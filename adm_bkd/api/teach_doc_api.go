package api

import (
	"adm_bkd/config"
	apiresult "adm_bkd/utils/api_result"
	errmgr "adm_bkd/utils/err_mgr"
	"adm_bkd/utils/storage"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
)

// ===============================
// 请求结构：目录定位（dept/date/teacher）
// ===============================
type reqTeachLocator struct {
	Dept    string `json:"dept"`
	Date    string `json:"date"`
	Teacher string `json:"teacher"`
}

type reqTeachDateList struct {
	Dept string `json:"dept"`
}

type reqTeachTeacherList struct {
	Dept string `json:"dept"`
	Date string `json:"date"`
}

type reqTeachFileList struct {
	Dept    string `json:"dept"`
	Date    string `json:"date"`
	Teacher string `json:"teacher"`
}

// TeachDeptList 列出所有部门目录
// Route: POST /api/teach/dept/list
func TeachDeptList(ctx iris.Context) {
	st := config.GlobalConfig.GetStorage()
	root := st.RootDir

	if _, err := os.Stat(root); err != nil {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_storage_root_not_found, nil))
		return
	}

	list, err := storage.ListDirs(root)
	if err != nil {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_storage_list_failed, nil))
		return
	}

	ctx.JSON(apiresult.NewAPIResult(errmgr.SUCCESS, list))
}

// TeachDateList 列出某部门下的日期目录
// Route: POST /api/teach/date/list
func TeachDateList(ctx iris.Context) {
	var req reqTeachDateList
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_http_imput_params_json_parse, nil))
		return
	}
	if strings.TrimSpace(req.Dept) == "" {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_http_input_params_empty, nil))
		return
	}

	st := config.GlobalConfig.GetStorage()
	deptPath, code := storage.JoinUnderRoot(st.RootDir, req.Dept)
	if code != errmgr.SUCCESS {
		ctx.JSON(apiresult.NewAPIResult(code, nil))
		return
	}

	list, err := storage.ListDirs(deptPath)
	if err != nil {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_storage_list_failed, nil))
		return
	}
	ctx.JSON(apiresult.NewAPIResult(errmgr.SUCCESS, list))
}

// TeachTeacherList 列出某部门某日期下的老师目录
// Route: POST /api/teach/teacher/list
func TeachTeacherList(ctx iris.Context) {
	var req reqTeachTeacherList
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_http_imput_params_json_parse, nil))
		return
	}
	if strings.TrimSpace(req.Dept) == "" || strings.TrimSpace(req.Date) == "" {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_http_input_params_empty, nil))
		return
	}

	st := config.GlobalConfig.GetStorage()
	p1, code := storage.JoinUnderRoot(st.RootDir, req.Dept, req.Date)
	if code != errmgr.SUCCESS {
		ctx.JSON(apiresult.NewAPIResult(code, nil))
		return
	}

	list, err := storage.ListDirs(p1)
	if err != nil {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_storage_list_failed, nil))
		return
	}
	ctx.JSON(apiresult.NewAPIResult(errmgr.SUCCESS, list))
}

// TeachFileList 返回某老师目录下文件清单（按类型目录）
// Route: POST /api/teach/file/list
func TeachFileList(ctx iris.Context) {
	var req reqTeachFileList
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_http_imput_params_json_parse, nil))
		return
	}
	if strings.TrimSpace(req.Dept) == "" || strings.TrimSpace(req.Date) == "" || strings.TrimSpace(req.Teacher) == "" {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_http_input_params_empty, nil))
		return
	}

	st := config.GlobalConfig.GetStorage()
	teacherDir, code := storage.JoinUnderRoot(st.RootDir, req.Dept, req.Date, req.Teacher)
	if code != errmgr.SUCCESS {
		ctx.JSON(apiresult.NewAPIResult(code, nil))
		return
	}

	data, err := storage.ListTeacherFiles(teacherDir)
	if err != nil {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_storage_list_failed, nil))
		return
	}

	ctx.JSON(apiresult.NewAPIResult(errmgr.SUCCESS, data))
}

// ===============================
// 上传接口
// ===============================

type reqTeachUpload struct {
	Dept     string `json:"dept"`
	Date     string `json:"date"`
	Teacher  string `json:"teacher"`
	FileType string `json:"file_type"` // pdf/docx/xlsx/img/others/video...
}

