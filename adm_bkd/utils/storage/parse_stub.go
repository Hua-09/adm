package storage

// ParseResult holds the structured content extracted from a document.
type ParseResult struct {
	DocID    string   `json:"doc_id"`
	Sections []string `json:"sections"`
	RawText  string   `json:"raw_text"`
}

// parseDoc is a stub that returns an empty ParseResult.
// Replace with a real parser (e.g., PDF / DOCX extraction) as needed.
func parseDoc(rootDir, id string) (ParseResult, error) {
	return ParseResult{DocID: id}, nil
}
