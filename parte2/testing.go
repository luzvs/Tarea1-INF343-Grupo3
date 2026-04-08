package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const baseURL = "http://10.10.28.17:8080"

const requests = 100

func measureGET() time.Duration {

	start := time.Now()

	resp, err := http.Get(baseURL + "/weapons")
	if err != nil {
		return 0
	}

	resp.Body.Close()

	return time.Since(start)
}

func measurePOST(name string) time.Duration {

	body := map[string]interface{}{
		"weapon_name": name,
		"stock":       100,
	}

	data, _ := json.Marshal(body)

	start := time.Now()

	resp, err := http.Post(baseURL+"/weapons", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return 0
	}

	resp.Body.Close()

	return time.Since(start)
}

func measurePATCH(name string) time.Duration {

	body := map[string]int{
		"stock": 1,
	}

	data, _ := json.Marshal(body)

	req, _ := http.NewRequest(
		"PATCH",
		baseURL+"/weapons/"+name,
		bytes.NewBuffer(data),
	)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		return 0
	}

	resp.Body.Close()

	return time.Since(start)
}

func main() {

	file, err := os.Create("testing_results.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var totalGET time.Duration
	var totalPOST time.Duration
	var totalPATCH time.Duration

	fmt.Println("Iniciando pruebas...")

	// GET
	for i := 0; i < requests; i++ {

		t := measureGET()

		totalGET += t

		file.WriteString(fmt.Sprintf("GET %d: %v\n", i+1, t))

		time.Sleep(50 * time.Millisecond)
	}

	// POST + PATCH sobre armas de prueba
	for i := 0; i < requests; i++ {

		name := fmt.Sprintf("test_weapon_%d", i)

		postTime := measurePOST(name)
		totalPOST += postTime

		file.WriteString(fmt.Sprintf("POST %d: %v\n", i+1, postTime))

		time.Sleep(50 * time.Millisecond)

		patchTime := measurePATCH(name)
		totalPATCH += patchTime

		file.WriteString(fmt.Sprintf("PATCH %d: %v\n", i+1, patchTime))

		time.Sleep(50 * time.Millisecond)
	}

	avgGET := totalGET / requests
	avgPOST := totalPOST / requests
	avgPATCH := totalPATCH / requests

	file.WriteString("\n=== PROMEDIOS ===\n")

	file.WriteString(fmt.Sprintf("GET promedio: %v\n", avgGET))
	file.WriteString(fmt.Sprintf("POST promedio: %v\n", avgPOST))
	file.WriteString(fmt.Sprintf("PATCH promedio: %v\n", avgPATCH))

	fmt.Println("Testing finalizado. Resultados guardados en testing_results.txt")
}
