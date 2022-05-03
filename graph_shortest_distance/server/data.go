package main

import (
	pb "github.com/firebearrex/graph-shortest-distance-grpc-server-go/graph_shortest_distance/proto"
)

var idHead int32

var graphStore map[int32]Graph

type Graph struct {
	totalVertices int32
	edges         []*pb.Edge
}
