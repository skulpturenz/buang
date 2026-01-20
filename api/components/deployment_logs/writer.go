package deploymentlogs

import (
	"bufio"
	"context"
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

// reason for intercepting this call although we add new lines when we upsert logs
// is that when it is buffered everything gets flushed at once so there are no new lines in this case
// since we're adding a new line here along with adding one at the query we might have two new lines where there should be one
// TODO: better way?
func (nw bufferedLogWriter) Write(p []byte) (n int, err error) {
	wrote := 0

	n, err = nw.bw.Write(p)
	if err != nil {
		return n, err
	}
	wrote += len(p)

	nl := []byte("\n")
	_, err = nw.bw.Write(nl)
	if err != nil {
		return wrote, err
	}
	wrote += len(nl)

	return wrote, err
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
