package partition_v1

import (
	"context"
	apic "github.com/antinvestor/apis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"

	"math"
)

const ctxKeyService = "partitionClientKey"

func defaultPartitionClientOptions() []apic.ClientOption {
	return []apic.ClientOption{
		apic.WithEndpoint("partitions.api.antinvestor.com:443"),
		apic.WithGRPCDialOption(grpc.WithDisableServiceConfig()),
		apic.WithGRPCDialOption(grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32))),
	}
}

func ToContext(ctx context.Context, partitionsClient *PartitionClient) context.Context {
	return context.WithValue(ctx, ctxKeyService, partitionsClient)
}

func FromContext(ctx context.Context) *PartitionClient {
	partitionsClient, ok := ctx.Value(ctxKeyService).(*PartitionClient)
	if !ok {
		return nil
	}

	return partitionsClient
}

// PartitionClient is a client for interacting with the partitions service API.
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
type PartitionClient struct {
	// gRPC connection to the service.
	clientConn *grpc.ClientConn

	// The gRPC API client.
	client PartitionServiceClient

	// The x-ant-* metadata to be sent with each request.
	xMetadata metadata.MD
}

// NewPartitionsClient creates a new partitions client.
/// The service that an application uses to access and manipulate partition information
func NewPartitionsClient(ctx context.Context, opts ...apic.ClientOption) (*PartitionClient, error) {
	clientOpts := defaultPartitionClientOptions()

	connPool, err := apic.DialConnection(ctx, append(clientOpts, opts...)...)
	if err != nil {
		return nil, err
	}
	cl := &PartitionClient{
		clientConn: connPool,
		client:     NewPartitionServiceClient(connPool),
	}

	cl.setClientInfo()

	return cl, nil
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (partCl *PartitionClient) Close() error {
	return partCl.clientConn.Close()
}

// setClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (partCl *PartitionClient) setClientInfo(keyval ...string) {
	kv := append([]string{"gl-go", apic.VersionGo()}, keyval...)
	kv = append(kv, "grpc", grpc.Version)
	partCl.xMetadata = metadata.Pairs("x-ai-api-client", apic.XAntHeader(kv...))
}

// NewTenant used to create a new tenant instance.
// This is a fairly static and infrequently used option that creates an almost physical data separation
// To allow the use of same databases in a multitentant fashion.
func (partCl *PartitionClient) NewTenant(ctx context.Context, name string,
	description string, props map[string]string) (*TenantObject, error) {

	profileCtx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	serviceClient := NewPartitionServiceClient(partCl.clientConn)

	request := TenantRequest{
		Name:        name,
		Description: description,
		Properties:  props,
	}

	return serviceClient.CreateTenant(profileCtx, &request)
}

// NewPartition Creates a further logical multitenant environment at a softer level.
// This separation at the partition level is enforced at the application level that is consuming the api.
func (partCl *PartitionClient) NewPartition(ctx context.Context, tenantId string, name string, description string,
	props map[string]string) (*PartitionObject, error) {
	return partCl.newPartition(ctx, tenantId, "", name, description, props)
}

// NewChildPartition partitions can have children, for example a bank can have multiple branches
func (partCl *PartitionClient) NewChildPartition(ctx context.Context, tenantId string, parentId string, name string,
	description string, props map[string]string) (*PartitionObject, error) {
	return partCl.newPartition(ctx, tenantId, parentId, name, description, props)
}

func (partCl *PartitionClient) newPartition(ctx context.Context, tenantId string,
	parentId string, name string, description string, props map[string]string) (*PartitionObject, error) {

	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	serviceClient := NewPartitionServiceClient(partCl.clientConn)

	request := PartitionCreateRequest{
		TenantId:    tenantId,
		ParentId:    parentId,
		Name:        name,
		Description: description,
		Properties:  props,
	}

	return serviceClient.CreatePartition(cancelCtx, &request)
}


func (partCl *PartitionClient) UpdatePartition(ctx context.Context, partitionId string,
	name string, description string, props map[string]string) (*PartitionObject, error) {

	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	serviceClient := NewPartitionServiceClient(partCl.clientConn)

	request := PartitionUpdateRequest{
		PartitionId:    partitionId,
		Name:        name,
		Description: description,
		Properties:  props,
	}

	return serviceClient.UpdatePartition(cancelCtx, &request)
}


