// Package jekyll manages `jekyll serve` and `jekyll build` subprocesses.
package jekyll

import (
	"errors"
	"io"
	"os"
	"os/exec"
)

// Jekyll controls the lifecycle of jekyll subprocesses.
type Jekyll struct {
	serverCmd *exec.Cmd
	buildCmd  *exec.Cmd
}

// NewJekyll creates a new Jekyll instance with no running subprocesses.
func NewJekyll() (j *Jekyll) {
	j = &Jekyll{nil, nil}
	return
}

// startCmd starts cmd and relays its stdout/stderr to the parent process.
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

// startRubyJekyll spawns `jekyll serve` if not already running.
func (j *Jekyll) startRubyJekyll() (err error) {
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

// Start launches `jekyll serve` in the background.
func (j *Jekyll) Start() (err error) {
	return j.startRubyJekyll()
}

// Stop kills the running `jekyll serve` process, if any.
func (j *Jekyll) Stop() {
	if j.serverCmd == nil {
		return
	}
	j.serverCmd.Process.Kill()
	j.serverCmd = nil
}

// Build runs `jekyll build` synchronously. Only one build can run at a time.
func (j *Jekyll) Build() (err error) {
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
