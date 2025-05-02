

from typing import List, Dict, Any

from pydantic import BaseModel

class Group(BaseModel):
    keys:str
    ids:List[str]




class Request(BaseModel):
    group_count:int
    keys:List[str]

