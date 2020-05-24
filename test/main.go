package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/icrowley/fake"
)

type StdOutListener struct {
	buf []byte
}

func (s *StdOutListener) Write(p []byte) (n int, err error) {
	s.buf = append(s.buf, p...)
	return len(p), nil
}

func (s *StdOutListener) GetStdOut() []string {
	lines := make([]string, 0)
	buf := make([]byte, 0)
	var n byte = 10
	for _, b := range s.buf {
		if b == n {
			lines = append(lines, string(buf))
			buf = make([]byte, 0)
			continue
		}
		buf = append(buf, b)
	}
	return lines
}

func main() {
	envSlice := make([]string, 100)
	args := make([]string, 12)
	var err error

	envDir, err := filepath.Abs("./test/env/")
	testProg, err := filepath.Abs("./test/test.sh")
	if err != nil {
		fmt.Println(err)
		return
	}

	args[0] = envDir
	args[1] = testProg

	var env, word string
	for i := 0; i < 100; i++ {
		word = fake.Word()
		env = fmt.Sprintf("%s=%s", strings.ToUpper(word), word)
		envSlice[i] = env
	}
	for i := 2; i < 12; i++ {
		args[i] = fake.Word()
	}

	envFiles := make([]string, 2)
	envFiles[0], err = makeEnvFile(fake.Word(), envSlice[0:50])
	envFiles[1], err = makeEnvFile(fake.Word(), envSlice[50:])
	if err != nil {
		fmt.Println(err)
	}
	defer removeEnvFiles(envFiles)

	progPath, err := filepath.Abs("./goenvdir")
	if err != nil {
		fmt.Println(err)
		return
	}

	listener := &StdOutListener{}
	cmd := exec.Command(progPath, args...)
	cmd.Stdout = listener
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	if isValid(listener.GetStdOut(), append(envSlice, args[2:]...)) {
		fmt.Println("Test passed")
	} else {
		fmt.Println("test failed")
	}
}

func isValid(output []string, contain []string) bool {
	for _, i := range contain {
		if !containLine(output, i) {
			fmt.Println(i)
			return false
		}
	}
	return true
}

func containLine(arr []string, line string) bool {
	for _, i := range arr {
		if i == line {
			return true
		}
	}

	return false
}

func makeEnvFile(fileName string, envSlice []string) (envPath string, err error) {
	envPath, err = filepath.Abs(fmt.Sprintf("./test/env/%s.env", fileName))
	if err != nil {
		return
	}
	file, err := os.Create(envPath)
	if err != nil {
		return
	}
	for _, s := range envSlice {
		_, err = file.WriteString(s + "\n")
		if err != nil {
			return
		}
	}

	return
}

func removeEnvFiles(envFiles []string) {
	for _, f := range envFiles {
		err := os.Remove(f)
		if err != nil {
			fmt.Println(err)
		}
	}
}
