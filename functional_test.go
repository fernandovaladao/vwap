package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"context"
	"time"

	dockertest "github.com/ory/dockertest/v3"
	docker "github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

func TestVwapEngine(t *testing.T) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		require.FailNow(t, fmt.Sprintf("could not connect to Docker: %v", err))
	}

	container, err := pool.RunWithOptions(&dockertest.RunOptions{Repository: "zero-hash-vwap", Tag: "v1.0.0", Cmd: []string{"BTC-USD"}})
	if err != nil {
		require.FailNow(t, fmt.Sprintf("could not start container: %v", err))
	}

	// start streaming logs
	containerLogs := streamLogs(pool, container)
	fmt.Printf("Logs are %s", containerLogs)
	
	t.Cleanup(func() {
		require.NoError(t, pool.Purge(container), "failed to remove container")
	})
} 

func streamLogs(pool *dockertest.Pool, container *dockertest.Resource) string {
	var b bytes.Buffer
	logsWriter := io.Writer(&b)
	opts := docker.LogsOptions{
		Context: context.TODO(),
		Stderr:      true,
		Stdout:      true,
		Follow:      true,
		Timestamps:  true,
		RawTerminal: true,
		Container: container.Container.ID,
		OutputStream: logsWriter,
	}
	go pool.Client.Logs(opts)
	time.Sleep(10*time.Second)
	return b.String()
}