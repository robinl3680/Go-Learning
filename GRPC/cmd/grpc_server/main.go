package main

import (
	"context"
	"fmt"
	pb "grpc/proto/tag"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedTagServiceServer
}

func (s *server) GetTags(ctx context.Context, in *pb.GetTagRequest) (*pb.GetTagResponse, error) {
	tagID := in.GetTagId()
	// Query the database for the tag with the given ID
	fmt.Println(tagID)
	// Return the tag name in the response
	return &pb.GetTagResponse{
		Tags: []*pb.Tag{{ Id: "1", Name: "Sample"}},
	}, nil
}

func main() {

	// Create the gRPC server and register the TagService server
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	// Start the server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Server listening on %v", lis.Addr())
	log.Printf("Starting gRPC server on port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
