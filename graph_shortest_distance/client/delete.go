package main

import (
	"context"
	pb "github.com/firebearrex/graph-shortest-distance-grpc-server-go/graph_shortest_distance/proto"
	"log"
)

// doDelete executes the client request
func doDelete(client pb.GraphServiceClient, id int32) {
	log.Println("Deleting the specified graph now...")

	res, err := client.Delete(context.Background(), &pb.DeleteRequest{
		Id: id,
	})

	// Error handling
	if err != nil {
		log.Fatalln("Unknown error")
	}

	if res.Result {
		log.Printf("The graph[id=%d] is succesfully deleted from the data store.\n", id)
	} else {
		log.Printf("The graph[id=%d] does not exist in the data store.\n", id)
	}
}
