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

var client = &http.Client{
	Timeout: 5 * time.Second,
}

func measureGET() int64 {

	start := time.Now()

	resp, err := client.Get(baseURL + "/weapons")
	if err != nil {
		return 0
	}

	resp.Body.Close()

	return time.Since(start).Microseconds()
}

func measurePOST(name string) int64 {

	body := map[string]interface{}{
		"weapon_name": name,
		"stock":       100,
	}

	data, _ := json.Marshal(body)

	start := time.Now()

	resp, err := client.Post(baseURL+"/weapons", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return 0
	}

	resp.Body.Close()

	return time.Since(start).Microseconds()
}

func measurePATCH(name string) int64 {

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

	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		return 0
	}

	resp.Body.Close()

	return time.Since(start).Microseconds()
}

func main() {

	file, err := os.Create("testing_results.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var totalGET int64
	var totalPOST int64
	var totalPATCH int64

	fmt.Println("Iniciando pruebas...")

	// GET
	for i := 0; i < requests; i++ {

		t := measureGET()

		totalGET += t

		file.WriteString(fmt.Sprintf("GET %d: %d us\n", i+1, t))

		time.Sleep(50 * time.Millisecond)
	}

	// POST + PATCH
	for i := 0; i < requests; i++ {

		name := fmt.Sprintf("test_weapon_%d", i)

		postTime := measurePOST(name)
		totalPOST += postTime

		file.WriteString(fmt.Sprintf("POST %d: %d us\n", i+1, postTime))

		time.Sleep(50 * time.Millisecond)

		patchTime := measurePATCH(name)
		totalPATCH += patchTime

		file.WriteString(fmt.Sprintf("PATCH %d: %d us\n", i+1, patchTime))

		time.Sleep(50 * time.Millisecond)
	}

	avgGET := totalGET / requests
	avgPOST := totalPOST / requests
	avgPATCH := totalPATCH / requests

	file.WriteString("\n=== PROMEDIOS ===\n")

	file.WriteString(fmt.Sprintf("GET promedio: %d us\n", avgGET))
	file.WriteString(fmt.Sprintf("POST promedio: %d us\n", avgPOST))
	file.WriteString(fmt.Sprintf("PATCH promedio: %d us\n", avgPATCH))

	fmt.Println("Testing finalizado. Resultados guardados en testing_results.txt")
}
