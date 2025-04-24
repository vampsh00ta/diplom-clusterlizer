from internal.model.сlusterizer import Clusterizer  # Имя файла с классом, без `.py`

sample_data = {
    "doc1": "The stock market crashed due to economic instability.",
    "doc2": "Investors are worried about inflation and rising interest rates.",
    "doc3": "New AI models are revolutionizing the tech industry.",
    "doc4": "Machine learning and neural networks are gaining popularity.",
    "doc5": "Climate change is causing more frequent natural disasters.",
    "doc6": "Global warming leads to extreme weather events worldwide.",
}

def main():

    clust = Clusterizer()
    result = clust.do(sample_data, group_count=3)

    for topic, ids in result.items():
        print(f"\nTopic: {topic}")
        for doc_id in ids:
            print(f" - {doc_id}")

if __name__ == "__main__":
    main()