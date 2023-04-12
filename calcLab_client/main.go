package main

import (
	"context"
	"log"
	"time"

	"calcLab2/grpc_api"

	"google.golang.org/grpc"
)

const (
	address     = "localhost:8083"
	defaultName = "world!"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := grpc_api.NewCalculateLabParamClient(conn)

	// Contact the server and print out its response.

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.CleanStage(ctx, &grpc_api.DNAcon{C: 100})

	if err != nil {
		log.Fatalf("could not calculate: %v", err)
	}

	log.Printf("Result: %d", r.GetS())
	//log.Printf("Next value: %s", r.GetId())

}
