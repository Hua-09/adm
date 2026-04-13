package storage

import (
	"os"
	"path/filepath"
	"time"
)

// TryLock 在 teacherDir 下创建 .lock 文件；成功返回 true
func TryLock(teacherDir string) (bool, error) {
	lockPath := filepath.Join(teacherDir, ".lock")
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		// 已存在 -> 认为被锁
		if os.IsExist(err) {
			return false, nil
		}
		return false, err
	}
	defer f.Close()

	_, _ = f.WriteString(time.Now().Format(time.RFC3339))
	return true, nil
}

func Unlock(teacherDir string) {
	_ = os.Remove(filepath.Join(teacherDir, ".lock"))
}
