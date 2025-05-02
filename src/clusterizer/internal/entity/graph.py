from typing import Dict, List, Any

from pydantic import BaseModel


class Node(BaseModel):
    id: str
    title: str
    cluster: int
    type:str

class Link(BaseModel):
    source: str
    target: str
    weight: float

class GraphData(BaseModel):
    directed: bool
    multigraph: bool
    graph: Dict[str, Any]
    nodes: List[Node]
    links: List[Link]




class Group(BaseModel):
    keys: str
    ids: List[str]