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

	nClassPath := make([]string, 0)
	filepath.Walk(basePath+"/lwjgl/", func(path string, info fs.FileInfo, err error) error {
		if strings.Contains(info.Name(), "lwjgl") {
			nClassPath = append(nClassPath, path)
		}
		return err
	})

	classPath := strings.Split(os.Getenv("CLASSPATH"), ":")

	for _, v := range classPath {
		if !strings.Contains(v, "lwjgl") {
			nClassPath = append(nClassPath, v)
		}
	}

	cmd := exec.Command(args[1], args[2:]...)
	cmd.Env = append(os.Environ(), "CLASSPATH="+strings.Join(nClassPath, ":"))
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
