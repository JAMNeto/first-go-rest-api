package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Produto struct {
	Id            int
	Nome          string
	Cor           string
	Preco         float64
	Estoque       int
	Codigo        string
	Poblicacao    bool
	DataDeCriacao string
}

type Produtos struct {
	Produtos []Produto
}

var db []Produto

func GetAll(ctx *gin.Context) {
	file, err := ioutil.ReadFile("theme.json")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Falha ao acessar os dados //readfile",
		})
		return
	}
	var Produtos Produtos
	err2 := json.Unmarshal(file, &Produtos)
	if err2 != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Falha ao acessar os dados//Unmarshal",
		})
		return
	}
	if len(db) == 0 {
		for _, prod := range Produtos.Produtos {
			db = append(db, prod)
		}
	}
	ctx.JSON(http.StatusOK, Produtos)
}

func GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	idConvert, _ := strconv.Atoi(id)
	for _, prod := range db {
		if prod.Id == idConvert {
			ctx.JSON(http.StatusOK, prod)
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{
		"message": "Id nao cadastrado no banco de dados",
	})
}

func main() {
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Seja bem vindo!",
		})
	})

	productsRouter := router.Group("/products")
	{
		productsRouter.GET("/", GetAll)
		productsRouter.GET("/:id", GetById)
	}

	router.Run()

	fmt.Println("API Running at PORT:8080 (default port)!")
}
