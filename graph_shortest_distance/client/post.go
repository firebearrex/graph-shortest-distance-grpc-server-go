package main

import (
	"context"
	pb "github.com/firebearrex/graph-shortest-distance-grpc-server-go/graph_shortest_distance/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func doPost(client pb.GraphServiceClient, totalVertices int32, edgesRaw [][2]int32) {
	log.Println("Posting new graph now...")

	edgesPb := make([]*pb.Edge, len(edgesRaw))

	for i := 0; i < len(edgesRaw); i++ {
		// src := edgesRaw[i][0]
		// if src >= totalVertices {
		// 	log.Fatalf("The node value [%d] is greater than or equal to the total number of nodes, "+
		// 		"meaning the node does not exist in the graph", src)
		// }

		edgesPb[i] = &pb.Edge{Src: edgesRaw[i][0], Dest: edgesRaw[i][1]}
	}

	res, err := client.Post(context.Background(), &pb.PostRequest{
		TotalVertices: totalVertices,
		Edges:         edgesPb,
	})

	// Error handling
	if err != nil {
		sts, ok := status.FromError(err)

		if ok {
			log.Printf("Error message from server: %v\n", sts.Message())
			log.Printf("Error code: %d\n", sts.Code())

			if sts.Code() == codes.InvalidArgument {
				log.Fatalf("Please check if the node values representing the edges are all valid.\n")
			}
		} else {
			log.Fatalf("A non gRPC error: %v\n", err)
		}

	}

	log.Printf("New graph posted to server successfully. Graph ID: %d\n", res.Result)
}
