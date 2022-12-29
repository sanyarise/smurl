package models

import "time"

// The internal structure of the Smurl object
type Smurl struct {
	CreatedAt time.Time
	ModifiedAt time.Time
	SmallURL string 
	LongURL  string
	AdminURL string
	IPInfo   []string
	Count    uint64
}
