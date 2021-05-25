package behaviourtest

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/cucumber/godog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/world-wide-coffee/identity-service/boundaries/grpc/v1/proto"
	grpc_service "github.com/world-wide-coffee/identity-service/boundaries/grpc/v1/server"
	"github.com/world-wide-coffee/identity-service/logic"
	"github.com/world-wide-coffee/identity-service/persistence/memdb"
)

//nolint: unused
type sharedState struct {
	listener        *bufconn.Listener
	requests        map[string]request
	lastRequestType string
}

type request struct {
	req   interface{}
	reply interface{}
	err   error
}

func (state *sharedState) aRunningIdentityService() error {
	const bufSize = 1024 * 1024

	persister := memdb.NewPersister()
	logic := logic.NewLogic().SetPersister(persister)
	identityService := grpc_service.NewServer().SetLogic(logic)

	grpcServer := grpc.NewServer()
	pb.RegisterIdentityServiceServer(grpcServer, identityService)

	state.listener = bufconn.Listen(bufSize)
	go func(lis net.Listener) {
		grpcServer.Serve(lis) //nolint: errcheck
	}(state.listener)

	return nil
}

func (state *sharedState) sendTheRequestToTheService(requestType, serviceName string) error {
	bufDialer := func(context.Context, string) (net.Conn, error) {
		return state.listener.Dial()
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, serviceName, grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		return err
	}

	defer conn.Close()

	client := pb.NewIdentityServiceClient(conn)

	request := state.requests[requestType]

	switch req := request.req.(type) {
	case *pb.SignUpRequest:
		request.reply, request.err = client.SignUp(ctx, req)
	case *pb.SignInRequest:
		request.reply, request.err = client.SignIn(ctx, req)
	case *pb.ChangePasswordRequest:
		request.reply, request.err = client.ChangePassword(ctx, req)
	default:
		return fmt.Errorf("unsupported request: %T", request.req)
	}

	state.requests[requestType] = request

	return nil
}

func (state *sharedState) iCreateARequest(requestType string) error {
	switch requestType {
	case "sign-up":
		state.requests[requestType] = request{req: &pb.SignUpRequest{}}
	case "sign-in":
		state.requests[requestType] = request{req: &pb.SignInRequest{}}
	case "change-password":
		state.requests[requestType] = request{req: &pb.ChangePasswordRequest{}}
	default:
		return fmt.Errorf("unsupported request: %q", requestType)
	}

	state.lastRequestType = requestType

	return nil
}

func (state *sharedState) getValue(value string) (string, error) {
	if !strings.HasPrefix(value, "{{") {
		return value, nil
	}

	key := value[2 : len(value)-2]

	splitKey := strings.SplitN(key, ".", 3)
	requestType := splitKey[0]

	request := state.requests[requestType]

	if strings.HasPrefix(splitKey[2], "error") {
		return state.getErrorValue(request, splitKey[2])
	}

	switch reply := request.reply.(type) {
	case *pb.SignUpReply:
		return state.getSignUpReplyValue(reply, splitKey[2])
	case *pb.SignInReply:
		return state.getSignInReplyValue(reply, splitKey[2])
	case *pb.ChangePasswordReply:
		return state.getChangePasswordReplyValue(reply, splitKey[2])
	default:
		return "", fmt.Errorf("unsupported request: %T, %q", request.req, requestType)
	}
}

func (state *sharedState) getErrorValue(req request, key string) (val string, err error) {
	switch key {
	case "error":
		if req.err != nil {
			val = req.err.Error()
		}
	case "error.code":
		if req.err != nil {
			val = status.Code(req.err).String()
		}
	case "error.message":
		if req.err != nil {
			s, _ := status.FromError(req.err)
			val = s.Message()
		}
	default:
		err = fmt.Errorf("unsupported argument: %q", key)
	}

	return
}

func (state *sharedState) setTo(key, value string) (err error) {
	requestType := state.lastRequestType

	value, err = state.getValue(value)
	if err != nil {
		return
	}

	request := state.requests[requestType]

	switch req := request.req.(type) {
	case *pb.SignUpRequest:
		err = state.populateSignUpRequest(req, key, value)
	case *pb.SignInRequest:
		err = state.populateSignInRequest(req, key, value)
	case *pb.ChangePasswordRequest:
		err = state.populateChangePasswordRequest(req, key, value)
	default:
		return fmt.Errorf("unsupported request: %T, %q", request.req, requestType)
	}

	return
}

func (state *sharedState) theShouldNotBe(key, expectedValue string) error {
	actualValue, err := state.extractReply(key)
	if err != nil {
		return err
	}

	if actualValue == expectedValue {
		return fmt.Errorf("should not be equal: %q == %q", expectedValue, actualValue)
	}

	return nil
}

func (state *sharedState) theShouldBe(key, expectedValue string) error {
	actualValue, err := state.extractReply(key)
	if err != nil {
		return err
	}

	if actualValue != expectedValue {
		return fmt.Errorf("should be equal: %q != %q", expectedValue, actualValue)
	}

	return nil
}

func (state *sharedState) extractReply(key string) (val string, err error) {
	valueKey := fmt.Sprintf("{{%s.reply.%s}}", state.lastRequestType, key)
	return state.getValue(valueKey)
}

//nolint: unused
func InitializeScenario(ctx *godog.ScenarioContext) {
	state := &sharedState{requests: make(map[string]request)}

	ctx.Step(`^a "signed-up" user: "([^"]*)"$`, state.aSignedUpUser)
	ctx.Step(`^a "signed-in" user: "([^"]*)"$`, state.aSignedInUser)
	ctx.Step(`^a running "Identity" service$`, state.aRunningIdentityService)
	ctx.Step(`^send the "([^"]*)" request to the "([^"]*)" service$`, state.sendTheRequestToTheService)
	ctx.Step(`^I "sign-in" with "([^"]*)"$`, state.iSignIn)
	ctx.Step(`^I create a "([^"]*)" request$`, state.iCreateARequest)
	ctx.Step(`^set "([^"]*)" to "([^"]*)"$`, state.setTo)
	ctx.Step(`^the "([^"]*)" should be "([^"]*)"$`, state.theShouldBe)
	ctx.Step(`^the "([^"]*)" should not be "([^"]*)"$`, state.theShouldNotBe)
}
