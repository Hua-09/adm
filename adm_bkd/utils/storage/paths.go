package storage

import (
	"path/filepath"
)

// docsDir returns the absolute path to the documents directory under rootDir.
func docsDir(rootDir string) string {
	return filepath.Join(rootDir, "docs")
}

// docPath returns the full file path for a document with the given id.
func docPath(rootDir, id string) string {
	return filepath.Join(docsDir(rootDir), id)
}
