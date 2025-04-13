from pydantic import BaseModel

class Token(BaseModel):
    access_token:str
    token_type:str = 'bearer'
class CreateUser(BaseModel):
    username:str
    email:str
    password:str
class AnonymousUserRead(BaseModel):
    id: int
    username: str

class OrderRead(BaseModel):
    id:int
    user_id:int
class UserRead(BaseModel):
    id: int
    email: str
    username: str
    is_active: bool = True
    is_superuser: bool = False
    is_verified: bool = False
    class Config:
        orm_mode = True

class User_Change_Email(BaseModel):
    email: str
    password: str

