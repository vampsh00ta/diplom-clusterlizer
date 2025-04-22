
from internal.config.config import AppConfig
from internal.converter.converter import Converter

import aio_pika
from aio_pika import IncomingMessage, Message
import logging
from typing import List


class RabbitMQServer:
    def __init__(
        self,
        config: AppConfig,
        logger: logging.Logger,
        rabbitmq_url: str,
        s3_client,
        convertor: Converter
    ):
        self.config = config
        self.logger = logger
        self.rabbitmq_url = rabbitmq_url
        self.s3_client = s3_client
        self.convertor = convertor
        self.connection = None
        self.channel = None
        self.queue = None

    async def start(self):
        self.connection = await aio_pika.connect_robust(self.rabbitmq_url)
        self.channel = await self.connection.channel()
        self.queue = await self.channel.declare_queue(
            self.config.rabbitmq.queue_name,
            durable=True
        )
        self.logger.info("RabbitMQ consumer started")
        await self.queue.consume(self.handle_message)

    async def stop(self):
        if self.connection:
            await self.connection.close()
            self.logger.info("RabbitMQ consumer stopped")

    def __get_files_from_s3(self, filenames: List[str]) -> List[str]:
        res = []
        for filename in filenames:
            self.logger.info(f"Fetching {filename} from S3...")

            response = self.s3_client.get_object(
                Bucket=self.config.s3.bucket,
                Key=filename
            )

            body = response["Body"].read()
            data = self.convertor.file_to_str(body)
            if data is None:
                self.logger.error(f"Unknown type {filename}, {len(response['Body'].read())} bytes")
            else:
                res.append(data)
                self.logger.info(f"Fetched {filename}, {len(response['Body'].read())} bytes")
        return res

    async def handle_message(self, message: IncomingMessage):
        async with message.process():  # auto ack after block if no exception
            try:
                raw_message = message.body
                filenames = self.convertor.byte_to_list_str(raw_message)

                self.logger.info(f"Received messages: {filenames}")

                files = self.__get_files_from_s3(filenames)
                summary_files_len = sum(len(f) for f in files)
                self.logger.info(f"Summary files lengths {summary_files_len}")

            except Exception as e:
                self.logger.error(f"Failed to process message: {e}")
                # optionally: requeue message
                await message.nack(requeue=True)

