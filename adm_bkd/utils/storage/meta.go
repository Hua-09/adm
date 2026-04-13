package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// DocMeta holds metadata for a stored document.
type DocMeta struct {
	ID          string    `json:"id"`
	OrigName    string    `json:"orig_name"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	SHA256      string    `json:"sha256"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

func metaPath(rootDir, id string) string {
	return filepath.Join(docsDir(rootDir), id+".meta.json")
}

func writeMeta(rootDir string, m DocMeta) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return atomicWrite(metaPath(rootDir, m.ID), data, 0o644)
}

func readMeta(rootDir, id string) (DocMeta, error) {
	data, err := os.ReadFile(metaPath(rootDir, id))
	if err != nil {
		return DocMeta{}, err
	}
	var m DocMeta
	return m, json.Unmarshal(data, &m)
}