// TeachFileUpload 上传文件并保存到对应分类目录下
// Route: POST /api/teach/file/upload
func TeachFileUpload(ctx iris.Context) {
	st := config.GlobalConfig.GetStorage()

	// multipart 表单字段
	dept := ctx.FormValue("dept")
	date := ctx.FormValue("date")
	teacher := ctx.FormValue("teacher")
	fileType := ctx.FormValue("file_type")

	if strings.TrimSpace(dept) == "" || strings.TrimSpace(date) == "" || strings.TrimSpace(teacher) == "" || strings.TrimSpace(fileType) == "" {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_http_input_params_empty, nil))
		return
	}
	if !storage.IsSafeName(fileType) {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_storage_path_invalid, nil))
		return
	}

	teacherDir, code := storage.JoinUnderRoot(st.RootDir, dept, date, teacher)
	if code != errmgr.SUCCESS {
		ctx.JSON(apiresult.NewAPIResult(code, nil))
		return
	}

	// 接收文件
	f, info, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_http_input_params_validate_error, nil))
		return
	}
	defer f.Close()

	// 大小限制
	maxBytes := int64(st.MaxUploadMB) * 1024 * 1024
	if info.Size > maxBytes {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_storage_file_too_large, nil))
		return
	}

	// 扩展名白名单（只做基础校验）
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(info.Filename), "."))
	if !storage.IsAllowedExt(ext, st.AllowExtList) {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_storage_ext_not_allowed, nil))
		return
	}

	// 保存到：{teacherDir}/{fileType}/
	targetDir := filepath.Join(teacherDir, fileType)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_storage_mkdir_failed, nil))
		return
	}

	filename := filepath.Base(info.Filename)
	if !storage.IsSafeName(filename) {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_storage_path_invalid, nil))
		return
	}
	dstAbs := filepath.Join(targetDir, filename)

	sha256Hex, saveErr := storage.SaveMultipartFileWithHash(dstAbs, f)
	if saveErr != nil {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_storage_save_failed, nil))
		return
	}

	// 更新 meta.json（无数据库）
	_ = storage.MetaUpsertFile(teacherDir, storage.MetaFileItem{
		RelPath:   filepath.ToSlash(filepath.Join(fileType, filename)),
		Sha256:    sha256Hex,
		Size:      info.Size,
		Ext:       ext,
		UpdatedAt: time.Now().Unix(),
	})

	resp := map[string]interface{}{
		"abs_path": dstAbs,
		"rel_path": filepath.ToSlash(filepath.Join(fileType, filename)),
		"sha256":   sha256Hex,
		"size":     info.Size,
		"ext":      ext,
	}
	ctx.JSON(apiresult.NewAPIResult(errmgr.SUCCESS, resp))
}

// ===============================
// 分析任务：创建/状态/读取结果
// ===============================

type reqAnalyzeCreate struct {
	Dept    string `json:"dept"`
	Date    string `json:"date"`
	Teacher string `json:"teacher"`

	// 可选：指定只分析某些文件（不传则分析所有支持类型文件）
	RelPaths []string `json:"rel_paths"`
}

type respAnalyzeStatus struct {
	Status      string `json:"status"`     // idle/running/success/failed
	LastError   string `json:"last_error"` // 错误信息
	UpdatedAt   int64  `json:"updated_at"`
	MdRelPath   string `json:"md_rel_path"`   // ai_result/综合分析.md
	JsonRelPath string `json:"json_rel_path"` // ai_result/各文件分析.json
}

