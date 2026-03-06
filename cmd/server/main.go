package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()

	log.Println("gRPC server running on 50051")

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}