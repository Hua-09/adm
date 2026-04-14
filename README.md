# adm backend

本仓库只包含两个后端服务：

- `adm_bkd/`：Go + Iris 教学资料后端（`/api/teach/*`）
- `ai_task_server/`：Python + FastAPI 模型任务服务（`/api/run_model`）

## 目录结构

```text
adm/
├─ README.md
├─ adm_bkd/
└─ ai_task_server/
```

## 1) 运行 Python AI 服务

### Linux / macOS

```bash
cd /home/runner/work/adm/adm/ai_task_server
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
export LLM_API_BASE="https://api.openai.com/v1"
export LLM_API_KEY="your_api_key"
export LLM_MODEL="gpt-4o-mini"
uvicorn ai_task_server:app --host 0.0.0.0 --port 6678
```

### Windows (PowerShell)

```powershell
cd C:\path\to\adm\ai_task_server
python -m venv .venv
.\.venv\Scripts\Activate.ps1
pip install -r requirements.txt
$env:LLM_API_BASE="https://api.openai.com/v1"
$env:LLM_API_KEY="your_api_key"
$env:LLM_MODEL="gpt-4o-mini"
uvicorn ai_task_server:app --host 0.0.0.0 --port 6678
```

## 2) 运行 Go 后端服务

默认配置文件：`adm_bkd/config/config.yaml`

- `storage.rootDir` 默认：`/data/teaching_repo`
- `storage.aiResultWaitSec`：等待 `ai_result/综合分析.md` 和 `ai_result/各文件分析.json` 的超时秒数
- `aiServer.apiUrl` 默认：`http://127.0.0.1:6678/api/run_model`

### Linux / macOS

```bash
mkdir -p /data/teaching_repo
cd /home/runner/work/adm/adm/adm_bkd
go mod tidy
go run .
```

### Windows (PowerShell)

```powershell
cd C:\path\to\adm\adm_bkd
# 建议把 config.yaml 的 storage.rootDir 改为本机可写目录，例如 D:/data/teaching_repo
go mod tidy
go run .
```

## 3) 主要接口

- `POST /api/teach/dept/list`
- `POST /api/teach/date/list`
- `POST /api/teach/teacher/list`
- `POST /api/teach/file/list`
- `POST /api/teach/file/upload` (multipart)
- `POST /api/teach/analyze/create`
- `POST /api/teach/analyze/status`
- `POST /api/teach/result/md/get`
- `POST /api/teach/result/json/get`
- `POST /api/teach/result/viz/get`

## 4) 解析依赖说明（可选）

`adm_bkd` 的解析逻辑会尝试：

- `pdftotext`（PDF）
- `pandoc`（DOCX）
- `excelize`（XLSX/XLS）

未安装外部命令时会跳过对应文件，服务仍可运行。
