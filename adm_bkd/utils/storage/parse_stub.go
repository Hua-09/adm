package storage

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ParseToTextStub：占位解析器
// - 将老师目录下（或指定 rel_paths）的文件“假解析”成 txt，写入 parsed/ 目录
// - 你后续可以替换为真正的 pdf/docx/xlsx 解析
func ParseToTextStub(teacherDir string, relPaths []string) error {
	parsedDir := filepath.Join(teacherDir, "parsed")
	if err := os.MkdirAll(parsedDir, 0o755); err != nil {
		return err
	}

	// 如果没指定文件列表，则不生成（由 Python 自行读取原文件也可以）
	if len(relPaths) == 0 {
		// 也可以在此扫描 pdf/docx/xlsx 目录；这里先保持最小实现
		return nil
	}

	for _, rel := range relPaths {
		rel = strings.TrimSpace(rel)
		if rel == "" {
			continue
		}

		// 只允许相对路径（不能包含 .. 或绝对路径）
		if strings.Contains(rel, "..") || strings.HasPrefix(rel, "/") || strings.Contains(rel, "\\") {
			continue
		}

		srcAbs := filepath.Join(teacherDir, filepath.FromSlash(rel))
		if _, err := os.Stat(srcAbs); err != nil {
			continue
		}

		txtName := filepath.Base(srcAbs) + ".txt"
		txtAbs := filepath.Join(parsedDir, txtName)

		content := []byte("stub parsed content\nfile=" + rel + "\nparsed_at=" + time.Now().Format(time.RFC3339) + "\n")
		_ = WriteFileAtomic(txtAbs, content, 0o644)
	}

	return nil
}
