package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// SaveMultipartFileWithHash 将 reader 内容写到 dstAbs，并返回 sha256
func SaveMultipartFileWithHash(dstAbs string, r io.Reader) (string, error) {
	if err := os.MkdirAll(filepathDir(dstAbs), 0o755); err != nil {
		return "", err
	}

	out, err := os.Create(dstAbs)
	if err != nil {
		return "", err
	}
	defer out.Close()

	h := sha256.New()
	w := io.MultiWriter(out, h)

	if _, err := io.Copy(w, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func filepathDir(p string) string {
	i := len(p) - 1
	for ; i >= 0; i-- {
		if p[i] == '/' || p[i] == '\\' {
			break
		}
	}
	if i <= 0 {
		return "."
	}
	return p[:i]
}
