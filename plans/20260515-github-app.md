# Description

Buang is a REST API which deploys docker compose services for short lived feature branches while they are open.
This facilitates collaboration by reducing friction when other members in a team need to verify changes which are being introduced.

Develop a GitHub application which deploys a feature branch when a pull request is opened in either draft or ready to review state
and tears down the deployment when the pull request is either closed or merged.

Each repository corresponds to a Buang project. Every time a deployment is created, it is created against the project.
So, a project must be created for all repositories the GitHub application has access to and all new projects after.
Creating a project requires a few details:
- The GitHub repository address
- Whether access to the repository requires authentication
- The username to use if authentication is required
- The password to use if authentication is required
- The path to the compose file which will be deployed

This means that our app will also need to notify the user that they must setup the repository before using 
our app by creating a check run with a link to the setup page. The user should then be able to follow the link and we should store the details they need.

Additionally, there are a few variables which are required for every deployment:
- BUANG_API_BASE_URL: The base URL (this is the same for all projects)
- BUANG_API_KEY: API key to use the Buang REST API (this is the same for all projects)
- BUANG_PROJECT_ID: The project id associated with this repository (this differs depending on the repository and should be stored by our bot)

# Tech stack

- Probot
- TypeScript
- PNPM
- Prettier
- ESLint with recommended rules (for both Probot and TypeScript along with any custom ones defined below which take precedence)
- HonoJS
   - Middlewares
      - oidc-auth
      - secure-headers
      - request-id
      - logger
      - timing
      - secure-headers
      - sentry
- React + Vite
   - shadcn/ui for components
   - React Hook Form for forms
   - yup to validate the form
   - React Query to make API requests
- PostgreSQL
   - node-postgres
      - To interact with our database
