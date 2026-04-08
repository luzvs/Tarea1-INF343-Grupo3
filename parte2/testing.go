package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const serverURL = "http://10.10.28.17:8080"
const totalTests = 100

func main() {

	file, err := os.Create("testing_results.txt")
	if err != nil {
		fmt.Println("Error creando archivo:", err)
		return
	}
	defer file.Close()

	fmt.Println("Iniciando pruebas...")

	var totalGET time.Duration
	var totalPOST time.Duration
	var totalPATCH time.Duration

	client := &http.Client{}

	for i := 1; i <= totalTests; i++ {

		start := time.Now()

		resp, err := client.Get(serverURL + "/weapons")
		if err != nil {
			fmt.Println("Error GET:", err)
			continue
		}

		io.ReadAll(resp.Body)
		resp.Body.Close()

		duration := time.Since(start)
		totalGET += duration

		ms := duration.Seconds() * 1000

		fmt.Fprintf(file, "GET %d: %.3f ms\n", i, ms)
	}

	for i := 1; i <= totalTests; i++ {

		weaponName := fmt.Sprintf("test_weapon_%d", i)

		// POST
		postBody := map[string]interface{}{
			"weapon_name": weaponName,
			"stock":       10,
		}

		jsonData, _ := json.Marshal(postBody)

		startPost := time.Now()

		reqPost, _ := http.NewRequest("POST", serverURL+"/weapons", bytes.NewBuffer(jsonData))
		reqPost.Header.Set("Content-Type", "application/json")

		respPost, err := client.Do(reqPost)
		if err != nil {
			fmt.Println("Error POST:", err)
			continue
		}

		io.ReadAll(respPost.Body)
		respPost.Body.Close()

		durationPost := time.Since(startPost)
		totalPOST += durationPost

		msPost := durationPost.Seconds() * 1000

		fmt.Fprintf(file, "POST %d: %.3f ms\n", i, msPost)

		// PATCH
		patchBody := map[string]interface{}{
			"stock": 1,
		}

		jsonPatch, _ := json.Marshal(patchBody)

		startPatch := time.Now()

		reqPatch, _ := http.NewRequest("PATCH", serverURL+"/weapons/"+weaponName, bytes.NewBuffer(jsonPatch))
		reqPatch.Header.Set("Content-Type", "application/json")

		respPatch, err := client.Do(reqPatch)
		if err != nil {
			fmt.Println("Error PATCH:", err)
			continue
		}

		io.ReadAll(respPatch.Body)
		respPatch.Body.Close()

		durationPatch := time.Since(startPatch)
		totalPATCH += durationPatch

		msPatch := durationPatch.Seconds() * 1000

		fmt.Fprintf(file, "PATCH %d: %.3f ms\n", i, msPatch)
	}

	avgGET := (totalGET.Seconds() * 1000) / totalTests
	avgPOST := (totalPOST.Seconds() * 1000) / totalTests
	avgPATCH := (totalPATCH.Seconds() * 1000) / totalTests

	fmt.Fprintf(file, "\n=== PROMEDIOS ===\n")
	fmt.Fprintf(file, "GET promedio: %.3f ms\n", avgGET)
	fmt.Fprintf(file, "POST promedio: %.3f ms\n", avgPOST)
	fmt.Fprintf(file, "PATCH promedio: %.3f ms\n", avgPATCH)

	fmt.Println("Testing finalizado. Resultados guardados en testing_results.txt")
}
