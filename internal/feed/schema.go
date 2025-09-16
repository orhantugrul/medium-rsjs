package feed

type Feed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Posts       []Post `json:"posts"`
}

type Post struct {
	Title      string    `json:"title"`
	Link       string    `json:"link"`
	Author     string    `json:"author"`
	Published  string    `json:"published"`
	Content    []Element `json:"content,omitempty"`
	Categories []string  `json:"categories"`
}

type Element struct {
	Tag        string      `json:"tag,omitempty"`
	Attributes []Attribute `json:"attributes,omitempty"`
	Value      string      `json:"value,omitempty"`
	Children   []Element   `json:"children,omitempty"`
}

type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
