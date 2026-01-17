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

*Buang* is a service to deploy any application in preview environments. Preview environments are temporary environments which allow you to test and validate your changes before merging it. Deploying APIs to preview environments usually requires managing complex infrastructure and K8s, with Buang all that's required to get started is a VM and a few GitHub actions.


# Development

1. Run the dev temporal server with:

```bash
task run-temporal-dev-server
```

2. Buang runs using an in memory SQLite database by default. Run it with:

```bash
task dev
```


# Documentation

Swagger documentation is available [here](http://buang.skulpture.xyz/docs/index.html)
