package storage

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// ParseToText 将 teacherDir 下的 pdf/docx/xlsx 解析为 txt 写入 parsed/ 目录。
// - relPaths 为空：自动扫描 pdf/docx/xlsx 目录下的所有文件
// - relPaths 非空：只解析指定相对路径文件（如 pdf/xxx.pdf）
func ParseToText(teacherDir string, relPaths []string) error {
	parsedDir := filepath.Join(teacherDir, "parsed")
	if err := os.MkdirAll(parsedDir, 0o755); err != nil {
		return err
	}

	targets, err := collectTargets(teacherDir, relPaths)
	if err != nil {
		return err
	}

	for _, rel := range targets {
		srcAbs := filepath.Join(teacherDir, filepath.FromSlash(rel))
		if _, err := os.Stat(srcAbs); err != nil {
			continue
		}

		ext := strings.ToLower(filepath.Ext(srcAbs))
		outName := filepath.Base(srcAbs) + ".txt"
		outAbs := filepath.Join(parsedDir, outName)

		var out []byte
		switch ext {
		case ".pdf":
			out, err = parsePDFToText(srcAbs)
		case ".docx":
			out, err = parseDocxToText(srcAbs)
		case ".xlsx", ".xls":
			out, err = parseXlsxToText(srcAbs)
		default:
			// 不支持就跳过
			continue
		}
		if err != nil {
			// 单文件失败不中断全局：你也可以改为 return err
			continue
		}

		header := fmt.Sprintf("parsed_at=%s\nsource=%s\n\n", time.Now().Format(time.RFC3339), rel)
		_ = WriteFileAtomic(outAbs, append([]byte(header), out...), 0o644)
	}

	return nil
}

func collectTargets(teacherDir string, relPaths []string) ([]string, error) {
	if len(relPaths) > 0 {
		out := make([]string, 0, len(relPaths))
		for _, rel := range relPaths {
			rel = strings.TrimSpace(rel)
			if rel == "" {
				continue
			}
			// 只允许相对路径（避免穿越）
			if strings.Contains(rel, "..") || strings.HasPrefix(rel, "/") || strings.Contains(rel, "\\") {
				continue
			}
			out = append(out, rel)
		}
		return out, nil
	}

	// 未指定：自动扫描 pdf/docx/xlsx 三个目录
	out := make([]string, 0)
	addDir := func(sub string) {
		dir := filepath.Join(teacherDir, sub)
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			out = append(out, filepath.ToSlash(filepath.Join(sub, e.Name())))
		}
	}
	addDir("pdf")
	addDir("docx")
	addDir("xlsx")
	return out, nil
}

// -----------------------
// PDF -> TXT（依赖系统命令 pdftotext）
// apt-get install poppler-utils
// -----------------------
func parsePDFToText(pdfAbs string) ([]byte, error) {
	// pdftotext -layout input.pdf -
	cmd := exec.Command("pdftotext", "-layout", pdfAbs, "-")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, errors.New("pdftotext failed: " + stderr.String())
	}
	return stdout.Bytes(), nil
}

// -----------------------
// DOCX -> TXT（依赖系统命令 pandoc，或你自行换 docx2txt）
// apt-get install pandoc
// -----------------------
func parseDocxToText(docxAbs string) ([]byte, error) {
	// pandoc input.docx -t plain
	cmd := exec.Command("pandoc", docxAbs, "-t", "plain")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, errors.New("pandoc failed: " + stderr.String())
	}
	return stdout.Bytes(), nil
}

// -----------------------
// XLSX/XLS -> TXT（Go 直接解析，依赖 excelize）
// -----------------------
func parseXlsxToText(xlsxAbs string) ([]byte, error) {
	f, err := excelize.OpenFile(xlsxAbs)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	sheets := f.GetSheetList()
	var buf strings.Builder
	for _, sh := range sheets {
		rows, err := f.GetRows(sh)
		if err != nil {
			continue
		}
		buf.WriteString("sheet=" + sh + "\n")
		for _, row := range rows {
			buf.WriteString(strings.Join(row, "\t"))
			buf.WriteString("\n")
		}
		buf.WriteString("\n")
	}
	return []byte(buf.String()), nil
}
