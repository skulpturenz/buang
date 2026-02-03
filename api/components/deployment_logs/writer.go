package deploymentlogs

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"skulpture/buang/app"
	"skulpture/buang/db/interfaces"
)

type DeploymentLogWriterParams struct {
	logs         []byte
	services     *app.ApplicationServices
	ProjectID    int64
	DeploymentID int64
}

type DeploymentLogWriteResult struct{}

func (p *DeploymentLogWriterParams) Exec(ctx context.Context, s *app.ApplicationServices) (*DeploymentLogWriteResult, error) {
	p.services = s
	q := *s.Queries

	_, err := q.UpsertDeploymentLog(ctx, interfaces.UpsertDeploymentLogParams{
		ProjectID:    p.ProjectID,
		DeploymentID: p.DeploymentID,
	})
	if err != nil {
		return nil, err
	}

	ret := DeploymentLogWriteResult{}

	return &ret, nil
}

func (p DeploymentLogWriterParams) Write(b []byte) (n int, err error) {
	q := *p.services.Queries

	log := string(b)
	_, err = q.UpsertDeploymentLog(context.Background(), interfaces.UpsertDeploymentLogParams{
		ProjectID:    p.ProjectID,
		DeploymentID: p.DeploymentID,
		Log:          &log,
	})
	if err != nil {
		return 0, err
	}

	return len(b), nil
}

type bufferedLogWriter struct {
	bw *bufio.Writer
}

func (nw bufferedLogWriter) Write(p []byte) (n int, err error) {
	n, err = fmt.Fprintf(nw.bw, "%v\n", string(p))
	if err != nil {
		return n, err
	}

	return n, nil
}

func (nw *bufferedLogWriter) Flush() error {
	return nw.bw.Flush()
}

func CreateBufferedWriter(w io.Writer) bufferedLogWriter {
	writer := bufferedLogWriter{
		bw: bufio.NewWriter(w),
	}

	return writer
}
