package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"

	pb "github.com/firebearrex/graph-shortest-distance-grpc-server-go/graph_shortest_distance/proto"
)

// TestServer_Post tests for successful posting actions
func TestServer_Post(t *testing.T) {
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

	tests := []struct {
		expected      int32
		totalVertices int32
		edgesPb       []*pb.Edge
	}{
		{
			expected:      0,
			totalVertices: 4,
			edgesPb: []*pb.Edge{
				{Src: 0, Dest: 1},
				{Src: 1, Dest: 2},
				{Src: 1, Dest: 3},
				{Src: 3, Dest: 0},
			},
		},
		{
			expected:      1,
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
			expected:      2,
			totalVertices: 3,
			edgesPb: []*pb.Edge{
				{Src: 0, Dest: 1},
			},
		},
	}

	for _, tt := range tests {
		res, err := client.Post(context.Background(), &pb.PostRequest{
			TotalVertices: tt.totalVertices,
			Edges:         tt.edgesPb,
		})

		if err != nil {
			t.Errorf("Post(%+v) got unexpected error", tt)
		}

		if res.Result != tt.expected {
			t.Errorf("Post(%+v) = %v, expected: %v", tt, res.Result, tt.expected)
		}
	}
}

// TestServer_PostInvalidInput tests for failed posting actions due to invalid arguments
func TestServer_PostInvalidInput(t *testing.T) {
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

	// Source node does not exist
	req1 := &pb.PostRequest{
		TotalVertices: 3,
		Edges: []*pb.Edge{
			{Src: 3, Dest: 1},
		},
	}

	res1, err1 := client.Post(context.Background(), req1)
	if err1 == nil {
		t.Fatal("Failed to catch expected error\n")
	} else if res1 != nil {
		t.Fatalf("Post(%+v) = %v, expected: nil", req1, res1.Result)
	}

	// Destination node does not exist
	req2 := &pb.PostRequest{
		TotalVertices: 3,
		Edges: []*pb.Edge{
			{Src: 1, Dest: 3},
		},
	}

	res2, err2 := client.Post(context.Background(), req2)
	if err2 == nil {
		t.Fatal("Failed to catch expected error\n")
	} else if res2 != nil {
		t.Fatalf("Post(%+v) = %v, expected: nil", req2, res2.Result)
	}

	// TotalVertices < 0
	req3 := &pb.PostRequest{
		TotalVertices: -1,
		Edges: []*pb.Edge{
			{Src: -1, Dest: -1},
		},
	}

	res3, err3 := client.Post(context.Background(), req3)
	if err3 == nil {
		t.Fatal("Failed to catch expected error\n")
	} else if res3 != nil {
		t.Fatalf("Post(%+v) = %v, expected: nil", req3, res3.Result)
	}

	// Src < 0
	req4 := &pb.PostRequest{
		TotalVertices: 3,
		Edges: []*pb.Edge{
			{Src: -1, Dest: 1},
		},
	}

	res4, err4 := client.Post(context.Background(), req4)
	if err4 == nil {
		t.Fatal("Failed to catch expected error\n")
	} else if res4 != nil {
		t.Fatalf("Post(%+v) = %v, expected: nil", req4, res4.Result)
	}

	// Dest < 0
	req5 := &pb.PostRequest{
		TotalVertices: 3,
		Edges: []*pb.Edge{
			{Src: 1, Dest: -1},
		},
	}

	res5, err5 := client.Post(context.Background(), req5)
	if err5 == nil {
		t.Fatal("Failed to catch expected error\n")
	} else if res5 != nil {
		t.Fatalf("Post(%+v) = %v, expected: nil", req5, res5.Result)
	}
}
