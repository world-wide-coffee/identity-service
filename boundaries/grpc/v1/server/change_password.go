package server

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/world-wide-coffee/identity-service/boundaries/grpc/v1/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordReply, error) {
	reply := pb.ChangePasswordReply{}

	userID, err := uuid.Parse(req.Id)
	if err != nil {
		err = status.New(codes.InvalidArgument, "invalid_id").Err()
		return &reply, err
	}

	err = s.logic.ChangePassword(ctx, userID, req.OldPassword, req.NewPassword)
	if err != nil {
		err = status.New(codes.Internal, err.Error()).Err()
		return &reply, err
	}

	return &reply, nil
}
