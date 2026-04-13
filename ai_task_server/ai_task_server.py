from fastapi import FastAPI
from pydantic import BaseModel, Field
from typing import List, Optional, Dict, Any
import os, json, time, glob

app = FastAPI(title="AI Task Server", version="1.0.0")


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

    # TODO: 在这里接入真正的大模型（OpenAI/通义/智谱/本地模型等）
    # 先 stub：产出可用的 md + json
    keywords = []
    for d in docs:
        words = [w for w in d["text"].replace("\n", " ").split(" ") if w.strip()]
        keywords.extend(words[:10])
    keywords = keywords[:30]

    md_lines = [
        f"# 综合分析（{req.dept}/{req.date}/{req.teacher}）",
        "",
        "## 摘要",
        f"- 文本文件数量：{len(docs)}",
        f"- 关键词示例：{', '.join(keywords[:10])}",
        "",
        "## 文件列表",
    ]
    for d in docs:
        md_lines.append(f"- {d['file']}")
    md_content = "\n".join(md_lines) + "\n"

    per_file = []
    for d in docs:
        per_file.append(
            {
                "file": d["file"],
                "summary": (d["text"][:200] + "...") if len(d["text"]) > 200 else d["text"],
                "keywords": keywords[:10],
            }
        )

    files_json = {
        "dept": req.dept,
        "date": req.date,
        "teacher": req.teacher,
        "generated_at": int(time.time()),
        "files": per_file,
    }

    md_path = os.path.join(req.result_dir, "综合分析.md")
    json_path = os.path.join(req.result_dir, "各文件分析.json")

    with open(md_path, "w", encoding="utf-8") as f:
        f.write(md_content)

    with open(json_path, "w", encoding="utf-8") as f:
        json.dump(files_json, f, ensure_ascii=False, indent=2)

    return {"ok": True, "md_path": md_path, "json_path": json_path}
