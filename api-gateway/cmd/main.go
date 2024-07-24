package main

import (
	"api-gateway/api"
	"api-gateway/grpc-client/order"
	"api-gateway/grpc-client/product"
	"api-gateway/grpc-client/user"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	uClient := user.UserApiConn(os.Getenv("user_url"))
	userHandler := api.NewUserClient(uClient)

	http.HandleFunc("POST /user/create", userHandler.CreateUser)
	http.HandleFunc("GET /user/read", userHandler.GetUserID)

	pClient := product.ProductApiConn(os.Getenv("product_url"))
	productHandler := api.NewProductClient(pClient)

	http.HandleFunc("POST /product/create", productHandler.CreateProduct)
	http.HandleFunc("GET /product/readid", productHandler.GetProductID)
	http.HandleFunc("GET /product/readall", productHandler.ListProducts)

	oClient := order.OrderApiConn(os.Getenv("order_url"))
	orderHandler := api.NewOrderClient(oClient)

	http.HandleFunc("POST /order/create", orderHandler.CreateOrder)
	http.HandleFunc("POST /order/createall", orderHandler.CreateOrders)
	http.HandleFunc("GET /order/read", orderHandler.GetOrder)

	log.Println("Server is running on port", os.Getenv("api"))
	http.ListenAndServe(os.Getenv("api"), nil)
}
	