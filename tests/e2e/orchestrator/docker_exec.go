package orchestrator

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func execCommand(dkrPool *dockertest.Pool, val *validator, cmd []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	exec, err := dkrPool.Client.CreateExec(docker.CreateExecOptions{
		Context:      ctx,
		AttachStdout: true,
		AttachStderr: true,
		Container:    val.dockerResource.Container.ID,
		User:         "root",
		Cmd:          cmd,
	})
	if err != nil {
		return err
	}

	var (
		outBuf bytes.Buffer
		errBuf bytes.Buffer
	)

	err = dkrPool.Client.StartExec(exec.ID, docker.StartExecOptions{
		Context:      ctx,
		Detach:       false,
		OutputStream: &outBuf,
		ErrorStream:  &errBuf,
	})

	if err != nil {
		return fmt.Errorf("exec command failed; stdout: %s, stderr: %s", outBuf.String(), errBuf.String())
	}

	fmt.Println(errBuf.String())
	fmt.Println(outBuf.String())
	return nil
}
