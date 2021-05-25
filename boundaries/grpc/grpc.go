package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	v1_pb "github.com/world-wide-coffee/identity-service/boundaries/grpc/v1/proto"
	v1_server "github.com/world-wide-coffee/identity-service/boundaries/grpc/v1/server"
	"github.com/world-wide-coffee/identity-service/logic"
)

const (
	port = ":50051"
)

type Server struct {
	v1 *v1_server.Server
}

func NewServer() *Server {
	return &Server{v1: v1_server.NewServer()}
}

func (s *Server) SetLogic(l *logic.Logic) *Server {
	s.v1.SetLogic(l)
	return s
}

func (s *Server) ListenAndServe() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen on tcp-port: %s: %w", port, err)
	}

	grpcServer := grpc.NewServer()
	v1_pb.RegisterIdentityServiceServer(grpcServer, s.v1)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
