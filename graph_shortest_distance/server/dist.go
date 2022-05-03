package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"math"

	pb "github.com/firebearrex/graph-shortest-distance-grpc-server-go/graph_shortest_distance/proto"
)

// Dist computes the shortest distance between the source node and destination node in the graph specified in the
// request. If the specified source or destination node does not exist in the graph,
// the server will send error accordingly.
func (*Server) Dist(ctx context.Context, req *pb.DistRequest) (*pb.DistResponse, error) {
	log.Printf("Dist was invoked with %v\n", req)

	graph, ok := graphStore[req.Id]

	// The graph does not exist in the data store
	if !ok {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("The graph[id=%d] does not exist in the data store", req.Id),
		)
	}

	totalVertices := graph.totalVertices
	edges := graph.edges

	// Error handling
	if req.Src >= totalVertices {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("The source node [%d] does not exist in the graph", req.Src),
		)
	}
	if req.Dest >= totalVertices {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("The source node [%d] does not exist in the graph", req.Dest),
		)
	}

	adjList := make([][]int32, totalVertices)
	for _, edge := range edges {
		src := edge.Src
		dest := edge.Dest
		adjList[src] = append(adjList[src], dest)
		adjList[dest] = append(adjList[dest], src)
	}

	shortestDistance := getShortestDistance(totalVertices, req.Src, req.Dest, adjList)

	return &pb.DistResponse{Result: shortestDistance}, nil
}

// getShortestDistance takes the total number of vertices, the source node, the destination node,
// as well as the adjacency list, and returns the shortest distance between those two nodes.
// The function uses BFS algorithm, since the graph is undirected and unweighted.
// The time complexity of this algorithm is O(V+E), where V represents the number of vertices in the graph,
// and E represents the number of edges in the graph.
func getShortestDistance(totalVertices int32, src int32, dest int32, adjList [][]int32) int32 {
	// The dist list records the shortest distance of each vertex to the source node
	dist := make([]int32, totalVertices)
	for i := 0; i < len(dist); i++ {
		dist[i] = math.MaxInt32
	}

	// A queue to maintain queue of vertices whose adjacency list is to be scanned
	var queue []int32

	// Boolean array visited[] which stores the information whether ith vertex is reached at least once in the
	// breadth first search
	visited := make([]bool, totalVertices)

	dist[src] = 0
	visited[src] = true
	queue = offer(queue, src)

	// BFS algorithm
	for len(queue) != 0 {
		nextNode := poll(&queue)
		destFound := false
		for i := 0; i < len(adjList[nextNode]); i++ {
			if !visited[adjList[nextNode][i]] {
				visited[adjList[nextNode][i]] = true
				dist[adjList[nextNode][i]] = dist[nextNode] + 1
				queue = offer(queue, adjList[nextNode][i])

				if adjList[nextNode][i] == dest {
					destFound = true
					break
				}
			}
		}
		if destFound {
			break
		}
	}

	return dist[dest]
}

// offer takes the queue and enqueue the given element
func offer(queue []int32, element int32) []int32 {
	queue = append(queue, element) // Offer element to the queue
	return queue
}

// poll takes the queue and slice off the first element and return the element
func poll(queue *[]int32) int32 {
	element := (*queue)[0] // Poll element from the queue
	*queue = (*queue)[1:]  // Slice off the element once it is dequeued.
	return element
}
