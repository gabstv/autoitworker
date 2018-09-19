package embedded

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

//go:generate bin2var -I au3worker.exe

// Exe contains the paths and the current process information (of au3worker.exe)
type Exe struct {
	TempDir     string
	TempExePath string
	Cmd         *exec.Cmd
}

// Remove temporary files, folders and running programs
func (ee *Exe) Remove() {
	if ee.Cmd != nil && ee.Cmd.Process != nil {
		ee.Cmd.Process.Kill()
	}
	os.RemoveAll(ee.TempDir)
}

// IsRunning checks if the process is still running
func (ee *Exe) IsRunning() bool {
	if ee.Cmd == nil {
		return false
	}
	if ee.Cmd.Process == nil {
		return false
	}
	_, err := os.FindProcess(ee.Cmd.Process.Pid)
	return err == nil
}

// Open creates a temporary directory, extracts the autoitworker executable
// and runs it.
func Open(httppath string) (*Exe, error) {
	dirpath, err := ioutil.TempDir("", "autoitworker")
	if err != nil {
		return nil, err
	}
	ee := &Exe{}
	ee.TempDir = dirpath
	//
	filepath := path.Join(dirpath, "au3worker.exe")
	tfw, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		os.RemoveAll(dirpath)
		return nil, fmt.Errorf("os.OpenFile: %s", err.Error())
	}
	buf := bytes.NewBuffer(auWorkerExe)
	_, err = io.Copy(tfw, buf)
	if err != nil {
		tfw.Close()
		os.RemoveAll(dirpath)
		return nil, err
	}
	buf.Reset()
	tfw.Close()
	ee.TempExePath = filepath
	cmdd := exec.Command(filepath, httppath)
	//cmdd.Stdout = os.Stdout
	ee.Cmd = cmdd
	err = ee.Cmd.Start()
	if err != nil {
		os.RemoveAll(dirpath)
		return nil, err
	}
	return ee, nil
}
