package main

import (
	"context"
	pb "github.com/firebearrex/graph-shortest-distance-grpc-server-go/graph_shortest_distance/proto"
	"log"
)

func doPost(client pb.GraphServiceClient) {
	log.Println("Posting new graph now...")

	res, err := client.Post(context.Background())
}
