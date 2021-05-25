package behaviourtest

import (
	"fmt"
	"strings"

	"github.com/cucumber/godog"

	pb "github.com/world-wide-coffee/identity-service/boundaries/grpc/v1/proto"
)

func (state *sharedState) iSignIn(arg1 string) godog.Steps {
	parts := strings.Split(arg1, ":")

	return godog.Steps{
		`I create a "sign-in" request`,
		fmt.Sprintf(`set "username" to "%s"`, parts[0]),
		fmt.Sprintf(`set "password" to "%s"`, parts[1]),
		`send the "sign-in" request to the "Identity" service`,
	}
}

func (state *sharedState) aSignedInUser(arg1 string) godog.Steps {
	return godog.Steps{
		fmt.Sprintf(`a "signed-up" user: "%s"`, arg1),
		fmt.Sprintf(`I "sign-in" with "%s"`, arg1),
		`the "error" should be ""`,
		`the "access_token" should not be ""`,
	}
}

func (state *sharedState) getSignInReplyValue(reply *pb.SignInReply, fieldKey string) (string, error) {
	switch fieldKey {
	case "access_token":
		return reply.GetAccessToken(), nil
	default:
		return "", fmt.Errorf("unsupported SignIn fieldKey: %q", fieldKey)
	}
}

func (state *sharedState) populateSignInRequest(req *pb.SignInRequest, key, value string) error {
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
