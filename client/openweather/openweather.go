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

func (o *OpenWeatherClient) Weather(lat, lon float64) (Weather, error) {
	url := "https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric"
	resp, err := http.Get(fmt.Sprintf(url, lat, lon, o.apiKey))
	if err != nil {
		return Weather{}, fmt.Errorf("failed to get weather: %w", err)
	}

	if resp.StatusCode != 200 {
		return Weather{}, fmt.Errorf("failed to get weather, status code: %d", resp.StatusCode)
	}

	var weatherResponse WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weatherResponse)
	if err != nil {
		return Weather{}, fmt.Errorf("failed to decode weather response: %w", err)
	}

	return Weather{
		Temp: weatherResponse.Main.Temp,
	}, nil
}
