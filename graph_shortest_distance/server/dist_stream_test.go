package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
	"testing"

	pb "github.com/firebearrex/graph-shortest-distance-grpc-server-go/graph_shortest_distance/proto"
)

// BenchmarkServer_DistStream serves as the performance testing
func BenchmarkServer_DistStream(b *testing.B) {
	idHead = 0
	graphStore = make(map[int32]Graph)

	ctx := context.Background()
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), creds)

	if err != nil {
		b.Fatalf("Failed to dial bufnet: %v", err)
	}

	defer conn.Close()
	client := pb.NewGraphServiceClient(conn)

	graphs := []struct {
		totalVertices int32
		edgesPb       []*pb.Edge
	}{
		{
			totalVertices: 4,
			edgesPb: []*pb.Edge{
				{Src: 0, Dest: 1},
				{Src: 1, Dest: 2},
				{Src: 1, Dest: 3},
				{Src: 3, Dest: 0},
			},
		},
		{
			totalVertices: 5,
			edgesPb: []*pb.Edge{
				{Src: 0, Dest: 1},
				{Src: 0, Dest: 4},
				{Src: 1, Dest: 2},
				{Src: 1, Dest: 3},
				{Src: 1, Dest: 4},
				{Src: 2, Dest: 3},
				{Src: 3, Dest: 4},
			},
		},
		{
			totalVertices: 8,
			edgesPb: []*pb.Edge{
				{Src: 0, Dest: 1},
				{Src: 0, Dest: 3},
				{Src: 1, Dest: 2},
				{Src: 3, Dest: 4},
				{Src: 3, Dest: 7},
				{Src: 4, Dest: 5},
				{Src: 4, Dest: 6},
				{Src: 4, Dest: 7},
				{Src: 5, Dest: 6},
				{Src: 6, Dest: 7},
			},
		},
		{
			totalVertices: 3,
			edgesPb: []*pb.Edge{
				{Src: 0, Dest: 1},
			},
		},
	}

	for _, graph := range graphs {
		_, err := client.Post(context.Background(), &pb.PostRequest{
			TotalVertices: graph.totalVertices,
			Edges:         graph.edgesPb,
		})

		if err != nil {
			b.Errorf("Post(%+v) got unexpected error", graph)
		}
	}

	// 0 1 3 0 0 2 1 0 2 1 4 1 1 3 0 1 4 2 2 3 7 2 2 0 2 6 0 2 1 6 2 2 6 3 0 1 3 0 2 3 2 1
	ids := []int32{0, 0, 1, 1, 1, 1, 2, 2, 2, 2, 2, 3, 3, 3}
	srcs := []int32{1, 0, 0, 4, 3, 4, 3, 2, 6, 1, 2, 0, 0, 2}
	dests := []int32{3, 2, 2, 1, 0, 2, 7, 0, 0, 6, 6, 1, 2, 1}

	for i := 0; i < b.N; i++ {
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
}
