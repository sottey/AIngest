# AIngest — AI-Ready Document Ingestion Framework (Go)

## Overview
**AIngest** is an open-source Go utility that ingests one or more files of varying formats  
(e.g., `.md`, `.txt`, `.rtf`, `.plist`, `.json`, `.go`, `.cs`, `.java`, etc.) and converts them into a single, structured, AI-friendly bundle — typically in JSON.

The output schema is **richly descriptive**, capturing file structure, metadata, and semantic meaning, designed for **large language model (LLM) ingestion** and **AI pipeline pre-processing**.

---

## Core Goals

- 🧩 **Multi-format ingestion** — parse and normalize diverse document and code formats.  
- 📦 **Unified schema output** — produce clean, structured JSON (or XML/SHON) with consistent semantics.  
- 🧠 **AI interpretability** — design data to be easily understood by LLMs (sections, types, metadata).  
- 📏 **Token awareness** — estimate token usage per file and total, with model-window guidance.  
- 🪶 **Lightweight & modular** — pure Go, CLI + library, plug-in friendly architecture.  
- 🔍 **Provenance & traceability** — retain metadata like relative path, hash, encoding, and file size.

---

## Architecture

```plaintext
cmd/
 └── aingest/          → CLI entrypoint (Cobra)
internal/
 ├── ingestors/        → Format-specific file parsers
 │     ├── md.go
 │     ├── txt.go
 │     ├── rtf.go
 │     ├── plist.go
 │     └── code.go
 ├── outputs/          → Output formatters (json.go, xml.go, shon.go)
 ├── schema/           → Shared Go structs for normalized data
 ├── bundler/          → Core orchestration layer
 ├── token/            → Token estimation logic
 └── utils/            → File I/O, MIME detection, hashing, etc.
```

---

## Data Flow

```plaintext
[ Input Files ]
     │
     ▼
  Detect Format
     │
     ▼
  Select Ingestor (via interface)
     │
     ▼
  Normalize content → DocumentData
     │
     ▼
  Aggregate into Bundle
     │
     ▼
  Serialize via Formatter → JSON / XML / SHON
```

---

## Interfaces

### Ingestor Interface
Each ingestor parses a specific file type.

```go
type Ingestor interface {
    CanIngest(path string) bool
    Ingest(path string) (*schema.DocumentData, error)
}
```

### Formatter Interface
Responsible for output serialization.

```go
type Formatter interface {
    Format(bundle *schema.Bundle) ([]byte, error)
    Extension() string
}
```

---

## Schema Design

### Bundle
```go
type Bundle struct {
    Version       string          `json:"version"`
    CreatedAt     time.Time       `json:"created_at"`
    SourceDir     string          `json:"source_dir,omitempty"`
    Description   string          `json:"description,omitempty"`
    Files         []FileBundle    `json:"files"`
    Summary       *SummaryInfo    `json:"summary"`
    Generator     string          `json:"generator"`
    TokenSettings map[string]any  `json:"token_settings,omitempty"`
}
```

### FileBundle
```go
type FileBundle struct {
    Path         string            `json:"path"`           // absolute
    RelativePath string            `json:"relative_path"`  // relative to bundle root
    Name         string            `json:"name"`
    Extension    string            `json:"extension"`
    MIMEType     string            `json:"mime_type"`
    FileType     string            `json:"file_type"`      // text, code, binary, data
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
```

### DocumentData
```go
type DocumentData struct {
    Title    string         `json:"title,omitempty"`
    Summary  string         `json:"summary,omitempty"`
    Sections []SectionData  `json:"sections"`
}
```

### SectionData
```go
type SectionData struct {
    ID         string            `json:"id"`
    Type       string            `json:"type"` // paragraph, code, list, table, etc.
    Heading    string            `json:"heading,omitempty"`
    Language   string            `json:"language,omitempty"` // for code
    Content    string            `json:"content,omitempty"`
    Items      []string          `json:"items,omitempty"`
    Table      [][]string        `json:"table,omitempty"`
    Attributes map[string]string `json:"attributes,omitempty"`
}
```

### SummaryInfo
```go
type SummaryInfo struct {
    TotalFiles     int     `json:"total_files"`
    TotalTokens    int     `json:"total_tokens"`
    TotalSize      int64   `json:"total_size"`
    AvgTokens      float64 `json:"avg_tokens"`
    LargestFile    string  `json:"largest_file"`
    LargestTokens  int     `json:"largest_tokens"`
    EstimatedModel string  `json:"estimated_model,omitempty"`
}
```

---

## Token Awareness

- `TokenEstimator` interface supports both fast and precise modes:
  - **Fast:** approximate word count × 0.75  
  - **Precise:** integrate with `tiktoken-go` or OpenAI’s tokenizer  
- CLI options:
  ```bash
  aingest run ./docs --token-mode fast
  aingest run ./docs --token-mode precise
  ```
- Reports per-file and total token usage with model-context warnings.

---

## CLI Examples

```bash
# Convert single file
aingest run myfile.rtf --out bundle.json

# Bundle multiple files recursively
aingest run ./docs --recursive --out report.json

# Output in XML instead of JSON
aingest run ./src --format xml

# Display token summary
aingest info bundle.json
```

---

## Example Output

```json
{
  "version": "1.0",
  "created_at": "2025-10-24T17:12:00Z",
  "generator": "aingest v0.1.0",
  "description": "Bundle for summarizing project documentation",
  "files": [
    {
      "relative_path": "docs/overview.md",
      "name": "overview.md",
      "extension": ".md",
      "mime_type": "text/markdown",
      "file_type": "text",
      "size_bytes": 4832,
      "token_count": 1231,
      "document": {
        "title": "Overview",
        "sections": [
          {
            "id": "intro",
            "type": "paragraph",
            "content": "This project defines the architecture for..."
          }
        ]
      }
    }
  ],
  "summary": {
    "total_files": 12,
    "total_tokens": 21232,
    "avg_tokens": 1769.3,
    "estimated_model": "gpt-4o-32k"
  }
}
```

---

## Future Enhancements
- **Embeddings-aware chunking**  
- **Server mode** (`aingest serve`) for HTTP ingestion  
- **Plugin system** for community parsers  
- **Model window awareness** — automatic chunking to fit target model  
- **Statistics mode** (`aingest stats`) for size and token reports  
- **Schema registry** for domain-specific extensions (e.g., “notes”, “logs”)

---

## Why It Matters
AI systems perform best with **structured, context-rich input** — not unorganized text blobs.  
By combining:
- **rich metadata**,  
- **token awareness**, and  
- **consistent semantic structure**,  

`AIngest` bridges the gap between raw data and intelligent AI ingestion.

---

## Repository
[github.com/sottey/aingest](https://github.com/sottey/aingest)

---

## License
MIT — free for personal, research, and commercial use.
