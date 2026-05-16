package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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

	sourceRoot, err := getSourceRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	src := client.Host().Directory(sourceRoot, dagger.HostDirectoryOpts{Exclude: []string{".git", "dist"}})

	switch os.Args[1] {
	case "lint":
		err = lint(ctx, client, src)
	case "test":
		err = test(ctx, client, src)
	case "build":
		err = build(ctx, client, src, filepath.Join(sourceRoot, "dist"))
	default:
		err = fmt.Errorf("unknown command %q", os.Args[1])
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getSourceRoot() (string, error) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to resolve current directory: %w", err)
	}

	return filepath.Clean(filepath.Join(workingDirectory, "..")), nil
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

func build(ctx context.Context, client *dagger.Client, src *dagger.Directory, distPath string) error {
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
			WithExec([]string{"go", "build", "-trimpath", "-ldflags=-s -w", "-o", "/tmp/next-ip", "./cmd/next-ip"})

		fileName := fmt.Sprintf("next-ip-%s", platformToArtifactSuffix(platform))
		artifacts = artifacts.WithFile(fileName, container.File("/tmp/next-ip"))
	}

	if _, err := artifacts.Export(ctx, filepath.Clean(distPath)); err != nil {
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
