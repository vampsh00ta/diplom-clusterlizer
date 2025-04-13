import asyncio
from typing import AsyncGenerator

from sqlalchemy import NullPool

from migrations.relationships import Base
from sqlalchemy.ext.asyncio import async_sessionmaker, create_async_engine, AsyncSession

DATABASE_URL = f"sqlite+aiosqlite:///./test.db"

engine = create_async_engine(DATABASE_URL,poolclass=NullPool)
TestingSessionLocal = async_sessionmaker(engine, expire_on_commit=False)



async def override_get_async_session() -> AsyncGenerator[AsyncSession, None]:
    async with TestingSessionLocal() as session:
        yield session





