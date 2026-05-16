package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"dagger.io/dagger"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./ci <lint|test|build>")
		os.Exit(1)
	}

	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to Dagger: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = client.Close() }()

	src := client.Host().Directory("..", dagger.HostDirectoryOpts{Exclude: []string{".git", "dist"}})

	switch os.Args[1] {
	case "lint":
		err = lint(ctx, client, src)
	case "test":
		err = test(ctx, client, src)
	case "build":
		err = build(ctx, client, src)
	default:
		err = fmt.Errorf("unknown command %q", os.Args[1])
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func baseContainer(client *dagger.Client, src *dagger.Directory) *dagger.Container {
	return client.Container().
		From("golang:1.24-alpine").
		WithEnvVariable("CGO_ENABLED", "0").
		WithMountedDirectory("/src", src).
		WithWorkdir("/src")
}

func lint(ctx context.Context, client *dagger.Client, src *dagger.Directory) error {
	_, err := baseContainer(client, src).
		WithExec([]string{"sh", "-c", "test -z \"$(gofmt -l .)\""}).
		WithExec([]string{"go", "vet", "./..."}).
		Sync(ctx)
	if err != nil {
		return fmt.Errorf("lint failed: %w", err)
	}

	return nil
}

func test(ctx context.Context, client *dagger.Client, src *dagger.Directory) error {
	_, err := baseContainer(client, src).
		WithExec([]string{"go", "test", "./..."}).
		Sync(ctx)
	if err != nil {
		return fmt.Errorf("tests failed: %w", err)
	}

	return nil
}

func build(ctx context.Context, client *dagger.Client, src *dagger.Directory) error {
	platforms := []string{"linux/amd64", "linux/arm64"}
	artifacts := client.Directory()

	for _, platform := range platforms {
		goos, goarch, found := strings.Cut(platform, "/")
		if !found {
			return fmt.Errorf("invalid platform %q", platform)
		}

		container := baseContainer(client, src).
			WithEnvVariable("GOOS", goos).
			WithEnvVariable("GOARCH", goarch).
			WithExec([]string{"go", "build", "-trimpath", "-ldflags=-s -w", "-o", "/tmp/get-next-ip", "./cmd/get-next-ip"})

		fileName := fmt.Sprintf("get-next-ip-%s", platformToArtifactSuffix(platform))
		artifacts = artifacts.WithFile(fileName, container.File("/tmp/get-next-ip"))
	}

	if _, err := artifacts.Export(ctx, "../dist"); err != nil {
		return fmt.Errorf("failed to export build artifacts: %w", err)
	}

	return nil
}

func platformToArtifactSuffix(platform string) string {
	switch platform {
	case "linux/amd64":
		return "linux-amd64"
	case "linux/arm64":
		return "linux-arm64"
	default:
		return platform
	}
}
