package main

import (
	"context"
	"log"

	pb "github.com/firebearrex/graph-shortest-distance-grpc-server-go/graph_shortest_distance/proto"
)

func (*Server) Post(ctx context.Context, req *pb.PostRequest) (*pb.PostResponse, error) {
	log.Printf("Post was invoked with %v\n", req)

	newGraph := Graph{
		totalVertices: req.TotalVertices,
		edges:         req.Edges,
	}
	currId := idHead
	graphStore[idHead] = newGraph
	idHead++

	return &pb.PostResponse{Result: currId}, nil
}
