import asyncio
import logging

import boto3
from fastapi import FastAPI
from dotenv import load_dotenv

import os
from fastapi.middleware.cors import CORSMiddleware
from internal.config.config import load_config
from internal.converter.converter import Converter

from internal.consumer.rabbitmq import RabbitMQServer

from internal.сlusterizer.graph_builder import ClusterGraphBuilder


load_dotenv()
cfg = load_config(os.getenv('CONFIG_PATH'))

app = FastAPI(title='main')
#loggers
logging.basicConfig(level=logging.INFO)



rabbitmq_logger = logging.getLogger("rabbitmq-consumer")
logging.basicConfig(level=logging.INFO)

converter_logger = logging.getLogger("converter")
converter = Converter(converter_logger)
s3_client = boto3.client(
    "s3"
)

graphbuilder = ClusterGraphBuilder()

rabbitmq_url = cfg.rabbitmq.url
rabbitmqServer = RabbitMQServer(
    config=cfg,
    logger=rabbitmq_logger,
    s3_client=s3_client,
    convertor=converter,
graphbuilder = graphbuilder)

origins = ["localhost:8000"]
app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.on_event("startup")
async def startup_event():
    # asyncio.create_task(kafkaServer.start())
    asyncio.create_task(rabbitmqServer.start())

@app.get("/health")
def health():
    return {"status": "ok"}




