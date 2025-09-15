package feed

type ElementType string

const (
	ElementText          ElementType = "text"
	ElementParagraph     ElementType = "paragraph"
	ElementHeading       ElementType = "heading"
	ElementLink          ElementType = "link"
	ElementImage         ElementType = "image"
	ElementFigure        ElementType = "figure"
	ElementFigcaption    ElementType = "figcaption"
	ElementStrong        ElementType = "strong"
	ElementEmphasis      ElementType = "emphasis"
	ElementUnorderedList ElementType = "unorderedList"
	ElementOrderedList   ElementType = "orderedList"
	ElementListItem      ElementType = "listItem"
	ElementBreak         ElementType = "break"
)

func isBlockElement(elementType ElementType) bool {
	if elementType == ElementHeading ||
		elementType == ElementParagraph ||
		elementType == ElementFigure ||
		elementType == ElementUnorderedList ||
		elementType == ElementOrderedList {
		return true
	}
	return false
}
