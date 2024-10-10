package main

import (
	proto "ChittyChat/grpc"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type ChittyChatServer struct {
	proto.UnimplementedChittyChatServer
	users          map[string]proto.ChittyChat_JoinServerServer
	server_lamport int32
}

func (cs *ChittyChatServer) SendMessage(ctx context.Context, in *proto.Message) (*proto.Empty, error) {
	return &proto.Empty{}, nil
}

func (cs *ChittyChatServer) Broadcast(message *proto.Message) {
	cs.server_lamport++
	for _, stream := range cs.users {
		stream.Send(message)
	}
}

func (cs *ChittyChatServer) JoinServer(client *proto.UserJoin, stream proto.ChittyChat_JoinServerServer) error {
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

	cs.Broadcast(message)

	return nil
}

func main() {
	cs := &ChittyChatServer{}
	cs.start_server()
}

func (cs *ChittyChatServer) start_server() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("Did not work")
	}

	proto.RegisterChittyChatServer(grpcServer, cs)

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalf("Did not work")
	}
}
