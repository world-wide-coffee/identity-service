package server

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/world-wide-coffee/identity-service/boundaries/grpc/v1/proto"
	"github.com/world-wide-coffee/identity-service/logic"
)

func (s *Server) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpReply, error) {
	id, err := s.logic.SignUp(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		switch {
		case errors.Is(err, logic.ErrorEmptyUsername{}):
			err = status.New(codes.InvalidArgument, "empty_username").Err()
		case errors.Is(err, logic.ErrorTooShortPassword{}):
			err = status.New(codes.InvalidArgument, "too_short_password").Err()
		default:
			err = status.New(codes.Internal, "").Err()
		}

		return &pb.SignUpReply{}, err
	}

	return &pb.SignUpReply{Id: id.String()}, nil
}
