package utils

import (
	"os"
	"time"
	"errors"
	"strings"
	"os/exec"
)
type Result struct {
	data []byte
	err  error
}

func myExec(c *exec.Cmd, ch chan bool, result *Result) {
	result.data, result.err = c.CombinedOutput()
	ch <- true
}
func BashCommand(cmdStr string, timeout int64) (string, error) {
	t := time.NewTimer(time.Duration(timeout) * time.Second)
	done := make(chan bool)

	c := exec.Command("/bin/bash", "-c", cmdStr)

	result := &Result{data: []byte{}, err: nil}

	go myExec(c, done, result)

	select {
	case <-done:
		return string(result.data), result.err
	case <-t.C:
		c.Process.Kill()
		result.err = errors.New("timeout")
		return string(result.data), result.err
	}
}

func BashComandNew(cmdPath,cmdStr string,timeout int64) (string, error)  {
	t := time.NewTimer(time.Duration(timeout) * time.Second)
	done := make(chan bool)

	c := exec.Command(cmdPath, "-c", cmdStr)

	result := &Result{data: []byte{}, err: nil}

	go myExec(c,done,result)

	select {
	case <-done:
		return string(result.data), result.err
	case <-t.C:
		c.Process.Kill()
		result.err = errors.New("timeout")
		return string(result.data), result.err
	}
}

func Command(cmd string, timeout int64, args []string) (string, error) {
	c := exec.Command(cmd, args...)
	t := time.NewTimer(time.Duration(timeout) * time.Second)
	done := make(chan bool)
	result := &Result{data: []byte{}, err: nil}

	go myExec(c, done, result)
	select {
	case <-done:
		return string(result.data), result.err
	case <-t.C:
		c.Process.Kill()
		result.err = errors.New("timeout")
		return string(result.data), result.err
	}
}

func StartCommand(cmdPath, cmdStr string) error {
	procAttr := &os.ProcAttr{
		Dir: cmdPath,
		Env: os.Environ(),
		Files: []*os.File{os.Stdin, nil, nil},
	}

	var params []string
	ss := strings.Split(cmdStr, " ")
	for _, s := range ss {
		params = append(params, s)
	}

	_, err := os.StartProcess(params[0], params, procAttr)

	if err != nil {
		return err
	}

	return nil
}
