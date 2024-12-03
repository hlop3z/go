package pathlib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Dict = map[string]interface{}
type Path struct {
	path string
}

// GetBaseDir returns the current working directory as the base directory.
func GetBaseDir() Path {
	dir, err := os.Getwd()
	if err != nil {
		// fmt.Println("Error getting current directory:", err)
		return NewPath(".")
	}
	return NewPath(dir)
}

// NewPath creates a new Path instance with the given directory string.
func NewPath(p string) Path {
	return Path{path: filepath.Clean(p)}
}

// String returns the name of the file.
func (p Path) Name() string {
	return filepath.Base(p.path)
}

// String returns the string representation of the Path.
func (p Path) String() string {
	return p.path
}

// Exists checks if the path exists on the filesystem.
func (p Path) Exists() bool {
	_, err := os.Stat(p.path)
	return err == nil
}

// IsAbsolute checks if the path is an absolute path.
func (p Path) IsAbsolute() bool {
	return filepath.IsAbs(p.path)
}

// Join joins the current path with another path segment.
func (p Path) Join(other string) Path {
	return Path{path: filepath.Join(p.path, other)}
}

// Parent returns the immediate parent directory of the current path.
func (p Path) Parent() Path {
	return Path{path: filepath.Dir(p.path)}
}

// Parents returns the parent directories up to the specified depth.
func (p Path) Parents(depth uint) Path {
	var parents []Path
	current := p

	depth++ // Increase depth to include the current path

	// Collect parents up to the specified depth
	for i := uint(0); i < depth; i++ {
		current = current.Parent()

		// Stop if we reach the root
		if current.path == "." || current.path == string(filepath.Separator) {
			break
		}

		parents = append(parents, current)
	}

	// Return the last collected parent or the current path if no parents are found
	return parents[len(parents)-1]
}

// Find searches for files matching the given pattern recursively
// and returns a slice of Path objects. It always returns a list, even if empty.
func (p Path) Find(patterns []string) map[string][]Path {
	dict := map[string][]Path{}
	for _, pattern := range patterns {
		dict[pattern] = p.FindOne(pattern)
	}
	return dict
}

// Find searches for files matching the given pattern recursively
// and returns a slice of Path objects. It always returns a list, even if empty.
func (p Path) FindOne(pattern string) []Path {
	var matches []Path

	// Walk through the directory structure
	err := filepath.Walk(p.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories
		if info.IsDir() {
			return nil
		}
		// Match the file against the pattern
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, NewPath(path))
		}
		return nil
	})

	// If there's an error during walking, print it but still return an empty slice
	if err != nil {
		// Log or handle the error here if needed
		fmt.Println("Error during walk:", err)
	}

	// Always return the list (empty or populated)
	return matches
}

// Mkdir creates the directory specified by the Path, including any necessary parent directories.
func (p Path) Mkdir() error {
	dirname := p.String()
	// Create the directory and any necessary parent directories
	err := os.MkdirAll(dirname, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	return nil
}

// Touch creates a file at the specified Path if it does not exist.
func (p Path) Touch(pathname string) error {
	dirname := p.String()
	filename := filepath.Join(dirname, pathname)
	// Check if the file exists
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// File does not exist, create it
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}
		defer file.Close()
	} else if err != nil {
		// Some other error (not just "file not exists")
		return fmt.Errorf("failed to check file status: %v", err)
	}
	return nil
}

// createPath creates necessary directories and files for the specified path.
func createPath(pathname string) Path {
	folder, file := splitPath(pathname)
	if file == "" || folder == "" && file == "" {
		folder = pathname
	}
	p := NewPath(folder)
	if folder != "" {
		p.Mkdir()
	}
	if file != "" {
		p.Touch(file)
		return NewPath(p.Join(file).String())
	}
	return p
}

// splitPath splits a given string into a folder and a file. If the path ends with "/",
// it considers the last part as a folder; otherwise, it's a file.
func splitPath(path string) (folder, file string) {
	urlSeparator := string(filepath.Separator)
	// Trim any trailing slashes for consistency
	trimmedPath := strings.TrimRight(path, urlSeparator)

	// Find the last occurrence of a slash
	lastSlashIndex := strings.LastIndex(trimmedPath, urlSeparator)

	// If there is a slash in the path, it means we have a folder and a file
	if lastSlashIndex != -1 {
		folder = trimmedPath[:lastSlashIndex]
		file = trimmedPath[lastSlashIndex+1:]
	} else {
		// If there is no slash, the entire path is treated as a file
		folder = ""
		file = trimmedPath
	}

	// If the path ends with "/", it means the last part is a folder, and there's no file
	if strings.HasSuffix(path, urlSeparator) {
		file = "" // If it's a directory, we leave file empty
	}

	return folder, file
}

// Create creates a new path, ensuring the necessary directories and files exist.
func (p Path) Create(pathname string) Path {
	path := p.Join(pathname)
	return createPath(path.String())
}

// Read reads file content
func (p Path) Read() interface{} {
	data, err := os.ReadFile(p.String())
	if err != nil {
		// fmt.Println("Error reading file:", err)
		return nil
	}
	return data

}

// Remove file from the folder
func (p Path) Delete() bool {
	if p.Exists() {
		err := os.RemoveAll(p.String())
		if err != nil {
			fmt.Println("Error Deleting path:", err)
			return false
		}
	}
	return true
}
