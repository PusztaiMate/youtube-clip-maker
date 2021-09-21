package utils

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func CheckEnvVars(vars ...string) error {
	var missing []string

	for _, v := range vars {
		if _, ok := os.LookupEnv(v); !ok {
			missing = append(missing, v)
		}
	}

	if len(missing) != 0 {
		return fmt.Errorf("missing environment variables: %s", strings.Join(missing, ", "))
	}
	return nil
}

func MakeDirs(paths ...string) error {
	for _, p := range paths {
		if !isDir(p) {
			p = path.Base(p)
		}

		err := os.MkdirAll(p, 0777)
		if err != nil {
			return fmt.Errorf("could not create source video dir ('%s'): '%s'", p, err)
		}

	}
	return nil
}

func isDir(path_ string) bool {
	fileInfo, err := os.Stat(path_)
	if err == nil {
		return fileInfo.IsDir()
	}

	if path.Ext(path_) == "" {
		return true
	}

	return false
}
