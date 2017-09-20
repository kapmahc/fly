package app

import (
	"math/rand"
	"os"
	"os/exec"
	"syscall"
	"time"
)

//Shell exec shell command
func Shell(cmd string, args ...string) error {
	bin, err := exec.LookPath(cmd)
	if err != nil {
		return err
	}
	return syscall.Exec(bin, append([]string{cmd}, args...), os.Environ())
}

//Random randome string
func Random(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	buf := make([]rune, n)
	for i := range buf {
		buf[i] = letters[rd.Intn(len(letters))]
	}
	return string(buf)
}
