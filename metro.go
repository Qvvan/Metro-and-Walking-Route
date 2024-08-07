package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type MetroObjectCollection struct {
	FeatureMember []struct {
		GeoObject struct {
			Name  string `json:"name"`
			Point struct {
				Pos string `json:"pos"`
			} `json:"Point"`
		} `json:"GeoObject"`
	} `json:"featureMember"`
}

func findNearestMetro(apiKey string, coords GeoObject) (string, GeoObject, error) {
	baseURL := "https://geocode-maps.yandex.ru/1.x/"
	params := url.Values{}
	params.Set("apikey", apiKey)
	params.Set("geocode", fmt.Sprintf("%f,%f", coords.Lon, coords.Lat))
	params.Set("kind", "metro")
	params.Set("format", "json")
	params.Set("results", "1")

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return "", GeoObject{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", GeoObject{}, err
	}

	var result struct {
		Response struct {
			GeoObjectCollection MetroObjectCollection `json:"GeoObjectCollection"`
		} `json:"response"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", GeoObject{}, err
	}

	if len(result.Response.GeoObjectCollection.FeatureMember) == 0 {
		return "", GeoObject{}, fmt.Errorf("не удалось найти ближайшее метро")
	}

	metro := result.Response.GeoObjectCollection.FeatureMember[0].GeoObject
	metroName := metro.Name

	coordsStr := strings.Split(metro.Point.Pos, " ")
	if len(coordsStr) != 2 {
		return "", GeoObject{}, fmt.Errorf("не удалось распарсить координаты метро")
	}

	lon, err := strconv.ParseFloat(coordsStr[0], 64)
	if err != nil {
		return "", GeoObject{}, fmt.Errorf("ошибка парсинга долготы метро: %v", err)
	}

	lat, err := strconv.ParseFloat(coordsStr[1], 64)
	if err != nil {
		return "", GeoObject{}, fmt.Errorf("ошибка парсинга широты метро: %v", err)
	}

	metroCoords := GeoObject{Lat: lat, Lon: lon}
	return metroName, metroCoords, nil
}
