package main

import (
	"context"
	pb "github.com/firebearrex/graph-shortest-distance-grpc-server-go/graph_shortest_distance/proto"
	"log"
)

func doPost(client pb.GraphServiceClient, totalVertices int32, edgesRaw [][2]int32) {
	log.Println("Posting new graph now...")

	edgesPb := make([]*pb.Edge, len(edgesRaw))

	for i := 0; i < len(edgesPb); i++ {
		edgesPb[i] = &pb.Edge{Src: edgesRaw[i][0], Dest: edgesRaw[i][1]}
	}

	res, err := client.Post(context.Background(), &pb.PostRequest{
		TotalVertices: totalVertices,
		Edges:         edgesPb,
	})

	if err != nil {
		log.Fatalf("Failed to post new graph: %v\n", err)
	}

	log.Printf("New graph posted to server successfully. Graph ID: %d\n", res.Result)
}
