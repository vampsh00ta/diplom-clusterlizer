from typing import List

from pydantic import BaseModel

class Group(BaseModel):
    keys:str
    ids:List[str]

class  ClusterizationRes(BaseModel):
    id:str
    groups:List[Group]


class Request(BaseModel):
    group_count:int
    keys:List[str]