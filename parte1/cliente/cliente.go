package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const serverURL = "http://localhost:8080"

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("========================================")
		fmt.Println("GESTIÓN DE ARMAMENTO")
		fmt.Println("========================================")
		fmt.Println("1. Ver Inventario")
		fmt.Println("2. Añadir armas")
		fmt.Println("3. Retirar armas")
		fmt.Println("4. Salir")
		fmt.Print("Selecciona una opción: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			viewInv()
		case "2":
			addWeapon(reader)
		case "3":
			removeWeapon(reader)
		case "4":
			return
		default:
			fmt.Println("Opción inválida.")
		}
	}
}

func viewInv() {
	resp, err := http.Get(serverURL + "/weapons")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func addWeapon(reader *bufio.Reader) {
	fmt.Print("Nombre del arma: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Cantidad: ")
	qtyStr, _ := reader.ReadString('\n')
	qtyStr = strings.TrimSpace(qtyStr)
	qty := atoi(qtyStr)

	body := map[string]interface{}{
		"weapon_name": name,
		"stock":       qty,
	}

	sendJSON("POST", "/weapons", body)
}

func removeWeapon(reader *bufio.Reader) {
	fmt.Print("¿Qué arma quieres descontar?: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Cantidad a retirar: ")
	qtyStr, _ := reader.ReadString('\n')
	qtyStr = strings.TrimSpace(qtyStr)
	qty := atoi(qtyStr)

	body := map[string]interface{}{
		"stock": qty,
	}

	sendJSON("PATCH", "/weapons/"+name, body)
}

func sendJSON(method, path string, data map[string]interface{}) {
	jsonData, _ := json.Marshal(data)

	req, _ := http.NewRequest(method, serverURL+path, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func atoi(s string) int {
	var n int
	fmt.Sscan(s, &n)
	return n
}