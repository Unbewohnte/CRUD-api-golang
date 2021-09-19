package randomdata

import (
	"encoding/json"
	"time"
)

// The `bridge` between input values and a record in db
type RandomData struct {
	DateCreated int64
	LastUpdated int64
	Title       string `json:"title"`
	Text        string `json:"text"`
}

// Create a new `RandomData`
func New(title string, text string) *RandomData {
	return &RandomData{
		DateCreated: time.Now().UTC().Unix(),
		LastUpdated: time.Now().UTC().Unix(),
		Title:       title,
		Text:        text,
	}
}

// Unmarshal `RandomData` from Json encoded bytes
func FromJson(jsonBytes []byte) (*RandomData, error) {
	var randomData RandomData
	err := json.Unmarshal(jsonBytes, &randomData)
	if err != nil {
		return nil, err
	}

	return &randomData, nil
}
