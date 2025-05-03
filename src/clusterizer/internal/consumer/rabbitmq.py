
import networkx as nx

from internal.config.config import AppConfig
from internal.converter.converter import Converter
from internal.entity.response import GraphRes
from internal.entity.graph import GraphData,Link,Node
from internal.entity.document import Document as DocumentEntity


from internal.сlusterizer.graph_builder import ClusterGraphBuilder


import aio_pika
from aio_pika import IncomingMessage, Message
import logging
from typing import List, Dict
from collections import defaultdict


class RabbitMQServer:
    def __init__(
        self,
        config: AppConfig,
        logger: logging.Logger,
        s3_client,
        convertor: Converter,
            graphbuilder:ClusterGraphBuilder
    ):
        self.exchange = None
        self.queue = None
        self.channel = None
        self.connection = None
        self.config = config
        self.logger = logger
        self.s3_client = s3_client
        self.convertor = convertor
        self.graphbuilder = graphbuilder
    

    async def start(self):
        self.connection = await aio_pika.connect_robust(self.config.rabbitmq.url)
        self.channel = await self.connection.channel()
        self.queue = await self.channel.declare_queue(
            self.config.rabbitmq.consumer.queue_name,
            durable=True
        )
      
        self.logger.info("RabbitMQ consumer started")
        await self.queue.consume(self.handle_message)

    async def stop(self):
        if self.connection:
            await self.connection.close()
            self.logger.info("RabbitMQ consumer stopped")

    def __get_files_from_s3(self, filenames: List[str]) -> Dict[str, DocumentEntity]:
        res: Dict[str, DocumentEntity] = defaultdict()
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
                continue

            self.logger.info(f"Fetched {filename}, {len(response['Body'].read())} bytes")
            res[filename] = data
        return res

    async def handle_message(self, message: IncomingMessage):
        try:
            raw_message = message.body
            req = self.convertor.parse_message(raw_message)

            self.logger.info(f"Received messages: {req.keys}")

            id_texts = self.__get_files_from_s3(req.keys)

            # clustered_texts = self.clusterizer.do(id_texts,group_count = req.group_count)
            clustered_texts = self.graphbuilder.build_cluster_graph(id_texts)

            id = req.keys[0].split("_")[0]
            res = self.__make_res(id,clustered_texts)
            await self.send_message(
                self.config.rabbitmq.producer.queue_name,
                res)

            self.logger.info(f"Message processed and clustered successfully for: {req}")
            await message.ack()

        except Exception as e:
            self.logger.error(f"Failed to process message: {e}")
            await message.nack(requeue=True)

    async def send_message(self, routing_key: str, data: GraphRes):
        try:
            body = data.json().encode()
            print(body)
            message = Message(body)
            await self.channel.default_exchange.publish(
                message,
                routing_key=routing_key
            )
            self.logger.info(f"Message published to {routing_key}")
        except Exception as e:
            self.logger.error(f"Failed to publish message: {e}")

    def __make_res(self,graph_id: str, graph: nx.Graph) -> GraphRes:
        nodes = []
        links = []

        for node_id, attr in graph.nodes(data=True):
            node = Node(
                id=node_id,
                title=attr.get("title", ""),
                cluster=attr.get("cluster", -1),
                type = attr.get("type", -1)
            )
            nodes.append(node)

        for source, target, attr in graph.edges(data=True):
            link = Link(
                source=source,
                target=target,
                weight=attr.get("weight", 1.0)
            )
            links.append(link)

        graph_data = GraphData(
            directed=nx.is_directed(graph),
            multigraph=graph.is_multigraph(),
            graph={},
            nodes=nodes,
            links=links
        )

        return GraphRes(id=graph_id, graph=graph_data)

