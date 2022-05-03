package main

import pb "github.com/firebearrex/graph-shortest-distance-grpc-server-go/graph_shortest_distance/proto"

type Server struct {
	pb.GraphServiceServer
}
