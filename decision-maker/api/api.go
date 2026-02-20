//go:generate docker run --rm -t -v ../../api:/spec -v .:/go redocly/cli build-docs openapi.yml -o /go/docs_gen.html
//go:generate docker run --rm -t -v ../../api:/spec -v .:/go redocly/cli bundle openapi.yml -o /go/openapi_gen.yml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest -config "api_config.yml" "../../api/openapi.yml"

package api

import (
	"context"
	"decisionMaker/persistence"
	"decisionMaker/service/device"
	"decisionMaker/service/workflow"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"net"
)

//go:embed openapi_gen.yml
var Spec []byte

//go:embed docs_gen.html
var Docs []byte

type Server struct {
	deviceService   *device.Service
	workflowService *workflow.Service
}

func NewServer(ds *device.Service, ws *workflow.Service) Server {
	return Server{
		deviceService:   ds,
		workflowService: ws,
	}
}

func (s Server) Ping(_ context.Context, _ PingRequestObject) (PingResponseObject, error) {
	return Ping200JSONResponse{Status: "ok"}, nil
}

func (s Server) GetDevice(ctx context.Context, request GetDeviceRequestObject) (GetDeviceResponseObject, error) {
	d, err := s.deviceService.GetById(ctx, request.Id)
	if errors.Is(err, persistence.ErrNotFound) {
		return GetDevice404JSONResponse{
			Message: err.Error(),
		}, nil

	}
	if err != nil {
		log.Print(err)
		return nil, err //TODO implement 500 error message with ref code
	}
	return GetDevice200JSONResponse{
		Id: d.Id,
		Ip: d.Ip.String(),
	}, err
}

func (s Server) GetDevices(ctx context.Context, _ GetDevicesRequestObject) (GetDevicesResponseObject, error) {
	devices, err := s.deviceService.GetAll(ctx)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	deviceStructSlice := make([]struct {
		Id string `json:"id"`
		Ip string `json:"ip"`
	}, len(devices))
	for i, d := range devices {
		deviceStructSlice[i] = struct {
			Id string `json:"id"`
			Ip string `json:"ip"`
		}{
			Id: d.Id,
			Ip: d.Ip.String(),
		}
	}
	return GetDevices200JSONResponse{
		Devices: &deviceStructSlice,
	}, nil
}

func (s Server) CreateDevice(ctx context.Context, request CreateDeviceRequestObject) (CreateDeviceResponseObject, error) {
	ip := net.ParseIP(request.Body.Ip)
	if ip == nil {
		return CreateDevice400JSONResponse{
			Message: fmt.Sprintf(`'%s' is not a valid ip`, request.Body.Ip),
		}, nil
	}
	d, err := s.deviceService.Create(ctx, ip)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return CreateDevice201JSONResponse{
		Id: d.Id,
		Ip: d.Ip.String(),
	}, nil
}

func (s Server) DeleteDevice(ctx context.Context, request DeleteDeviceRequestObject) (DeleteDeviceResponseObject, error) {
	err := s.deviceService.Delete(ctx, request.Id)
	if errors.Is(err, persistence.ErrNotFound) {
		return DeleteDevice404JSONResponse{Message: err.Error()}, nil
	}
	return DeleteDevice204Response{}, nil
}

func (s Server) CreateWorkflow(ctx context.Context, request CreateWorkflowRequestObject) (CreateWorkflowResponseObject, error) {
	name := request.Body.Name
	if name == "" {
		return CreateWorkflow400JSONResponse{
			Message: "name field must not be empty",
		}, nil
	}
	w, err := s.workflowService.Create(ctx, name)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return CreateWorkflow201JSONResponse{
		Devices: nil,
		Id:      w.Id,
		Name:    w.Name,
	}, nil
}

func (s Server) DeleteWorkflow(ctx context.Context, request DeleteWorkflowRequestObject) (DeleteWorkflowResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetWorkflow(ctx context.Context, request GetWorkflowRequestObject) (GetWorkflowResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetWorkflows(ctx context.Context, request GetWorkflowsRequestObject) (GetWorkflowsResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) UpdateWorkflow(ctx context.Context, request UpdateWorkflowRequestObject) (UpdateWorkflowResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) AssociateWorkflowDevices(ctx context.Context, request AssociateWorkflowDevicesRequestObject) (AssociateWorkflowDevicesResponseObject, error) {
	//TODO implement me
	panic("implement me")
}
