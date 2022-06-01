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

type request struct {
	Id            int     `json:"id"`
	Nome          string  `json:"nome"`
	Cor           string  `json:"cor"`
	Preco         float64 `json:"preco"`
	Estoque       int     `json:"estoque"`
	Codigo        string  `json:"codigo"`
	Publicacao    bool    `json:"publicacao"`
	DataDeCriacao string  `json:"dataDeCriacao"`
}

var db []Produto

func GenerateID() int {
	id := len(db) + 1
	return id
}

func GetAll(ctx *gin.Context) {
	if len(db) > 0 {
		ctx.JSON(http.StatusOK, db)
		return
	}
	file, err := ioutil.ReadFile("theme.json")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Falha ao acessar os dados //readfile",
		})
		return
	}
	var Produtos []Produto
	err2 := json.Unmarshal(file, &Produtos)
	if err2 != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Falha ao acessar os dados//Unmarshal",
		})
		return
	}
	if len(db) == 0 {
		for _, prod := range Produtos {
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

func InsertNewProduct(ctx *gin.Context) {
	var req Produto
	token := ctx.GetHeader("token")
	if token != "123456" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "token inválido",
		})
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if req.Nome == "" || req.Cor == "" || req.Codigo == "" || req.Estoque == 0 || req.Preco == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "dados inválidos",
		})
		return
	}
	req.Id = GenerateID()
	db = append(db, req)
	ctx.JSON(http.StatusOK, req)
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
		productsRouter.POST("/insert", InsertNewProduct)
	}

	router.Run()

	fmt.Println("API Running at PORT:8080 (default port)!")
}
