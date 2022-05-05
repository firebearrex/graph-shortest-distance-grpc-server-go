package main

import (
	"context"
	pb "github.com/firebearrex/graph-shortest-distance-grpc-server-go/graph_shortest_distance/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
)

// doDistStream executes the client request
func doDistStream(client pb.GraphServiceClient, ids []int32, srcs []int32, dests []int32) {
	log.Println("Processing multiple shortest distance requests now...")

	// Parameter validation
	if len(ids) != len(srcs) || len(srcs) != len(dests) {
		log.Fatalln("Invalid length of input parameters. Number of parameters for graph ID, source node, " +
			"and destination node must equal.")
	}

	reqLen := len(ids)
	stream, err := client.DistStream(context.Background())

	if err != nil {
		log.Fatalf("Error while opening stream: %v\n", err)
	}

	streamSync := make(chan struct{})

	go func() {
		for i := 0; i < reqLen; i++ {
			req := pb.DistRequest{
				Id:   ids[i],
				Src:  srcs[i],
				Dest: dests[i],
			}
			log.Printf("Sending request: %+v\n", req)
			stream.Send(&req)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				sts, ok := status.FromError(err)

				if ok {
					log.Printf("Error message from server: %v\n", sts.Message())
					log.Printf("Error code: %d\n", sts.Code())

					if sts.Code() == codes.InvalidArgument {
						log.Fatalf("Please check if the specified source node or destination node exist in the graph.\n")
					} else if sts.Code() == codes.NotFound {
						log.Fatalf("Please check if the graph ID is correct.\n")
					}
				} else {
					log.Fatalf("A non gRPC error: %v\n", err)
				}

				break
			}

			if res.Result == math.MaxInt32 {
				log.Printf("The source node [%d] and destination node [%d] in graph[id=%d] are not connected.\n",
					res.Src, res.Dest, res.Id)
			} else {
				log.Printf("The shortest distance between node [%d] and node [%d] in graph[id=%d] is: %d\n",
					res.Src, res.Dest, res.Id, res.Result)
			}
		}
		close(streamSync)
	}()

	<-streamSync
}
