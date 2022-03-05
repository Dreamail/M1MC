package main

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args
	proPath, _ := filepath.Abs(args[0])
	basePath := filepath.Dir(proPath)
	isEnv := os.Getenv("CLASSPATH") != ""
	classPathIndex := 0

	nClassPath := make([]string, 0)
	filepath.Walk(basePath+"/libraries/", func(path string, info fs.FileInfo, err error) error {
		if strings.Contains(info.Name(), "lwjgl") || strings.Contains(info.Name(), "java-objc-bridge") {
			nClassPath = append(nClassPath, path)
		}
		return err
	})

	classPathStr := ""

	if isEnv {
		classPathStr = os.Getenv("CLASSPATH")
	} else {
		for i, arg := range args {
			if arg == "-cp" {
				classPathIndex = i + 1
				classPathStr = args[i+1]
				break
			}
		}
	}

	for _, v := range strings.Split(classPathStr, ":") {
		if !strings.Contains(v, "lwjgl") && !strings.Contains(v, "java-objc-bridge") {
			nClassPath = append(nClassPath, v)
		}
	}

	args[classPathIndex] = strings.Join(nClassPath, ":")
	cmd := exec.Command(args[1], args[2:]...)
	if isEnv {
		cmd.Env = append(os.Environ(), "CLASSPATH=")
	}
	cmd.Stdout = os.Stdout

	//println(strings.Join(nClassPath, ":"))

	/*var stdout bytes.Buffer

	err := cmd.Run()
	println(string(stdout.Bytes()))
	if err != nil {
		println(err.Error())
		return
	}*/

	cmd.Run()
}
