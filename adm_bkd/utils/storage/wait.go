package storage

import (
	"os"
	"time"
)

// WaitFiles 等待 files 全部存在且大小>0，直到超时
func WaitFiles(file1, file2 string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if fileReady(file1) && fileReady(file2) {
			return true
		}
		time.Sleep(500 * time.Millisecond)
	}
	return false
}

func fileReady(p string) bool {
	st, err := os.Stat(p)
	if err != nil {
		return false
	}
	return !st.IsDir() && st.Size() > 0
}
