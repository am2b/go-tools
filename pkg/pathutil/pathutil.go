package pathutil

import (
    "os"
    "path/filepath"
    "strings"
)

//扩展路径中的~,$HOME,${HOME},"${HOME}"
func ExpandPath(path string) string {
    //expand ~ to home
    if strings.HasPrefix(path, "~") {
        if home, err := os.UserHomeDir(); err == nil {
            path = strings.Replace(path, "~", home, 1)
        }
    }

    //expand environment variables
    path = os.ExpandEnv(path)

    //convert to absolute path
    if abs, err := filepath.Abs(path); err == nil {
        path = abs
    }

    return path
}
