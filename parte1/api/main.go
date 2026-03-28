package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Weapon struct {
	ID         int    `json:"id"`
	WeaponName string `json:"weapon_name"`
	Stock      int    `json:"stock"`
}

var weapons = map[string]*Weapon{}
var nextID = 1

func AddWeapon(name string, stock int) error {
	if _, exists := weapons[name]; exists {
		return errors.New("El arma ya existe")
	}

	weapons[name] = &Weapon{
		ID:         nextID,
		WeaponName: name,
		Stock:      stock,
	}
	nextID++
	return nil
}

func GetAllWeapons() []*Weapon {
	list := []*Weapon{}
	for _, w := range weapons {
		list = append(list, w)
	}
	return list
}

func RevStock(name string, qty int) error {
	w, exists := weapons[name]
	if !exists {
		return errors.New("Arma no encontrada")
	}

	if w.Stock < qty {
		return errors.New("No hay suficiente stock")
	}

	w.Stock -= qty
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
	r := gin.Default()

	r.GET("/weapons", GetWeaponsh)
	r.POST("/weapons", AddWeaponh)
	r.PATCH("/weapons/:weapon_name", DesWeaponh)

	r.Run(":8080")
}