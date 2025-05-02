import json
import logging
from io import BytesIO
import pdfplumber
from docx import Document

from internal.entity.request import Request
from internal.entity.document import Document as DocumentEntity


class Converter:
    def __init__(self,logger:logging.Logger):
        self.logger = logger
    def file_to_str(self, data:bytes)->DocumentEntity:
        self.logger.info(f"convert file to str")

        res = None
        if data[:4] == b"%PDF":
            try:
                with pdfplumber.open(BytesIO(data)) as pdf:
                    text = "\n".join([
                        page.extract_text() or "" for page in pdf.pages
                    ])
                    res = DocumentEntity(
                        text = text,
                        type = "pdf"
                    )
            except Exception as e:
                self.logger.error(f"[PDF ERROR] {e}")
        elif data[:2] == b"PK":
            try:
                doc = Document(BytesIO(data))
                text =  "\n".join([
                    para.text for para in doc.paragraphs if para.text.strip()
                ])
                res = DocumentEntity(
                    text=text,
                    type="docx"
                )
            except Exception as e:
                self.logger.error(f"[DOCX ERROR] {e}")
        return res

    def byte_to_list_str(self, data: bytes) -> list[str]:
        try:
            filenames = json.loads(data.decode("utf-8"))
            if isinstance(filenames, list):
                return filenames
            elif isinstance(filenames, dict) and "files" in filenames:
                return filenames["files"]
            else:
                raise ValueError("Unexpected message format")
        except Exception as e:
            return []

    def parse_message(self,raw_message: bytes) -> Request:
        try:
            data = json.loads(raw_message.decode("utf-8"))
            request = Request(**data)
            return request
        except Exception as e:
            raise ValueError(f"Failed to parse message: {e}")


