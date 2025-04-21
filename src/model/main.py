import asyncio
import logging

import boto3
from aiokafka import AIOKafkaConsumer
from fastapi import FastAPI

import os
from fastapi.middleware.cors import CORSMiddleware
from internal.config.config import load_config
from internal.converter.converter import Converter

from internal.kafka.kafka import KafkaServer





cfg = load_config(os.getenv('CONFIG_PATH'))
app = FastAPI(title='main')

#loggers
logging.basicConfig(level=logging.INFO)

kafka_logger = logging.getLogger("kafka-consumer")
logging.basicConfig(level=logging.INFO)

converter_logger = logging.getLogger("converter")
converter = Converter(converter_logger)

s3_client = boto3.client(
    "s3"
)
kafka_consumer = AIOKafkaConsumer(
    cfg.kafka.topic,
    bootstrap_servers=   cfg.kafka.bootstrap_servers,
    group_id= cfg.kafka.group_id,
    max_poll_interval_ms=cfg.kafka.max_poll_interval_ms,
    enable_auto_commit=False,
)


kafkaServer = KafkaServer(cfg,
                          kafka_logger,
                          kafka_consumer,
                          s3_client,
                          converter)


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
    asyncio.create_task(kafkaServer.start())

@app.get("/health")
def health():
    return {"status": "ok"}