// TeachAnalyzeCreate 创建分析任务：parsed + ai_result
// Route: POST /api/teach/analyze/create
func TeachAnalyzeCreate(ctx iris.Context) {
	var req reqAnalyzeCreate
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_http_imput_params_json_parse, nil))
		return
	}
	if strings.TrimSpace(req.Dept) == "" || strings.TrimSpace(req.Date) == "" || strings.TrimSpace(req.Teacher) == "" {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_http_input_params_empty, nil))
		return
	}

	st := config.GlobalConfig.GetStorage()
	teacherDir, code := storage.JoinUnderRoot(st.RootDir, req.Dept, req.Date, req.Teacher)
	if code != errmgr.SUCCESS {
		ctx.JSON(apiresult.NewAPIResult(code, nil))
		return
	}

	// 并发锁：同一目录同一时刻只允许一个分析任务
	if st.EnableLock {
		ok, lockErr := storage.TryLock(teacherDir)
		if lockErr != nil {
			ctx.JSON(apiresult.NewAPIResult(errmgr.Err_storage_write_failed, nil))
			return
		}
		if !ok {
			ctx.JSON(apiresult.NewAPIResult(errmgr.Err_storage_lock_conflict, nil))
			return
		}
		defer storage.Unlock(teacherDir)
	}

	// 更新状态 running
	_ = storage.MetaSetStatus(teacherDir, "running", "")

	// 1) parsed：这里先做 stub（后续可接真实解析器）
	// 说明：你可以把 pdf/docx/xlsx 转文本的逻辑放在 storage.ParseToText(...) 中
	if err := storage.ParseToText(teacherDir, req.RelPaths); err != nil {
		_ = storage.MetaSetStatus(teacherDir, "failed", err.Error())
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_storage_write_failed, nil))
		return
	}
	// 2) 转发 AI 服务（沿用你现有 aiServer.apiUrl）
	aiSrv := config.GlobalConfig.GetAiServer()

	// 这里组装给 Python 的 payload：为了兼容你现有 Python 服务，这里先提供通用字段
	payload := map[string]interface{}{
		"dept":        req.Dept,
		"date":        req.Date,
		"teacher":     req.Teacher,
		"teacher_dir": teacherDir,
		"parsed_dir":  filepath.Join(teacherDir, "parsed"),
		"result_dir":  filepath.Join(teacherDir, "ai_result"),
		"rel_paths":   req.RelPaths,
	}

	_, code2 := storage.ForwardToAiServer(aiSrv.ApiUrl, payload)
	if code2 != errmgr.SUCCESS {
		_ = storage.MetaSetStatus(teacherDir, "failed", errmgr.ErrStr(code2))
		ctx.JSON(apiresult.NewAPIResult(code2, nil))
		return
	}

	// 3) 写入 ai_result：兼容两种情况
	// A) Python 已经在 result_dir 内写好了 md/json -> Go 只校验存在
	// B) Python 返回内容，Go 落盘
	if err := os.MkdirAll(filepath.Join(teacherDir, "ai_result"), 0o755); err != nil {
		_ = storage.MetaSetStatus(teacherDir, "failed", err.Error())
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_storage_mkdir_failed, nil))
		return
	}

	// Python 负责直接落盘：ai_result/综合分析.md 与 ai_result/各文件分析.json
	mdPath := filepath.Join(teacherDir, "ai_result", "综合分析.md")
	jsonPath := filepath.Join(teacherDir, "ai_result", "各文件分析.json")

	// 允许 Python 异步写入：这里简单轮询等待一小段时间（可调）
	waitSec := st.AiResultWaitSec
	if waitSec <= 0 {
		waitSec = 120
	}
	waitOk := storage.WaitFiles(mdPath, jsonPath, time.Duration(waitSec)*time.Second)
	if !waitOk {
		_ = storage.MetaSetStatus(teacherDir, "failed", "AI结果文件未在超时内生成")
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_ai_forward_resp_parse, nil))
		return
	}

	_ = storage.MetaSetStatus(teacherDir, "success", "")

	ret := map[string]interface{}{
		"status":        "success",
		"md_rel_path":   "ai_result/综合分析.md",
		"json_rel_path": "ai_result/各文件分析.json",
	}
	ctx.JSON(apiresult.NewAPIResult(errmgr.SUCCESS, ret))
}

// TeachAnalyzeStatus 读取 meta.json 返回状态
// Route: POST /api/teach/analyze/status
func TeachAnalyzeStatus(ctx iris.Context) {
	var req reqTeachLocator
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_http_imput_params_json_parse, nil))
		return
	}
	if strings.TrimSpace(req.Dept) == "" || strings.TrimSpace(req.Date) == "" || strings.TrimSpace(req.Teacher) == "" {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_http_input_params_empty, nil))
		return
	}

	st := config.GlobalConfig.GetStorage()
	teacherDir, code := storage.JoinUnderRoot(st.RootDir, req.Dept, req.Date, req.Teacher)
	if code != errmgr.SUCCESS {
		ctx.JSON(apiresult.NewAPIResult(code, nil))
		return
	}

	meta, err := storage.MetaLoad(teacherDir)
	if err != nil {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_storage_read_failed, nil))
		return
	}

	resp := respAnalyzeStatus{
		Status:      meta.Status,
		LastError:   meta.LastError,
		UpdatedAt:   meta.UpdatedAt,
		MdRelPath:   "ai_result/综合分析.md",
		JsonRelPath: "ai_result/各文件分析.json",
	}
	ctx.JSON(apiresult.NewAPIResult(errmgr.SUCCESS, resp))
}