- austenite (https://github.com/ezzatron/austenite)
   - Validate environment variables required by the node app

# Scope
1. Gather repository and user details
- When the app is installed, create a check run with a link to the setup page for every repository it has access to
   - Visiting the setup page should allow the user to login and after logging in, they should be redirected to a page to setup their repository
   - If the user has not registered, first capture their login details (email and password) along with details of where their Buang server is hosted:
      - Buang API base URL
         - Verify that this URL is reachable by pinging it (ignore CORS, we just want to see if the host is reachable)
      - Buang API key
         - Do not store this in cleartext. Encrypt the Buang API key with a secret key defined by the environment variable (`BUANG_GHA_SECRET_KEY`) and decode it
   - Store the users details in a Postgres table
      - All database changes should be applied via migrations
- Listen for new repositories and repositories which have been made private:
   - For new repositories, create a check run with a link to the setup page
   - For repositories which were previously public but now private, create a check run with a link to update the details associated with that project
   - The setup link should include the repository URL so that when users are presented with the form to setup the repository,
     the input for the repository URL is autofilled
- Capture details associated with the repository
   - Display a form to capture:
      - The repository URL
         - Required
      - Whether the repository requires authentication
         - Required
      - If authentication is required, the username to authenticate with
         - Conditionally required
      - If authentication is required, the password to authenticate with
        Display helper text which reminds the user to use a Personal Access Token (PAT) instead of their real password
         - Conditionally required
      - Environment variables. This will be a list of key value pairs which are ideally represented by a table which allows rows to be added and deleted
         - Conditional
      - Wait for workflow run
         - Conditional
         - Some projects will require a build and push workflow to complete before it is ready to deploy. In these cases the deployment
           shouldn't happen immediately, it should only occur after the workflow run is successful
   - When the form is submitted, we should:
      - Create a Buang project:
         - POST `${BUANG_API_BASE_URL}/api/v1/project` with:
            - Headers:
               - X-API-KEY: `${BUANG_API_KEY}`
            - Body:
               - composePath: The path to the compose path to deploy
               - requiresAuthn: Whether the repository requires authentication
               - username: If authentication is required, the username to authenticate with
               - password: If authentication is required, the password to authenticate with
      - Store the environment variables for the Buang project
        - Store it as an encrypted JSON string in the database. Encrypt it with a secret key defined by the environment variable (`BUANG_GHA_SECRET_KEY`) and decode it when we need to use it while creating deployments
- Notes:
   - The implementation will probably require us to setup a Hono app and then use Probot as a middleware. The React application should also be deployed with the same Hono app
   - Plan appropriate tables to store user details and repository details
   - Users should also be able to authenticate via OAuth with the HonoJS `oidc-auth` middleware

2. Create a deployment
- For each repository, listen for either: when a new pull request is created or updated (either active or draft) OR when a workflow run is successful for the feature branch
- When a new pull request is created, we should make a request to the Buang API (for the user):
   - POST `${BUANG_API_BASE_URL}/api/v1/project/${BUANG_PROJECT_ID}/deployment?waitForDeployment=true` with:
      - Headers:
         - X-API-KEY: `${BUANG_API_KEY}`
      - Body:
         - branch: The feature branch ref
         - sha: The commit sha to deploy
         - serviceEntrypoint: The path to the compose file to deploy
         - env: Any environment variables for the project
            - The environment variables for the project are stored when the user sets up the project with the GitHub app. Remember that it is encrypted and can be decrypted with `BUANG_GHA_SECRET_KEY`

3. Tear down a deployment
- For each repository, listen for when a pull request is closed or merged
- When a pull request is closed or merged, we should make a request to the Buang API (for the user):
   - DELETE `${BUANG_API_BASE_URL}/api/v1/project/${BUANG_PROJECT_ID}/branch` with:
      - Headers:
         - X-API-KEY: `${BUANG_API_KEY}`
      - Body:
         - branch: The feature branch ref

# Verification
- Write unit tests where it makes sense to

# Related
- Implement this under `github-app`. Setup devcontainer for it with the Dockerfile:
   - ```
        FROM node:lts AS base

        RUN apt update && apt upgrade -y && \
            apt install vim bash-completion -y && \
            curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-$(dpkg --print-architecture).deb -o migrate.deb && \
            dpkg -i migrate.deb
        RUN git config --global pager.branch false
        RUN echo "if [ -f /etc/bash_completion ]; then\n\t. /etc/bash_completion\nfi" >> /etc/bash.bashrc

        FROM base

        RUN bash -c "corepack enable pnpm && pnpm config set store-dir /home/node/.local/share/pnpm/store"
     ```
   - Add a `migrate` script in the node project: `migrate -source file://$(pwd)/migrations/ -database $PG_CONNECTION_STRING up`
- Buang APIs can be found in:
   - Generally: `api/handlers`
   - POST `${BUANG_API_BASE_URL}/api/v1/project/${BUANG_PROJECT_ID}/deployment`: `api/handlers/deployments/create_deployment.go`
   - DELETE `${BUANG_API_BASE_URL}/api/v1/project/${BUANG_PROJECT_ID}/branch`: `api/handlers/deployments/buang_deployment.go`
- When designing the tables, it is important to note:
   - Users table:
      - Record if they have registered with an email and password or if they have used OAuth
      - DO NOT STORE PASSWORDS IN PLAIN TEXT
   - Repositories:
      - Keep track of the Buang project
      - Keep track of the environment variables for the repository which will be required for deployments
         - DO NOT STORE ENVIRONMENT VARIABLES IN PLAIN TEXT
         - Environment variables are stored as a JSON string
- Prettier configuration
   - ```
        {
            "plugins": ["prettier-plugin-organize-imports"],
            "tabWidth": 4,
            "useTabs": true,
            "semi": true,
            "singleQuote": false,
            "quoteProps": "as-needed",
            "trailingComma": "all",
            "bracketSpacing": true,
            "bracketSameLine": true,
            "arrowParens": "avoid",
            "endOfLine": "lf",
            "embeddedLanguageFormatting": "auto",
            "printWidth": 80
        }
     ```
- ESLint configuration:
   - ```
        {
            files: ["**/*.{ts,tsx}"],
            languageOptions: {
                ecmaVersion: "latest",
                sourceType: "module",
                parser: tseslint.parser,
                parserOptions: {
                    projectService: true,
                    tsconfigRootDir: import.meta.dirname,
                },
            },
            rules: {
                eqeqeq: "error",
                "no-duplicate-imports": "error",
                "import/no-namespace": "error",
                "import/no-default-export": "error",
                "@typescript-eslint/no-unused-vars": [
                    "error",
                    {
                        args: "all",
                        argsIgnorePattern: "^_",
                        caughtErrors: "all",
                        caughtErrorsIgnorePattern: "^_",
                        destructuredArrayIgnorePattern: "^_",
                        varsIgnorePattern: "^_",
                        ignoreRestSiblings: true,
                    },
                ],
                "@typescript-eslint/no-empty-object-type": "warn",
                "@typescript-eslint/no-explicit-any": "off",
                "import/no-unresolved": "off",
            },
        },
     ```
- Prefer anonymous functions over function declarations
- Use meaningful variable names but keep it terse
