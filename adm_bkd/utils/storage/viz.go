package storage

import (
	"regexp"
	"strings"
)

// BuildVizFromMarkdown：非常轻量的启发式解析
// - 提取关键词（简单：取 ##/# 标题、以及出现频率较���的词）
// - 提取数字（形如 123 / 12.34 / 50%）
// 你后续可以：让 AI 直接输出结构化 JSON，就不需要启发式解析。
func BuildVizFromMarkdown(md string) map[string]interface{} {
	lines := strings.Split(md, "\n")

	// 标题做关键词
	keywords := make([]string, 0)
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if strings.HasPrefix(ln, "#") {
			kw := strings.TrimSpace(strings.TrimLeft(ln, "#"))
			if kw != "" {
				keywords = append(keywords, kw)
			}
		}
		if len(keywords) >= 20 {
			break
		}
	}

	// 数值抽取
	reNum := regexp.MustCompile(`\b\d+(\.\d+)?%?\b`)
	nums := make([]string, 0)
	for _, m := range reNum.FindAllString(md, -1) {
		nums = append(nums, m)
		if len(nums) >= 200 {
			break
		}
	}

	return map[string]interface{}{
		"keywords": keywords,
		"numbers":  nums,
		// 给前端一个可扩展结构：nodes/edges（流程图/关系图）
		"graph": map[string]interface{}{
			"nodes": []interface{}{},
			"edges": []interface{}{},
		},
	}
}
