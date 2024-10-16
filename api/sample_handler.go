package api

import (
	"context"

	"connectrpc.com/connect"

	v1 "github.com/maybemaby/sk-connect/gen/proto/api/v1"
)

type SampleHandler struct{}

func (h *SampleHandler) SampleMethod(ctx context.Context, req *connect.Request[v1.SampleRequest]) (*connect.Response[v1.SampleResponse], error) {

	name := req.Msg.Name

	return connect.NewResponse(&v1.SampleResponse{
		Message: "Hello " + name,
	}), nil
}
