package behaviourtest

import (
	"fmt"
	"strings"

	"github.com/cucumber/godog"

	pb "github.com/world-wide-coffee/identity-service/boundaries/grpc/v1/proto"
)

func (state *sharedState) aSignedUpUser(arg1 string) godog.Steps {
	parts := strings.Split(arg1, ":")

	return godog.Steps{
		`I create a "sign-up" request`,
		fmt.Sprintf(`set "username" to "%s"`, parts[0]),
		fmt.Sprintf(`set "password" to "%s"`, parts[1]),
		`send the "sign-up" request to the "Identity" service`,
	}
}

func (state *sharedState) getSignUpReplyValue(reply *pb.SignUpReply, fieldKey string) (string, error) {
	switch fieldKey {
	case "id":
		return reply.GetId(), nil
	default:
		return "", fmt.Errorf("unsupported SignUp fieldKey: %q", fieldKey)
	}
}

func (state *sharedState) populateSignUpRequest(req *pb.SignUpRequest, key, value string) error {
	switch key {
	case "username":
		req.Username = value
	case "password":
		req.Password = value
	default:
		return fmt.Errorf("unsupported fieldKey: %q", key)
	}

	return nil
}
