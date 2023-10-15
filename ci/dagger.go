package main

import (
	"context"
	"dagger.io/dagger"
	"fmt"
	"os"
	"strconv"
)

func main() {
	ctx := context.Background()

	// initialize Dagger client
	client, err := dagger.Connect(ctx,
		dagger.WithLogOutput(os.Stdout),
	)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	sqlInit := client.Host().Directory(".").File("init.sql")
	src := client.Host().Directory(".")
	pgPort := 5444

	database := client.Container().
		From("postgres:15.4-bookworm").
		WithEnvVariable("POSTGRES_USER", "admin").
		WithEnvVariable("POSTGRES_DB", "postgres").
		WithEnvVariable("POSTGRES_PASSWORD", "admin").
		WithMountedFile("/docker-entrypoint-initdb.d/init.sql", sqlInit).
		WithExec([]string{"postgres"}).
		WithExposedPort(5432)

	_, err = database.Endpoint(ctx, dagger.ContainerEndpointOpts{Port: pgPort})
	if err != nil {
		panic(err.Error())
	}

	outputWithCrash, err := client.Container().
		From("golang:1.21.3-bookworm").
		WithServiceBinding("postgres", database).
		WithEnvVariable("DB_HOST", "postgres").
		WithEnvVariable("DB_PORT", strconv.Itoa(pgPort)).
		WithEnvVariable("DB_PASSWORD", "admin").
		WithEnvVariable("DB_USER", "briskport_user").
		WithEnvVariable("DB_NAME", "postgres").
		WithEnvVariable("DB_SCHEMA", "briskport").
		WithDirectory("/src", src, dagger.ContainerWithDirectoryOpts{
			Exclude: []string{"ci/"},
		}).
		WithWorkdir("/src").
		WithExec([]string{"go", "run", "cmd/main.go", "--must-panic"}).
		Stdout(ctx)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("debug outputWithCrash = ", outputWithCrash)

	outputWithoutCrash, err := client.Container().
		From("golang:1.21.3-bookworm").
		WithServiceBinding("postgres", database).
		WithEnvVariable("DB_HOST", "postgres").
		WithEnvVariable("DB_PORT", strconv.Itoa(pgPort)).
		WithEnvVariable("DB_PASSWORD", "admin").
		WithEnvVariable("DB_USER", "briskport_user").
		WithEnvVariable("DB_NAME", "postgres").
		WithEnvVariable("DB_SCHEMA", "briskport").
		WithDirectory("/src", src, dagger.ContainerWithDirectoryOpts{
			Exclude: []string{"ci/"},
		}).
		WithWorkdir("/src").
		WithExec([]string{"go", "run", "cmd/main.go"}).
		Stdout(ctx)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("debug outputWithoutCrash = ", outputWithoutCrash)
}
