from typing import Dict
from sentence_transformers import SentenceTransformer
from typing import List
from collections import defaultdict
from sklearn.cluster import KMeans
from sklearn.feature_extraction.text import TfidfVectorizer
from internal.consumer.entity import Group, ClusterizationRes

class  Clusterizer:
    def __init__(self):
        self.model = SentenceTransformer('all-MiniLM-L6-v2')

    def __get_embeddings(self, texts: List[str]):
        return self.model.encode(texts)

    def __generate_topic_name(self, texts: List[str]) -> str:
        vectorizer = TfidfVectorizer(stop_words='english', max_features=5)
        vectorizer.fit(texts)
        return " ".join(vectorizer.get_feature_names_out()).capitalize()

    def do(self, id_texts: Dict[str, str], group_count: int = 1) -> ClusterizationRes:
        ids = list(id_texts.keys())
        texts = list(id_texts.values())
        embeddings = self.__get_embeddings(texts)
        kmeans = KMeans(n_clusters=group_count, random_state=42)
        labels = kmeans.fit_predict(embeddings)

        clusters = defaultdict(list)
        texts_in_clusters = defaultdict(list)

        for idx, label in enumerate(labels):
            clusters[label].append(ids[idx])
            texts_in_clusters[label].append(texts[idx])

        group_results = []
        for cluster_id in clusters:
            topic_name = self.__generate_topic_name(texts_in_clusters[cluster_id])
            group_results.append(
                Group(keys=topic_name, ids=clusters[cluster_id])
            )

        doc_id = ids[0].split("_")[0] if ids else "unknown"

        return ClusterizationRes(id=doc_id, groups=group_results)


