

from typing import List

from pydantic import BaseModel

class Group(BaseModel):
    keys:str
    ids:List[str]




class Request(BaseModel):
    keys:List[str]

