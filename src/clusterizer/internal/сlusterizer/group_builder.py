import json
from typing import Dict

import networkx as nx
from networkx.readwrite import json_graph
from sentence_transformers import SentenceTransformer
from typing import List
from collections import defaultdict
from sklearn.cluster import KMeans
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import cosine_similarity

from internal.entity.graph import Group
from internal.entity.response import ClusterizationRes
from internal.entity.document import Document as DocumentEntity

from keybert import KeyBERT

class  Groupbuilder:
    def __init__(self):
        self.model = SentenceTransformer('paraphrase-multilingual-MiniLM-L12-v2')  # Поддержка рус/англ
        self.keyword_model = KeyBERT(model=self.model)

    def __get_embeddings(self, texts: List[str]):
        return self.model.encode(texts)

    def __generate_topic_name(self, texts: List[str]) -> str:
        try:
            combined_text = " ".join(texts)
            keywords = self.keyword_model.extract_keywords(
                combined_text,
                keyphrase_ngram_range=(1, 3),
                stop_words=None,  # или 'russian'
                top_n=5,
                use_maxsum=True,
                nr_candidates=20
            )
            if not keywords:
                return "Untitled group"
            # Возвращаем только самую сильную фразу
            return keywords[0][0].capitalize()
        except Exception:
            return "Untitled group"

    def __graph_to_json(self,graph: nx.Graph) -> str:
        data = json_graph.node_link_data(graph)
        return json.dumps(data, indent=2, ensure_ascii=False)
    def __build_semantic_graph(self, id_texts: Dict[str, DocumentEntity], threshold: float = 0.75) -> str:
        try:
            ids = list(id_texts.keys())

            texts =[]
            for document in id_texts.values():
                texts.append(document.text)
            # Эмбеддинги текстов
            embeddings = self.model.encode(texts)

            # Косинусное сходство
            similarity_matrix = cosine_similarity(embeddings)

            # Построение графа
            graph = nx.Graph()

            # Добавляем вершины
            for idx, doc_id in enumerate(ids):
                graph.add_node(doc_id, label=doc_id)

            # Добавляем рёбра при достаточной близости
            for i in range(len(ids)):
                for j in range(i + 1, len(ids)):
                    sim = similarity_matrix[i][j]
                    if sim >= threshold:
                        graph.add_edge(ids[i], ids[j], weight=float(sim))

            return self.__graph_to_json(graph)

        except Exception as e:
            print(f"Ошибка при построении графа: {e}")
            return ""
    def do(self, id_texts: Dict[str, DocumentEntity], group_count: int = 1) -> ClusterizationRes:
        ids = list(id_texts.keys())
        texts = []
        for document in id_texts.values():
            texts.append(document.text)
        embeddings = self.__get_embeddings(texts)

        kmeans = KMeans(n_clusters=group_count, random_state=42)
        labels = kmeans.fit_predict(embeddings)

        clusters = defaultdict(list)
        texts_in_clusters = defaultdict(list)

        for idx, label in enumerate(labels):
            clusters[label].append(ids[idx])
            texts_in_clusters[label].append(texts[idx])
        # print(self.__build_semantic_graph(id_texts))
        group_results = []
        for cluster_id in clusters:
            topic_name = self.__generate_topic_name(texts_in_clusters[cluster_id])
            group_results.append(
                Group(keys=topic_name, ids=clusters[cluster_id])
            )

        doc_id = ids[0].split("_")[0] if ids else "unknown"

        return ClusterizationRes(id=doc_id, groups=group_results)


