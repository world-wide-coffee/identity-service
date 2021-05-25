package server

import (
	"context"
	"errors"

	pb "github.com/world-wide-coffee/identity-service/boundaries/grpc/v1/proto"
	"github.com/world-wide-coffee/identity-service/logic"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInReply, error) {
	token, err := s.logic.SignIn(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		switch {
		case errors.Is(err, logic.ErrorWrongUsername{}),
			errors.Is(err, logic.ErrorWrongPassword{}):
			err = status.New(codes.InvalidArgument, "wrong_username_or_password").Err()
		default:
			err = status.New(codes.Internal, "").Err()
		}

		return &pb.SignInReply{}, err
	}

	return &pb.SignInReply{AccessToken: token}, nil
}
