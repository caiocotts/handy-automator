//go:generate redocly build-docs /api/openapi.yml -o /workspace/api/docs_gen.html
//go:generate redocly bundle /api/openapi.yml -o /workspace/api/openapi_gen.yml
//go:generate oapi-codegen -config "/workspace/api/api_config.yml" "/api/openapi.yml"

package api

import (
	"context"
	"decisionMaker/persistence"
	"decisionMaker/service/auth"
	"decisionMaker/service/device"
	"decisionMaker/service/user"
	"decisionMaker/service/workflow"
	_ "embed"
	"errors"
	"fmt"
)

const internalErrorMessage = "internal server error"

//go:embed openapi_gen.yml
var Spec []byte

//go:embed docs_gen.html
var Docs []byte

type Server struct {
	deviceService   *device.Service
	userService     *user.Service
	workflowService *workflow.Service
	authService     *auth.Service
}

func NewServer(ds *device.Service, us *user.Service, ws *workflow.Service, as *auth.Service) Server {
	return Server{
		deviceService:   ds,
		userService:     us,
		workflowService: ws,
		authService:     as,
	}
}

func (s Server) Ping(context.Context, PingRequestObject) (PingResponseObject, error) {
	return Ping200JSONResponse{Status: "ok"}, nil
}

func (s Server) CreateUser(ctx context.Context, request CreateUserRequestObject) (CreateUserResponseObject, error) {
	u, err := s.userService.Register(ctx, request.Body.Username, request.Body.Password)
	if errors.Is(err, user.ErrPasswordTooLong) || errors.Is(err, persistence.ErrUsernameAlreadyTaken) {
		return CreateUser400JSONResponse{
			Message: err.Error(),
		}, nil
	}

	if err != nil {
		refCode := logWithRef(err, "CreateUser")
		return CreateUser500JSONResponse{
			Message: internalErrorMessage,
			Ref:     refCode,
		}, nil
	}

	return CreateUser201JSONResponse{
		Id:       u.Id,
		Username: u.Username,
	}, nil
}

func (s Server) LoginUser(ctx context.Context, request LoginUserRequestObject) (LoginUserResponseObject, error) {
	u, at, err := s.authService.Login(ctx, request.Body.Username, request.Body.Password)
	if errors.Is(err, auth.ErrInvalidCredentials) {
		return LoginUser401Response{}, nil
	}

	if err != nil {
		refCode := logWithRef(err, "LoginUser")
		return LoginUser500JSONResponse{
			Message: internalErrorMessage,
			Ref:     refCode,
		}, nil
	}

	return LoginUser200JSONResponse{
		UserId:       u.Id,
		Username:     u.Username,
		AccessToken:  at,
		RefreshToken: *u.RefreshToken,
	}, nil
}

func (s Server) LoginUserWithFace(context.Context, LoginUserWithFaceRequestObject) (LoginUserWithFaceResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) RefreshAccessToken(ctx context.Context, _ RefreshAccessTokenRequestObject) (RefreshAccessTokenResponseObject, error) {
	at, err := s.authService.Refresh(ctx)
	if err != nil {
		refCode := logWithRef(err, "RefreshAccessToken")
		return RefreshAccessToken500JSONResponse{
			Message: internalErrorMessage,
			Ref:     refCode,
		}, nil
	}

	return RefreshAccessToken200JSONResponse{
		AccessToken: at,
	}, nil
}

func (s Server) DeleteUser(ctx context.Context, request DeleteUserRequestObject) (DeleteUserResponseObject, error) {
	err := s.userService.Delete(ctx, request.Id)
	if errors.Is(err, persistence.ErrNotFound) {
		return DeleteUser404JSONResponse{Message: err.Error()}, nil
	}

	if err != nil {
		refCode := logWithRef(err, "DeleteUser")
		return DeleteUser500JSONResponse{
			Message: internalErrorMessage,
			Ref:     refCode,
		}, nil
	}

	return DeleteUser204Response{}, nil
}

func (s Server) CreateDevice(ctx context.Context, request CreateDeviceRequestObject) (CreateDeviceResponseObject, error) {
	d, err := s.deviceService.Create(ctx, request.Body.Ip)
	if errors.Is(err, device.ErrInvalidIPFormat) {
		return CreateDevice400JSONResponse{
			Message: fmt.Sprintf(`'%s' is not a valid ip`, request.Body.Ip),
		}, nil
	}

	if err != nil {
		refCode := logWithRef(err, "CreateDevice")
		return CreateDevice500JSONResponse{
			Message: internalErrorMessage,
			Ref:     refCode,
		}, nil
	}

	return CreateDevice201JSONResponse{
		Id: d.Id,
		Ip: d.Ip.String(),
	}, nil
}

func (s Server) GetDevice(ctx context.Context, request GetDeviceRequestObject) (GetDeviceResponseObject, error) {
	d, err := s.deviceService.GetById(ctx, request.Id)
	if errors.Is(err, persistence.ErrNotFound) {
		return GetDevice404JSONResponse{
			Message: err.Error(),
		}, nil
	}

	if err != nil {
		refCode := logWithRef(err, "GetDevice")
		return GetDevice500JSONResponse{
			Message: internalErrorMessage,
			Ref:     refCode,
		}, nil
	}

	return GetDevice200JSONResponse{
		Id: d.Id,
		Ip: d.Ip.String(),
	}, nil
}

