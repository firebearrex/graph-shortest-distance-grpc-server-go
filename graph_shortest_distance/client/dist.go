package main

import (
	"context"
	pb "github.com/firebearrex/graph-shortest-distance-grpc-server-go/graph_shortest_distance/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"math"
)

func doDist(client pb.GraphServiceClient, id int32, src int32, dest int32) {
	log.Println("Calculating shortest distance ")

	res, err := client.Dist(context.Background(), &pb.DistRequest{
		Id:   id,
		Src:  src,
		Dest: dest,
	})

	if err != nil {
		sts, ok := status.FromError(err)

		if ok {
			log.Printf("Error message from server: %v\n", sts.Message())
			log.Printf("Error code: %d\n", sts.Code())

			if sts.Code() == codes.InvalidArgument {
				log.Println("Check if the specified source node or destination node exist in the graph.")
				return
			}
		} else {
			log.Fatalf("A non gRPC error: %v\n", err)
		}
	}

	if res.Result == math.MaxInt32 {
		log.Printf("The source node [%d] and destination node [%d] are not connected.\n", src, dest)
	} else {
		log.Printf("The shortest distance between node [%d] and node [%d] in graph with id=%d is: %d", src, dest, id,
			res.Result)
	}
}
