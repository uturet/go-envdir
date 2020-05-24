package pkg

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// Main execution line
func Run(args []string) error {
	dir, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}

	prog := args[1]
	if len(args) > 2 {
		args = args[2:]
	}

	envSlice, err := ReadDir(dir)
	if err != nil {
		return err
	}

	return RunCmd(prog, args, envSlice)
}

// Read dir and extract env variables from all *.env files
func ReadDir(dir string) (envSlice []string, err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	var envs []string
	for _, file := range files {
		if path.Ext(file.Name()) == ".env" {
			envs, err = readFile(path.Join(dir, file.Name()))
			if err != nil {
				return
			}
			envSlice = append(envSlice, envs...)
		}
	}

	return
}

// Run command with env variables
func RunCmd(prog string, args []string, env []string) error {
	cmd := exec.Command(prog, args...)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	return cmd.Run()
}

// Read env file
func readFile(filename string) (envSlice []string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	return Parse(file)
}

// Extract all env variables from file
func Parse(r io.Reader) (envSlice []string, err error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		trimmedLine := strings.TrimSpace(scanner.Text())
		if len(trimmedLine) == 0 || strings.HasPrefix(trimmedLine, "#") {
			continue
		}
		envSlice = append(envSlice, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return
	}

	return
}
