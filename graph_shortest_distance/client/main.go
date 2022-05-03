package main

import (
	"flag"
	pb "github.com/firebearrex/graph-shortest-distance-grpc-server-go/graph_shortest_distance/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strconv"
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

	flag.Parse()

	switch *method {
	case "post":
		// Parse the inputs
		args := flag.Args()
		if len(args) < 1 {
			log.Fatalln("Insufficient number of arguments")
		} else if (len(args)-1)%2 != 0 {
			log.Fatalln("Make sure the number of values to represent the edges is even (in pairs)")
		}

		totalVertices, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			log.Fatalf("Invalid input: %s\n", args[0])
		}

		var edgesRaw [][2]int32 = make([][2]int32, (len(args)-1)/2)
		for i := 1; i < len(args); i++ {
			src, err := strconv.ParseInt(args[i], 10, 32)
			if err != nil {
				log.Fatalf("Invalid input: %s\n", args[i])
			}
			i++

			dst, err := strconv.ParseInt(args[i], 10, 32)
			if err != nil {
				log.Fatalf("Invalid input: %s\n", args[i])
			}

			edgesRaw[i/2-1][0] = int32(src)
			edgesRaw[i/2-1][1] = int32(dst)
		}

		doPost(client, int32(totalVertices), edgesRaw)
	case "dist":
		log.Println("Dist method has not been implemented yet")
		// Parse the inputs
		// doDist(client)
	case "delete":
		log.Println("Delete method has not been implemented yet")
		// Parse the inputs
		// doDelete(client)
	default:
		log.Fatalf("%s is not a valid method", *method)
	}
}
