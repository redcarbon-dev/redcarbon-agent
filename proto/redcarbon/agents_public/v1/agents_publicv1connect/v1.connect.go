// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: redcarbon/agents_public/v1/v1.proto

package agents_publicv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	http "net/http"
	v1 "pkg.redcarbon.ai/proto/redcarbon/agents_public/v1"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// AgentsPublicAPIsV1SrvName is the fully-qualified name of the AgentsPublicAPIsV1Srv service.
	AgentsPublicAPIsV1SrvName = "redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// AgentsPublicAPIsV1SrvHZProcedure is the fully-qualified name of the AgentsPublicAPIsV1Srv's HZ
	// RPC.
	AgentsPublicAPIsV1SrvHZProcedure = "/redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv/HZ"
	// AgentsPublicAPIsV1SrvIngestIncidentProcedure is the fully-qualified name of the
	// AgentsPublicAPIsV1Srv's IngestIncident RPC.
	AgentsPublicAPIsV1SrvIngestIncidentProcedure = "/redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv/IngestIncident"
	// AgentsPublicAPIsV1SrvFetchAgentConfigurationProcedure is the fully-qualified name of the
	// AgentsPublicAPIsV1Srv's FetchAgentConfiguration RPC.
	AgentsPublicAPIsV1SrvFetchAgentConfigurationProcedure = "/redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv/FetchAgentConfiguration"
	// AgentsPublicAPIsV1SrvFetchAgentRequestsProcedure is the fully-qualified name of the
	// AgentsPublicAPIsV1Srv's FetchAgentRequests RPC.
	AgentsPublicAPIsV1SrvFetchAgentRequestsProcedure = "/redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv/FetchAgentRequests"
	// AgentsPublicAPIsV1SrvSubmitAgentResponseProcedure is the fully-qualified name of the
	// AgentsPublicAPIsV1Srv's SubmitAgentResponse RPC.
	AgentsPublicAPIsV1SrvSubmitAgentResponseProcedure = "/redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv/SubmitAgentResponse"
)

// AgentsPublicAPIsV1SrvClient is a client for the redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv
// service.
type AgentsPublicAPIsV1SrvClient interface {
	HZ(context.Context, *connect.Request[v1.HZRequest]) (*connect.Response[v1.HZResponse], error)
	IngestIncident(context.Context, *connect.Request[v1.IngestIncidentRequest]) (*connect.Response[v1.IngestIncidentResponse], error)
	FetchAgentConfiguration(context.Context, *connect.Request[v1.FetchAgentConfigurationRequest]) (*connect.Response[v1.FetchAgentConfigurationResponse], error)
	FetchAgentRequests(context.Context, *connect.Request[v1.FetchAgentRequestsRequest]) (*connect.Response[v1.FetchAgentRequestsResponse], error)
	SubmitAgentResponse(context.Context, *connect.Request[v1.SubmitAgentResponseRequest]) (*connect.Response[v1.SubmitAgentResponseResponse], error)
}

// NewAgentsPublicAPIsV1SrvClient constructs a client for the
// redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv service. By default, it uses the Connect
// protocol with the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed
// requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewAgentsPublicAPIsV1SrvClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) AgentsPublicAPIsV1SrvClient {
	baseURL = strings.TrimRight(baseURL, "/")
	agentsPublicAPIsV1SrvMethods := v1.File_redcarbon_agents_public_v1_v1_proto.Services().ByName("AgentsPublicAPIsV1Srv").Methods()
	return &agentsPublicAPIsV1SrvClient{
		hZ: connect.NewClient[v1.HZRequest, v1.HZResponse](
			httpClient,
			baseURL+AgentsPublicAPIsV1SrvHZProcedure,
			connect.WithSchema(agentsPublicAPIsV1SrvMethods.ByName("HZ")),
			connect.WithClientOptions(opts...),
		),
		ingestIncident: connect.NewClient[v1.IngestIncidentRequest, v1.IngestIncidentResponse](
			httpClient,
			baseURL+AgentsPublicAPIsV1SrvIngestIncidentProcedure,
			connect.WithSchema(agentsPublicAPIsV1SrvMethods.ByName("IngestIncident")),
			connect.WithClientOptions(opts...),
		),
		fetchAgentConfiguration: connect.NewClient[v1.FetchAgentConfigurationRequest, v1.FetchAgentConfigurationResponse](
			httpClient,
			baseURL+AgentsPublicAPIsV1SrvFetchAgentConfigurationProcedure,
			connect.WithSchema(agentsPublicAPIsV1SrvMethods.ByName("FetchAgentConfiguration")),
			connect.WithClientOptions(opts...),
		),
		fetchAgentRequests: connect.NewClient[v1.FetchAgentRequestsRequest, v1.FetchAgentRequestsResponse](
			httpClient,
			baseURL+AgentsPublicAPIsV1SrvFetchAgentRequestsProcedure,
			connect.WithSchema(agentsPublicAPIsV1SrvMethods.ByName("FetchAgentRequests")),
			connect.WithClientOptions(opts...),
		),
		submitAgentResponse: connect.NewClient[v1.SubmitAgentResponseRequest, v1.SubmitAgentResponseResponse](
			httpClient,
			baseURL+AgentsPublicAPIsV1SrvSubmitAgentResponseProcedure,
			connect.WithSchema(agentsPublicAPIsV1SrvMethods.ByName("SubmitAgentResponse")),
			connect.WithClientOptions(opts...),
		),
	}
}

