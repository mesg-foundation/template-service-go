package main

import (
	"context"
	"log"

	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
)

var healtCheckData = map[string]*types.Value{
	"x": {
		Kind: &types.Value_BoolValue{
			BoolValue: true,
		},
	},
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to mesg server")

	estream, err := client.StreamExecution()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("created execution stream")
	defer estream.CloseSend()

	if err := client.CreateEvent("started", healtCheckData); err != nil {
		log.Fatal(err)
	}
	log.Println("emitted healt-check event")
	processExecutions(client, estream)
}

func processExecutions(c *Client, stream pb.Execution_StreamClient) {
	for {
		exec, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("received execution %s %s\n", exec.TaskKey, exec.Hash)
		req := &pb.UpdateExecutionRequest{
			Hash: exec.Hash,
		}

		if exec.Inputs.Fields["foo"].GetStringValue() != "hello" &&
			exec.Inputs.Fields["bar"].GetStringValue() != "world" {
			req.Result = &pb.UpdateExecutionRequest_Error{
				Error: "invalid inputs",
			}
		} else {
			req.Result = &pb.UpdateExecutionRequest_Outputs{
				Outputs: &types.Struct{
					Fields: map[string]*types.Value{
						"message": {
							Kind: &types.Value_StringValue{
								StringValue: "Hello world is valid",
							},
						},
					},
				},
			}
		}

		if _, err := c.ExecutionClient.Update(context.Background(), req); err != nil {
			log.Fatal(err)
		}
	}
}