func (s Server) GetDevices(ctx context.Context, _ GetDevicesRequestObject) (GetDevicesResponseObject, error) {
	devices, err := s.deviceService.GetAll(ctx)
	if err != nil {
		refCode := logWithRef(err, "GetDevices")
		return GetDevices500JSONResponse{
			Message: internalErrorMessage,
			Ref:     refCode,
		}, nil
	}

	deviceStructSlice := make([]struct {
		Id   string  `json:"id"`
		Ip   string  `json:"ip"`
		Name *string `json:"name,omitempty"`
		Type *string `json:"type,omitempty"`
	}, len(devices))

	for i, d := range devices {
		deviceStructSlice[i] = struct {
			Id   string  `json:"id"`
			Ip   string  `json:"ip"`
			Name *string `json:"name,omitempty"`
			Type *string `json:"type,omitempty"`
		}{
			Id: d.Id,
			Ip: d.Ip.String(),
		}
	}

	return GetDevices200JSONResponse{
		Devices: &deviceStructSlice,
	}, nil
}

func (s Server) DeleteDevice(ctx context.Context, request DeleteDeviceRequestObject) (DeleteDeviceResponseObject, error) {
	err := s.deviceService.Delete(ctx, request.Id)
	if errors.Is(err, persistence.ErrNotFound) {
		return DeleteDevice404JSONResponse{
			Message: err.Error(),
		}, nil
	}

	if err != nil {
		refCode := logWithRef(err, "DeleteDevice")
		return DeleteDevice500JSONResponse{
			Message: internalErrorMessage,
			Ref:     refCode,
		}, nil
	}

	return DeleteDevice204Response{}, nil
}

func (s Server) CreateWorkflow(ctx context.Context, request CreateWorkflowRequestObject) (CreateWorkflowResponseObject, error) {
	w, err := s.workflowService.Create(ctx, request.Body.Name)
	if errors.Is(err, workflow.ErrEmptyName) {
		return CreateWorkflow400JSONResponse{
			Message: "name field must not be empty",
		}, nil
	}
	if err != nil {
		refCode := logWithRef(err, "CreateWorkflow")
		return CreateWorkflow500JSONResponse{
			Message: internalErrorMessage,
			Ref:     refCode,
		}, nil
	}

	return CreateWorkflow201JSONResponse{
		Devices: nil,
		Id:      w.Id,
		Name:    w.Name,
		UserId:  w.UserId,
	}, nil
}

func (s Server) GetWorkflow(ctx context.Context, request GetWorkflowRequestObject) (GetWorkflowResponseObject, error) {
	w, err := s.workflowService.GetById(ctx, request.Id)
	if errors.Is(err, persistence.ErrNotFound) {
		return GetWorkflow404JSONResponse{
			Message: err.Error(),
		}, nil
	}

	if err != nil {
		refCode := logWithRef(err, "GetWorkflow")
		return GetWorkflow500JSONResponse{
			Message: internalErrorMessage,
			Ref:     refCode,
		}, nil
	}

	return GetWorkflow200JSONResponse{
		Id:   w.Id,
		Name: w.Name,
	}, nil
}

func (s Server) GetWorkflows(ctx context.Context, _ GetWorkflowsRequestObject) (GetWorkflowsResponseObject, error) {
	workflows, err := s.workflowService.GetAll(ctx)
	if err != nil {
		refCode := logWithRef(err, "GetWorkflows")
		return GetWorkflows500JSONResponse{
			Message: internalErrorMessage,
			Ref:     refCode,
		}, nil
	}

	workflowStructSlice := make([]struct {
		Devices *[]struct {
			Id   string  `json:"id"`
			Ip   string  `json:"ip"`
			Name *string `json:"name,omitempty"`
			Type *string `json:"type,omitempty"`
		} `json:"devices,omitempty"`
		Id     string `json:"id"`
		Name   string `json:"name"`
		UserId string `json:"userId"`
	}, len(workflows))

	for i, w := range workflows {
		workflowStructSlice[i] = struct {
			Devices *[]struct {
				Id   string  `json:"id"`
				Ip   string  `json:"ip"`
				Name *string `json:"name,omitempty"`
				Type *string `json:"type,omitempty"`
			} `json:"devices,omitempty"`
			Id     string `json:"id"`
			Name   string `json:"name"`
			UserId string `json:"userId"`
		}{Devices: nil, Id: w.Id, Name: w.Name, UserId: w.UserId}
	}

	return GetWorkflows200JSONResponse{
		Workflows: &workflowStructSlice,
	}, nil
}

func (s Server) UpdateWorkflow(context.Context, UpdateWorkflowRequestObject) (UpdateWorkflowResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) DeleteWorkflow(context.Context, DeleteWorkflowRequestObject) (DeleteWorkflowResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) AssociateWorkflowDevices(context.Context, AssociateWorkflowDevicesRequestObject) (AssociateWorkflowDevicesResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) TriggerWorkflow(context.Context, TriggerWorkflowRequestObject) (TriggerWorkflowResponseObject, error) {
	//TODO implement me
	panic("implement me")
}
