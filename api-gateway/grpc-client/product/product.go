package product

import (
	product_proto "api-gateway/protos/product-proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ProductApiConn(product_url string) product_proto.ProductServiceClient {

	conn, err := grpc.NewClient(product_url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Product-servicega ulanishda xatolik...")
	}
	product := product_proto.NewProductServiceClient(conn)
	return product
}
