package main

import (
	proto "ChittyChat/grpc"
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working")
	}

	client := proto.NewChittyChatClient(conn)

	ts, err := client.SendMessage(context.Background(), &proto.Message{Msg: "Hello World!", Lamport: 1})

}
