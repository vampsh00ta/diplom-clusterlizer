import logging
import re
from typing import Dict, List, Optional, Tuple
from internal.сlusterizer.stop_words import stopwords
import networkx as nx
import numpy as np
import hdbscan
import yake
from networkx.readwrite import json_graph
from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity
from internal.entity.document import Document as DocumentEntity
from keybert import KeyBERT
from sklearn.metrics import silhouette_score, davies_bouldin_score, silhouette_samples
import matplotlib.pyplot as plt

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
        if enable_logging:
            logging.basicConfig(level=logging.INFO)
        self.logger = logging.getLogger(__name__)

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

        self.keybert = KeyBERT(keyword_model) if keyword_model else None





    def preprocess_text(self, text: str) -> str:
        """
        Очищает, лемматизирует и удаляет стоп-слова из текста.
        """
        text = text.lower()
        text = re.sub(r'https?://\S+|www\.\S+', '', text)
        text = re.sub(r'[^а-яё\s]', ' ', text)
        text = re.sub(r'\s+', ' ', text).strip()

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
        sim_matrix = cosine_similarity(embeddings)

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

    def evaluate_cluster_quality(
            self,
            embeddings: np.ndarray,
            labels: np.ndarray,
            visualize: bool = True,
            figsize: Tuple[int, int] = (10, 6)
    ) -> Dict[str, float]:
        """
        Отдельный метод для оценки качества кластеризации.
        Вычисляет Silhouette Score и Davies-Bouldin Index.
        При visualize=True строит силиэтный график.
        """
        metrics: Dict[str, float] = {}
        valid_idx = labels >= 0
        if len(set(labels[valid_idx])) > 1:
            silhouette_avg = silhouette_score(embeddings[valid_idx], labels[valid_idx])
            metrics['silhouette_score'] = silhouette_avg
            if visualize:
                sample_silhouette_values = silhouette_samples(embeddings[valid_idx], labels[valid_idx])
                y_lower = 10
                plt.figure(figsize=figsize)
                for cluster in np.unique(labels[valid_idx]):
                    ith_vals = sample_silhouette_values[labels[valid_idx] == cluster]
                    ith_vals.sort()
                    size = ith_vals.shape[0]
                    y_upper = y_lower + size
                    plt.fill_betweenx(
                        np.arange(y_lower, y_upper),
                        0, ith_vals
                    )
                    plt.text(-0.05, y_lower + 0.5 * size, str(cluster))
                    y_lower = y_upper + 10
                plt.title('Silhouette plot for various clusters')
                plt.xlabel('Silhouette coefficient values')
                plt.ylabel('Cluster label')
                plt.axvline(x=silhouette_avg, color='red', linestyle='--')
                plt.show()
        else:
            metrics['silhouette_score'] = float('nan')
            if visualize:
                self.logger.warning("Not enough clusters for silhouette plot.")
        if len(set(labels)) > 1:
            db_index = davies_bouldin_score(embeddings, labels)
            metrics['davies_bouldin'] = db_index
        else:
            metrics['davies_bouldin'] = float('nan')
        return metrics


