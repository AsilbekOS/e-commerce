package main

import (
	"log"
	"net"
	"os"
	"productservice/database"
	"productservice/intercepter"
	"productservice/logs"
	"productservice/proto"
	"productservice/server/service"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lgs := logs.GetLogger("logs/logger.log")
	db := database.OpenDb("postgres", os.Getenv("postgres_url"))
	defer db.Close()

	lis, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		log.Fatalf("Server listening error - %v", err)
	}
	defer lis.Close()

	server := service.NewProductServer(db, lgs)
	s := grpc.NewServer(grpc.StreamInterceptor(intercepter.StreamInterceptor), grpc.UnaryInterceptor(intercepter.UnaryInterceptor))
	proto.RegisterProductServiceServer(s, server)
	reflection.Register(s)

	log.Println("Server is listening on port: 8888")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error to server - %v", err)
	}
}
