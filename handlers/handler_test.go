package handlers

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
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
	}{
		{[]string{}, []string{}, 0, []EnvRepairFlag{NvimFolder, EnvimLuaFile, PluginsFolder}},
		{[]string{}, []string{".nvim"}, 3, []EnvRepairFlag{}},
    {[]string{".nvim"}, []string{".plugins"}, 1, []EnvRepairFlag{EnvimLuaFile}},
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
			ce.Execute(nil)

			if len(ce.errors) != tc.expectedErrCount {
				t.Errorf("Expected number of errors: %d, returned number of errors: %d", tc.expectedErrCount, len(ce.errors))
			}

			require.ElementsMatchf(t,
				ce.GetRepairFlags(),
				tc.expectedFlags,
				"Expected flags: %s, returned flags: %s",
				strings.Join(transliterate(tc.expectedFlags), ", "),
				strings.Join(transliterate(ce.GetRepairFlags()), ", "))
		})
	}
}
