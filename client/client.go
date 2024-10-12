package main

import (
	proto "ChittyChat/grpc"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var mu sync.Mutex

func main() {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working")
	}
	var wg sync.WaitGroup

	msgCh := make(chan *proto.Message)
	inputCh := make(chan string)
	outputCh := make(chan string)

	currentUser, _ := user.Current()

	client := proto.NewChittyChatClient(conn)
	user := &proto.UserJoin{
		Name:    currentUser.Name,
		Lamport: 0,
	}

	stream, _ := client.JoinServer(context.Background(), user)
	wg.Add(2)

	go Listen(stream, msgCh, &wg)

	go ReadInput(inputCh, &wg)

	go PrintMessages(outputCh)

	for {
		select {
		case msg := <-msgCh:
			mu.Lock()
			str := fmt.Sprintf("Message: %s: %s (Lamport: %d)", msg.Username, msg.Msg, msg.Lamport)
			outputCh <- str
			mu.Unlock()

		case input, ok := <-inputCh:
			if !ok {
				close(outputCh)
				fmt.Println("Exiting...")
				return
			}
			msg1 := &proto.Message{
				Username: currentUser.Name,
				Msg:      input,
				Lamport:  0,
			}
			_, _ = client.SendMessage(context.Background(), msg1)

		}
	}
	wg.Wait()
}

	for {
	}
}

	for {
		}
	}
}
