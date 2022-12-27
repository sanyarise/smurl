package models

// The internal structure of the Smurl object
type Smurl struct {
	SmallURL string 
	LongURL  string
	AdminURL string
	IPInfo   []string
	Count    uint64
}
