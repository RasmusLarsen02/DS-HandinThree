package main

import (
	proto "ChittyChat/grpc"
	"context"
	"fmt"
	"log"
	"os/user"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working")
	}

	currentUser, _ := user.Current()

	client := proto.NewChittyChatClient(conn)
	user := &proto.UserJoin{
		Name:    currentUser.Name,
		Lamport: 0,
	}

	stream, _ := client.JoinServer(context.Background(), user)
	go Listen(stream)

	time.Sleep(1000)
	msg1 := &proto.Message{
		Username: currentUser.Username,
		Msg:      "Hello Every One!",
		Lamport:  0,
	}
	_, _ = client.SendMessage(context.Background(), msg1)

	for {
	}
}

func Listen(stream grpc.ServerStreamingClient[proto.Message]) {
	for {
		msg, _ := stream.Recv()
		if msg == nil {
			continue
		}
		fmt.Printf("%s : %s (at time %d)", msg.Username, msg.Msg, msg.Lamport)
	}
}
