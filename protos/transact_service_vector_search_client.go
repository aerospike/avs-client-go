package protos

import grpc "google.golang.org/grpc"

// Wrapping this because mockgen doesn't support generating mocks
// for type aliases. The generated protobuf code uses type aliases
// THIS IS ONLY USED FOR TESTING
type wrappedTransactService_VectorSearchClient interface {
	grpc.ServerStreamingClient[Neighbor]
}
