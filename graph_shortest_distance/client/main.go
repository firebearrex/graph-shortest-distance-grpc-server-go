package main

import (
	"flag"
	pb "github.com/firebearrex/graph-shortest-distance-grpc-server-go/graph_shortest_distance/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var addr string = "0.0.0.0:50051"

func main() {
	method := flag.String("method", "dist", "Specify one of the following methods to use with the "+
		"client: post/dist/delete. Default usage is 'dist' for querying the shortest distance between two graph nodes.")

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Did not connect: %v\n", err)
	}

	defer conn.Close()

	client := pb.NewGraphServiceClient(conn)

	switch *method {
	case "post":
		// Parse the inputs
		doPost(client)
	case "dist":
		// Parse the inputs
		//doDist(client)
	case "delete":
		// Parse the inputs
		//doDelete(client)
	default:
		log.Fatalf("%s is not a valid method", *method)
	}
}
