package main

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"

	pb "github.com/firebearrex/graph-shortest-distance-grpc-server-go/graph_shortest_distance/proto"
)

func (*Server) DistStream(stream pb.GraphService_DistStreamServer) error {
	log.Println("DistStream was invoked")

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
		}

		graph, ok := graphStore[req.Id]

		// The graph does not exist in the data store
		if !ok {
			return status.Errorf(
				codes.NotFound,
				fmt.Sprintf("The graph[id=%d] does not exist in the data store", req.Id),
			)
		}

		totalVertices := graph.totalVertices
		edges := graph.edges

		// Parameter validation
		if req.Src < 0 {
			return status.Errorf(
				codes.InvalidArgument,
				fmt.Sprintf("Invalid source node: %d. Must not be negative.", req.Src),
			)
		}
		if req.Dest < 0 {
			return status.Errorf(
				codes.InvalidArgument,
				fmt.Sprintf("Invalid destination node: %d. Must not be negative.", req.Dest),
			)
		}
		if req.Src >= totalVertices {
			return status.Errorf(
				codes.InvalidArgument,
				fmt.Sprintf("The source node [%d] does not exist in the graph", req.Src),
			)
		}
		if req.Dest >= totalVertices {
			return status.Errorf(
				codes.InvalidArgument,
				fmt.Sprintf("The destination node [%d] does not exist in the graph", req.Dest),
			)
		}

		// Build adjacency list
		adjList := make([][]int32, totalVertices)
		for _, edge := range edges {
			src := edge.Src
			dest := edge.Dest
			adjList[src] = append(adjList[src], dest)
			adjList[dest] = append(adjList[dest], src)
		}

		shortestDistance := getShortestDistance(totalVertices, req.Src, req.Dest, adjList)

		err = stream.Send(&pb.DistStreamResponse{
			Result: shortestDistance,
			Id:     req.Id,
			Src:    req.Src,
			Dest:   req.Dest,
		})

		if err != nil {
			log.Fatalf("Error while sending data to client: %v\n", err)
		}
	}
}
