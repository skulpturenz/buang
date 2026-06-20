package main

//go:generate go tool sqlc generate
//go:generate go tool swag init -d ./,./handlers/projects
//go:generate go tool swag fmt
