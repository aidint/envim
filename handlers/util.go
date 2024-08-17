package handlers

import (
	"fmt"
	"os"
	"path/filepath"
)

// This go file contains utility functions and utility handlers

// statPath function: Check if a path exists as the given type
// Arguments:
// path: string - the path to check
// isDir: bool - whether the path should be a directory
// Returns:
// bool: whether the path exists as the:
//   - directory if isDir is true
//   - file if isDir is false
//
// bool: whether the path is creatable
func statPath(p string, isDir bool) (bool, bool) {
  p = filepath.Clean(p)
	if info, err := os.Stat(p); os.IsNotExist(err) {
    return false, true
	} else if err == nil {
		return info.IsDir() == isDir, false
  }
  return false, false
}
// createFolder
func createFolder(path string) (bool, error) {
	isFolder, isCreatable := statPath(path, true)
	if !isCreatable {
		if isFolder {
      return true, fmt.Errorf("Create Folder %s error: Folder already exists", path)
		} else {
      return false, fmt.Errorf("Create Folder %s error: A file exists with the same name", path)
		}
	}

	if err := os.MkdirAll(path, 0755); err != nil {
    return false, fmt.Errorf("Create Folder %s error: %s", path, err.Error())
	}
	return true, nil
}
