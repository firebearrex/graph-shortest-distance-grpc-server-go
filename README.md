# graph-shortest-distance-grpc-server-go

**Author:** Cheng Zhao

**Email:** cz.rex@outlook.com

## About/Overview
This program is a Golang microservice built on gRPC API framework and HTTP/2. It relies on Protocol Buffers for 
server interface definition. The server offers functionalities for posting a graph representation, computing the 
shortest distance between two specified nodes, and deleting a previously posted graph from the data store.

## List of Features
* Post a graph, returning an ID to be used in subsequent operations
* Get the shortest path between two vertices in a previously posted graph
* Delete a graph from the server

## How to Run

* ### Compile the protocol definition
    * To generate the Protocol Buffer client/server interface boilerplate code, run the following command from the root directory:  
      `protoc -Igraph_shortest_distance/proto --go_opt=module=github.com/firebearrex/graph-shortest-distance-grpc-server-go --go_out=. --go-grpc_opt=module=github.com/firebearrex/graph-shortest-distance-grpc-server-go --go-grpc_out=. graph_shortest_distance/proto/*.proto`
* ### Build client/server executables
  * The client/server executables can be found under _bin/graph_shortest_distance/_ directory by running the 
    following command:  
    `go build -o bin/graph_shortest_distance/server ./graph_shortest_distance/server; go build -o bin/graph_shortest_distance/client ./graph_shortest_distance/client`
* ### Run the server
  * When server executable has been generated from the previous step, run the following command from the root 
    directory to start the server, so it can respond to client requests:  
    `./bin/graph_shortest_distance/server`
* ### Run the client
  * When client executable has been generated from the previous step, run the following command from the root
    directory to use the client to trigger the desired method with appropriate arguments required by that method:  
    `./bin/graph_shortest_distance/client -method=[<post>/<dist>/<delete> | default=dist] [args]`  
    Refer to the next section __How to Use the Program__ for more information regarding the program arguments. 

## How to Use the Program

* ### Post a graph
  * For posting a graph, the arguments are numerical values to represent the following attributes:
    * Total number of vertices
    * Pairs of nodes as an edge (e.g. 1 2 represents an edge from node 1 to node 2)
  * The first argument represents the total number of vertices, while the following arguments represent the edges. 
    The node values representing the edges must be in pairs, i.e. the number of arguments for representing the edges 
    must be even.
  * The following example post a new graph having 4 vertices in total and edges 0-1, 1-2, 1-3, 3-0:  
    `./bin/graph_shortest_distance/client -method=post 4 0 1 1 2 1 3 3 0`
  * After running the command, the program will respond with a prompt to show the newly posted graph's ID number. 
    This ID number can be used for computing the shortest distance of two nodes or deleting the associated graph.
  * If there is an error, the corresponding message will be prompted.

* ### Compute the shortest distance between two nodes
  * #### Unary gRPC API call for one single request
    * If this method only receives 3 arguments, it will trigger the unary call for this single request.
    * For computing the shortest distance of two nodes, the arguments are numerical values to represent the following
      attributes:
      * The graph's ID which is queried on
      * The source node
      * The destination node
    * The arguments should be specified following the same order as the order above, i.e. first argument - graph ID; 
      second argument - source node; third argument - destination node
    * The following example computes the shortest distance between node 1 and 3 in the graph with ID equal to 0:  
      `./bin/graph_shortest_distance/client -method=dist 0 1 3`
  * #### Bi-Directional gRPC streaming for multiple requests
    * If this method receives more than 3 arguments, it will trigger the bi-directional streaming for these multiple 
      requests.
    * The arguments should come as pairs of 3 for making a single request, and follow the same attributes order.
  * You can omit the [-method=dist] part as the default method to be used on this program is _dist_.
  * After running the command, the program will respond with a prompt to show the shortest distance between the two 
    nodes which are queried on.
  * If there is an error, the corresponding message will be prompted.

* ### Delete a graph
  * For deleting a graph, the arguments are numerical values to represent the following attributes:
    * The graph's ID which is to be deleted
  * The following example deletes the graph whose ID is 0:  
    `./bin/graph_shortest_distance/client -method=delete 0`
  * After running the command, the program will prompt accordingly depending on whether the graph has been 
    successfully deleted, or it does not exist in the server's data store. 

## Example Runs
An example run can be found under the res/ folder. The example run showcases the typical use cases of using 
this program, including:
* Using an incorrect method name
* Invalid arguments are given
* Posting new graphs
* Making queries for computing the shortest distance of two nodes
* Deleting a graph
* The requested graph does not exist

## Tests
* Test files can be found in the server's directory.
* Both `dist_test` and `dist_stream_test` contains performance testing for computing the shortest distances.

## Assumptions
* The graph nodes are represented as numerical values. If there are N vertices in the graph, then the values 0, 1, 2,
  ... , N - 1 represent each of nodes in this graph.
