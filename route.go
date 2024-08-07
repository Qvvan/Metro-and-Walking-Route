package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RouteRequest struct {
	Points []struct {
		Lat  float64 `json:"lat"`
		Lon  float64 `json:"lon"`
		Type string  `json:"type"`
	} `json:"points"`
	Transport string `json:"transport"`
	Params    struct {
		Pedestrian struct {
			UseInstructions bool `json:"use_instructions"`
		} `json:"pedestrian"`
	} `json:"params"`
	Filters []string `json:"filters"`
}

type RouteResponse struct {
	Result []struct {
		UITotalDuration string `json:"ui_total_duration"`
	} `json:"result"`
}

func getWalkingDuration(apiKey string, start, end GeoObject) (string, error) {
	url := "http://routing.api.2gis.com/routing/7.0.0/global?key=" + apiKey

	request := RouteRequest{
		Points: []struct {
			Lat  float64 `json:"lat"`
			Lon  float64 `json:"lon"`
			Type string  `json:"type"`
		}{
			{Lat: start.Lat, Lon: start.Lon, Type: "walking"},
			{Lat: end.Lat, Lon: end.Lon, Type: "walking"},
		},
		Transport: "pedestrian",
		Params: struct {
			Pedestrian struct {
				UseInstructions bool `json:"use_instructions"`
			} `json:"pedestrian"`
		}{
			Pedestrian: struct {
				UseInstructions bool `json:"use_instructions"`
			}{UseInstructions: true},
		},
		Filters: []string{"dirt_road", "ferry", "highway"},
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result RouteResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result.Result) == 0 {
		return "", fmt.Errorf("не удалось получить время пешего маршрута")
	}

	return result.Result[0].UITotalDuration, nil
}
