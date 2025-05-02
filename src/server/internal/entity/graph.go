package entity

type Node struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Cluster int    `json:"cluster"`
	Type    string `json:"type"`
}

type Link struct {
	Source string  `json:"source"`
	Target string  `json:"target"`
	Weight float64 `json:"weight"`
}

type GraphData struct {
	Directed   bool                   `json:"directed"`
	Multigraph bool                   `json:"multigraph"`
	Graph      map[string]interface{} `json:"graph"` // Пустой объект или дополнительные поля
	Nodes      []Node                 `json:"nodes"`
	Links      []Link                 `json:"links"`
}

type DocumentSaverReq struct {
	ID    string    `json:"id"`
	Graph GraphData `json:"graph"`
}
