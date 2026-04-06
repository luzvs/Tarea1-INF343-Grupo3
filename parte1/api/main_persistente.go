package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Weapon struct {
	ID         int    `json:"id"`
	WeaponName string `json:"weapon_name"`
	Stock      int    `json:"stock"`
}

var db *sql.DB

func connectDB() {
	connStr := "host=10.10.28.19 port=5432 user=admin password=1234 dbname=inventario sslmode=disable"

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error al conectar:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("No se pudo conectar a la BD:", err)
	}

	log.Println("✅ Conectado a PostgreSQL")
}

func AddWeapon(name string, stock int) error {
	_, err := db.Exec("INSERT INTO weapons (weapon_name, stock) VALUES ($1, $2)", name, stock)
	if err != nil {
		return errors.New("El arma ya existe o error en BD")
	}
	return nil
}

func GetAllWeapons() []*Weapon {
	rows, err := db.Query("SELECT id, weapon_name, stock FROM weapons")
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	var list []*Weapon

	for rows.Next() {
		w := &Weapon{}
		err := rows.Scan(&w.ID, &w.WeaponName, &w.Stock)
		if err != nil {
			log.Println(err)
			continue
		}
		list = append(list, w)
	}

	return list
}

func RevStock(name string, qty int) error {
	var stock int

	err := db.QueryRow("SELECT stock FROM weapons WHERE weapon_name=$1", name).Scan(&stock)
	if err != nil {
		return errors.New("Arma no encontrada")
	}

	if stock < qty {
		return errors.New("No hay suficiente stock")
	}

	_, err = db.Exec("UPDATE weapons SET stock = stock - $1 WHERE weapon_name=$2", qty, name)
	if err != nil {
		return err
	}

	return nil
}

func GetWeaponsh(c *gin.Context) {
	c.JSON(http.StatusOK, GetAllWeapons())
}

func AddWeaponh(c *gin.Context) {
	var body struct {
		WeaponName string `json:"weapon_name"`
		Stock      int    `json:"stock"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json invalido"})
		return
	}

	if err := AddWeapon(body.WeaponName, body.Stock); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Arma añadida"})
}

func DesWeaponh(c *gin.Context) {
	name := c.Param("weapon_name")

	var body struct {
		Stock int `json:"stock"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json invalido"})
		return
	}

	if err := RevStock(name, body.Stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock actualizado"})
}

func main() {
	connectDB()

	r := gin.Default()

	r.GET("/weapons", GetWeaponsh)
	r.POST("/weapons", AddWeaponh)
	r.PATCH("/weapons/:weapon_name", DesWeaponh)

	r.Run("0.0.0.0:8080")
}
