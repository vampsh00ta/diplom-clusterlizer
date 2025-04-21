import logging
import time
from typing import List

import boto3
from aiokafka import AIOKafkaConsumer, TopicPartition
from internal.config.config import AppConfig
from internal.converter.converter import Converter

# from internal.model.model import MLModel





class KafkaServer:
    def __init__(self,
             config: AppConfig,
            logger: logging.Logger,
             consumer: AIOKafkaConsumer,
             s3_client,
            convertor:Converter):
        self.config = config
        self.consumer = consumer
        self.s3_client = s3_client
        self.convertor = convertor
        self.logger = logger

    async def start(self):

        await self.consumer.start()
        self.logger.info("Kafka consumer started")
        await self.consume_loop()

    async def stop(self):
        if self.consumer:
            await self.consumer.stop()
            self.logger.info("Kafka consumer stopped")

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
                self.logger.info(f"Fetched {filename}, {len( response['Body'].read())} bytes")
        return res

    async def consume_loop(self):
        try:
            async for msg in self.consumer:
                await self.handle_message(msg)
        except Exception as e:
            self.logger.error(f"Error in consume loop: {e}")
        finally:
            await self.stop()

    async def handle_message(self, msg):
        try:
            raw_message = msg.value
            filenames = self.convertor.byte_to_list_str(raw_message)

            self.logger.info(f"Received messages: {filenames}")

            files = self.__get_files_from_s3(filenames)
            summary_files_len = 0
            for file in files:
                summary_files_len+=len(file)
            self.logger.info(f"Summary files lengths {summary_files_len}")
            # time.sleep(10)
            tp = TopicPartition(msg.topic, msg.partition)
            await self.consumer.commit({tp: msg.offset + 1})
            self.logger.info(f"Offset committed for {filenames}")

        except Exception as e:
            self.logger.error(f"Failed to process message: {e}")
