package common

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	COMMON_DIR = "c:\\monkeydb2\\data"
)

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		dir = "."
	}
	return FixPath(dir)
}

func FixPath(s string) string {
	if runtime.GOOS == "windows" {
		return s
	}
	return strings.Replace(s, "\\", "/", -1)
}

func init() {
	if runtime.GOOS != "windows" {
		COMMON_DIR = "/var/local/monkeydb2/data"
	}
}
