package handlers

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"strings"
	"testing"
)

func touchFoldersAndFiles(t *testing.T, folders []string, files []string) string {
	dir := t.TempDir()
	for _, folderName := range folders {
		if err := os.MkdirAll(path.Join(dir, folderName), 0755); err != nil {
			t.Fatalf("Can't create folder %s", folderName)
		}
	}

	for _, fileName := range files {
		if _, err := os.OpenFile(path.Join(dir, fileName), os.O_CREATE, 0644); err != nil {
			t.Fatalf("Can't create file %s", fileName)
		}
	}

	return dir
}

func transliterate[S ~[]E, E fmt.Stringer](l S) []string {
	var res []string
	for _, val := range l {
		res = append(res, val.String())
	}
	return res
}

func TestCheckEnvironmentHandler(t *testing.T) {

	var tests = []struct {
		folderList       []string
		fileList         []string
		expectedErrCount int
		expectedFlags    []EnvRepairFlag
		expectedState    HandlerState
	}{
		// Case 1
		{
			[]string{},
			[]string{},
			0,
			[]EnvRepairFlag{
				NvimFolder,
				EnvimLuaFile,
				PluginsFolder,
			},
			HandlerSuccessState,
		},
		// Case 2
		{
			[]string{},
			[]string{
				".nvim",
			},
			3,
			[]EnvRepairFlag{},
			HandlerErrorState,
		},
		// Case 3
		{
			[]string{
				".nvim",
			},
			[]string{},
			0,
			[]EnvRepairFlag{
				EnvimLuaFile,
				PluginsFolder,
			},
			HandlerSuccessState,
		},
		// Case 4
		{
			[]string{
				".nvim",
				path.Join(".nvim", "envim.lua"),
			},
			[]string{},
			1,
			[]EnvRepairFlag{
				PluginsFolder,
			},
			HandlerErrorState,
		},
		// Case 5
		{
			[]string{
				".nvim",
			},
			[]string{
				path.Join(".nvim", "envim.lua"),
			},
			0,
			[]EnvRepairFlag{
				PluginsFolder,
			},
			HandlerSuccessState,
		},
		// Case 6
		{
			[]string{
				".nvim",
			},
			[]string{
				path.Join(".nvim", "envim.lua"),
				path.Join(".nvim", "plugins"),
			},
			1,
			[]EnvRepairFlag{},
			HandlerErrorState,
		},
		// Case 7
		{
			[]string{
				".nvim",
				path.Join(".nvim", "plugins"),
			},
			[]string{
				path.Join(".nvim", "envim.lua"),
			},
			0,
			[]EnvRepairFlag{},
			HandlerSuccessState,
		},
	}

	for _, tc := range tests {
		name := strings.Join(append(tc.folderList, tc.fileList...), ", ")
		if name == "" {
			name = "Empty"
		} else {
			name = fmt.Sprintf("{ %s }", name)
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			dir := touchFoldersAndFiles(t, tc.folderList, tc.fileList)
			ce := CheckEnvironment{Path: dir}
			require.Equal(t,
				ce.GetType(),
				CheckEnvironmentType,
				"Expected type: %s, returned type: %s",
				CheckEnvironmentType,
				ce.GetType())

			require.Equal(t,
				ce.DependsOn(),
				[]HandlerType{},
				"Expected dependencies: %s, returned dependencies: %s",
				[]HandlerType{},
				ce.DependsOn())

			ce.Execute(nil)

			if len(ce.GetErrors()) != tc.expectedErrCount {
				t.Errorf("Expected number of errors: %d, returned number of errors: %d", tc.expectedErrCount, len(ce.GetErrors()))
			}

			require.ElementsMatchf(t,
				ce.GetRepairFlags(),
				tc.expectedFlags,
				"Expected flags: %s, returned flags: %s",
				strings.Join(transliterate(tc.expectedFlags), ", "),
				strings.Join(transliterate(ce.GetRepairFlags()), ", "))

			require.True(t,
				ce.ShouldProceed() == (ce.GetState() == HandlerSuccessState),
				"ShouldProceed() and GetState() should be in sync. If the state is error, we should not proceed, "+
					"Otherwise we should.")

		})
	}
}

func TestRepairEnvironmentHandler(t *testing.T) {
	tests := []struct {
		// folders to be created in the environment before execution
		folders    []string
		testFlags []EnvRepairFlag
	}{
		{
      folders: []string{},
      testFlags: []EnvRepairFlag{NvimFolder, EnvimLuaFile, PluginsFolder},
		},
		{
      folders: []string{".nvim"},
      testFlags: []EnvRepairFlag{EnvimLuaFile, PluginsFolder},
		},
	}

	for _, tc := range tests {
		temp := []string{}
		for _, flag := range tc.testFlags {
			temp = append(temp, flag.String())
		}
		name := strings.Join(temp, ", ")
    dir := touchFoldersAndFiles(t, tc.folders, []string{})

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ce := &CheckEnvironment{rflags: tc.testFlags, Path: dir}
			re := RepairEnvironment{}
			re.Execute(map[HandlerType]Handler{CheckEnvironmentType: ce})
			require.Equal(t, HandlerSuccessState, re.GetState(), "Expected state: %s, returned state: %s", HandlerSuccessState, re.GetState())
		})
	}
}
