package user

import (
	user_proto "api-gateway/protos/user-proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func UserApiConn(user_url string) user_proto.UserServiceClient {

	conn, err := grpc.NewClient(user_url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("User-servicega ulanishda xatolik...")
	}
	user := user_proto.NewUserServiceClient(conn)
	return user
}
