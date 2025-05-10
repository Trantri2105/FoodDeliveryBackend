package models

const MapboxBaseURL = "https://api.mapbox.com"

type MapboxGeocodingResponse struct {
	Features []struct {
		Center []float64 `json:"center"`
	} `json:"features"`
}

type MapboxDirectionsResponse struct {
	Routes []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
		Geometry string  `json:"geometry"`
	} `json:"routes"`
}
