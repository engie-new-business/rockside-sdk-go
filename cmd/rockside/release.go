// +build ignore

package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var builds = map[string][]string{
	"darwin":  {"amd64"},
	"linux":   {"amd64"},
	"windows": {"amd64"},
}

func main() {
	log.SetFlags(0)

	var wg sync.WaitGroup

	for osname, archs := range builds {
		for _, arch := range archs {
			wg.Add(1)
			go func(o, a string) {
				defer wg.Done()
				if err := buildAndPackage(o, a); err != nil {
					log.Fatal("%s\n", err)
				}
			}(osname, arch)
		}
	}

	wg.Wait()
}

func buildAndPackage(osname, arch string) error {
	builddir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(builddir)

	binName := "rockside"
	switch osname {
	case "windows":
		binName = "rockside.exe"
	}

	artefactPath := filepath.Join(builddir, binName)

	var stderr, stdout bytes.Buffer
	cmd := exec.Command("go", "build", "-trimpath", "-o", artefactPath)
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOARCH=%s", arch), fmt.Sprintf("GOOS=%s", osname))
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build: %v\n%s", err, stderr.Bytes())
	}

	zipFile, err := os.OpenFile(fmt.Sprintf("%s-%s-%s.zip", strings.Split(binName, ".")[0], osname, arch), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	w := zip.NewWriter(zipFile)

	f, err := w.Create(binName)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(artefactPath)
	if err != nil {
		return err
	}

	if _, err = f.Write(content); err != nil {
		return err
	}

	return w.Close()
}
