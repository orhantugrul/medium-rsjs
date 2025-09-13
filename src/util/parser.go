package util

import (
	"encoding/xml"
	"strings"
	"time"
)

// Document represents the root structure of an RSS/XML feed document.
// It contains a single Channel which holds all the feed information and items.
type Document struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title      string   `xml:"title"`
	Link       string   `xml:"link"`
	Creator    string   `xml:"creator"`
	PubDate    string   `xml:"pubDate"`
	Content    string   `xml:"encoded"`
	Categories []string `xml:"category"`
}

// ParseDocument parses raw XML/RSS feed data and returns a structured Document.
// It processes the XML data, cleans up text content by removing CDATA tags and
// fixing common encoding issues, and normalizes date formats to RFC3339.
//
// Parameters:
//   - data: Raw byte slice containing the XML/RSS feed data
//
// Returns:
//   - *Document: A parsed and cleaned Document structure, or nil if parsing fails
//   - error: Any error that occurred during parsing or processing
//
// The function performs the following processing steps:
//   - Unmarshals the XML data into the Document structure
//   - Cleans text content by removing CDATA tags and fixing encoding issues
//   - Normalizes publication dates to RFC3339 format
//   - Returns a fully processed Document ready for use
func ParseDocument(data []byte) (*Document, error) {
	var rss Document

	err := xml.Unmarshal(data, &rss)
	if err != nil {
		return nil, err
	}

	items := make([]Item, 0, len(rss.Channel.Items))
	for _, item := range rss.Channel.Items {
		items = append(items, Item{
			Title:      extractValue(item.Title),
			Link:       item.Link,
			Creator:    extractValue(item.Creator),
			PubDate:    extractDate(item.PubDate),
			Content:    extractValue(item.Content),
			Categories: item.Categories,
		})
	}

	return &Document{
		Channel: Channel{
			Title:       extractValue(rss.Channel.Title),
			Description: extractValue(rss.Channel.Description),
			Link:        rss.Channel.Link,
			Items:       items,
		},
	}, nil
}

// extractValue cleans and normalizes text content from RSS feed data.
// It removes CDATA tags, fixes common encoding issues, and trims whitespace.
//
// Parameters:
//   - text: Raw text string that may contain CDATA tags and encoding issues
//
// Returns:
//   - string: Cleaned and normalized text content
//
// The function handles:
//   - CDATA tag removal (<![CDATA[...]]>)
//   - Common encoding issues (smart quotes, non-breaking spaces)
//   - Whitespace trimming
func extractValue(text string) string {
	text = strings.TrimPrefix(text, "<![CDATA[")
	text = strings.TrimSuffix(text, "]]>")

	text = strings.ReplaceAll(text, "â€™", "'")
	text = strings.ReplaceAll(text, "â€œ", "\"")
	text = strings.ReplaceAll(text, "â€", "\"")
	text = strings.ReplaceAll(text, "Â", "")

	return strings.TrimSpace(text)
}

// extractDate parses and normalizes date strings from RSS feeds to RFC3339 format.
// It attempts to parse the input date using multiple common RSS date formats
// and returns a standardized RFC3339 formatted date string.
//
// Parameters:
//   - date: Raw date string from RSS feed (various formats supported)
//
// Returns:
//   - string: Date in RFC3339 format, or current time if parsing fails
//
// Supported date formats:
//   - RFC1123Z (Mon, 02 Jan 2006 15:04:05 -0700)
//   - RFC1123 (Mon, 02 Jan 2006 15:04:05 MST)
//   - RFC3339 (2006-01-02T15:04:05Z07:00)
//   - Custom format (Mon, 2 Jan 2006 15:04:05 GMT)
func extractDate(date string) string {
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC3339,
		"Mon, 2 Jan 2006 15:04:05 GMT",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, date); err == nil {
			return t.Format(time.RFC3339)
		}
	}

	return time.Now().Format(time.RFC3339) // fallback
}
