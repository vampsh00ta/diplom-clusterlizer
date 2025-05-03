from typing import List

from pydantic import BaseModel

from internal.entity.graph import Group, GraphData


class  ClusterizationRes(BaseModel):
    id:str
    groups:List[Group]




class  GraphRes(BaseModel):
    id:str
    graph:GraphData