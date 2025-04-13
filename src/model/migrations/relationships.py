import enum
from datetime import datetime,timedelta

from sqlalchemy import MetaData, Table, Column, Integer, String, TIMESTAMP, ForeignKey, Boolean,DateTime
from sqlalchemy.orm import DeclarativeBase, relationship
class Base(DeclarativeBase):
    metadata = MetaData()
# class Base(DeclarativeBase):
#     metadata = MetaData()


item_order = Table(
    "item_order",
    Base.metadata,
    Column("item_id", ForeignKey("items.id"), primary_key=True),
    Column("order_id", ForeignKey("orders.id"), primary_key=True),
)

item_category = Table(
    "item_category_ref",
    Base.metadata,
    Column("item_id", ForeignKey("items.id"), primary_key=True),
    Column("category_id", ForeignKey("categories.id"), primary_key=True),
)

user_item_recommendation = Table(
    "user_item_ref_recommendation",
    Base.metadata,
    Column("item_id", ForeignKey("items.id"), primary_key=True),
    Column("user_id", ForeignKey("customers.id"), primary_key=True),
)

user_item = Table(
    "user_item_ref",
    Base.metadata,
    Column("item_id", ForeignKey("items.id"), primary_key=True),
    Column("user_id", ForeignKey("customers.id"), primary_key=True),
)

