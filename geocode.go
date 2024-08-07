package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type GeoObject struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type GeoObjectCollection struct {
	Items []struct {
		Point GeoObject `json:"point"`
	} `json:"items"`
}

func geocode(apiKey, address string) (GeoObject, error) {
	baseURL := "https://catalog.api.2gis.com/3.0/items/geocode"
	params := url.Values{}
	params.Set("q", address)
	params.Set("fields", "items.point")
	params.Set("key", apiKey)

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return GeoObject{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GeoObject{}, err
	}

	var result struct {
		Result GeoObjectCollection `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return GeoObject{}, err
	}

	if len(result.Result.Items) == 0 {
		return GeoObject{}, fmt.Errorf("не удалось найти координаты для адреса")
	}

	return result.Result.Items[0].Point, nil
}
