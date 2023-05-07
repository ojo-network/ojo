package orchestrator

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/ory/dockertest/v3/docker"
)

// ExecCommand executes a command on the validator container
// and outputs the stdout and stderr to the console
func (o *Orchestrator) ExecCommand(cmd []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	validator := o.chain.validators[1]

	exec, err := o.dkrPool.Client.CreateExec(docker.CreateExecOptions{
		Context:      ctx,
		AttachStdout: true,
		AttachStderr: true,
		Container:    validator.dockerResource.Container.ID,
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

	err = o.dkrPool.Client.StartExec(exec.ID, docker.StartExecOptions{
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
