from fastapi import FastAPI
from pydantic import BaseModel, Field
from typing import List, Optional, Dict, Any
import os, json, time, glob
import requests

app = FastAPI(title="AI Task Server", version="1.0.0")

def call_llm(api_base: str, api_key: str, model: str, prompt: str) -> str:
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
    r = requests.post(url, headers=headers, json=payload, timeout=120)
    r.raise_for_status()
    data = r.json()
    return data["choices"][0]["message"]["content"]
    
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

        # =========================
    # 1) 组装 Prompt（控制长度：只取每个文件前 N 字）
    # =========================
    max_chars_per_file = 8000
    prompt_parts = []
    prompt_parts.append("请对以下教学材料进行分析，并输出：\n"
                        "1) 综合摘要\n"
                        "2) 关键点（条目化）\n"
                        "3) 可视化数据（用JSON代码块输出，字段包含 keywords(数组)、numbers(数组)、relations(数组，元素含source/target/label)）\n"
                        "要求：先输出Markdown正文，再输出一个 ```json ... ``` 的可视化JSON代码块。\n\n")

    for d in docs:
        text = d["text"][:max_chars_per_file]
        prompt_parts.append(f"【文件】{d['file']}\n{text}\n\n")

    prompt = "".join(prompt_parts)

    # =========================
    # 2) 调用大模型：生成综合Markdown
    # =========================
    api_base = os.getenv("LLM_API_BASE", "https://api.openai.com/v1")
    api_key = os.getenv("LLM_API_KEY", "")
    model = os.getenv("LLM_MODEL", "gpt-4o-mini")

    if not api_key:
        # 没配置 key 就直接报错，避免你以为调用成功
        raise RuntimeError("LLM_API_KEY is empty")

    md_content = call_llm(api_base=api_base, api_key=api_key, model=model, prompt=prompt)

    # =========================
    # 3) 各文件分析.json：先做一个最简版（后续可让模型也输出每文件JSON）
    # =========================
    files_json = {
        "dept": req.dept,
        "date": req.date,
        "teacher": req.teacher,
        "generated_at": int(time.time()),
        "llm": {"api_base": api_base, "model": model},
        "files": [
            {
                "file": d["file"],
                "text_len": len(d["text"]),
            } for d in docs
        ],
    }
    md_path = os.path.join(req.result_dir, "综合分析.md")
    json_path = os.path.join(req.result_dir, "各文件分析.json")

    with open(md_path, "w", encoding="utf-8") as f:
        f.write(md_content)

    with open(json_path, "w", encoding="utf-8") as f:
        json.dump(files_json, f, ensure_ascii=False, indent=2)

    return {"ok": True, "md_path": md_path, "json_path": json_path}
