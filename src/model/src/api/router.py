
from fastapi import APIRouter, Depends, Request
from sqlalchemy import Column, Integer, Float, JSON, DateTime
from sqlalchemy.sql import func
from pydantic import BaseModel
class Test(BaseModel):
    test: str



router = APIRouter(tags=["test"])

@router.get("/test")
async def api_metrics():
    return {"test":"test"}
