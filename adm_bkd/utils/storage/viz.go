package storage

import (
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// DocSummary is a lightweight view of a stored document.
type DocSummary struct {
	ID         string    `json:"id"`
	OrigName   string    `json:"orig_name"`
	Size       int64     `json:"size"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// SaveUpload stores the uploaded multipart file and returns its assigned path.
func SaveUpload(rootDir string, fh *multipart.FileHeader) (string, error) {
	id := uuid.NewString()
	dest := docPath(rootDir, id)

	src, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	data := make([]byte, fh.Size)
	if _, err := src.Read(data); err != nil {
		return "", err
	}

	if err := atomicWrite(dest, data, 0o644); err != nil {
		return "", err
	}

	sha, err := hashFile(dest)
	if err != nil {
		return "", err
	}

	meta := DocMeta{
		ID:          id,
		OrigName:    fh.Filename,
		ContentType: fh.Header.Get("Content-Type"),
		Size:        fh.Size,
		SHA256:      sha,
		UploadedAt:  time.Now().UTC(),
	}
	return dest, writeMeta(rootDir, meta)
}

// ListDocs returns a summary of all stored documents.
func ListDocs(rootDir string) ([]DocSummary, error) {
	dir := docsDir(rootDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []DocSummary{}, nil
		}
		return nil, err
	}

	var docs []DocSummary
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) == ".json" {
			continue
		}
		id := e.Name()
		m, err := readMeta(rootDir, id)
		if err != nil {
			continue
		}
		docs = append(docs, DocSummary{
			ID:         m.ID,
			OrigName:   m.OrigName,
			Size:       m.Size,
			UploadedAt: m.UploadedAt,
		})
	}
	return docs, nil
}

// DeleteDoc removes the document file and its metadata.
func DeleteDoc(rootDir, id string) error {
	return withFileLock(docPath(rootDir, id), func() error {
		if err := os.Remove(docPath(rootDir, id)); err != nil && !os.IsNotExist(err) {
			return err
		}
		if err := os.Remove(metaPath(rootDir, id)); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	})
}
