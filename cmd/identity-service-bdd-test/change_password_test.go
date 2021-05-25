package behaviourtest

import (
	"fmt"

	pb "github.com/world-wide-coffee/identity-service/boundaries/grpc/v1/proto"
)

func (state *sharedState) getChangePasswordReplyValue(reply *pb.ChangePasswordReply, fieldKey string) (string, error) {
	switch fieldKey {
	default:
		return "", fmt.Errorf("unsupported ChangePassword fieldKey: %q", fieldKey)
	}
}

func (state *sharedState) populateChangePasswordRequest(req *pb.ChangePasswordRequest, key, value string) error {
	switch key {
	case "id":
		req.Id = value
	case "access_token":
		req.AccessToken = value
	case "new_password":
		req.NewPassword = value
	case "old_password":
		req.OldPassword = value
	default:
		return fmt.Errorf("unsupported fieldKey: %q", key)
	}

	return nil
}
