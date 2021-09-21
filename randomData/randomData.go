package randomdata

import (
	"encoding/json"
)

// The `bridge` between input values and a record in db
type RandomData struct {
	ID          uint   `json:"ID"`
	Title       string `json:"title"`
	Text        string `json:"text"`
	DateCreated int64  `json:"date_created"`
	LastUpdated int64  `json:"last_updated"`
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

// Convert struct to json bytes
func (rd *RandomData) ToJson() ([]byte, error) {
	bytes, err := json.Marshal(rd)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// Convert all rdatas to []byte
func ToJsonAll(rdatas []*RandomData, addNewlines bool) ([]byte, error) {
	var rdatasBytes []byte
	for _, rdata := range rdatas {
		jsonBytes, err := rdata.ToJson()
		if err != nil {
			return nil, err
		}
		rdatasBytes = append(rdatasBytes, jsonBytes...)

		if addNewlines {
			rdatasBytes = append(rdatasBytes, 10)
		}
	}

	return rdatasBytes, nil
}
