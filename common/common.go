package common

import (
    "runtime"
    "path/filepath"
    "strings"
    "os"
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