package storage

import (
	"os"
	"path/filepath"
	"sort"
)

// ListDirs 列出某目录下的一级子目录名
func ListDirs(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0)
	for _, e := range entries {
		if e.IsDir() {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out, nil
}

type TeacherFileList struct {
	PDF    []string `json:"pdf"`
	DOCX   []string `json:"docx"`
	XLSX   []string `json:"xlsx"`
	IMG    []string `json:"img"`
	OTHERS []string `json:"others"`
	PARSED []string `json:"parsed"`
	AIRES  []string `json:"ai_result"`
}

// ListTeacherFiles 扫描老师目录下的文件清单
func ListTeacherFiles(teacherDir string) (TeacherFileList, error) {
	var out TeacherFileList
	var err error

	out.PDF, _ = listFiles(filepath.Join(teacherDir, "pdf"))
	out.DOCX, _ = listFiles(filepath.Join(teacherDir, "docx"))
	out.XLSX, _ = listFiles(filepath.Join(teacherDir, "xlsx"))
	out.IMG, _ = listFiles(filepath.Join(teacherDir, "img"))
	out.OTHERS, _ = listFiles(filepath.Join(teacherDir, "others"))
	out.PARSED, _ = listFiles(filepath.Join(teacherDir, "parsed"))
	out.AIRES, _ = listFiles(filepath.Join(teacherDir, "ai_result"))

	_ = err
	return out, nil
}

func listFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		// 目录不存在不算错误
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	out := make([]string, 0)
	for _, e := range entries {
		if !e.IsDir() {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out, nil
}
