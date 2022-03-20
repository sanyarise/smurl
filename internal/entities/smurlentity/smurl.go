package smurlentity

//Внутренняя структура объекта Smurl
type Smurl struct {
	SmallURL string `json:"small_url,omitempty"`
	LongURL  string `json:"long_url,omitempty"`
	AdminURL string `json:"admin_url,omitempty"`
	IPInfo   string `json:"ip_info,omitempty"`
	Count    string `json:"count,omitempty"`
}
