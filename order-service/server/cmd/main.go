package main //port:9999
import (
	"log"
	"net"
	"orderservice/database"
	"orderservice/intercepter"
	"orderservice/logs"
	"orderservice/proto"
	"orderservice/server/service"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lgs := logs.GetLogger("logs/logger.log")
	db := database.OpenDb("postgres", os.Getenv("postgres_url"))

	lis, err := net.Listen("tcp", "localhost:8001")
	if err != nil {
		log.Fatalf("Server listening error - %v", err)
	}
	defer lis.Close()

	server := service.NewOrderServer(db, lgs)
	s := grpc.NewServer(grpc.StreamInterceptor(intercepter.StreamInterceptor), grpc.UnaryInterceptor(intercepter.UnaryInterceptor))
	proto.RegisterOrderServiceServer(s, server)

	reflection.Register(s)

	log.Println("Server is listening on port: 8001")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error to server")
	}
}
