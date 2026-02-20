//go:generate docker run --rm -t -v ../../api:/spec -v .:/go redocly/cli build-docs openapi.yml -o /go/docs_gen.html
//go:generate docker run --rm -t -v ../../api:/spec -v .:/go redocly/cli bundle openapi.yml -o /go/openapi_gen.yml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest -config "api_config.yml" "../../api/openapi.yml"

package api

import (
	"context"
	"decisionMaker/persistence"
	"decisionMaker/service/device"
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
	deviceService *device.Service
}

func NewServer(ds *device.Service) Server {
	return Server{
		deviceService: ds,
	}
}

func (s Server) GetDevice(ctx context.Context, request GetDeviceRequestObject) (GetDeviceResponseObject, error) {
	d, err := s.deviceService.GetDeviceById(ctx, request.Id)
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

func (s Server) PostDevice(ctx context.Context, request PostDeviceRequestObject) (PostDeviceResponseObject, error) {
	ip := net.ParseIP(request.Body.Ip)
	if ip == nil {
		return PostDevice400JSONResponse{
			Message: fmt.Sprintf(`'%s' is not a valid ip`, request.Body.Ip),
		}, nil
	}
	d, err := s.deviceService.RegisterDevice(ctx, ip)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return PostDevice201JSONResponse{
		Id: d.Id,
		Ip: d.Ip.String(),
	}, nil
}

func (s Server) GetPing(_ context.Context, _ GetPingRequestObject) (GetPingResponseObject, error) {
	return GetPing200JSONResponse{
		Status: "ok",
	}, nil
}
