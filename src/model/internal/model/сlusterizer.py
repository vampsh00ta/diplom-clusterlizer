from typing import List

from sentence_transformers import SentenceTransformer
from sklearn.cluster import KMeans
from sklearn.feature_extraction.text import TfidfVectorizer
from collections import defaultdict

# # Примерные данные
# sample_papers = [
#     {"title": "Quantum Entanglement in Particle Physics", "abstract": "Study on entangled states and their behavior under quantum field interactions."},
#     {"title": "Astrophysical Observations of Black Holes", "abstract": "Investigating black hole spin and mass from distant galaxies."},
#     {"title": "Photosynthesis in Extreme Environments", "abstract": "How bacteria perform photosynthesis in hydrothermal vents."},
#     {"title": "CRISPR Gene Editing Techniques", "abstract": "Improving precision in genetic modifications using CRISPR-Cas9."},
#     {"title": "String Theory and Multiverse Models", "abstract": "Evaluating the implications of string theory in multiverse predictions."},
#     {"title": "Protein Folding Simulations with AI", "abstract": "Use of machine learning to predict protein structures."},
#     {"title": "Quantum Computing for Cryptography", "abstract": "Exploring quantum algorithms to solve RSA encryption."},
#     {"title": "Deep Space Telescope Imaging", "abstract": "New methods in imaging galaxies with deep space telescopes."},
#     {"title": "RNA Sequencing in Cancer Research", "abstract": "Analysis of gene expression in tumor cells using RNA-seq."},
#     {"title": "Theoretical Models in Dark Matter", "abstract": "Simulations and predictions related to dark matter behavior."},
# ]
from typing import List
from collections import defaultdict
from sklearn.cluster import KMeans
from sklearn.feature_extraction.text import TfidfVectorizer
from sentence_transformers import SentenceTransformer

class Clusturilzer:
    def __init__(self):
        self.model = SentenceTransformer('all-MiniLM-L6-v2')

    def __get_embeddings(self, texts: List[str]):
        return self.model.encode(texts)

    def __generate_topic_name(self, texts: List[str]) -> str:
        vectorizer = TfidfVectorizer(stop_words='english', max_features=5)
        vectorizer.fit(texts)
        return " ".join(vectorizer.get_feature_names_out()).capitalize()

    def do(self, texts: List[str], group_count: int) -> dict:
        embeddings = self.__get_embeddings(texts)

        kmeans = KMeans(n_clusters=group_count, random_state=42)
        labels = kmeans.fit_predict(embeddings)

        clusters = defaultdict(list)
        for i, label in enumerate(labels):
            clusters[label].append(texts[i])

        result = dict()
        for cluster_id, texts_in_cluster in clusters.items():
            topic_name = self.__generate_topic_name(texts_in_cluster)
            result[topic_name] = texts_in_cluster

        return result

# # 1. Получаем семантические векторы с помощью SBERT
# model = SentenceTransformer('all-MiniLM-L6-v2')
# texts = [f"{p['title']}. {p['abstract']}" for p in sample_papers]
# embeddings = model.encode(texts)
#
# # 2. Кластеризация: задаём нужное число кластеров (например, 3)
# kmeans = KMeans(n_clusters=3, random_state=42)
# labels = kmeans.fit_predict(embeddings)
#
# # 3. Группируем работы по кластерам
# clusters = defaultdict(list)
# for i, label in enumerate(labels):
#     clusters[label].append(sample_papers[i])
#
# # 4. Функция для генерации названия темы через TF-IDF
# def generate_topic_name(papers):
#     corpus = [p["title"] + " " + p["abstract"] for p in papers]
#     vectorizer = TfidfVectorizer(stop_words='english', max_features=5)
#     vectorizer.fit(corpus)
#     # Объединяем ключевые слова через пробел
#     return " ".join(vectorizer.get_feature_names_out()).capitalize()
#
# # 5. Генерируем HTML с визуальным разделением групп кружками
# html = """<html>
#   <head>
#     <meta charset="utf-8">
#     <style>
#        body {
#          font-family: Arial, sans-serif;
#        }
#        .group-container {
#          display: inline-block;
#          margin: 20px;
#          width: 250px;
#          height: 250px;
#          border: 2px solid #444;
#          border-radius: 50%;
#          vertical-align: top;
#          padding: 20px;
#          text-align: center;
#          overflow: auto;
#        }
#        .group-title {
#          font-weight: bold;
#          margin-bottom: 10px;
#          font-size: 16px;
#        }
#        .paper-title {
#          font-size: 14px;
#          margin: 5px 0;
#        }
#     </style>
#   </head>
#   <body>
# """
#
# for cluster_id, papers in clusters.items():
#     topic_name = generate_topic_name(papers)
#     html += f'<div class="group-container">\n'
#     html += f'<div class="group-title">Тема: {topic_name}</div>\n'
#     for p in papers:
#         html += f'<div class="paper-title">{p["title"]}</div>\n'
#     html += "</div>\n"
#
# html += """
#   </body>
# </html>
# """
#
# # Сохраняем HTML в файл
# with open("../../grouped_papers.html", "w", encoding="utf-8") as f:
#     f.write(html)
#
# print("HTML-страница 'grouped_papers.html' успешно создана.")
