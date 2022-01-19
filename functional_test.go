// +build functional

package main

import (
	"fmt"
	"testing"
	"context"

	dockertest "github.com/ory/dockertest/v3"
	docker "github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

func TestVwapEngine(t *testing.T) {
	pool, err := dockertest.NewPool("")
	require.NoError(t, err, "could not connect to Docker")

	container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "zero-hash-vwap",
		Tag:        "v1.0.0",
	})
	fmt.Printf("%v", container.Container.ID)
	opts := docker.LogsOptions{
		Context: context.TODO(),
		Stderr:      true,
		Stdout:      true,
		Follow:      false,
		Timestamps:  true,
		RawTerminal: true,
		Container: container.Container.ID,
	}
	fmt.Printf("%v", pool.Client.Logs(opts))
	require.NoError(t, err, "could not start container")

	t.Cleanup(func() {
		require.NoError(t, pool.Purge(container), "failed to remove container")
	})
}