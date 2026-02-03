## [1.0.0-canary.0] - 2026-02-03

### 🚀 Features

- Add queries to create, select and update deployments
- Add deployments component
- Add routes for deployments
- Add docker, fs and git components
- Add workers, activities, update schema
- Add git authn, update taskfile, other small fixes
- Buang all active deployments when a new one occurs
- List deployments and projects req validation
- Housekeeping worker
- Add buang branch endpoint, register housekeeping worker
- Add endpoint to find project by repository, enforce only one active repository at a time, add github actions examples (untested)
- Add FindDeploymentById endpoint, update actions to leave a link to the deployed environment
- Docker prune every 2 days
- Add support for dbos
- Recover from dbos workflow failures
- Update deployment status to error after recovering from panic
- More panic handling
- Allow connecting to dbos conductor
- Deployment logs
- Allow blocking until deployment completes
- Get deployment logs endpoint
- Allow waiting for deployment to complete when retrieving logs
- Git clone logs
- Provide BUANG_DEPLOYMENT_PATH when composing up
- Add queries DeleteProject and UpdateProject
- Add diagnostic log table
- Bootstrap
- Capture all panics from workflows
- Prune diagnostic logs
- Parse stack
- Enable sticky cookies
- DeleteProject, UpdateProject
- Stream deployment logs
- Fallback to waiting for deployment to complete if streaming is unsupported
- Docker stats
- Sort docker stats
- Add version route
- Add image and image id to stats response
- Enable assertions in dev, ensure all singletons are initialized
- Return container status and state for stats, GET /diagnostics/version -> GET /diagnostics/info - add tz and uptime information
- Update compose up build options to push, pass writer for build progress
- Add fields for containers in DurableExecutorConfiguration

### 🐛 Bug Fixes

- Issues with sqlite, issues with SelectProjectsDesc
- Add req validations, buang all active deployments by branch, schema updates, component updates
- Issues with deployment and validation
- Update ServiceEntrypoint validation
- Setting env vars
- Add buang_branch, fix deployment issues
- Update clone configuration
- Pg SelectProjectsDesc
- Run housekeeping worker
- Dbos issues
- Housekeeping workflow DI
- Panic recovery and improve
- Handle compose down active deployments errors better. spamming the execute button 60+ times surfaces it. services are still running but compose file is gone: ostrich, find really old services and kill them or maybe not using temp dirs is more suitable
- Housekeeping workflow nil ptr, concurrent deployments breaking 1 active deployment constraint
- Status code
- 1 active deployment per branch for every project
- Housekeeping temporal workflow
- Stringify stack before saving to db, fix issue with sqlite CreateDeployment
- Buang_deployment activity panic if clonePath is nil
- Nil envs
- Handle nil map cases at db wrappers
- DeleteProject, UpdateProject swagger annotations, update create_deployment query to check if project has been deleted
- GetDeploymentLogs returns too early when deployment is already in progress but not completed
- Update isDeploying check
- Deployment router paths
- On delete cascade
- Update cpu usage calculation
- MemUsage calculation, used total_inactive_file vs inactive_file
- Sqlite in memory getting cleared
- Compose logs not included with temporal
- Retries and handle context timeouts
- Increase step retries

### 🚜 Refactor

- Update deployment routes, int temporal and http
- Assertions, update Taskfile
- Don't treat panics as errors, the application is setup wrong if it panics. iirc if compose tries to open a directory and it does not exist then we have a panic, but i think we shouldn't consider it as a deployment failure because the data is right, the issue is that we're using /tmp for files which are temporary but expected to be there because we don't try to retrieve the files again if its not once we have initially done it
- More err handling
- Sqlc named params for pg, update queries UpdateDeployment and UpdateDeploymentStatus to check for project id
- Rename sqlite constraints
- Envs
- Standardize diagnostic logs a little
- Dberrors
- Tidy GetDeploymentLog and PollDeploymentLog
- Simplify and standardize workflow panic recovery
- Reorganize handlers
- Update swagger grouping
- Collect
- Update DI - inject docker client, inject gorilla schema from main
- Scope envs, allow configuring traefik ports
- Make durable executor fields in ApplicationServices private
- Replace conditionals with polymorphism
- Ensure clients implement the Workflows interface
- Improve sagas for concurrency
- Improve DI for services which don't have global scope

### 📚 Documentation

- Swagger annotations
- Update CreateDeployment swagger
- Add TODO

### 🧪 Testing

- New project with deployment
- Check sha and branch of checkout
- Check compose service after deployment
- Check envs are set correctly
- Listing projects, deployments, paging, buang branch
- Docker stats. and fix docker stats sorting, we were sorting it in ascending and when we returned the response we were iterating from newest so we end up with a response which was in descending. also update sorting tests in list_deployments_rest and list_projects_test
- Supress testcontainer logs
- Run tests in parallel, 45s max, under 30 most of the time. slightly flaky but more when runnign in watch mode. there are orphan containers
- Retry failing tests
- Add TemporalPg test config
- Assert busiest and idlest containers
- Remove time.Sleep, better logging for retry
- TemporalPg for TestListProjectsExcludesDeleted, TestSearchParams and TestDeletedProjects
- Assert proxying, fix issue with compensations not running, increase timeouts
- Cleanup deployments, refactor compensations
- Remove timeouts
- Basic fitness tests
- Fitness test assert proxy
- Add deploy workflow failure test when trying to compose up, fix issues with deploy workflow and add TODOs, refactor to allow proxying per request services
- Buang deployment errors, fix issues with dbos workflow, fix buang deployment activity return

### ⚙️ Miscellaneous Tasks

- Initial commit
- Setup project
- Db models
- Setup db connections, di, slog, refactor
- Update schema
- Add Taskfile
- Update swagger endpoint security, update bundled migrations
- Task generate-all
- Add infrastructure
- Update README
- Add demo to README
- Update README
- Add deploy-prod task and prod compose file
- Update README
- Task generate-all
- Update README
- Add configuration.md docs
- Add git-chglog
- Update docker compose to add health check and update volume mount points, still untested
- Panic if BUANG_VERSION is not set in production
- Docker in docker for buang
- Update Dockerfile
- Setup utils for integration tests
- Update taskfile
- Update logger configuration
- Fix Dockerfile
- Add test github action workflow
- Update gh actions workflow to cache go package downloads by default
- Fix cache dependency path
- Authn github test-services action
- Authenticate install task step in test services gh workflow
- Migrate to git-cliff
- Add release workflow
