package order

import (
	order_proto "api-gateway/protos/order-proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func OrderApiConn(order_url string) order_proto.OrderServiceClient {

	conn, err := grpc.NewClient(order_url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Product-servicega ulanishda xatolik...")
	}
	order := order_proto.NewOrderServiceClient(conn)
	return order
}