func (partCl *PartitionClient) CreatePartitionRole(ctx context.Context, partitionId string,
	name string, props map[string]string) (*PartitionRoleObject, error) {

	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	serviceClient := NewPartitionServiceClient(partCl.clientConn)

	request := PartitionRoleCreateRequest{
		Name: name,
		PartitionId:  partitionId,
		Properties:  props,
	}

	return serviceClient.CreatePartitionRole(cancelCtx, &request)
}

func (partCl *PartitionClient) RemovePartitionRole(ctx context.Context, partitionRoleId string) (*RemoveResponse, error) {

	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	serviceClient := NewPartitionServiceClient(partCl.clientConn)

	request := PartitionRoleRemoveRequest{
		PartitionRoleId: partitionRoleId,
	}

	return serviceClient.RemovePartitionRole(cancelCtx, &request)
}

func (partCl *PartitionClient) ListPartitionRoles(ctx context.Context, partitionId string) (*PartitionRoleListResponse, error) {

	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	serviceClient := NewPartitionServiceClient(partCl.clientConn)

	partitionRoleRequest := PartitionRoleListRequest{
		PartitionId:  partitionId,
	}

	return serviceClient.ListPartitionRoles(cancelCtx, &partitionRoleRequest)
}


// NewPage a partition has a provision to store custom pages that can be shown to users later.
// These pages can include signup or customer specified customized pictures
func (partCl *PartitionClient) NewPage(ctx context.Context, partitionId string, name string, html string) (*PageObject, error) {

	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	serviceClient := NewPartitionServiceClient(partCl.clientConn)

	request := PageCreateRequest{
		Name:        name,
		Html: html,
		PartitionId:  partitionId,
	}

	return serviceClient.CreatePage(cancelCtx, &request)
}

// GetPage simple way to quickly pull custom pages accessed by clients of a partition
func (partCl *PartitionClient) GetPage(ctx context.Context, partitionId string, name string) (*PageObject, error) {

	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	serviceClient := NewPartitionServiceClient(partCl.clientConn)

	request := PageGetRequest{
		Name:        name,
		PartitionId:  partitionId,
	}

	return serviceClient.GetPage(cancelCtx, &request)
}

func (partCl *PartitionClient) CreateAccess(ctx context.Context, partitionId string, profileId string) (*AccessObject, error) {

	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	serviceClient := NewPartitionServiceClient(partCl.clientConn)

	request := AccessCreateRequest{
		ProfileId: profileId,
		PartitionId:  partitionId,
	}

	return serviceClient.CreateAccess(cancelCtx, &request)
}

func (partCl *PartitionClient) RemoveAccess(ctx context.Context, accessId string) (*RemoveResponse, error) {

	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	serviceClient := NewPartitionServiceClient(partCl.clientConn)

	request := AccessRemoveRequest{
		AccessId: accessId,
	}

	return serviceClient.RemoveAccess(cancelCtx, &request)
}

func (partCl *PartitionClient) GetAccess(ctx context.Context, partitionId string, profileId string) (*AccessObject, error) {

	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	serviceClient := NewPartitionServiceClient(partCl.clientConn)

	request := AccessGetRequest{
		ProfileId: profileId,
		PartitionId:  partitionId,
	}

	return serviceClient.GetAccess(cancelCtx, &request)
}


func (partCl *PartitionClient) CreateAccessRole(ctx context.Context, accessId string, partitionRoleId string) (*AccessRoleObject, error) {

	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	serviceClient := NewPartitionServiceClient(partCl.clientConn)

	request := AccessRoleCreateRequest{
		AccessId: accessId,
		PartitionRoleId:  partitionRoleId,
	}

	return serviceClient.CreateAccessRole(cancelCtx, &request)
}

func (partCl *PartitionClient) RemoveAccessRole(ctx context.Context, accessRoleId string) (*RemoveResponse, error) {

	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	serviceClient := NewPartitionServiceClient(partCl.clientConn)

	request := AccessRoleRemoveRequest{
		AccessRoleId: accessRoleId,
	}

	return serviceClient.RemoveAccessRole(cancelCtx, &request)
}

func (partCl *PartitionClient) ListAccess(ctx context.Context, accessId string) (*AccessRoleListResponse, error) {

	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	serviceClient := NewPartitionServiceClient(partCl.clientConn)

	request := AccessRoleListRequest{
		AccessId: accessId,
	}

	return serviceClient.ListAccessRoles(cancelCtx, &request)
}