// agentsPublicAPIsV1SrvClient implements AgentsPublicAPIsV1SrvClient.
type agentsPublicAPIsV1SrvClient struct {
	hZ                      *connect.Client[v1.HZRequest, v1.HZResponse]
	ingestIncident          *connect.Client[v1.IngestIncidentRequest, v1.IngestIncidentResponse]
	fetchAgentConfiguration *connect.Client[v1.FetchAgentConfigurationRequest, v1.FetchAgentConfigurationResponse]
	fetchAgentRequests      *connect.Client[v1.FetchAgentRequestsRequest, v1.FetchAgentRequestsResponse]
	submitAgentResponse     *connect.Client[v1.SubmitAgentResponseRequest, v1.SubmitAgentResponseResponse]
}

// HZ calls redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv.HZ.
func (c *agentsPublicAPIsV1SrvClient) HZ(ctx context.Context, req *connect.Request[v1.HZRequest]) (*connect.Response[v1.HZResponse], error) {
	return c.hZ.CallUnary(ctx, req)
}

// IngestIncident calls redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv.IngestIncident.
func (c *agentsPublicAPIsV1SrvClient) IngestIncident(ctx context.Context, req *connect.Request[v1.IngestIncidentRequest]) (*connect.Response[v1.IngestIncidentResponse], error) {
	return c.ingestIncident.CallUnary(ctx, req)
}

// FetchAgentConfiguration calls
// redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv.FetchAgentConfiguration.
func (c *agentsPublicAPIsV1SrvClient) FetchAgentConfiguration(ctx context.Context, req *connect.Request[v1.FetchAgentConfigurationRequest]) (*connect.Response[v1.FetchAgentConfigurationResponse], error) {
	return c.fetchAgentConfiguration.CallUnary(ctx, req)
}

// FetchAgentRequests calls redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv.FetchAgentRequests.
func (c *agentsPublicAPIsV1SrvClient) FetchAgentRequests(ctx context.Context, req *connect.Request[v1.FetchAgentRequestsRequest]) (*connect.Response[v1.FetchAgentRequestsResponse], error) {
	return c.fetchAgentRequests.CallUnary(ctx, req)
}

// SubmitAgentResponse calls redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv.SubmitAgentResponse.
func (c *agentsPublicAPIsV1SrvClient) SubmitAgentResponse(ctx context.Context, req *connect.Request[v1.SubmitAgentResponseRequest]) (*connect.Response[v1.SubmitAgentResponseResponse], error) {
	return c.submitAgentResponse.CallUnary(ctx, req)
}

// AgentsPublicAPIsV1SrvHandler is an implementation of the
// redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv service.
type AgentsPublicAPIsV1SrvHandler interface {
	HZ(context.Context, *connect.Request[v1.HZRequest]) (*connect.Response[v1.HZResponse], error)
	IngestIncident(context.Context, *connect.Request[v1.IngestIncidentRequest]) (*connect.Response[v1.IngestIncidentResponse], error)
	FetchAgentConfiguration(context.Context, *connect.Request[v1.FetchAgentConfigurationRequest]) (*connect.Response[v1.FetchAgentConfigurationResponse], error)
	FetchAgentRequests(context.Context, *connect.Request[v1.FetchAgentRequestsRequest]) (*connect.Response[v1.FetchAgentRequestsResponse], error)
	SubmitAgentResponse(context.Context, *connect.Request[v1.SubmitAgentResponseRequest]) (*connect.Response[v1.SubmitAgentResponseResponse], error)
}

