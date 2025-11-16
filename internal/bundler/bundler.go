package bundler

import (
	"os"
	"path/filepath"
	"time"

	"github.com/sottey/aingest/internal/schema"
)

// Bundler handles directory walking and bundle creation.
type Bundler struct {
	RootPath  string
	Recursive bool
}

// NewBundler creates a new Bundler instance.
func NewBundler(root string, recursive bool) *Bundler {
	return &Bundler{RootPath: root, Recursive: recursive}
}

// BuildBundle scans the root directory and builds a metadata-only bundle.
func (b *Bundler) BuildBundle() (*schema.Bundle, error) {
	var files []schema.FileBundle

	err := filepath.Walk(b.RootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories (unless root)
		if info.IsDir() {
			if !b.Recursive && path != b.RootPath {
				return filepath.SkipDir
			}
			return nil
		}

		// Gather file metadata
		rel, _ := filepath.Rel(b.RootPath, path)
		file := schema.FileBundle{
			Path:         path,
			RelativePath: rel,
			Name:         info.Name(),
			Extension:    filepath.Ext(path),
			MIMEType:     detectMIME(path),
			FileType:     classifyType(filepath.Ext(path)),
			SizeBytes:    info.Size(),
			CreatedAt:    info.ModTime(),
			ModifiedAt:   info.ModTime(),
			TokenCount:   0, // Placeholder until token estimation is added
		}

		files = append(files, file)
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Build bundle summary
	totalSize := int64(0)
	for _, f := range files {
		totalSize += f.SizeBytes
	}

	summary := &schema.SummaryInfo{
		TotalFiles: len(files),
		TotalSize:  totalSize,
	}

	bundle := &schema.Bundle{
		Version:     "1.0",
		CreatedAt:   time.Now(),
		SourceDir:   b.RootPath,
		Files:       files,
		Summary:     summary,
		Generator:   "aingest v0.1.0",
		Description: "Generated bundle (metadata only)",
	}

	return bundle, nil
}

// detectMIME returns a simple MIME guess based on extension.
func detectMIME(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".md":
		return "text/markdown"
	case ".txt":
		return "text/plain"
	case ".json":
		return "application/json"
	case ".go":
		return "text/x-go"
	case ".rtf":
		return "application/rtf"
	default:
		return "application/octet-stream"
	}
}

// classifyType returns a broad category for a file.
func classifyType(ext string) string {
	switch ext {
	case ".go", ".js", ".ts", ".py", ".cs", ".java":
		return "code"
	case ".json", ".yaml", ".yml", ".toml", ".plist", ".csv":
		return "data"
	case ".md", ".txt", ".rtf":
		return "text"
	default:
		return "binary"
	}
}
