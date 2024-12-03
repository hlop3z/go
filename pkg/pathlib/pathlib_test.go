package pathlib

import (
	"reflect"
	"testing"
)

// TestGetBaseDir verifies the behavior of GetBaseDir.
// It checks if the base directory name matches the expected value.
func TestGetBaseDir(t *testing.T) {
	expectedName := "pathlib"
	file_name := GetBaseDir().Name()
	if file_name != expectedName {
		t.Fatalf("Expected base directory name %v, but got %v", expectedName, file_name)
	}
}

// TestParents verifies the functionality of Parents.
// It ensures the parent path name matches the expected value.
func TestParents(t *testing.T) {
	expectedName := "pkg"
	path := GetBaseDir()
	parent_name := path.Parents(0).Name()
	if parent_name != expectedName {
		t.Fatalf("Expected parent directory name %v, but got %v", expectedName, parent_name)
	}
}

// TestCreate verifies the behavior of the Create method.
// It ensures that a new file is created with the correct name.
func TestCreate(t *testing.T) {
	expectedName := "base.json"
	path := GetBaseDir()
	outpath := path.Create("go_demo_data/templates/base.json")
	outpath_name := outpath.Name()
	if outpath_name != expectedName {
		t.Fatalf("Expected file name %v, but got %v", expectedName, outpath_name)
	}
}

// TestRead verifies the behavior of the Read method.
// It ensures the data read from a file is not empty.
func TestRead(t *testing.T) {
	filePath := "go_demo_data/templates/base.json"
	path := GetBaseDir()
	path.Create(filePath)
	file := path.Join(filePath)
	jsonData := file.Read()
	if reflect.DeepEqual(jsonData, make(map[string]interface{})) {
		t.Fatalf("Expected non-empty data, but got empty data from file %v", file.String())
	}
}

// TestDelete verifies the behavior of the Delete method.
// It ensures a directory or file is deleted successfully.
func TestDelete(t *testing.T) {
	dirPath := "go_demo_data/"
	path := GetBaseDir()
	success := path.Join(dirPath).Delete()
	if !success {
		t.Fatalf("Failed to delete directory %v", dirPath)
	}
}

// TestStringExistsAbsolute verifies directory operations and attributes.
// It checks existence, absolute path, and cleanup.
func TestStringExistsAbsolute(t *testing.T) {
	expectedName := "templates"
	path := GetBaseDir()
	baseDir := path.Join("data")
	tmplDir := baseDir.Join("templates")

	if tmplDir.Name() != expectedName {
		t.Fatalf("Expected directory name %v, but got %v", expectedName, tmplDir.Name())
	}
	if tmplDir.Exists() {
		t.Fatalf("Directory %v already exists, but it shouldn't", tmplDir.String())
	}
	if !tmplDir.IsAbsolute() {
		t.Fatalf("Path %v is not absolute", tmplDir.String())
	}

	// Clean up
	path.Join("data/").Delete()
}

// TestSearchRecursively verifies the behavior of Find for recursive file searches.
// It ensures the correct number of files are matched by patterns.
func TestSearchRecursively(t *testing.T) {
	path := GetBaseDir()
	files := path.Parent().Find([]string{"*.go", "*.py"})
	expectedCount := 2
	if len(files["*.go"]) != expectedCount {
		t.Fatalf("Expected %v files matching '*.go', but found %v", expectedCount, len(files["*.go"]))
	}
}