// TeachResultMdGet 读取综合分析.md
// Route: POST /api/teach/result/md/get
func TeachResultMdGet(ctx iris.Context) {
	var req reqTeachLocator
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_http_imput_params_json_parse, nil))
		return
	}
	if strings.TrimSpace(req.Dept) == "" || strings.TrimSpace(req.Date) == "" || strings.TrimSpace(req.Teacher) == "" {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_http_input_params_empty, nil))
		return
	}

	st := config.GlobalConfig.GetStorage()
	teacherDir, code := storage.JoinUnderRoot(st.RootDir, req.Dept, req.Date, req.Teacher)
	if code != errmgr.SUCCESS {
		ctx.JSON(apiresult.NewAPIResult(code, nil))
		return
	}

	p := filepath.Join(teacherDir, "ai_result", "综合分析.md")
	b, err := os.ReadFile(p)
	if err != nil {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_storage_read_failed, nil))
		return
	}

	ctx.JSON(apiresult.NewAPIResult(errmgr.SUCCESS, map[string]interface{}{
		"md_rel_path": "ai_result/综合分析.md",
		"content":     string(b),
	}))
}

// TeachResultJsonGet 读取各文件分析.json
// Route: POST /api/teach/result/json/get
func TeachResultJsonGet(ctx iris.Context) {
	var req reqTeachLocator
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_http_imput_params_json_parse, nil))
		return
	}
	if strings.TrimSpace(req.Dept) == "" || strings.TrimSpace(req.Date) == "" || strings.TrimSpace(req.Teacher) == "" {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_http_input_params_empty, nil))
		return
	}

	st := config.GlobalConfig.GetStorage()
	teacherDir, code := storage.JoinUnderRoot(st.RootDir, req.Dept, req.Date, req.Teacher)
	if code != errmgr.SUCCESS {
		ctx.JSON(apiresult.NewAPIResult(code, nil))
		return
	}

	p := filepath.Join(teacherDir, "ai_result", "各文件分析.json")
	b, err := os.ReadFile(p)
	if err != nil {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_storage_read_failed, nil))
		return
	}

	var obj interface{}
	_ = json.Unmarshal(b, &obj)

	ctx.JSON(apiresult.NewAPIResult(errmgr.SUCCESS, map[string]interface{}{
		"json_rel_path": "ai_result/各文件分析.json",
		"data":          obj,
	}))
}

// TeachResultVizGet 返回可绘图结构化数据（缓存到 viz.json）
// Route: POST /api/teach/result/viz/get
func TeachResultVizGet(ctx iris.Context) {
	var req reqTeachLocator
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_http_imput_params_json_parse, nil))
		return
	}
	if strings.TrimSpace(req.Dept) == "" || strings.TrimSpace(req.Date) == "" || strings.TrimSpace(req.Teacher) == "" {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_http_input_params_empty, nil))
		return
	}

	st := config.GlobalConfig.GetStorage()
	teacherDir, code := storage.JoinUnderRoot(st.RootDir, req.Dept, req.Date, req.Teacher)
	if code != errmgr.SUCCESS {
		ctx.JSON(apiresult.NewAPIResult(code, nil))
		return
	}

	vizPath := filepath.Join(teacherDir, "ai_result", "viz.json")
	if b, err := os.ReadFile(vizPath); err == nil && len(b) > 0 {
		var obj interface{}
		_ = json.Unmarshal(b, &obj)
		ctx.JSON(apiresult.NewAPIResult(errmgr.SUCCESS, obj))
		return
	}

	mdPath := filepath.Join(teacherDir, "ai_result", "综合分析.md")
	mdBytes, err := os.ReadFile(mdPath)
	if err != nil {
		ctx.JSON(apiresult.NewAPIResult(errmgr.Err_storage_read_failed, nil))
		return
	}

	vizObj := storage.BuildVizFromMarkdown(string(mdBytes))

	if st.EnableCache {
		_ = os.MkdirAll(filepath.Join(teacherDir, "ai_result"), 0o755)
		b, _ := json.MarshalIndent(vizObj, "", "  ")
		_ = storage.WriteFileAtomic(vizPath, b, 0o644)
	}

	ctx.JSON(apiresult.NewAPIResult(errmgr.SUCCESS, vizObj))
}
