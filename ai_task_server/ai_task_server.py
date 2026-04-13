from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field
from typing import List, Optional, Dict, Any
import os, json, time, glob, re
import requests

app = FastAPI(title="AI Task Server", version="1.0.0")


def call_llm(api_base: str, api_key: str, model: str, prompt: str) -> str:
    """
    OpenAI 兼容 /chat/completions 调用
    - 适用于 OpenAI 官方 / 以及多数兼容网关（需自行填 api_base）
    """
    url = api_base.rstrip("/") + "/chat/completions"
    headers = {"Authorization": f"Bearer {api_key}", "Content-Type": "application/json"}
    payload = {
        "model": model,
        "messages": [
            {"role": "system", "content": "你是一个文档分析助手，请输出结构化总结。"},
            {"role": "user", "content": prompt},
        ],
        "temperature": 0.2,
    }
    r = requests.post(url, headers=headers, json=payload, timeout=180)
    r.raise_for_status()
    data = r.json()
    return data["choices"][0]["message"]["content"]


def extract_json_block(text: str) -> Optional[dict]:
    """
    从模型输出中提取 ```json ... ``` 代码块
    返回 dict 或 None
    """
    m = re.search(r"```json\s*(\{.*?\})\s*```", text, flags=re.S | re.I)
    if not m:
        return None
    raw = m.group(1)
    try:
        return json.loads(raw)
    except Exception:
        return None


class RunModelReq(BaseModel):
    dept: str = Field(..., description="部门，如 自动化系")
    date: str = Field(..., description="日期目录，如 2025-12-04")
    teacher: str = Field(..., description="老师目录名，如 张三老师")

    teacher_dir: str = Field(..., description="老师目录绝对路径")
    parsed_dir: str = Field(..., description="解析后的txt目录绝对路径")
    result_dir: str = Field(..., description="输出结果目录绝对路径")

    rel_paths: Optional[List[str]] = Field(default=None)


@app.post("/api/run_model")
def run_model(req: RunModelReq) -> Dict[str, Any]:
    """
    约定（与你 Go 端一致）：
    - 必须在 req.result_dir 内写出：
      1) 综合分析.md（必须非空）
      2) 各文件分析.json（必须非空）
    - Go 端只检查文件是否生成，不依赖本接口返回体结构
    """
    os.makedirs(req.result_dir, exist_ok=True)

    # 读取 parsed/*.txt
    txt_files = sorted(glob.glob(os.path.join(req.parsed_dir, "*.txt")))
    docs = []
    for p in txt_files:
        try:
            with open(p, "r", encoding="utf-8", errors="ignore") as f:
                docs.append({"file": os.path.basename(p), "text": f.read()})
        except Exception:
            continue

    if not docs:
        docs = [{"file": "empty.txt", "text": "No parsed text found. Please check parsed_dir."}]

    # 组 prompt（限制长度）
    max_chars_per_file = int(os.getenv("MAX_CHARS_PER_FILE", "8000"))
    prompt_parts = []
    prompt_parts.append(
        "请对以下教学材料进行分析，并输出Markdown，包含：\n"
        "1) 综合摘要\n"
        "2) 关键点（条目化）\n"
        "3) 风险/问题\n"
        "4) 建议\n"
        "最后请输出一个 ```json``` 代码块，JSON字段包含：\n"
        "- keywords: string[]\n"
        "- numbers: string[]\n"
        "- relations: {source:string,target:string,label:string}[]\n"
        "注意：json 必须是严格可解析的 JSON。\n\n"
    )
    for d in docs:
        text = d["text"][:max_chars_per_file]
        prompt_parts.append(f"【文件】{d['file']}\n{text}\n\n")
    prompt = "".join(prompt_parts)

    # 从环境变量读取 LLM 配置
    api_base = os.getenv("LLM_API_BASE", "https://api.openai.com/v1")
    api_key = os.getenv("LLM_API_KEY", "")
    model = os.getenv("LLM_MODEL", "gpt-4o-mini")

    if not api_key.strip():
        raise HTTPException(status_code=500, detail="LLM_API_KEY is empty (set env var before starting server)")

    # 调用大模型生成 Markdown
    try:
        md_content = call_llm(api_base=api_base, api_key=api_key, model=model, prompt=prompt)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"call_llm failed: {e}")

    # 从 md 中尝试提取 json 代码块（可选，但强烈建议）
    viz = extract_json_block(md_content) or {"keywords": [], "numbers": [], "relations": []}

    # 生成 各文件分析.json（包含：文件清单 + LLM信息 + viz）
    files_json = {
        "dept": req.dept,
        "date": req.date,
        "teacher": req.teacher,
        "generated_at": int(time.time()),
        "llm": {"api_base": api_base, "model": model},
        "files": [{"file": d["file"], "text_len": len(d["text"])} for d in docs],
        "viz": viz,
    }

    md_path = os.path.join(req.result_dir, "综合分析.md")
    json_path = os.path.join(req.result_dir, "各文件分析.json")

    with open(md_path, "w", encoding="utf-8") as f:
        f.write(md_content if md_content.endswith("\n") else (md_content + "\n"))

    with open(json_path, "w", encoding="utf-8") as f:
        json.dump(files_json, f, ensure_ascii=False, indent=2)

    # 额外写一个 viz.json（可选）：给前端/Go 直接读取更方便
    viz_path = os.path.join(req.result_dir, "viz.json")
    with open(viz_path, "w", encoding="utf-8") as f:
        json.dump(viz, f, ensure_ascii=False, indent=2)

    return {"ok": True, "md_path": md_path, "json_path": json_path, "viz_path": viz_path}
