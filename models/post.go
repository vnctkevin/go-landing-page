package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// ContentType defines the type of content block
type ContentType string

const (
	ContentTypeHeading1  ContentType = "h1"
	ContentTypeHeading2  ContentType = "h2"
	ContentTypeHeading3  ContentType = "h3"
	ContentTypeParagraph ContentType = "p"
	ContentTypeImage     ContentType = "img"
	ContentTypeCaption   ContentType = "caption"
	ContentTypeCode      ContentType = "code"
)

// ContentBlock represents a single block of content
type ContentBlock struct {
	Type    ContentType `json:"type"`
	Text    string      `json:"text,omitempty"`    // For headings, paragraphs, captions
	Src     string      `json:"src,omitempty"`     // For images
	Alt     string      `json:"alt,omitempty"`     // For images
	Caption string      `json:"caption,omitempty"` // Optional caption for images if nested, or standalone
	Code    string      `json:"code,omitempty"`    // For code blocks
}

// PostContent is a slice of ContentBlocks that implements sql.Scanner and driver.Valuer
type PostContent []ContentBlock

// Value simply returns the JSON-encoded representation of the struct.
func (pc PostContent) Value() (driver.Value, error) {
	return json.Marshal(pc)
}

func (pc *PostContent) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &pc)
}

type BlogPost struct {
	gorm.Model
	Title   string      `json:"title"`
	Content PostContent `json:"content" gorm:"type:jsonb"`
	Author  string      `json:"author"`
	Slug    string      `json:"slug"`
}
