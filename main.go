package main

import (
	"archive/zip"
	"github.com/cavaliergopher/grab/v3"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func main() {
	getLibs()

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
	cmd.Stderr = os.Stderr

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

func getLibs() {
	// what i need to get is lwjgl lwjgl-glfw lwjgl-jemalloc lwjgl-tinyfd lwjgl-stb lwjgl-opengl lwjgl-openal(without native) openal objc
	lwjgls := []string{"lwjgl", "lwjgl-glfw", "lwjgl-jemalloc", "lwjgl-tinyfd", "lwjgl-stb", "lwjgl-opengl", "lwjgl-openal"}
	metaUrls := []string{"https://repo1.maven.org/maven2/ca/weblite/java-objc-bridge/maven-metadata.xml", "https://repo1.maven.org/maven2/org/lwjgl/lwjgl/maven-metadata.xml"}

	lwjglVersion := ""
	ObjCVersion := ""

	/*client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}*/

	for _, v := range []string{"libraries", "natives", "temp"} {
		err := os.Mkdir(v, 0770)
		if err != nil {
			if !strings.Contains(err.Error(), "file exists") {
				log.Fatal(err)
			}
		}
	}

	for _, v := range metaUrls {
		resp, err := http.Get(v)
		if err != nil {
			log.Fatal(err)
		}

		xmlbytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()

		xmlstr := string(xmlbytes)
		versionIndex := strings.Index(xmlstr, "<latest>") + 8
		versionOutdex := strings.Index(xmlstr, "</latest>")

		version := string([]rune(xmlstr)[versionIndex:versionOutdex])
		if strings.Contains(v, "lwjgl") {
			lwjglVersion = version
		} else {
			ObjCVersion = version
		}
	}

	for _, v := range lwjgls { // process lwjgl
		jarUrl := "https://repo1.maven.org/maven2/org/lwjgl/" + v + "/" + lwjglVersion + "/" + v + "-" + lwjglVersion + ".jar"
		nativeUrl := "https://repo1.maven.org/maven2/org/lwjgl/" + v + "/" + lwjglVersion + "/" + v + "-" + lwjglVersion + "-natives-macos-arm64.jar"

		wg.Add(1)
		go func() {
			defer wg.Done()

			_, err := grab.Get("libraries/", jarUrl)
			if err != nil {
				log.Fatal(err)
			}
		}()

		//if v == "lwjgl-openal" { // TODO using openal-soft, according to GameParrot/minecraft-mac-window-fix
		//	continue
		//}

		wg.Add(1)
		go func() {
			defer wg.Done()

			resp, err := grab.Get("temp/", nativeUrl)
			if err != nil {
				log.Fatal(err)
			}

			jarZip, err := zip.OpenReader(resp.Filename)
			if err != nil {
				log.Fatal(err)
			}

			for _, v := range jarZip.File {
				if strings.Contains(v.Name, "dylib") {
					dylibZip, err := v.Open()
					if err != nil {
						log.Fatal(err)
					}

					dylib, err := os.Create("natives/" + strings.Split(v.Name, "/")[strings.Count(v.Name, "/")])
					if err != nil {
						log.Fatal(err)
					}

					io.Copy(dylib, dylibZip)
					dylibZip.Close()
					dylib.Close()

					break
				}
			}

			jarZip.Close()
		}()
	}

	wg.Add(1)
	go func() { // process ojbc-bridge
		defer wg.Done()

		jarUrl := "https://repo1.maven.org/maven2/ca/weblite/java-objc-bridge/" + ObjCVersion + "/" + "java-objc-bridge-" + ObjCVersion + ".jar"

		resp, err := grab.Get("temp/", jarUrl)
		if err != nil {
			log.Fatal(err)
		}

		jarZip, err := zip.OpenReader(resp.Filename)
		if err != nil {
			log.Fatal(err)
		}

		dylib, err := os.Create("natives/" + "libjcocoa.dylib")
		if err != nil {
			log.Fatal(err)
		}

		dylibZip, err := jarZip.Open("libjcocoa.dylib")

		//TODO remove native in jar

		io.Copy(dylib, dylibZip)
		jarZip.Close()
		dylib.Close()
	}()

	wg.Wait()
	os.RemoveAll("temp/")
}
