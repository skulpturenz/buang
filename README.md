STATUS: WIP

# Get started

1. Deploy the compose services in `api/docker-compose.yml`

2. Run this CURL command using Postman or a Terminal to create a repository for your project

```bash
curl --request POST \
  --url http://localhost:65277/api/v1/project \
  --header 'content-type: application/json' \
  --header 'x-api-key: supersecureapikey' \
  --data '{
  "composePath": "<location of the compose file from the root of your repository>",
  "password": "<repository authentication (if any)>",
  "repository": "<github repository http url>",
  "requiresAuthn": true, // if any
  "username": "<repository authentication (if any)>"
}'
```

3. Run this CURL command using Postman or a Terminal to create a preview deployment

```bash
curl --request POST \
  --url http://localhost:65277/api/v1/project/<projectId>/deployment \
  --header 'content-type: application/json' \
  --header 'x-api-key: supersecureapikey' \
  --data '{
  "branch": "<the preview branch name>",
  "env": {
    "HELLO": "WORLD",
    // ... any other env variables for your service
  },
  "serviceEntrypoint": "<hostname:port of the service which listens for outside connections>",
  "sha": "<the commit hash to deploy>"
}'
```

4. Services deployed! 🚀

*Note: The compose services you'd like to deploy should use the `buang` network


# About

*Buang* is a service to deploy any application in preview environments.
Preview environments are temporary environments which allow you to test and validate your changes before merging it. 
Deploying APIs to preview environments usually requires managing complex infrastructure and K8s, with Buang all that's required to get started is a VM and a few GitHub actions.

Buang deploys services using Docker Compose and dynamically updates a Traefik instance to route traffic to the preview services.


# Requirements

Buang aims to be a deploy and forget service which is fairly low maintenance, we don't want to spend more time fixing issues with the preview server than developing.
To this end, Buang employs the use of durable executors such as [Temporal](https://temporal.io/) or [DBOS](https://docs.dbos.dev/) so that it is resilient to most failures. Durable executors make it convenient for a service to recover from failure from its last successful point. Buang also prunes unused containers and images periodically to ensure the system does not run out of storage.

Using Buang with Temporal requires a Temporal deployment by either using [Temporal Cloud](https://temporal.io/cloud) or [self-hosting](https://docs.temporal.io/self-hosted-guide). Temporal can be used with either Postgres or SQLite. Buang runs its Temporal workers on a goroutine each to keep things simple instead of separate processes.

DBOS runs in-process so all that's required to use Buang is a Postgres database. DBOS does not support SQLite.

Buang runs using an in memory SQLite database by default with Temporal.


# Development

## Temporal

1. Run the dev temporal server with:

```bash
task run-temporal-dev-server
```

2. Buang runs using an in memory SQLite database by default. Run it with:

```bash
task dev
```

## DBOS

1. Run the dev Postgres server with:

```bash
task run-pg
```

2. Set the correct environment variables

```bash
export DB_CONNECTION_STRING="postgresql://postgres:mysecretpassword@localhost/buang?sslmode=disable"
export DB_TYPE="postgres"
export DURABLE_EXECUTOR="dbos"
```

3. Run Buang:

```bash
task dev
```

# Documentation

Swagger documentation is available [here](http://buang.skulpture.xyz/docs/index.html)


# Demo

https://github.com/user-attachments/assets/6237db24-73f1-4df1-83d1-034b6f5f899c


# Limitations & tradeoffs

Buang is designed for single tenancy, using Docker Compose to do the heavy lifting means that the only way to scale is vertically.
As it is meant to be hosted on your own infrastructure for temporary deployments,
it is also light on security: tokens for private repository access are stored in plain text and any environment variables required for the deployment are sent in the request to deploy services.

Ensure that any private repository access tokens are readonly with limited scope and assume that the preview environment will be compromised.

Each deployment gets a unique path instead of a subdomain to minimize the configuration required to deploy Buang, so if the service depends on the base URL, it can be derived by checking the `X-Forwarded-Prefix` header.
