package handlers

import (
	"fmt"
	"log"
	"os"
)

// This go file contains utility functions and utility handlers

// StatPath function: Check if a path exists as the given type
// Arguments:
// path: string - the path to check
// isDir: bool - whether the path should be a directory
// Returns:
// bool: whether the path exists as the:
//   - directory if isDir is true
//   - file if isDir is false
//
// bool: whether the path is creatable
func StatPath(path string, isDir bool) (bool, bool) {
	if info, err := os.Stat(path); os.IsNotExist(err) {
		return false, true
	} else {
		return info.IsDir() == isDir, false
	}
}

// CreateFolder

type CreateFolder struct {
	state      HandlerState
	errors     []error
	FolderName string
}

func (cf *CreateFolder) GetType() HandlerType {
	return CreateFolderType
}

func (cf *CreateFolder) GetState() HandlerState {
	return cf.state
}

func (cf *CreateFolder) GetErrors() []error {
	return cf.errors
}

func (cf *CreateFolder) Execute(state map[HandlerType]Handler) {
	if cf.state != HandlerNotStarted {
		log.Panic("Cannot execute a handler that has already been executed.")
	}

	isFolder, isCreatable := StatPath(cf.FolderName, true)
	if !isCreatable {
		var error string
		if isFolder {
			error = "Folder already exists"
			cf.state = HandlerSuccess
		} else {
			error = "A file exists with the same name"
			cf.state = HandlerError
		}
		cf.errors = append(cf.errors, fmt.Errorf("Create Folder %s error: %s", cf.FolderName, error))
		return
	}

	if err := os.MkdirAll(cf.FolderName, 0755); err != nil {
		cf.errors = append(cf.errors, fmt.Errorf("Create Folder %s error: %s", cf.FolderName, err.Error()))
		cf.state = HandlerError
		return
	}
	cf.state = HandlerSuccess
	return
}

func (cf *CreateFolder) ShouldProceed() bool {
	return cf.state == HandlerSuccess
}
