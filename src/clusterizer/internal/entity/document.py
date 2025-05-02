
from typing import Dict, List, Any

from pydantic import BaseModel


class Document(BaseModel):
    text:str
    type:str