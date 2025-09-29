package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/ahmad-masud/Kivi/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewKVClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Put
	if _, err := c.Put(ctx, &pb.PutRequest{Key: "hello", Value: "world"}); err != nil {
		log.Fatalf("Put error: %v", err)
	}
	fmt.Println("Put hello")

	// Get
	gresp, err := c.Get(ctx, &pb.GetRequest{Key: "hello"})
	if err != nil {
		log.Fatalf("Get error: %v", err)
	}
	if gresp.Found {
		fmt.Printf("Get hello => %s\n", gresp.Value)
	} else {
		fmt.Println("Get hello => not found")
	}

	// List
	lresp, err := c.List(ctx, &pb.ListRequest{})
	if err != nil {
		log.Fatalf("List error: %v", err)
	}
	fmt.Println("List keys:")
	for _, p := range lresp.Pairs {
		fmt.Printf(" - %s: %s\n", p.Key, p.Value)
	}

	// Delete
	if _, err := c.Delete(ctx, &pb.DeleteRequest{Key: "hello"}); err != nil {
		log.Fatalf("Delete error: %v", err)
	}
	fmt.Println("Deleted hello")
}
