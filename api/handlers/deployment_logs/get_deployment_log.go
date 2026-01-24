package deploymentlogs

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	deploymentlogs "skulpture/buang/components/deployment_logs"
	"skulpture/buang/components/deployments"
	dberrors "skulpture/buang/db/db_errors"
	enumsdeploymentstatus "skulpture/buang/enums/deployment_status"
	workflowwrappers "skulpture/buang/workers/workflow_wrappers"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/negrel/assert"
)

type GetDeploymentLogRequest struct {
	ProjectId    int64 `schema:"-"`
	DeploymentId int64 `schema:"-"`
	Stream       *bool `schema:"stream,default:true"`
}

type GetDeploymentLogResult = string

var errStreamingUnsupported = errors.New("streaming unsupported")

// @summary	Get deployment logs
// @tags		api.v1, deployment-logs
// @security	ApiKeyAuth
// @param		projectId		path		int		required	"Project ID"
// @param		deploymentId	path		int		required	"Deployment ID"
// @param		stream			query		bool	false		"Stream logs"
// @success	200				{object}	string
// @success	204				{object}	nil
// @failure	400				{object}	any
// @failure	401
// @failure	500	{object}	string
// @router		/project/{projectId}/deployment/{deploymentId}/logs [get]
func GetDeploymentLog(s app.ApplicationServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req GetDeploymentLogRequest

		projectIdParam := chi.URLParam(r, "projectId")
		projectId, err := strconv.Atoi(projectIdParam)
		if err != nil {
			slog.ErrorContext(r.Context(), "get deployment logs", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.ProjectId = int64(projectId)

		deploymentIdParam := chi.URLParam(r, "deploymentId")
		deploymentId, err := strconv.Atoi(deploymentIdParam)
		if err != nil {
			slog.ErrorContext(r.Context(), "get deployment logs", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.DeploymentId = int64(deploymentId)

		err = s.GorillaSchemaDecoder.Decode(&req, r.URL.Query())
		if err != nil {
			slog.ErrorContext(r.Context(), "get deployment logs", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		d := deployments.FindDeploymentByIdParams{
			ProjectId: req.ProjectId,
			ID:        req.DeploymentId,
		}

		dply, dplyErr := d.Exec(r.Context(), &s)
		if dplyErr != nil && !dberrors.IsNoRows(dplyErr) {
			slog.ErrorContext(r.Context(), "get deployment logs", "err", dplyErr.Error())
			http.Error(w, dplyErr.Error(), http.StatusInternalServerError)
			return
		}

		p := deploymentlogs.GetDeploymentLogParams{
			ProjectId:    int64(projectId),
			DeploymentId: int64(deploymentId),
		}

		res, err := p.Exec(r.Context(), &s)
		if err != nil && !dberrors.IsNoRows(err) {
			slog.ErrorContext(r.Context(), "get deployment logs", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		isDeploying := func() bool {
			if dplyErr != nil {
				return dberrors.IsNoRows(dplyErr)
			}

			if err != nil {
				return dberrors.IsNoRows(err)
			}

			if dply.Deployment.GetStatus() == int16(enumsdeploymentstatus.New) {
				return true
			}

			if dply.Deployment.GetStatus() == int16(enumsdeploymentstatus.Deploying) {
				return true
			}

			return false
		}

		var streamErr error
		if isDeploying() && *req.Stream {
			streamErr = streamDeploymentLogs(r.Context(), s, w, req)

			return
		}

		if isDeploying() && (!*req.Stream || errors.Is(streamErr, errStreamingUnsupported)) {
			const WAIT_TIMEOUT = 5 * time.Minute
			waitCtx, cancelWaitCtx := context.WithTimeout(r.Context(), WAIT_TIMEOUT)

			watchForDeployment(waitCtx, cancelWaitCtx, s, req)

			if errors.Is(context.DeadlineExceeded, waitCtx.Err()) {
				slog.ErrorContext(r.Context(), "get deployment logs", "err", waitCtx.Err().Error())
				http.Error(w, waitCtx.Err().Error(), http.StatusInternalServerError)

				return
			}

			res, err = p.Exec(r.Context(), &s)
			if err != nil {
				slog.ErrorContext(r.Context(), "get deployment logs", "err", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		assert.True(dply.Deployment.GetStatus() == int16(enumsdeploymentstatus.Deployed))

		if res.Logs == nil {
			w.WriteHeader(http.StatusNoContent)

			return
		}

		err = app.WriteText(w, *res.Logs, http.StatusOK)
		if err != nil {
			slog.ErrorContext(r.Context(), "get deployment logs", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func streamDeploymentLogs(ctx context.Context, s app.ApplicationServices, w http.ResponseWriter, req GetDeploymentLogRequest) error {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return errStreamingUnsupported
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	p := deploymentlogs.PollDeploymentLogParams{
		ProjectId:    int64(req.ProjectId),
		DeploymentId: int64(req.DeploymentId),
	}

	const STREAM_TIMEOUT = 5 * time.Minute

	streamCtx, cancelStreamCtx := context.WithTimeout(ctx, STREAM_TIMEOUT)

	previousLogs := ""

	go watchForDeployment(streamCtx, cancelStreamCtx, s, req)

	for {
		select {
		case <-streamCtx.Done():
			if errors.Is(context.DeadlineExceeded, streamCtx.Err()) && previousLogs == "" {
				slog.ErrorContext(streamCtx, "get deployment logs", "err", streamCtx.Err().Error())
				http.Error(w, streamCtx.Err().Error(), http.StatusInternalServerError)

				return streamCtx.Err()
			}

			// won't write to body if logs are empty because len(currentLogs) == len(previousLogs)
			// so set correct headers
			if previousLogs == "" {
				w.WriteHeader(http.StatusNoContent)

				return nil
			}

			return nil
		case log, ok := <-p.Exec(streamCtx, &s):
			if !ok { // poll deployment logs result channel closed
				cancelStreamCtx()
			} else {
				newLogs := log.Logs

				cl := len(newLogs)
				pl := len(previousLogs)

				assert.True(cl >= pl, "expected deployment log to be append only")

				if cl > pl {
					app.WriteText(w, newLogs[pl:], http.StatusOK) // why we expect it to be append only: so that we can skip the diffing and just slice it
					flusher.Flush()

					previousLogs = newLogs
				}
			}
		}
	}
}

func watchForDeployment(ctx context.Context, cancelCtx context.CancelFunc, s app.ApplicationServices, req GetDeploymentLogRequest) {
	const HAS_DEPLOYED_POLL_INTERVAL = 50 * time.Millisecond

	ticker := time.NewTicker(HAS_DEPLOYED_POLL_INTERVAL)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			wp := workflowwrappers.HasDeployedParams{
				ProjectId:    req.ProjectId,
				DeploymentId: req.DeploymentId,
			}

			err := wp.Exec(ctx, s)
			if err == nil {
				cancelCtx()
				return
			}
		}
	}
}
