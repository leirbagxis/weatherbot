package openweather

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenWeatherClient struct {
	apiKey string
}

func New(apiKey string) *OpenWeatherClient {
	return &OpenWeatherClient{
		apiKey: apiKey,
	}
}

func (o *OpenWeatherClient) Coordinates(city string) (Coordinates, error) {
	url := "http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=5&appid=%s"
	resp, err := http.Get(fmt.Sprintf(url, city, o.apiKey))

	if err != nil {
		return Coordinates{}, fmt.Errorf("failed to get coordinates: %w", err)
	}

	if resp.StatusCode != 200 {
		return Coordinates{}, fmt.Errorf("failed to get coordinates, status code: %d", resp.StatusCode)
	}

	var coordinateResponse []CoordinateResponse
	err = json.NewDecoder(resp.Body).Decode(&coordinateResponse)
	if err != nil {
		return Coordinates{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(coordinateResponse) == 0 {
		return Coordinates{}, fmt.Errorf("no coordinates found for city: %s", city)
	}

	return Coordinates{
		Lat: coordinateResponse[0].Lat,
		Lon: coordinateResponse[0].Lon,
	}, nil

}
