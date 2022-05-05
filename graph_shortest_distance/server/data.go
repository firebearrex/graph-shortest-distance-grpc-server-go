package main

import (
	pb "github.com/firebearrex/graph-shortest-distance-grpc-server-go/graph_shortest_distance/proto"
)

// idHead is used to keep track of the next ID to assign to the next new graph
var idHead int32

// graphStore serves as an in-memory data store to keep the ID -> graph key-value pairs
var graphStore map[int32]Graph

// Graph is the abstract structure for representing a graph
type Graph struct {
	totalVertices int32
	edges         []*pb.Edge
}
