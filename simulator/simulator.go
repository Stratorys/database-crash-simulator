package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	setupPostgres(ctx, cli)

	runGolangApp(ctx, cli)
}

func setupPostgres(ctx context.Context, cli *client.Client) {
	// Pull Postgres Image
	_, err := cli.ImagePull(ctx, "postgres:15.4-bookworm", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	// For setting up the exposed port
	exposedPort := nat.Port(os.Getenv("EXPOSED_PORT") + "/tcp")
	portSet := nat.PortSet{
		exposedPort: struct{}{},
	}

	// Set up the mount
	initSqlAbsolutePath, err := filepath.Abs("init.sql")
	if err != nil {
		log.Printf("cannot get absolute path for init.sql file")
	}
	mnt := mount.Mount{
		Type:   mount.TypeBind,
		Source: initSqlAbsolutePath, // Path to your SQL file on the host
		Target: "/docker-entrypoint-initdb.d/init.sql",
	}

	// Start Postgres Container
	pgContainer, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: "postgres:15",
			Env: []string{
				"POSTGRES_USER=admin",
				"POSTGRES_PASSWORD=admin",
				"POSTGRES_DB=postgres",
			},
			ExposedPorts: portSet,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{mnt},
		},
		nil,
		nil,
		"",
	)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, pgContainer.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}

func runGolangApp(ctx context.Context, cli *client.Client) {
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: "golang:1.21.3-bookworm",
			Cmd:   []string{"go", "run", "cmd/main.go"},
			Tty:   true,
		},
		nil,
		nil,
		nil,
		"",
	)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	fmt.Println("debug statusCh = ", <-statusCh)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case status := <-statusCh:
		if status.StatusCode != 0 { // Non-zero exit code implies it crashed
			fmt.Println("The Golang app container crashed. Checking logs...")

			// Fetch the logs to check the reason for crash
			out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
			if err != nil {
				panic(err)
			}
			var builder strings.Builder
			stdCopy, err := stdcopy.StdCopy(&builder, &builder, out)
			if err != nil {
				return
			}
			logs := builder.String()
			fmt.Println("debug stdCopy = ", stdCopy)

			// Check if the logs contain database connection issues
			if strings.Contains(logs, "database connection error") || strings.Contains(logs, "unable to connect to database") {
				fmt.Println("Detected database connection issue. Handling...")

				// Here you can add logic to check PostgreSQL health, restart it, etc.
				// Example: Restarting PostgreSQL container
				// cli.ContainerRestart(ctx, "YourPostgresContainerID", nil)
			}
		}
	}
}
