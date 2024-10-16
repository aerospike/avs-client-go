package protos

import grpc "google.golang.org/grpc"

// Wrapping this because mockgen doesn't support generating mocks
// for type aliases. The generated protobuf code uses type aliases
// THIS IS ONLY USED FOR TESTING
//
//nolint:revive,stylecheck // This is a generated type, ignore _ in name
type wrappedTransactService_VectorSearchClient interface {
	grpc.ServerStreamingClient[Neighbor]
}
