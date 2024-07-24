package main

import (
	"log"
	"net"
	"os"
	"userservice/database"
	"userservice/intercepter"
	"userservice/logs"
	"userservice/proto"
	"userservice/server/service"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lgs := logs.GetLogger("logs/logger.log")
	db := database.OpenDb("postgres", os.Getenv("postgres_url"))
	defer db.Close()

	lis, err := net.Listen("tcp", "localhost:7070")
	if err != nil {
		lgs.Println("Server listening error -", err)
		log.Fatalf("Server listening error - %v", err)
	}
	defer lis.Close()

	server := service.NewUserServer(db, lgs)
	s := grpc.NewServer(grpc.StreamInterceptor(intercepter.StreamInterceptor), grpc.UnaryInterceptor(intercepter.UnaryInterceptor))
	proto.RegisterUserServiceServer(s, server)
	reflection.Register(s)

	log.Println("Server is listening to port: 7070")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error to Server - %v", err)
	}
}
