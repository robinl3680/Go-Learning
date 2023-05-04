package main

import (
	"context"
	"fmt"
	pb "grpc/proto/tag"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Set up a connection to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	// Create a new client for the TagService
	client := pb.NewTagServiceClient(conn)

	// Call the GetTag method to retrieve the tag with ID 1
	req := &pb.GetTagRequest{TagId: "1"}
	resp, err := client.GetTags(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling GetTag: %v", err)
	}
	// Print the tag name returned by the server
	fmt.Printf("Tag name: %s\n", resp.Tags[0])
}
