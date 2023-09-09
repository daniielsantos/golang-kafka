package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/daniielsantos/dss/internal/infra/web"

	"github.com/daniielsantos/dss/internal/infra/akafka"
	"github.com/daniielsantos/dss/internal/usecase"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/daniielsantos/dss/internal/infra/repository"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:1122@tcp(db:3306)/dss")
	if err != nil {
		fmt.Printf("FAlhouu ", err)
		panic(err)
	}
	defer db.Close()

	repository := repository.NewProductRepository(db)
	createProductUseCase := usecase.NewCreateProductUseCase(repository)
	listProductUseCase := usecase.NewListProductsUseCase(repository)

	productHandlers := web.NewProductHandlers(createProductUseCase, listProductUseCase)

	r := chi.NewRouter()
	r.Post("/products", productHandlers.CreateProductHandler)
	r.Get("/products", productHandlers.ListProductsHandler)
	r.Get("/health", productHandlers.Health)

	go http.ListenAndServe(":8000", r)
	fmt.Printf("Rodando\n")
	msgChan := make(chan *kafka.Message)
	go akafka.Consume([]string{"product"}, "broker:29092", msgChan)
	fmt.Printf("Passou\n")
	for msg := range msgChan {
		dto := usecase.CreateProductInputDto{}
		err := json.Unmarshal(msg.Value, &dto)
		if err != nil {
			continue
		}
		_, err = createProductUseCase.Execute(dto)
	}
}
