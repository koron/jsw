package jekyll

import (
	"errors"
	"io"
	"os"
	"os/exec"
)

type Jekyll struct {
	serverCmd *exec.Cmd
	buildCmd  *exec.Cmd
}

func NewJekyll() (j *Jekyll) {
	j = &Jekyll{nil, nil}
	return
}

func startCmd(cmd *exec.Cmd) (err error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return
	}
	err = cmd.Start()
	if err != nil {
		return
	}
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	return
}

func (j *Jekyll) Start() (err error) {
	if j.serverCmd != nil {
		return errors.New("jekyll serve is running already.")
	}
	cmd := exec.Command("jekyll", "serve")
	err = startCmd(cmd)
	if err != nil {
		return
	}
	j.serverCmd = cmd
	return
}

func (j *Jekyll) Stop() {
	if j.serverCmd == nil {
		return
	}
	j.serverCmd.Process.Kill()
	j.serverCmd = nil
}

func (j *Jekyll) Build() (err error) {
	// FIXME:
	if j.buildCmd != nil {
		return errors.New("jekyll build is running already")
	}
	j.buildCmd = exec.Command("jekyll", "build")
	err = startCmd(j.buildCmd)
	if err != nil {
		j.buildCmd = nil
		return
	}
	err = j.buildCmd.Wait()
	j.buildCmd = nil
	return
}
