import logging
from typing import List, Dict
from internal.entity.document import Document as DocumentEntity
from internal.converter.converter import Converter

from internal.config.config import S3 as S3Config

from collections import defaultdict

class S3:
    def __init__(self,
                 s3_client,
                  logger: logging.Logger,
                 config:S3Config,
                 convertor:Converter,
                 ):
         self.s3_client = s3_client
         self.logger = logger
         self.config = config
         self.convertor = convertor

    def get_files_by_ids(self, ids:List[str]):
        res: Dict[str, DocumentEntity] = defaultdict()
        print(1)
        for id in ids:
            self.logger.info(f"Fetching {id} from S3...")

            response = self.s3_client.get_object(
                Bucket=self.config.bucket,
                Key=id
            )

            body = response["Body"].read()
            data = self.convertor.file_to_str(body)
            if data is None:
                self.logger.error(f"Unknown type {id}, {len(response['Body'].read())} bytes")
                continue

            self.logger.info(f"Fetched {id}, {len(response['Body'].read())} bytes")
            res[id] = data
        print(1)

        return res