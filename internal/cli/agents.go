package cli

import (
	"net/http"
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
		f.agentClients[p.Name] = agents_publicv1connect.NewAgentsPublicAPIsV1SrvClient(http.DefaultClient, p.Profile.Host)
	}

	return f.agentClients[p.Name]
}
