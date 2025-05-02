import json
import logging
import re
from typing import Dict, List, Optional
from internal.сlusterizer.stop_words import stopwords
import networkx as nx
import numpy as np
import hdbscan
import yake
from networkx.readwrite import json_graph
from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity
import inspect
from internal.entity.document import Document as DocumentEntity

if not hasattr(inspect, 'getargspec'):
    inspect.getargspec = inspect.getfullargspec

import pymorphy3
from razdel import tokenize



class ClusterGraphBuilder:
    """
    Класс для кластеризации текстовых документов и построения графа сходства.

    Шаги:
      1. Предобработка текста (очистка, лемматизация, удаление стоп-слов).
      2. Получение эмбеддингов через SentenceTransformer.
      3. Кластеризация через HDBSCAN.
      4. Построение графа NetworkX на основе косинусного сходства.
      5. Генерация названий с помощью YAKE/KeyBERT.
      6. Экспорт графа в JSON формате node-link.
    """

    def __init__(
        self,
        embed_model: str = 'ai-forever/sbert_large_nlu_ru',
        keyword_model: Optional[str] = None,
        yake_top: int = 3,
        yake_ngram: int = 1,
        hdbscan_min_cluster_size: int = 2,
        hdbscan_min_samples: int = 1,
        sim_threshold: float = 0.75,
        enable_logging: bool = True
    ):
        # Настройка логирования
        if enable_logging:
            logging.basicConfig(level=logging.INFO)
        self.logger = logging.getLogger(__name__)

        # Инициализация моделей
        self.model = SentenceTransformer(embed_model)
        self.yake_extractor = yake.KeywordExtractor(
            lan="ru", n=yake_ngram, top=yake_top
        )

        # Параметры кластеризации и графа
        self.cluster_params = dict(
            min_cluster_size=hdbscan_min_cluster_size,
            min_samples=hdbscan_min_samples
        )
        self.sim_threshold = sim_threshold

        # Подготовка морфологического анализатора
        self.morph = pymorphy3.MorphAnalyzer()
        self.logger.info("Russian preprocessing is enabled.")




    def preprocess_text(self, text: str) -> str:
        """
        Очищает, лемматизирует и удаляет стоп-слова из текста.
        Если отсутствуют razdel или pymorphy2, выполняется базовая очистка.
        """
        # Нормализация и удаление лишних символов
        text = text.lower()
        text = re.sub(r'https?://\S+|www\.\S+', '', text)
        text = re.sub(r'[^а-яё\s]', ' ', text)
        text = re.sub(r'\s+', ' ', text).strip()



        # Токенизация и лемматизация
        lemmas = []
        for token in tokenize(text):
            word = token.text
            if word in stopwords or len(word) < 2:
                continue
            parse = self.morph.parse(word)[0]
            lemma = parse.normal_form
            if lemma not in stopwords:
                lemmas.append(lemma)
        return ' '.join(lemmas)

    def _get_embeddings(self, texts: List[str]) -> np.ndarray:
        """Получение эмбеддингов для списка текстов."""
        self.logger.info("Computing embeddings for %d documents", len(texts))
        return self.model.encode(texts, show_progress_bar=False)

    def generate_titles(self, id_texts: Dict[str, str]) -> Dict[str, str]:
        """Генерация названий для документов по ключевым словам."""
        titles: Dict[str, str] = {}
        for doc_id, text in id_texts.items():
            if not text.strip():
                titles[doc_id] = "Без названия"
                continue
            try:
                keywords = self.yake_extractor.extract_keywords(text)
                if keywords:
                    keyword, _ = min(keywords, key=lambda x: x[1])
                    titles[doc_id] = keyword.capitalize()
                else:
                    titles[doc_id] = "Без названия"
            except Exception as e:
                self.logger.error("Error generating title for %s: %s", doc_id, e)
                titles[doc_id] = "Без названия"
        return titles

    def export_graph_to_json(self, graph: nx.Graph) -> dict:
        """Экспорт графа в JSON node-link формат."""
        data = json_graph.node_link_data(graph)
        # data.
        return data

    def build_cluster_graph(
        self,
        id_texts: Dict[str, DocumentEntity],
        preprocess: bool = True,
        threshold: Optional[float] = None,
        **cluster_kwargs
    ) -> nx.Graph:
        """
        Основной метод: строит граф кластеров по входным текстам.

        Args:
            id_texts: Словарь id -> текст.
            preprocess: Применять ли предобработку.
            threshold: Порог косинусного сходства для ребер.
            cluster_kwargs: Переопределение параметров HDBSCAN.

        Returns:
            nx.Graph с атрибутами 'cluster' и 'title' на узлах и весами на ребрах.
        """
        ids = list(id_texts.keys())
        raw_texts = []
        for document in id_texts.values():
            raw_texts.append(document.text)


        self.logger.info("Applying preprocessing to texts...")
        texts = [self.preprocess_text(t) for t in raw_texts]


        # Встраивания и кластеризация
        embeddings = self._get_embeddings(texts)
        params = {**self.cluster_params, **cluster_kwargs}
        clusterer = hdbscan.HDBSCAN(**params)
        labels = clusterer.fit_predict(embeddings)

        # Матрица сходства
        sim_matrix = cosine_similarity(embeddings)  # импортируйте из sklearn

        graph = nx.Graph()
        titles = self.generate_titles({id_: text for id_, text in zip(ids, texts)})

        # Добавляем узлы
        for idx, doc_id in enumerate(ids):
            graph.add_node(
                doc_id,
                cluster=int(labels[idx]),
                title=titles.get(doc_id, "Без названия"),
                type  = id_texts[doc_id].type

            )

        # Добавляем ребра
        thr = threshold if threshold is not None else self.sim_threshold
        self.logger.info("Building edges with similarity threshold %.2f", thr)
        n = len(ids)
        for i in range(n):
            for j in range(i + 1, n):
                sim = float(sim_matrix[i, j])
                if sim >= thr:
                    graph.add_edge(ids[i], ids[j], weight=sim)

        self.logger.info("Generated %d nodes and %d edges", graph.number_of_nodes(), graph.number_of_edges())
        return graph


# if __name__ == '__main__':
#     # Пример использования
#     sample = {
#         'paper1': "Исследование глубокого обучения в задачах компьютерного зрения.",
#         'paper2': "Обзор методов кластеризации временных рядов.",
#         'paper3': "Применение нейронных сетей в анализе текста",
#     }
#     builder = ClusterGraphBuilder()
#     graph = builder.build_cluster_graph(sample)
#     print(builder.export_graph_to_json(graph))
