package schema

import "time"

type Bundle struct {
	Version       string         `json:"version"`
	CreatedAt     time.Time      `json:"created_at"`
	SourceDir     string         `json:"source_dir,omitempty"`
	Description   string         `json:"description,omitempty"`
	Files         []FileBundle   `json:"files"`
	Summary       *SummaryInfo   `json:"summary"`
	Generator     string         `json:"generator"`
	TokenSettings map[string]any `json:"token_settings,omitempty"`
}

type FileBundle struct {
	Path         string            `json:"path"`
	RelativePath string            `json:"relative_path"`
	Name         string            `json:"name"`
	Extension    string            `json:"extension"`
	MIMEType     string            `json:"mime_type"`
	FileType     string            `json:"file_type"`
	Language     string            `json:"language,omitempty"`
	Encoding     string            `json:"encoding,omitempty"`
	SizeBytes    int64             `json:"size_bytes"`
	Hash         string            `json:"hash,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	ModifiedAt   time.Time         `json:"modified_at"`
	TokenCount   int               `json:"token_count"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	Document     *DocumentData     `json:"document,omitempty"`
}

type DocumentData struct {
	Title    string        `json:"title,omitempty"`
	Summary  string        `json:"summary,omitempty"`
	Sections []SectionData `json:"sections"`
}

type SectionData struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	Heading    string            `json:"heading,omitempty"`
	Language   string            `json:"language,omitempty"`
	Content    string            `json:"content,omitempty"`
	Items      []string          `json:"items,omitempty"`
	Table      [][]string        `json:"table,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

type SummaryInfo struct {
	TotalFiles     int     `json:"total_files"`
	TotalTokens    int     `json:"total_tokens"`
	TotalSize      int64   `json:"total_size"`
	AvgTokens      float64 `json:"avg_tokens"`
	LargestFile    string  `json:"largest_file"`
	LargestTokens  int     `json:"largest_tokens"`
	EstimatedModel string  `json:"estimated_model,omitempty"`
}
