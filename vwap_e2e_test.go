// +build e2e

package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"regexp"
	"testing"
	"time"

	dockertest "github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVwapEngine(t *testing.T) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		require.FailNow(t, fmt.Sprintf("could not connect to Docker: %v", err))
	}

	container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "zero-hash-vwap",
		Tag:        "v1.0.0",
		Cmd:        []string{"BTC-USD"},
	})
	if err != nil {
		require.FailNow(t, fmt.Sprintf("could not start container: %v", err))
	}

	containerLogs := readContainerLogs(pool, container)
	matched, err := regexp.MatchString(`trade_pair=BTC-USD vwap=[0-9]+`, containerLogs.String())
	assert.Nil(t, err)
	assert.True(t, matched)

	t.Cleanup(func() {
		require.NoError(t, pool.Purge(container), "failed to remove container")
	})
}

func readContainerLogs(pool *dockertest.Pool, container *dockertest.Resource) bytes.Buffer {
	var b bytes.Buffer
	logsWriter := io.Writer(&b)
	opts := docker.LogsOptions{
		Context:      context.TODO(),
		Stderr:       true,
		Stdout:       true,
		Follow:       true,
		Timestamps:   true,
		RawTerminal:  true,
		Container:    container.Container.ID,
		OutputStream: logsWriter,
	}
	go pool.Client.Logs(opts)
	time.Sleep(5 * time.Second)
	return b
}
