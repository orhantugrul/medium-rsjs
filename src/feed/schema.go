package feed

type Feed struct {
	Title string `json:"title"`
	Link  string `json:"link"`
	Posts []Post `json:"posts"`
}

type Post struct {
	Title      string   `json:"title"`
	Link       string   `json:"link"`
	Author     string   `json:"author"`
	Published  string   `json:"published"`
	Content    string   `json:"content"`
	Categories []string `json:"categories"`
}
