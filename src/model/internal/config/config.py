import yaml
from pydantic import BaseModel
from pathlib import Path

class Kafka(BaseModel):
    bootstrap_servers: str = "localhost:9092"
    topic: str = "files-topic"
    group_id: str = "ml-model-consumer"
    max_poll_interval_ms: int = 100

class RabbitMQBase(BaseModel):
    queue_name: str = "ml_jobs"
    exchange: str = "diplom"
class RabbiqMQ(BaseModel):
    url:str = "amqp://guest:guest@localhost/"
    consumer:RabbitMQBase
    producer:RabbitMQBase


class S3(BaseModel):
    bucket: str

class AppConfig(BaseModel):
    rabbitmq: RabbiqMQ
    kafka: Kafka
    s3: S3

def load_config(path: str = "config/config.yaml") -> AppConfig:
    with open(Path(path), "r") as f:
        data = yaml.safe_load(f)
    return AppConfig(**data)