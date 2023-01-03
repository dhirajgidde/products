package main

import (
	"context"
	"log"
	"time"
	"google.golang.org/grpc"
	messager "grpc_ex.com/v1/productsProto"
)




func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := messager.NewMessageReceiverClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.EnableProducts(ctx, &messager.Product{Title : "Prodct1", SKU : "SKU1", AccountCode : "RRR"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf(r.GetResponseMessage())
}