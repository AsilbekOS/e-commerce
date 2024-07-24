package internal

import (
	"context"
	"log"
	product "orderservice/proto-prd"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Conn(id int32) *product.ProductResponse {
	conn, err := grpc.NewClient(":8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Connection error", err)
	}
	clinet := product.NewProductServiceClient(conn)

	resp, err := clinet.GetProductID(context.Background(), &product.GetProductRequest{Id: id})
	if err != nil {
		log.Fatal("Bunday product mavjud emas! - ", err)
	}
	return resp
}
