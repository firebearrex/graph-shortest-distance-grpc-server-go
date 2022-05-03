package main

import (
	"context"
	"log"

	pb "github.com/firebearrex/graph-shortest-distance-grpc-server-go/graph_shortest_distance/proto"
)

func (*Server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	log.Printf("Delete was invoked with: %v\n", req)

	_, ok := graphStore[req.Id]
	delete(graphStore, req.Id)

	return &pb.DeleteResponse{Result: ok}, nil
}
