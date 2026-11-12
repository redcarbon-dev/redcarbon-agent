package cli

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	"pkg.redcarbon.ai/internal/build"
	"pkg.redcarbon.ai/internal/config"
	"pkg.redcarbon.ai/proto/redcarbon/agents_public/v1/agents_publicv1connect"
)

type ClientFactory struct {
	agentClients map[string]agents_publicv1connect.AgentsPublicAPIsV1SrvClient
}

func NewClientFactory() ClientFactory {
	return ClientFactory{
		agentClients: make(map[string]agents_publicv1connect.AgentsPublicAPIsV1SrvClient),
	}
}

func (f *ClientFactory) GetAgentClientForProfile(p config.Profile) agents_publicv1connect.AgentsPublicAPIsV1SrvClient {
	if _, ok := f.agentClients[p.Name]; !ok {
		f.agentClients[p.Name] = agents_publicv1connect.NewAgentsPublicAPIsV1SrvClient(
			http.DefaultClient,
			p.Profile.Host,
			connect.WithInterceptors(userAgentInterceptor),
		)
	}

	return f.agentClients[p.Name]
}

var userAgentInterceptor connect.UnaryInterceptorFunc = func(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		req.Header().Set("User-Agent", build.UserAgent())
		return next(ctx, req)
	}
}
