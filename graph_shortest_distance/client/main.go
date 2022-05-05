package main

import (
	"flag"
	pb "github.com/firebearrex/graph-shortest-distance-grpc-server-go/graph_shortest_distance/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strconv"
)

var addr = "0.0.0.0:50051"

func main() {
	method := flag.String("method", "dist", "Specify one of the following methods to use with the "+
		"client: post/dist/delete.\n"+
		"post - post a new graph. The first argument is the total number of vertices, "+
		"followed by a sequence of node values for representing [src -> dest] pairs.\n"+
		"dist = compute the shortest distance between two nodes.")

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Did not connect: %v\n", err)
	}

	defer conn.Close()

	client := pb.NewGraphServiceClient(conn)

	flag.Parse()
	args := flag.Args()

	switch *method {
	case "post":
		// Parse the inputs
		if len(args) < 1 {
			log.Fatalln("Insufficient number of arguments")
		} else if (len(args)-1)%2 != 0 {
			log.Fatalln("Make sure the number of values to represent the edges is even (in pairs)")
		}

		totalVertices, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			log.Fatalf("Invalid input: %s\n", args[0])
		}

		var edgesRaw = make([][2]int32, (len(args)-1)/2)
		for i := 1; i < len(args); i++ {
			src, err := strconv.ParseInt(args[i], 10, 32)
			if err != nil {
				log.Fatalf("Invalid input: %s\n", args[i])
			}
			i++

			dest, err := strconv.ParseInt(args[i], 10, 32)
			if err != nil {
				log.Fatalf("Invalid input: %s\n", args[i])
			}

			edgesRaw[i/2-1][0] = int32(src)
			edgesRaw[i/2-1][1] = int32(dest)
		}

		// Do the posting action
		doPost(client, int32(totalVertices), edgesRaw)
	case "dist":
		// Parse the inputs
		if len(args) > 3 {
			if len(args)%3 != 0 {
				log.Fatalln("Incorrect number of arguments. " +
					"Arguments must be given as pairs of 3 consecutive numerical values.")
			}

			var ids = make([]int32, len(args)/3)
			var srcs = make([]int32, len(args)/3)
			var dests = make([]int32, len(args)/3)
			for i := 0; i < len(args); i += 3 {
				id, err := strconv.ParseInt(args[i], 10, 32)
				if err != nil {
					log.Fatalf("Invalid input: %s\n", args[i])
				}

				src, err := strconv.ParseInt(args[i+1], 10, 32)
				if err != nil {
					log.Fatalf("Invalid input: %s\n", args[i+1])
				}

				dest, err := strconv.ParseInt(args[i+2], 10, 32)
				if err != nil {
					log.Fatalf("Invalid input: %s\n", args[i+2])
				}

				ids[i/3] = int32(id)
				srcs[i/3] = int32(src)
				dests[i/3] = int32(dest)
			}

			doDistStream(client, ids, srcs, dests)
		} else if len(args) == 3 {
			id, err := strconv.ParseInt(args[0], 10, 32)
			if err != nil {
				log.Fatalf("Invalid input: %s\n", args[0])
			}

			src, err := strconv.ParseInt(args[1], 10, 32)
			if err != nil {
				log.Fatalf("Invalid input: %s\n", args[1])
			}

			dest, err := strconv.ParseInt(args[2], 10, 32)
			if err != nil {
				log.Fatalf("Invalid input: %s\n", args[2])
			}

			doDist(client, int32(id), int32(src), int32(dest))
		} else {
			log.Fatalf("The [dist] method accepts 3 or more numeral arguments\n")
		}
	case "delete":
		// Parse the inputs
		if len(args) != 1 {
			log.Fatalf("The [delete] method accepts 1 numeral argument exactly\n")
		} else {
			id, err := strconv.ParseInt(args[0], 10, 32)
			if err != nil {
				log.Fatalf("Invalid input: %s\n", args[0])
			}

			doDelete(client, int32(id))
		}

	default:
		log.Fatalf("%s is not a valid method", *method)
	}
}
