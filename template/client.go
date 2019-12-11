package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	types "github.com/mesg-foundation/engine/protobuf/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	// env variables for configure client.
	envMesgEndpoint     = "MESG_ENDPOINT"
	envMesgInstanceHash = "MESG_INSTANCE_HASH"
	envMesgRunnerHash   = "MESG_RUNNER_HASH"
)

// Client is a struct to wrap all mesg exposed protobuf API.
type Client struct {
	// all clients registered by mesg server.
	pb.EventClient
	pb.ExecutionClient

	// instance hash that could be used in api calls.
	InstanceHash hash.Hash

	// runner hash that could be used in api calls.
	RunnerHash hash.Hash
}

// NewClient creates a new client from env variables supplied by mesg engine.
func NewClient() (*Client, error) {
	endpoint := os.Getenv(envMesgEndpoint)
	if endpoint == "" {
		return nil, fmt.Errorf("client: server address env(%s) is empty", envMesgEndpoint)
	}

	instanceHash, err := hash.Decode(os.Getenv(envMesgInstanceHash))
	if err != nil {
		return nil, fmt.Errorf("client: error with instance hash env(%s): %s", envMesgInstanceHash, err)
	}

	runnerHash, err := hash.Decode(os.Getenv(envMesgRunnerHash))
	if err != nil {
		return nil, fmt.Errorf("client: error with runner hash env(%s): %s", envMesgRunnerHash, err)
	}

	dialoptions := []grpc.DialOption{
		// Keep alive prevents Docker network to drop TCP idle connections after 15 minutes.
		// See: https://forum.mesg.com/t/solution-summary-for-docker-dropping-connections-after-15-min/246
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: 5 * time.Minute, // 5 minutes is the minimun time of gRPC enforcement policy.
		}),
		grpc.WithInsecure(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, endpoint, dialoptions...)
	if err != nil {
		return nil, fmt.Errorf("client: connection error: %s", err)
	}

	return &Client{
		EventClient:     pb.NewEventClient(conn),
		ExecutionClient: pb.NewExecutionClient(conn),
		InstanceHash:    instanceHash,
		RunnerHash:      runnerHash,
	}, nil
}

// StreamExecution returns stream that recives executions with status in-progress.
func (c *Client) StreamExecution() (pb.Execution_StreamClient, error) {
	stream, err := c.ExecutionClient.Stream(context.Background(), &pb.StreamExecutionRequest{
		Filter: &pb.StreamExecutionRequest_Filter{
			Statuses:     []execution.Status{execution.Status_InProgress},
			ExecutorHash: c.RunnerHash,
		},
	})
	if err != nil {
		return nil, err
	}
	if err := acknowledgement.WaitForStreamToBeReady(stream); err != nil {
		return nil, err
	}
	return stream, nil
}

// CreateEvent is common wrapper for create events.
func (c *Client) CreateEvent(key string, data map[string]*types.Value) error {
	_, err := c.EventClient.Create(context.Background(), &pb.CreateEventRequest{
		InstanceHash: c.InstanceHash,
		Key:          key,
		Data: &types.Struct{
			Fields: data,
		},
	})
	return err
}
