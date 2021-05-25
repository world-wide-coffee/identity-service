package main

import (
	"context"
	"log"
	"time"

	pb "github.com/world-wide-coffee/identity-service/boundaries/grpc/v1/proto"
	"google.golang.org/grpc"
)

const (
	address  = "localhost:50051"
	name     = "world"
	password = "h3ll0w0r1d"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewIdentityServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if _, err := c.SignUp(ctx, &pb.SignUpRequest{Username: name, Password: password}); err != nil {
		log.Fatalf("could not sign up: %v", err)
	}

	r, err := c.SignIn(ctx, &pb.SignInRequest{Username: name, Password: password})
	if err != nil {
		log.Fatalf("could not sign in: %v", err)
	}

	log.Printf("Signed In: %s", r.GetAccessToken())
}
