package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math"
	"testing"

	pb "github.com/firebearrex/graph-shortest-distance-grpc-server-go/graph_shortest_distance/proto"
)

// TestServer_Dist tests for computing the shortest distance
func TestServer_Dist(t *testing.T) {
	idHead = 0
	graphStore = make(map[int32]Graph)

	ctx := context.Background()
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), creds)

	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	defer conn.Close()
	client := pb.NewGraphServiceClient(conn)

	graphs := []struct {
		totalVertices int32
		edgesPb       []*pb.Edge
	}{
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
			t.Errorf("Post(%+v) got unexpected error", graph)
		}
	}

	tests := []struct {
		expected int32
		id       int32
		src      int32
		dest     int32
	}{
		{
			expected: 2,
			id:       0,
			src:      3,
			dest:     6,
		},
		{
			expected: 3,
			id:       0,
			src:      0,
			dest:     5,
		},
		{
			expected: 5,
			id:       0,
			src:      6,
			dest:     2,
		},
		{
			expected: 1,
			id:       1,
			src:      4,
			dest:     1,
		},
		{
			expected: 2,
			id:       1,
			src:      0,
			dest:     3,
		},
		{
			expected: 1,
			id:       2,
			src:      0,
			dest:     1,
		},
		{
			expected: math.MaxInt32,
			id:       2,
			src:      1,
			dest:     2,
		},
		{
			expected: 0,
			id:       0,
			src:      1,
			dest:     1,
		},
	}

	for _, tt := range tests {
		res, err := client.Dist(context.Background(), &pb.DistRequest{
			Id:   tt.id,
			Src:  tt.src,
			Dest: tt.dest,
		})

		if err != nil {
			t.Errorf("Dist(%+v) got unexpected error", tt)
		}

		if res.Result != tt.expected {
			t.Errorf("Dist(%+v) = %v, expected: %v", tt, res.Result, tt.expected)
		}
	}
}

// TestServer_DistInvalidInput tests for invalid parameters
func TestServer_DistInvalidInput(t *testing.T) {
	idHead = 0
	graphStore = make(map[int32]Graph)

	ctx := context.Background()
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), creds)

	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	defer conn.Close()
	client := pb.NewGraphServiceClient(conn)

	graphs := []struct {
		totalVertices int32
		edgesPb       []*pb.Edge
	}{
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
			t.Errorf("Post(%+v) got unexpected error", graph)
		}
	}

	// The queried graph does not exist
	req1 := &pb.DistRequest{
		Id:   3,
		Src:  0,
		Dest: 0,
	}

	res1, err1 := client.Dist(context.Background(), req1)
	if err1 == nil {
		t.Fatal("Failed to catch expected error\n")
	} else if res1 != nil {
		t.Fatalf("Dist(%+v) = %v, expected: nil", req1, res1.Result)
	}

	// Source node does not exist
	req2 := &pb.DistRequest{
		Id:   1,
		Src:  5,
		Dest: 0,
	}

	res2, err2 := client.Dist(context.Background(), req2)
	if err2 == nil {
		t.Fatal("Failed to catch expected error\n")
	} else if res2 != nil {
		t.Fatalf("Dist(%+v) = %v, expected: nil", req2, res2.Result)
	}

	// Destination node does not exist
	req3 := &pb.DistRequest{
		Id:   1,
		Src:  0,
		Dest: 5,
	}

	res3, err3 := client.Dist(context.Background(), req3)
	if err3 == nil {
		t.Fatal("Failed to catch expected error\n")
	} else if res3 != nil {
		t.Fatalf("Dist(%+v) = %v, expected: nil", req3, res3.Result)
	}

	// Src < 0
	req4 := &pb.DistRequest{
		Id:   2,
		Src:  -1,
		Dest: 0,
	}

	res4, err4 := client.Dist(context.Background(), req4)
	if err4 == nil {
		t.Fatal("Failed to catch expected error\n")
	} else if res4 != nil {
		t.Fatalf("Dist(%+v) = %v, expected: nil", req4, res4.Result)
	}

	// Dest < 0
	req5 := &pb.DistRequest{
		Id:   2,
		Src:  0,
		Dest: -1,
	}

	res5, err5 := client.Dist(context.Background(), req5)
	if err5 == nil {
		t.Fatal("Failed to catch expected error\n")
	} else if res5 != nil {
		t.Fatalf("Dist(%+v) = %v, expected: nil", req5, res5.Result)
	}
}

// BenchmarkServer_Dist serves as the performance testing
func BenchmarkServer_Dist(b *testing.B) {
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

	tests := []struct {
		distance int32
		id       int32
		src      int32
		dest     int32
	}{
		{
			distance: 0,
			id:       0,
			src:      2,
			dest:     2,
		},
		{
			distance: 1,
			id:       0,
			src:      2,
			dest:     1,
		},
		{
			distance: 2,
			id:       0,
			src:      2,
			dest:     0,
		},
		{
			distance: 3,
			id:       0,
			src:      2,
			dest:     3,
		},
		{
			distance: 4,
			id:       0,
			src:      2,
			dest:     7,
		},
		{
			distance: 5,
			id:       0,
			src:      2,
			dest:     6,
		},
	}

	for _, tt := range tests {
		b.Run(fmt.Sprintf("shortest_distance_%d", tt.distance), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				client.Dist(context.Background(), &pb.DistRequest{
					Id:   tt.id,
					Src:  tt.src,
					Dest: tt.dest,
				})
			}
		})
	}
}
