package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type MetaFileItem struct {
	RelPath   string `json:"rel_path"`
	Sha256    string `json:"sha256"`
	Size      int64  `json:"size"`
	Ext       string `json:"ext"`
	UpdatedAt int64  `json:"updated_at"`
}

type Meta struct {
	Status    string         `json:"status"`     // idle/running/success/failed
	LastError string         `json:"last_error"`
	Files     []MetaFileItem `json:"files"`
	UpdatedAt int64          `json:"updated_at"`
}

func metaPath(teacherDir string) string {
	return filepath.Join(teacherDir, "meta.json")
}

func MetaLoad(teacherDir string) (*Meta, error) {
	p := metaPath(teacherDir)
	b, err := os.ReadFile(p)
	if err != nil {
		// 不存在则返回默认
		if os.IsNotExist(err) {
			return &Meta{Status: "idle", Files: make([]MetaFileItem, 0), UpdatedAt: time.Now().Unix()}, nil
		}
		return nil, err
	}

	var m Meta
	if err := json.Unmarshal(b, &m); err != nil {
		// 解析失败也给默认，避免全挂
		return &Meta{Status: "idle", Files: make([]MetaFileItem, 0), UpdatedAt: time.Now().Unix()}, nil
	}
	if m.Files == nil {
		m.Files = make([]MetaFileItem, 0)
	}
	return &m, nil
}

func MetaSave(teacherDir string, m *Meta) error {
	m.UpdatedAt = time.Now().Unix()
	b, _ := json.MarshalIndent(m, "", "  ")
	return WriteFileAtomic(metaPath(teacherDir), b, 0o644)
}

func MetaSetStatus(teacherDir, status, lastErr string) error {
	m, err := MetaLoad(teacherDir)
	if err != nil {
		return err
	}
	m.Status = status
	m.LastError = lastErr
	return MetaSave(teacherDir, m)
}

func MetaUpsertFile(teacherDir string, item MetaFileItem) error {
	m, err := MetaLoad(teacherDir)
	if err != nil {
		return err
	}

	found := false
	for i := range m.Files {
		if m.Files[i].RelPath == item.RelPath {
			m.Files[i] = item
			found = true
			break
		}
	}
	if !found {
		m.Files = append(m.Files, item)
	}
	return MetaSave(teacherDir, m)
}
