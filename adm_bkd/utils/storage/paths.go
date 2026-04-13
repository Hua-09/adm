package storage

import (
	errmgr "adm_bkd/utils/err_mgr"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// allowName：允许中文、字母、数字、下划线、横线、点、空格
var reSafeName = regexp.MustCompile(`^[\p{Han}a-zA-Z0-9_\-\. ]+$`)

func IsSafeName(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	// 禁止路径分隔符与冒号
	if strings.Contains(s, "/") || strings.Contains(s, "\\") || strings.Contains(s, ":") {
		return false
	}
	// 禁止 ..（简单挡穿越）
	if strings.Contains(s, "..") {
		return false
	}
	return reSafeName.MatchString(s)
}

// JoinUnderRoot 将 path segments 拼到 root 下，并做“必须在 root 内”的校验
func JoinUnderRoot(root string, parts ...string) (string, int) {
	rootClean := filepath.Clean(root)
	if _, err := os.Stat(rootClean); err != nil {
		// root 可能还不存在，视情况你也可以允许自动创建；这里严格一点
		return "", errmgr.Err_storage_root_not_found
	}

	for _, p := range parts {
		if !IsSafeName(p) {
			return "", errmgr.Err_storage_path_invalid
		}
	}

	joined := filepath.Join(append([]string{rootClean}, parts...)...)
	joinedClean := filepath.Clean(joined)

	// 前缀校验：必须仍在 root 下
	rel, err := filepath.Rel(rootClean, joinedClean)
	if err != nil {
		return "", errmgr.Err_storage_path_invalid
	}
	if rel == "." {
		return joinedClean, errmgr.SUCCESS
	}
	if strings.HasPrefix(rel, "..") {
		return "", errmgr.Err_storage_path_invalid
	}

	return joinedClean, errmgr.SUCCESS
}

func IsAllowedExt(ext string, allow []string) bool {
	ext = strings.ToLower(strings.TrimSpace(ext))
	if ext == "" {
		return false
	}
	for _, a := range allow {
		if ext == strings.ToLower(strings.TrimSpace(a)) {
			return true
		}
	}
	return false
}
