package feed

import (
	"encoding/xml"
	"errors"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// ParseMediumFeed parses a Medium RSS feed from raw XML data and returns a structured
// Feed object containing all posts with their parsed content.
//
// The function expects valid RSS XML data from Medium's feed endpoint. It extracts
// the channel metadata (title, description, link) and processes each post item,
// including parsing the HTML-encoded content into structured elements.
//
// Parameters:
//   - data: Raw XML bytes from the Medium RSS feed
//
// Returns:
//   - *Feed: Structured feed object with metadata and parsed posts
//   - error: Parsing error if XML is malformed or content cannot be processed
func ParseMediumFeed(data []byte) (*Feed, error) {
	var document struct {
		Channel struct {
			Title       string `xml:"title"`
			Description string `xml:"description"`
			Link        string `xml:"link"`
			Items       []struct {
				Title      string   `xml:"title"`
				Link       string   `xml:"link"`
				Creator    string   `xml:"creator"`
				PubDate    string   `xml:"pubDate"`
				Content    string   `xml:"encoded"`
				Categories []string `xml:"category"`
			} `xml:"item"`
		} `xml:"channel"`
	}

	err := xml.Unmarshal(data, &document)
	if err != nil {
		return nil, err
	}

	feed := Feed{
		Title:       document.Channel.Title,
		Description: document.Channel.Description,
		Link:        document.Channel.Link,
		Posts:       []Post{},
	}

	posts := []Post{}
	for _, item := range document.Channel.Items {
		elements, err := parseContent(item.Content)
		if err != nil {
			return nil, err
		}

		published, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			return nil, err
		}

		posts = append(posts, Post{
			Title:      item.Title,
			Link:       item.Link,
			Author:     item.Creator,
			Published:  published.Format(time.RFC3339),
			Content:    elements,
			Categories: item.Categories,
		})
	}

	feed.Posts = posts
	return &feed, nil
}

// parseContent parses HTML-encoded content from a Medium post and converts it into
// a structured representation of elements.
//
// This function takes the raw HTML content from a Medium post's <content:encoded>
// field and parses it into a tree of Element structures. It handles HTML fragments
// by creating a virtual root div element and parsing the content as a fragment.
//
// The function validates that the content is not empty and uses the html package
// to parse the HTML structure, then recursively processes each node to build
// the element tree.
//
// Parameters:
//   - content: Raw HTML string from the post's content field
//
// Returns:
//   - []Element: Slice of parsed elements representing the content structure
//   - error: Parsing error if content is empty or HTML is malformed
func parseContent(content string) ([]Element, error) {
	if strings.TrimSpace(content) == "" {
		return nil, errors.New("document is empty")
	}

	root := &html.Node{Type: html.ElementNode, DataAtom: atom.Div, Data: "div"}
	nodes, err := html.ParseFragment(strings.NewReader(content), root)
	if err != nil {
		return nil, err
	}

	elements := []Element{}
	for _, node := range nodes {
		children := parseElement(node)
		elements = append(elements, children...)
	}

	return elements, nil
}

// parseElement converts an HTML node into an Element structure, handling both
// text nodes and element nodes recursively.
//
// This function is the core of the HTML parsing logic. It processes each HTML node
// and converts it into the corresponding Element structure. For text nodes, it
// creates a simple element with just the text value. For element nodes, it extracts
// the tag name, attributes, and recursively processes all child nodes.
//
// The function normalizes tag names to lowercase and trims whitespace from text
// content to ensure consistent output.
//
// Parameters:
//   - node: HTML node to convert to Element structure
//
// Returns:
//   - []Element: Slice containing the converted element (always single element)
func parseElement(node *html.Node) []Element {
	if node.Type == html.TextNode {
		value := strings.TrimSpace(node.Data)
		return []Element{{Value: value}}
	}

	tag := strings.ToLower(node.Data)
	attributes := parseAttributes(node)
	children := parseChildren(node)

	return []Element{{
		Tag:        tag,
		Attributes: attributes,
		Children:   children,
	}}
}

// parseChildren recursively processes all child nodes of an HTML element and
// converts them into Element structures.
//
// This function traverses the HTML node tree starting from the first child and
// processes each sibling node recursively. It's used to build the complete
// element hierarchy for complex HTML structures with nested elements.
//
// The function uses a simple linked-list traversal pattern to visit each child
// node in order, ensuring that the resulting element structure maintains the
// original HTML hierarchy.
//
// Parameters:
//   - node: HTML node whose children should be processed
//
// Returns:
//   - []Element: Slice of all child elements in document order
func parseChildren(node *html.Node) []Element {
	elements := []Element{}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		elements = append(elements, parseElement(child)...)
	}
	return elements
}

// parseAttributes extracts all HTML attributes from a node and converts them
// into Attribute structures.
//
// This function processes the attribute list of an HTML element and creates
// corresponding Attribute structures. It preserves both the attribute names
// and values exactly as they appear in the original HTML, including any
// namespace prefixes or special characters.
//
// The function handles all standard HTML attributes as well as custom data
// attributes and other non-standard attributes that might be present in
// Medium's content.
//
// Parameters:
//   - node: HTML node whose attributes should be extracted
//
// Returns:
//   - []Attribute: Slice of all attributes found on the node
func parseAttributes(node *html.Node) []Attribute {
	attributes := []Attribute{}
	for _, attribute := range node.Attr {
		attributes = append(attributes, Attribute{
			Name:  attribute.Key,
			Value: attribute.Val,
		})
	}
	return attributes
}
