package tread

import (
	"os"
	"os/exec"
	"io"
)

func Exec(r io.Reader, cmd string, args ...string) (io.Reader, error) {
	ecmd := exec.Command(cmd, args...)
	ecmd.Stdin = r
	ecmd.Stderr = os.Stderr
	out, err := ecmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	ecmd.Start()
	return out, nil
}

func MustExec(r io.Reader, cmd string, args ...string) io.Reader {
	out, err := Exec(r, cmd, args...)
	if err != nil {
		panic(err)
	}
	return out
}
