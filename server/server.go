package main

import (
	proto "ChittyChat/grpc"
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type ChittyChatServer struct {
	proto.UnimplementedChittyChatServer
	users          map[string]proto.ChittyChat_JoinServerServer
	server_lamport int32
}

func (cs *ChittyChatServer) SendMessage(ctx context.Context, in *proto.Message) (*proto.Empty, error) {
	cs.server_lamport++
	in.Lamport = cs.server_lamport
	LogMessage(in)
	cs.Broadcast(in)
	return &proto.Empty{}, nil
}

func (cs *ChittyChatServer) Broadcast(message *proto.Message) {
	for _, stream := range cs.users {
		stream.Send(message)
	}
}

func LogMessage(msg *proto.Message) {
	log.Printf("Lamport: %d, %s: %s", msg.Lamport, msg.Username, msg.Msg)
}

func (cs *ChittyChatServer) JoinServer(client *proto.UserJoin, stream proto.ChittyChat_JoinServerServer) error {
	fmt.Println("User Joined")
	if cs.users[client.Name] != nil {
		log.Fatalf("User already exists")
	}
	cs.users[client.Name] = stream
	cs.server_lamport++

	message := &proto.Message{
		Username: client.Name,
		Msg:      "User Joined ",
		Lamport:  cs.server_lamport,
	}

	LogMessage(message)
	cs.Broadcast(message)

	for {
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (cs *ChittyChatServer) LeaveServer(ctx context.Context, client *proto.UserLeave) (*proto.Empty, error) {
	fmt.Println("User Left")
	if cs.users[client.Name] == nil {
		log.Fatalf("User doesn't exist")
	}

	cs.server_lamport++

	delete(cs.users, client.Name)
	message := &proto.Message{
		Username: client.Name,
		Msg:      "User Left ",
		Lamport:  cs.server_lamport,
	}

	cs.Broadcast(message)

	for {
		time.Sleep(1 * time.Second)
	}
	return &proto.Empty{}, nil

}

func main() {

	cs := &ChittyChatServer{
		users:          make(map[string]proto.ChittyChat_JoinServerServer),
		server_lamport: 0,
	}
	cs.start_server()
}

func (cs *ChittyChatServer) start_server() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", "192.168.1.237:5050")
	if err != nil {
		log.Fatalf("Did not work")
	}

	proto.RegisterChittyChatServer(grpcServer, cs)

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalf("Did not work")
	}
}
