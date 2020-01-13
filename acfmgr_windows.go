// +build windows

package acfmgr

import (
        "golang.org/x/sys/windows/registry"
        "os"
        "os/user"
        "path/filepath"
        "runtime"
        "strings"
)

// expandPathO takes the path as a string and attempts to massage
// it into something the ioutils can use. Should support situations
// where people try to put stuff like ~ and $HOME in the path
func expandPathO(path string, usr *user.User) (expandedPath string, err error) {
        dir := usr.HomeDir
        prefixTildeNix := "~/"
        prefixTildeWin := "~\\"
        // fmt.Printf("checking path `%s` for prefix `%s`\n", path, prefix)
        if strings.HasPrefix(path, prefixTildeNix) || strings.HasPrefix(path, prefixTildeWin) {
                // Use strings.HasPrefix so we don't match paths like
                // "/something/~/something/"
                expandedPath = filepath.Join(dir, path[2:])
        } else {
                expandedPath = path
        }
        if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
                if strings.HasPrefix(expandedPath, "$") {
                        // need to extract and expand env var
                        expandedPath = os.ExpandEnv(expandedPath)
                }
        } else if runtime.GOOS == "windows" {
                if strings.HasPrefix(expandedPath, "%") {
                        expandedPath, err = registry.ExpandString(expandedPath)
                        if err != nil {
                                return expandedPath, err
                        }
                }
        }
        expandedPath, err = filepath.Abs(expandedPath)
        return expandedPath, err
}

