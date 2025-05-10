package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"thanhnt208/delivery-service/internal/models"
)

type MapboxClient struct {
	apiKey string
}

func NewMapboxClient(apiKey string) *MapboxClient {
	return &MapboxClient{apiKey: apiKey}
}

func (m *MapboxClient) GeocodeAddress(address string) ([]float64, error) {
	encodedAddress := url.QueryEscape(address)
	geocodingURL := fmt.Sprintf("%s/geocoding/v5/mapbox.places/%s.json?access_token=%s",
		models.MapboxBaseURL, encodedAddress, m.apiKey)

	resp, err := http.Get(geocodingURL)
	if err != nil {
		return nil, fmt.Errorf("geocoding request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read geocoding response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding API returned status %d: %s", resp.StatusCode, string(body))
	}

	var geocodingResp models.MapboxGeocodingResponse
	if err := json.Unmarshal(body, &geocodingResp); err != nil {
		return nil, fmt.Errorf("failed to parse geocoding response: %w", err)
	}

	if len(geocodingResp.Features) == 0 {
		return nil, fmt.Errorf("no results found for address: %s", address)
	}

	return geocodingResp.Features[0].Center, nil
}

func (m *MapboxClient) GetDirections(fromCoords, toCoords []float64) (float64, float64, string, error) {
	directionsURL := fmt.Sprintf("%s/directions/v5/mapbox/driving/%f,%f;%f,%f?access_token=%s&geometries=polyline",
		models.MapboxBaseURL, fromCoords[0], fromCoords[1], toCoords[0], toCoords[1], m.apiKey)

	resp, err := http.Get(directionsURL)
	if err != nil {
		return 0, 0, "", fmt.Errorf("directions request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, "", fmt.Errorf("failed to read directions response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, 0, "", fmt.Errorf("directions API returned status %d: %s", resp.StatusCode, string(body))
	}

	var directionsResp models.MapboxDirectionsResponse
	if err := json.Unmarshal(body, &directionsResp); err != nil {
		return 0, 0, "", fmt.Errorf("failed to parse directions response: %w", err)
	}

	if len(directionsResp.Routes) == 0 {
		return 0, 0, "", fmt.Errorf("no route found between the specified locations")
	}

	route := directionsResp.Routes[0]
	return route.Distance, route.Duration, route.Geometry, nil
}
