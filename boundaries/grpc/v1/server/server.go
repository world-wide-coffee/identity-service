package server

import (
	pb "github.com/world-wide-coffee/identity-service/boundaries/grpc/v1/proto"
	"github.com/world-wide-coffee/identity-service/logic"
)

type Server struct {
	pb.UnimplementedIdentityServiceServer
	logic *logic.Logic
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) SetLogic(l *logic.Logic) *Server {
	s.logic = l
	return s
}
