package main

import (
	proto "ChittyChat/grpc"
	"context"
	"fmt"
	"log"
	"os/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working")
	}

	currentUser, err := user.Current()

	client := proto.NewChittyChatClient(conn)
	user := &proto.UserJoin{
		Name:    currentUser.Name,
		Lamport: 0,
	}

	stream, err := client.JoinServer(context.Background(), user)

	for {
		msg, _ := stream.Recv()
		if msg == nil {
			continue
		}
		fmt.Printf("%s : %s (at time %d)", msg.Username, msg.Msg, msg.Lamport)
	}

}
