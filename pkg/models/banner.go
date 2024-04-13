package models

import (
	"github.com/microcosm-cc/bluemonday"
	"strings"
	"time"
)

type Content struct {
	Title string `json:"title"   valid:"required"`
	Text  string `json:"text"    valid:"required"`
	URL   string `json:"url"     valid:"required"`
}

type Banner struct {
	BannerID  uint64    `json:"banner_id"    valid:"required"`
	TagIDs    []uint64  `json:"tag_ids"      valid:"required"`
	FeatureID uint64    `json:"feature_id"   valid:"required"`
	Content   Content   `json:"content"      valid:"required"`
	IsActive  bool      `json:"is_active"    valid:"required"`
	CreatedAt time.Time `json:"created_at"   valid:"required"`
	UpdatedAt time.Time `json:"updated_at"   valid:"optional"`
}

type PreBanner struct {
	TagIDs    []uint64 `json:"tag_ids"      valid:"required"`
	FeatureID uint64   `json:"feature_id"   valid:"required"`
	Content   Content  `json:"content"      valid:"required"`
	IsActive  bool     `json:"is_active"    valid:"required"`
}

func (b *PreBanner) Trim() {
	b.Content.Title = strings.TrimSpace(b.Content.Title)
	b.Content.URL = strings.TrimSpace(b.Content.URL)
	b.Content.Text = strings.TrimSpace(b.Content.Text)
}

func (c *Content) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()

	c.Title = sanitizer.Sanitize(c.Title)
	c.Text = sanitizer.Sanitize(c.Text)
	c.URL = sanitizer.Sanitize(c.URL)
}

func (b *Banner) Sanitize() {
	b.Content.Sanitize()
}
