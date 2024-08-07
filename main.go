package main

import (
	"fmt"
	"log"
)

func main() {
	apiKey2GIS := "d5a782a5-2577-48da-9efc-88ee73e965a2"
	apiKeyYandex := "b49495eb-e11e-43a2-b369-a010a8b43f65"
	address := "Москва, Садовническая 36"

	coords, err := geocode(apiKey2GIS, address)
	if err != nil {
		log.Fatalf("Ошибка геокодирования: %v", err)
	}

	metroName, metroCoords, err := findNearestMetro(apiKeyYandex, coords)
	if err != nil {
		log.Fatalf("Ошибка поиска метро: %v", err)
	}

	fmt.Printf("Ближайшее метро к адресу '%s': %s\n", address, metroName)

	duration, err := getWalkingDuration(apiKey2GIS, metroCoords, coords)
	if err != nil {
		log.Fatalf("Ошибка расчета маршрута: %v", err)
	}

	fmt.Printf("Время пешего перехода до адреса: %s\n", duration)
}