// NewAgentsPublicAPIsV1SrvHandler builds an HTTP handler from the service implementation. It
// returns the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewAgentsPublicAPIsV1SrvHandler(svc AgentsPublicAPIsV1SrvHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	agentsPublicAPIsV1SrvMethods := v1.File_redcarbon_agents_public_v1_v1_proto.Services().ByName("AgentsPublicAPIsV1Srv").Methods()
	agentsPublicAPIsV1SrvHZHandler := connect.NewUnaryHandler(
		AgentsPublicAPIsV1SrvHZProcedure,
		svc.HZ,
		connect.WithSchema(agentsPublicAPIsV1SrvMethods.ByName("HZ")),
		connect.WithHandlerOptions(opts...),
	)
	agentsPublicAPIsV1SrvIngestIncidentHandler := connect.NewUnaryHandler(
		AgentsPublicAPIsV1SrvIngestIncidentProcedure,
		svc.IngestIncident,
		connect.WithSchema(agentsPublicAPIsV1SrvMethods.ByName("IngestIncident")),
		connect.WithHandlerOptions(opts...),
	)
	agentsPublicAPIsV1SrvFetchAgentConfigurationHandler := connect.NewUnaryHandler(
		AgentsPublicAPIsV1SrvFetchAgentConfigurationProcedure,
		svc.FetchAgentConfiguration,
		connect.WithSchema(agentsPublicAPIsV1SrvMethods.ByName("FetchAgentConfiguration")),
		connect.WithHandlerOptions(opts...),
	)
	agentsPublicAPIsV1SrvFetchAgentRequestsHandler := connect.NewUnaryHandler(
		AgentsPublicAPIsV1SrvFetchAgentRequestsProcedure,
		svc.FetchAgentRequests,
		connect.WithSchema(agentsPublicAPIsV1SrvMethods.ByName("FetchAgentRequests")),
		connect.WithHandlerOptions(opts...),
	)
	agentsPublicAPIsV1SrvSubmitAgentResponseHandler := connect.NewUnaryHandler(
		AgentsPublicAPIsV1SrvSubmitAgentResponseProcedure,
		svc.SubmitAgentResponse,
		connect.WithSchema(agentsPublicAPIsV1SrvMethods.ByName("SubmitAgentResponse")),
		connect.WithHandlerOptions(opts...),
	)
	return "/redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case AgentsPublicAPIsV1SrvHZProcedure:
			agentsPublicAPIsV1SrvHZHandler.ServeHTTP(w, r)
		case AgentsPublicAPIsV1SrvIngestIncidentProcedure:
			agentsPublicAPIsV1SrvIngestIncidentHandler.ServeHTTP(w, r)
		case AgentsPublicAPIsV1SrvFetchAgentConfigurationProcedure:
			agentsPublicAPIsV1SrvFetchAgentConfigurationHandler.ServeHTTP(w, r)
		case AgentsPublicAPIsV1SrvFetchAgentRequestsProcedure:
			agentsPublicAPIsV1SrvFetchAgentRequestsHandler.ServeHTTP(w, r)
		case AgentsPublicAPIsV1SrvSubmitAgentResponseProcedure:
			agentsPublicAPIsV1SrvSubmitAgentResponseHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedAgentsPublicAPIsV1SrvHandler returns CodeUnimplemented from all methods.
type UnimplementedAgentsPublicAPIsV1SrvHandler struct{}

func (UnimplementedAgentsPublicAPIsV1SrvHandler) HZ(context.Context, *connect.Request[v1.HZRequest]) (*connect.Response[v1.HZResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv.HZ is not implemented"))
}

func (UnimplementedAgentsPublicAPIsV1SrvHandler) IngestIncident(context.Context, *connect.Request[v1.IngestIncidentRequest]) (*connect.Response[v1.IngestIncidentResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv.IngestIncident is not implemented"))
}

func (UnimplementedAgentsPublicAPIsV1SrvHandler) FetchAgentConfiguration(context.Context, *connect.Request[v1.FetchAgentConfigurationRequest]) (*connect.Response[v1.FetchAgentConfigurationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv.FetchAgentConfiguration is not implemented"))
}

func (UnimplementedAgentsPublicAPIsV1SrvHandler) FetchAgentRequests(context.Context, *connect.Request[v1.FetchAgentRequestsRequest]) (*connect.Response[v1.FetchAgentRequestsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv.FetchAgentRequests is not implemented"))
}

func (UnimplementedAgentsPublicAPIsV1SrvHandler) SubmitAgentResponse(context.Context, *connect.Request[v1.SubmitAgentResponseRequest]) (*connect.Response[v1.SubmitAgentResponseResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("redcarbon.agents_public.v1.AgentsPublicAPIsV1Srv.SubmitAgentResponse is not implemented"))
}
