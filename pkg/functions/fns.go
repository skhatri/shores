package functions

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func ListFiles(rootPath string, ext string) []string {
	files := make([]string, 0)
	walkerFn := func(path string, fi os.FileInfo, err error) error {
		if err == nil && !fi.IsDir() {
			if filepath.Ext(path) == ext && strings.HasSuffix(path, ext) {
				files = append(files, path)
			}
		}
		return err
	}
	filepath.Walk(rootPath, walkerFn)
	return files
}

func UnmarshalFile(file string, t interface{}) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.New(fmt.Sprintf("file: [%s], error: [%v]", file, err))
	}
	err = UnmarshalYaml(content, t)
	if err != nil {
		return errors.New(fmt.Sprintf("file: [%s], error: [%v]", file, err))
	}
	return nil
}

func UnmarshalYaml(content []byte, t interface{}) error {
	err := yaml.Unmarshal(content, t)
	if err != nil {
		return errors.New(fmt.Sprintf("error unmarshalling to %T, with error %v", t, err))
	}
	return nil
}
