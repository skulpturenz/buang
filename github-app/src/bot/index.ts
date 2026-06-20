import { type Probot } from "probot";
import { registerInstallationHandlers } from "./installation.js";
import { registerPullRequestHandlers } from "./pull-request.js";
import { registerRepositoryHandlers } from "./repository.js";
import { registerWorkflowRunHandlers } from "./workflow-run.js";

export const probotApp = (app: Probot) => {
	registerInstallationHandlers(app);
	registerRepositoryHandlers(app);
	registerPullRequestHandlers(app);
	registerWorkflowRunHandlers(app);
};
