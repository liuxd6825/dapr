package grpc

import (
	"context"
	runtimev1pb "github.com/liuxd6825/dapr/pkg/proto/runtime/v1"
	"strings"
)

type CommandRuntime interface {
	AppHealthChanged(ctx context.Context)
}
type CommandName string

const (
	CommandInit = "init"
	CommandHelp = "help"
)

func (a *api) Command(ctx context.Context, request *runtimev1pb.CommandRequest) (*runtimev1pb.CommandResponse, error) {
	resp := &runtimev1pb.CommandResponse{}
	if request == nil {
		return resp, nil
	}
	cmd := strings.Trim(strings.ToLower(request.Command), " ")
	switch cmd {
	case CommandHelp:
		{
			resp.Data = make(map[string]string)
			resp.Data["init"] = "重新初始化Dapr"
			resp.Data["help"] = "帮助说明"
		}
	case CommandInit:
		{
			a.commandRuntime.AppHealthChanged(ctx)
			return resp, nil
		}
	}
	return resp, nil
}
