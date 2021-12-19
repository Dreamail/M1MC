package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args
	proPath, _ := filepath.Abs(args[0])
	basePath := filepath.Dir(proPath)

	classPath := strings.Split(os.Getenv("CLASSPATH"), ":")
	nClassPath := []string{basePath + "/lwjglfat.jar"}

	for _, v := range classPath {
		if !strings.Contains(v, "lwjgl") {
			nClassPath = append(nClassPath, v)
		}
	}

	/*for i := range args {
		if strings.Contains(args[i], "-cp") {
			classPath := strings.Split(args[i+1], ":")

			for j := range classPath {
				if strings.Contains(classPath[j], "lwjgl") {
					classPath = append(classPath[:j], classPath[j+1:]...)
				}
			}

			classPath = append(classPath, basePath+"lwjglfat.jar")
			//println(strings.Join(classPath, " "))

			args[i+1] = strings.Join(classPath, ":")

			break
		}
	}*/

	//println(strings.Join(args, " "))

	//println(strings.Join(classPath, ":"))

	cmd := exec.Command(args[1], args[2:]...)
	cmd.Env = append(os.Environ(), "CLASSPATH="+strings.Join(nClassPath, ":"))
	cmd.Stdout = os.Stdout

	/*var err error

	var wg sync.WaitGroup

	go func() {
		wg.Add(1)
		err = cmd.Run()
		wg.Done()
	}()
	//println(string(stdout.Bytes()))
	wg.Wait()
	if err != nil {
		println(err.Error())
		return
	}*/

	println(strings.Join(nClassPath, ":"))

	var stdout bytes.Buffer

	err := cmd.Run()
	println(string(stdout.Bytes()))
	if err != nil {
		println(err.Error())
		return
	}
}
