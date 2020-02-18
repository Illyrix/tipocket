package creator

import (
	"context"

	clusterTypes "github.com/pingcap/tipocket/pkg/cluster/types"
	"github.com/pingcap/tipocket/pkg/core"
)

// PocketCreator create pocket instances
type PocketCreator struct{}

// PocketClient runs pocket
type PocketClient struct{}

// Create client
func (PocketCreator) Create(node clusterTypes.ClientNode) core.Client {
	return PocketClient{}
}

// SetUp sets up the client.
func (PocketClient) SetUp(ctx context.Context, nodes []clusterTypes.ClientNode, idx int) error {
	return nil
}

// TearDown tears down the client.
func (PocketClient) TearDown(ctx context.Context, nodes []clusterTypes.ClientNode, idx int) error {
	return nil
}

// Invoke invokes a request to the database.
func (PocketClient) Invoke(ctx context.Context, node clusterTypes.ClientNode, r interface{}) interface{} {
	return nil
}

// NextRequest generates a request for latter Invoke.
func (PocketClient) NextRequest() interface{} {
	return nil
}

// DumpState the database state(also the model's state)
func (PocketClient) DumpState(ctx context.Context) (interface{}, error) {
	return nil, nil
}

// Start runs self scheduled cases
func (PocketClient) Start(ctx context.Context, cfg interface{}, dsns []string) error {
	// upstream, downstream := dsns[0], dsns[1]
	return nil
}
