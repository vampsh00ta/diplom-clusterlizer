import asyncio

from fastapi import FastAPI


from src.api.router  import router
from fastapi.middleware.cors import CORSMiddleware



app = FastAPI(title='main')

app.include_router(router)


origins = ["localhost:8000"]
app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,

    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.on_event('startup')
async  def startup_event_producer():

    loop = asyncio.get_event_loop()



