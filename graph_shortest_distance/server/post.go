package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"

	pb "github.com/firebearrex/graph-shortest-distance-grpc-server-go/graph_shortest_distance/proto"
)

func (*Server) Post(ctx context.Context, req *pb.PostRequest) (*pb.PostResponse, error) {
	log.Printf("Post was invoked with: %v\n", req)

	totalVertices := req.TotalVertices
	edges := req.Edges

	// Parameter validation
	if totalVertices < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Invalid total number of vertices: %d. Must not be negative.", totalVertices),
		)
	}

	for _, edge := range edges {
		if edge.Src < 0 {
			return nil, status.Errorf(
				codes.InvalidArgument,
				fmt.Sprintf("Invalid edge node: %d. Must not be negative.", edge.Src),
			)
		}
		if edge.Dest < 0 {
			return nil, status.Errorf(
				codes.InvalidArgument,
				fmt.Sprintf("Invalid edge node: %d. Must not be negative.", edge.Dest),
			)
		}
		if edge.Src >= totalVertices {
			return nil, status.Errorf(
				codes.InvalidArgument,
				fmt.Sprintf("The node value [%d] is greater than or equal to the total number of nodes, "+
					"meaning the node does not exist in the graph", edge.Src),
			)
		}
		if edge.Dest >= totalVertices {
			return nil, status.Errorf(
				codes.InvalidArgument,
				fmt.Sprintf("The node value [%d] is greater than or equal to the total number of nodes, "+
					"meaning the node does not exist in the graph", edge.Dest),
			)
		}
	}

	// Saving the graph
	newGraph := Graph{
		totalVertices: totalVertices,
		edges:         edges,
	}
	currId := idHead
	graphStore[idHead] = newGraph
	idHead++

	return &pb.PostResponse{Result: currId}, nil
}
