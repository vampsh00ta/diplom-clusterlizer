

import pytest
from main import app
from migrations.relationships import Base
from tests.db.database import override_get_async_session, engine
from utils.database import get_async_session

from httpx import AsyncClient

app.dependency_overrides[get_async_session] =   override_get_async_session

client = AsyncClient(app=app, base_url='http://127.0.0.1:8000/auth/')

@pytest.fixture
def anyio_backend():
    return 'asyncio'


@pytest.mark.anyio
async def test_register_success():
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)

    username = "op1um_test"
    password = "op1um_test"
    email = 'opium_test'
    response  = await client.post('/register',json ={'username':username,'password':password,"email":email} )
    assert response.status_code == 201
    assert response.json()['token_type'] == 'bearer'


@pytest.mark.anyio
async def test_login_success():
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)

    username = "op1um_test"
    password = "op1um_test"
    response = await client.post('/login', data={'username': username, 'password': password})
    assert response.status_code == 200
    assert response.json()['token_type'] == 'bearer'

@pytest.mark.anyio
async def test_login_fake_user():
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)

    username = "fake"
    password = "fake"
    response = await client.post('/login', data={'username': username, 'password': password})
    assert response.status_code == 401

@pytest.mark.anyio
async def test_get_curr_user():
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
    username = "op1um_test"
    response = await client.get('/user')
    assert response.status_code == 200
    assert response.json()['username'] == username