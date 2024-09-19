package avs

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/aerospike/avs-client-go/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestInsert_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	mockConn := &connection{
		transactClient: mockTransactClient,
	}

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	// Prepare the expected PutRequest
	expectedPutRequest := &protos.PutRequest{
		Key: &protos.Key{
			Namespace: "testNamespace",
			Set:       nil,
			Value: &protos.Key_StringValue{
				StringValue: "testKey",
			},
		},
		WriteType: protos.WriteType_INSERT_ONLY,
		Fields: []*protos.Field{
			{
				Name:  "field1",
				Value: &protos.Value{Value: &protos.Value_StringValue{StringValue: "value1"}},
			},
		},
		IgnoreMemQueueFull: false,
	}

	// Set up expectations for transactClient.Put()
	mockTransactClient.
		EXPECT().
		Put(gomock.Any(), gomock.Any()).
		Return(&emptypb.Empty{}, nil).
		Do(func(ctx context.Context, in *protos.PutRequest, opts ...grpc.CallOption) {
			assert.Equal(t, expectedPutRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := "testKey"
	recordData := map[string]interface{}{"field1": "value1"}
	ignoreMemQueueFull := false

	err = client.Insert(ctx, namespace, set, key, recordData, ignoreMemQueueFull)

	assert.NoError(t, err)
}

func TestInsert_FailsGettingConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(nil, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := "testKey"
	recordData := map[string]interface{}{"field1": "value1"}
	ignoreMemQueueFull := false

	err = client.Insert(ctx, namespace, set, key, recordData, ignoreMemQueueFull)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError(failedToInsertRecord, fmt.Errorf("foo")))
}

func TestInsert_FailsConvertingKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	mockConn := &connection{
		transactClient: mockTransactClient,
	}

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := struct{}{}
	recordData := map[string]interface{}{"field1": "value1"}
	ignoreMemQueueFull := false

	err = client.Insert(ctx, namespace, set, key, recordData, ignoreMemQueueFull)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError(failedToInsertRecord, fmt.Errorf("unsupported key type: struct {}")))
}

func TestInsert_FailsConvertingFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	mockConn := &connection{
		transactClient: mockTransactClient,
	}

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := "key"
	recordData := map[string]interface{}{"field1": struct{}{}}
	ignoreMemQueueFull := false

	err = client.Insert(ctx, namespace, set, key, recordData, ignoreMemQueueFull)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError(failedToInsertRecord, fmt.Errorf("error converting field value for key 'field1': unsupported value type: struct {}")))
}

func TestInsert_FailsPutRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	mockConn := &connection{
		transactClient: mockTransactClient,
	}

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	mockTransactClient.
		EXPECT().
		Put(gomock.Any(), gomock.Any()).
		Return(&emptypb.Empty{}, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := "key"
	recordData := map[string]interface{}{"field1": "value1"}
	ignoreMemQueueFull := false

	err = client.Insert(ctx, namespace, set, key, recordData, ignoreMemQueueFull)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError(failedToInsertRecord, fmt.Errorf("foo")))
}

func TestUpdate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock connProvider
	mockConnProvider := NewMockconnProvider(ctrl)
	// Create a mock transactClient
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	// Create a mock connection
	mockConn := &connection{
		transactClient: mockTransactClient,
	}

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	// Prepare the expected PutRequest
	expectedPutRequest := &protos.PutRequest{
		// You need to fill this with expected values
		Key: &protos.Key{
			Namespace: "testNamespace",
			Set:       nil,
			Value: &protos.Key_StringValue{
				StringValue: "testKey",
			},
		},
		WriteType: protos.WriteType_UPDATE_ONLY,
		Fields: []*protos.Field{
			{
				Name:  "field1",
				Value: &protos.Value{Value: &protos.Value_StringValue{StringValue: "value1"}},
			},
		},
		IgnoreMemQueueFull: false,
	}

	// Set up expectations for transactClient.Put()
	mockTransactClient.
		EXPECT().
		Put(gomock.Any(), gomock.Any()).
		Return(&emptypb.Empty{}, nil).
		Do(func(ctx context.Context, in *protos.PutRequest, opts ...grpc.CallOption) {
			// Optionally, you can assert that req matches expectedPutRequest
			assert.Equal(t, expectedPutRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := "testKey"
	recordData := map[string]interface{}{"field1": "value1"}
	ignoreMemQueueFull := false

	// Call the method under test
	err = client.Update(ctx, namespace, set, key, recordData, ignoreMemQueueFull)

	// Assert no error occurred
	assert.NoError(t, err)
}

func TestReplace_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock connProvider
	mockConnProvider := NewMockconnProvider(ctrl)
	// Create a mock transactClient
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	// Create a mock connection
	mockConn := &connection{
		transactClient: mockTransactClient,
	}

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	// Prepare the expected PutRequest
	expectedPutRequest := &protos.PutRequest{
		// You need to fill this with expected values
		Key: &protos.Key{
			Namespace: "testNamespace",
			Set:       nil,
			Value: &protos.Key_StringValue{
				StringValue: "testKey",
			},
		},
		WriteType: protos.WriteType_UPSERT,
		Fields: []*protos.Field{
			{
				Name:  "field1",
				Value: &protos.Value{Value: &protos.Value_StringValue{StringValue: "value1"}},
			},
		},
		IgnoreMemQueueFull: false,
	}

	// Set up expectations for transactClient.Put()
	mockTransactClient.
		EXPECT().
		Put(gomock.Any(), gomock.Any()).
		Return(&emptypb.Empty{}, nil).
		Do(func(ctx context.Context, in *protos.PutRequest, opts ...grpc.CallOption) {
			// Optionally, you can assert that req matches expectedPutRequest
			assert.Equal(t, expectedPutRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := "testKey"
	recordData := map[string]interface{}{"field1": "value1"}
	ignoreMemQueueFull := false

	// Call the method under test
	err = client.Upsert(ctx, namespace, set, key, recordData, ignoreMemQueueFull)

	// Assert no error occurred
	assert.NoError(t, err)
}

func TestGet_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	mockConn := &connection{
		transactClient: mockTransactClient,
	}

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	// Prepare the expected PutRequest
	expectedGetRequest := &protos.GetRequest{
		Key: &protos.Key{
			Namespace: "testNamespace",
			Set:       nil,
			Value: &protos.Key_StringValue{
				StringValue: "testKey",
			},
		},
		Projection: &protos.ProjectionSpec{
			Include: &protos.ProjectionFilter{
				Type: protos.ProjectionType_ALL,
			},
			Exclude: &protos.ProjectionFilter{
				Type: protos.ProjectionType_NONE,
			},
		},
	}

	protosRecord := &protos.Record{
		Fields: []*protos.Field{
			{Name: "field1", Value: &protos.Value{Value: &protos.Value_StringValue{StringValue: "value1"}}},
		},
		Metadata: &protos.Record_AerospikeMetadata{
			AerospikeMetadata: &protos.AerospikeRecordMetadata{
				Generation: 10,
				Expiration: 1,
			},
		},
	}

	expTime := AerospikeEpoch.Add(time.Second * 1)

	expectedRecord := &Record{
		Data: map[string]any{
			"field1": "value1",
		},
		Generation: 10,
		Expiration: &expTime,
	}

	mockTransactClient.
		EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(protosRecord, nil).
		Do(func(ctx context.Context, in *protos.GetRequest, opts ...grpc.CallOption) {
			assert.Equal(t, expectedGetRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := "testKey"

	record, err := client.Get(ctx, namespace, set, key, nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, expectedRecord, record)
}

func TestGet_FailsGettingConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(nil, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := "testKey"

	_, err = client.Get(ctx, namespace, set, key, nil, nil)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError(failedToGetRecord, fmt.Errorf("foo")))
}

func TestGet_FailsConvertingKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	mockConn := &connection{
		transactClient: mockTransactClient,
	}

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := struct{}{}

	_, err = client.Get(ctx, namespace, set, key, nil, nil)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError(failedToGetRecord, fmt.Errorf("unsupported key type: struct {}")))
}

func TestGet_FailsGetRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	mockConn := &connection{
		transactClient: mockTransactClient,
	}

	mockTransactClient.
		EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("foo"))

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := "key"

	_, err = client.Get(ctx, namespace, set, key, nil, nil)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError(failedToGetRecord, fmt.Errorf("foo")))
}

func TestDelete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	mockConn := &connection{
		transactClient: mockTransactClient,
	}

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	// Prepare the expected PutRequest
	expectedDeleteRequest := &protos.DeleteRequest{
		Key: &protos.Key{
			Namespace: "testNamespace",
			Set:       nil,
			Value: &protos.Key_StringValue{
				StringValue: "testKey",
			},
		},
	}

	mockTransactClient.
		EXPECT().
		Delete(gomock.Any(), gomock.Any()).
		Return(&emptypb.Empty{}, nil).
		Do(func(ctx context.Context, in *protos.DeleteRequest, opts ...grpc.CallOption) {
			assert.Equal(t, expectedDeleteRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := "testKey"

	err = client.Delete(ctx, namespace, set, key)

	assert.NoError(t, err)
}

func TestDelete_FailsGettingConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(nil, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := "testKey"

	err = client.Delete(ctx, namespace, set, key)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError(failedToDeleteRecord, fmt.Errorf("foo")))
}

func TestDelete_FailsConvertingKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	mockConn := &connection{
		transactClient: mockTransactClient,
	}

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := struct{}{}

	err = client.Delete(ctx, namespace, set, key)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError(failedToDeleteRecord, fmt.Errorf("unsupported key type: struct {}")))
}

func TestDelete_FailsDeleteRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	mockConn := &connection{
		transactClient: mockTransactClient,
	}

	mockTransactClient.
		EXPECT().
		Delete(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("foo"))

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := "key"

	err = client.Delete(ctx, namespace, set, key)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError(failedToDeleteRecord, fmt.Errorf("foo")))
}

func TestExists_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	mockConn := &connection{
		transactClient: mockTransactClient,
	}

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	// Prepare the expected PutRequest
	expectedExistsRequest := &protos.ExistsRequest{
		Key: &protos.Key{
			Namespace: "testNamespace",
			Set:       nil,
			Value: &protos.Key_StringValue{
				StringValue: "testKey",
			},
		},
	}

	mockTransactClient.
		EXPECT().
		Exists(gomock.Any(), gomock.Any()).
		Return(&protos.Boolean{
			Value: true,
		}, nil).
		Do(func(ctx context.Context, in *protos.ExistsRequest, opts ...grpc.CallOption) {
			assert.Equal(t, expectedExistsRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := "testKey"

	exists, err := client.Exists(ctx, namespace, set, key)

	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestExists_FailsGettingConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(nil, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := "testKey"

	_, err = client.Exists(ctx, namespace, set, key)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError(failedToCheckRecordExists, fmt.Errorf("foo")))
}

func TestExists_FailsConvertingKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	mockConn := &connection{
		transactClient: mockTransactClient,
	}

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := struct{}{}

	_, err = client.Exists(ctx, namespace, set, key)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError(failedToCheckRecordExists, fmt.Errorf("unsupported key type: struct {}")))
}

func TestExists_FailsDeleteRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	mockConn := &connection{
		transactClient: mockTransactClient,
	}

	mockTransactClient.
		EXPECT().
		Exists(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("foo"))

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := "key"

	_, err = client.Exists(ctx, namespace, set, key)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError(failedToCheckRecordExists, fmt.Errorf("foo")))
}

func TestIsIndexed_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	mockConn := &connection{
		transactClient: mockTransactClient,
	}

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	// Prepare the expected PutRequest
	expectedIsIndexedRequest := &protos.IsIndexedRequest{
		Key: &protos.Key{
			Namespace: "testNamespace",
			Set:       nil,
			Value: &protos.Key_StringValue{
				StringValue: "testKey",
			},
		},
		IndexId: &protos.IndexId{
			Namespace: "testNamespace",
			Name:      "testIndex",
		},
	}

	mockTransactClient.
		EXPECT().
		IsIndexed(gomock.Any(), gomock.Any()).
		Return(&protos.Boolean{
			Value: true,
		}, nil).
		Do(func(ctx context.Context, in *protos.IsIndexedRequest, opts ...grpc.CallOption) {
			assert.Equal(t, expectedIsIndexedRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := "testKey"
	indexName := "testIndex"

	exists, err := client.IsIndexed(ctx, namespace, set, indexName, key)

	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestIsIndexed_FailsGettingConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(nil, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := "testKey"
	indexName := "testIndex"

	_, err = client.IsIndexed(ctx, namespace, set, indexName, key)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError(failedToCheckIsIndexed, fmt.Errorf("foo")))
}

func TestIsIndexed_FailsConvertingKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	mockConn := &connection{
		transactClient: mockTransactClient,
	}

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := struct{}{}
	indexName := "testIndex"

	_, err = client.IsIndexed(ctx, namespace, set, indexName, key)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError(failedToCheckIsIndexed, fmt.Errorf("unsupported key type: struct {}")))
}

func TestIsIndexed_FailsDeleteRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	mockConn := &connection{
		transactClient: mockTransactClient,
	}

	mockTransactClient.
		EXPECT().
		IsIndexed(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("foo"))

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	var set *string = nil
	key := "key"
	indexName := "testIndex"

	_, err = client.IsIndexed(ctx, namespace, set, indexName, key)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError(failedToCheckIsIndexed, fmt.Errorf("foo")))
}

func TestVectorSearchFloat32_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	mockConn := &connection{
		transactClient: mockTransactClient,
	}
	mockVectorSearchClient := protos.NewMockTransactService_VectorSearchClient(ctrl)

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	// Prepare the expected PutRequest
	expectedVectorSearchFloat32Request := &protos.VectorSearchRequest{
		Index: &protos.IndexId{
			Namespace: "testNamespace",
			Name:      "testIndex",
		},
		QueryVector: &protos.Vector{
			Data: &protos.Vector_FloatData{
				FloatData: &protos.FloatData{
					Value: []float32{1.0, 2.0, 3.0},
				},
			},
		},
		Limit: 7,
		SearchParams: &protos.VectorSearchRequest_HnswSearchParams{
			HnswSearchParams: &protos.HnswSearchParams{
				Ef: GetUint32Ptr(8),
			},
		},
		Projection: &protos.ProjectionSpec{
			Include: &protos.ProjectionFilter{
				Type: protos.ProjectionType_ALL,
			},
			Exclude: &protos.ProjectionFilter{
				Type: protos.ProjectionType_NONE,
			},
		},
	}

	mockTransactClient.
		EXPECT().
		VectorSearch(gomock.Any(), gomock.Any()).
		Return(mockVectorSearchClient, nil).
		Do(func(ctx context.Context, in *protos.VectorSearchRequest, opts ...grpc.CallOption) {
			assert.Equal(t, expectedVectorSearchFloat32Request, in)
		})

	vectorCounter := 0
	mockVectorSearchClient.
		EXPECT().
		Recv().
		AnyTimes().
		DoAndReturn(
			func() (*protos.Neighbor, error) {
				vectorCounter++

				if vectorCounter == 4 {
					return nil, io.EOF
				}

				return &protos.Neighbor{
					Key: &protos.Key{
						Namespace: "testNamespace",
						Set:       GetStrPtr("testSet"),
						Value: &protos.Key_StringValue{
							StringValue: fmt.Sprintf("key-%d", vectorCounter),
						},
					},
					Record: &protos.Record{
						Fields: []*protos.Field{
							{
								Name: "field1",
								Value: &protos.Value{
									Value: &protos.Value_StringValue{
										StringValue: "value1",
									},
								},
							},
						},
						Metadata: &protos.Record_AerospikeMetadata{
							AerospikeMetadata: &protos.AerospikeRecordMetadata{
								Generation: uint32(vectorCounter),
								Expiration: uint32(vectorCounter),
							},
						},
					},
					Distance: float32(vectorCounter),
				}, nil
			},
		)

	expectedNeighbors := []*Neighbor{
		{
			Record: &Record{
				Data: map[string]any{
					"field1": "value1",
				},
				Generation: uint32(1),
				Expiration: GetTimePtr(AerospikeEpoch.Add(time.Second * 1)),
			},
			Set:       GetStrPtr("testSet"),
			Key:       "key-1",
			Namespace: "testNamespace",
			Distance:  float32(1),
		},
		{
			Record: &Record{
				Data: map[string]any{
					"field1": "value1",
				},
				Generation: uint32(2),
				Expiration: GetTimePtr(AerospikeEpoch.Add(time.Second * 2)),
			},
			Set:       GetStrPtr("testSet"),
			Key:       "key-2",
			Namespace: "testNamespace",
			Distance:  float32(2),
		},
		{
			Record: &Record{
				Data: map[string]any{
					"field1": "value1",
				},
				Generation: uint32(3),
				Expiration: GetTimePtr(AerospikeEpoch.Add(time.Second * 3)),
			},
			Set:       GetStrPtr("testSet"),
			Key:       "key-3",
			Namespace: "testNamespace",
			Distance:  float32(3),
		},
	}

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	indexName := "testIndex"
	vector := []float32{1.0, 2.0, 3.0}
	limit := uint32(7)
	searchParams := &protos.HnswSearchParams{
		Ef: GetUint32Ptr(8),
	}

	neighbors, err := client.VectorSearchFloat32(ctx, namespace, indexName, vector, uint32(limit), searchParams, nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, len(expectedNeighbors), len(neighbors))
	for i, _ := range neighbors {
		assert.EqualExportedValues(t, expectedNeighbors[i], neighbors[i])
	}
}

func TestVectorSearchFloat32_FailsGettingConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(nil, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	indexName := "testIndex"
	vector := []float32{1.0, 2.0, 3.0}
	limit := uint32(7)
	searchParams := &protos.HnswSearchParams{
		Ef: GetUint32Ptr(8),
	}

	_, err = client.VectorSearchFloat32(ctx, namespace, indexName, vector, uint32(limit), searchParams, nil, nil)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to search for vector", fmt.Errorf("foo")))
}

func TestVectorSearchFloat32_FailsVectorSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	mockConn := &connection{
		transactClient: mockTransactClient,
	}

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	mockTransactClient.
		EXPECT().
		VectorSearch(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	indexName := "testIndex"
	vector := []float32{1.0, 2.0, 3.0}
	limit := uint32(7)
	searchParams := &protos.HnswSearchParams{
		Ef: GetUint32Ptr(8),
	}

	_, err = client.VectorSearchFloat32(ctx, namespace, indexName, vector, uint32(limit), searchParams, nil, nil)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to search for vector", fmt.Errorf("foo")))
}

func TestVectorSearchFloat32_FailedToRecvAllNeighbors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	mockConn := &connection{
		transactClient: mockTransactClient,
	}
	mockVectorSearchClient := protos.NewMockTransactService_VectorSearchClient(ctrl)

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	mockTransactClient.
		EXPECT().
		VectorSearch(gomock.Any(), gomock.Any()).
		Return(mockVectorSearchClient, nil)

	mockVectorSearchClient.
		EXPECT().
		Recv().
		AnyTimes().
		DoAndReturn(
			func() (*protos.Neighbor, error) {
				return nil, fmt.Errorf("foo")
			},
		)

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	indexName := "testIndex"
	vector := []float32{1.0, 2.0, 3.0}
	limit := uint32(7)
	searchParams := &protos.HnswSearchParams{
		Ef: GetUint32Ptr(8),
	}

	_, err = client.VectorSearchFloat32(ctx, namespace, indexName, vector, uint32(limit), searchParams, nil, nil)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to receive all neighbors", fmt.Errorf("foo")))
}

type unknowKeyValue struct{}

func (u *unknowKeyValue) isKey_Value() {}

func TestVectorSearchFloat32_FailedToConvertNeighbor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockTransactClient := protos.NewMockTransactServiceClient(ctrl)
	mockConn := &connection{
		transactClient: mockTransactClient,
	}
	mockVectorSearchClient := protos.NewMockTransactService_VectorSearchClient(ctrl)

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	mockTransactClient.
		EXPECT().
		VectorSearch(gomock.Any(), gomock.Any()).
		Return(mockVectorSearchClient, nil)

	mockVectorSearchClient.
		EXPECT().
		Recv().
		AnyTimes().
		DoAndReturn(
			func() (*protos.Neighbor, error) {
				return &protos.Neighbor{
					Key: &protos.Key{Value: protos.NewMockisKey_Value(ctrl)},
				}, nil
			},
		)

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	indexName := "testIndex"
	vector := []float32{1.0, 2.0, 3.0}
	limit := uint32(7)
	searchParams := &protos.HnswSearchParams{
		Ef: GetUint32Ptr(8),
	}

	_, err = client.VectorSearchFloat32(ctx, namespace, indexName, vector, uint32(limit), searchParams, nil, nil)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to convert neighbor", fmt.Errorf("error converting neighbor: unsupported key value type: *protos.MockisKey_Value")))
}

func TestWaitForIndexCompletion_SuccessAfterZeroCountReturnedTwice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	expectedGetStatusRequest := &protos.IndexStatusRequest{
		IndexId: &protos.IndexId{
			Namespace: "testNamespace",
			Name:      "testIndex",
		},
	}

	indexStatus := &protos.IndexStatusResponse{
		UnmergedRecordCount: 0,
	}

	mockIndexClient.
		EXPECT().
		GetStatus(gomock.Any(), gomock.Any()).
		Times(3).
		Return(indexStatus, nil).
		Do(func(ctx context.Context, in *protos.IndexStatusRequest, opts ...grpc.CallOption) {
			assert.Equal(t, expectedGetStatusRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	indexName := "testIndex"

	err = client.WaitForIndexCompletion(ctx, namespace, indexName, time.Millisecond*1)

	assert.NoError(t, err)
}

func TestWaitForIndexCompletion_SuccessAfterNonZeroUnmergedCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	expectedGetStatusRequest := &protos.IndexStatusRequest{
		IndexId: &protos.IndexId{
			Namespace: "testNamespace",
			Name:      "testIndex",
		},
	}

	indexStatusCount := 0
	mockIndexClient.
		EXPECT().
		GetStatus(gomock.Any(), gomock.Any()).
		Times(2).
		DoAndReturn(func(ctx context.Context, in *protos.IndexStatusRequest, opts ...grpc.CallOption) (*protos.IndexStatusResponse, error) {
			assert.Equal(t, expectedGetStatusRequest, in)
			indexStatusCount++

			if indexStatusCount == 1 {
				return &protos.IndexStatusResponse{
					UnmergedRecordCount: 1,
				}, nil
			} else {
				return &protos.IndexStatusResponse{
					UnmergedRecordCount: 0,
				}, nil
			}

		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	indexName := "testIndex"

	err = client.WaitForIndexCompletion(ctx, namespace, indexName, time.Millisecond*1)

	assert.NoError(t, err)
}

func TestWaitForIndexCompletion_FailToGetRandomConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	indexName := "testIndex"

	err = client.WaitForIndexCompletion(ctx, namespace, indexName, time.Millisecond*1)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to wait for index completion", fmt.Errorf("foo")))
}

func TestWaitForIndexCompletion_FailGetStatusCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	mockIndexClient.
		EXPECT().
		GetStatus(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	namespace := "testNamespace"
	indexName := "testIndex"

	err = client.WaitForIndexCompletion(ctx, namespace, indexName, time.Millisecond*1)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to wait for index completion", fmt.Errorf("foo")))
}

func TestWaitForIndexCompletion_FailTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	// Set up expectations for connProvider.GetRandomConn()
	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	mockIndexClient.
		EXPECT().
		GetStatus(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	namespace := "testNamespace"
	indexName := "testIndex"

	err = client.WaitForIndexCompletion(ctx, namespace, indexName, time.Second*1)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to wait for index completion", fmt.Errorf("context deadline exceeded")))
}

func TestIndexCreateFromIndexDef_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Times(2).
		Return(mockConn, nil)

	indexDef := &protos.IndexDefinition{
		Id: &protos.IndexId{
			Namespace: "testNamespace",
			Name:      "testIndex",
		},
	}
	expectedIndexCreateRequest := &protos.IndexCreateRequest{
		Definition: indexDef,
	}

	mockIndexClient.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Do(func(ctx context.Context, in *protos.IndexCreateRequest, opts ...grpc.CallOption) {
			assert.Equal(t, expectedIndexCreateRequest, in)
		})

	mockIndexClient.
		EXPECT().
		GetStatus(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()

	err = client.IndexCreateFromIndexDef(ctx, indexDef)

	assert.NoError(t, err)
}

func TestIndexCreateFromIndexDef_FailGetConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, fmt.Errorf("foo"))

	indexDef := &protos.IndexDefinition{
		Id: &protos.IndexId{
			Namespace: "testNamespace",
			Name:      "testIndex",
		},
	}

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()

	err = client.IndexCreateFromIndexDef(ctx, indexDef)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to create index from definition", fmt.Errorf("foo")))
}

func TestIndexCreateFromIndexDef_FailCreateCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	mockIndexClient.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("foo"))

	indexDef := &protos.IndexDefinition{
		Id: &protos.IndexId{
			Namespace: "testNamespace",
			Name:      "testIndex",
		},
	}

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()

	err = client.IndexCreateFromIndexDef(ctx, indexDef)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to create index from definition", fmt.Errorf("foo")))
}

func TestIndexUpdate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	expectedIndexUpdateRequest := &protos.IndexUpdateRequest{
		IndexId: &protos.IndexId{
			Namespace: "testNamespace",
			Name:      "testIndex",
		},
		Labels: map[string]string{
			"foo": "bar",
		},
		Update: &protos.IndexUpdateRequest_HnswIndexUpdate{
			HnswIndexUpdate: &protos.HnswIndexUpdate{
				MaxMemQueueSize: GetUint32Ptr(10),
			},
		},
	}

	mockIndexClient.
		EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Do(func(ctx context.Context, in *protos.IndexUpdateRequest, opts ...grpc.CallOption) {
			assert.Equal(t, expectedIndexUpdateRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testNamespace := "testNamespace"
	testIndex := "testIndex"
	testMetadata := map[string]string{
		"foo": "bar",
	}
	hnswParams := &protos.HnswIndexUpdate{
		MaxMemQueueSize: GetUint32Ptr(10),
	}

	err = client.IndexUpdate(ctx, testNamespace, testIndex, testMetadata, hnswParams)

	assert.NoError(t, err)
}

func TestIndexUpdate_FailGetConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testNamespace := "testNamespace"
	testIndex := "testIndex"
	testMetadata := map[string]string{
		"foo": "bar",
	}
	hnswParams := &protos.HnswIndexUpdate{
		MaxMemQueueSize: GetUint32Ptr(10),
	}

	err = client.IndexUpdate(ctx, testNamespace, testIndex, testMetadata, hnswParams)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to update index", fmt.Errorf("foo")))
}

func TestIndexUpdate_FailUpdateCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	mockIndexClient.
		EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("bar"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testNamespace := "testNamespace"
	testIndex := "testIndex"
	testMetadata := map[string]string{
		"foo": "bar",
	}
	hnswParams := &protos.HnswIndexUpdate{
		MaxMemQueueSize: GetUint32Ptr(10),
	}

	err = client.IndexUpdate(ctx, testNamespace, testIndex, testMetadata, hnswParams)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to update index", fmt.Errorf("bar")))
}

func TestIndexDrop_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Times(2).
		Return(mockConn, nil)

	expectedIndexDropRequest := &protos.IndexDropRequest{
		IndexId: &protos.IndexId{
			Namespace: "testNamespace",
			Name:      "testIndex",
		},
	}

	mockIndexClient.
		EXPECT().
		Drop(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Do(func(ctx context.Context, in *protos.IndexDropRequest, opts ...grpc.CallOption) {
			assert.Equal(t, expectedIndexDropRequest, in)
		})

	mockIndexClient.
		EXPECT().
		GetStatus(gomock.Any(), gomock.Any()).
		Return(nil, status.Errorf(codes.NotFound, "foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testNamespace := "testNamespace"
	testIndex := "testIndex"

	err = client.IndexDrop(ctx, testNamespace, testIndex)

	assert.NoError(t, err)
}

func TestIndexDrop_FailGetConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testNamespace := "testNamespace"
	testIndex := "testIndex"

	err = client.IndexDrop(ctx, testNamespace, testIndex)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to drop index", fmt.Errorf("foo")))
}

func TestIndexDrop_FailDropCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	mockIndexClient.
		EXPECT().
		Drop(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("bar"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testNamespace := "testNamespace"
	testIndex := "testIndex"

	err = client.IndexDrop(ctx, testNamespace, testIndex)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to drop index", fmt.Errorf("bar")))
}

func TestIndexList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	expectedIndexListRequest := &protos.IndexListRequest{
		ApplyDefaults: GetBoolPtr(true),
	}

	expectedIndexDefs := &protos.IndexDefinitionList{
		Indices: []*protos.IndexDefinition{
			{
				Id: &protos.IndexId{
					Namespace: "testNamespace0",
					Name:      "testIndex0",
				},
			},
			{
				Id: &protos.IndexId{
					Namespace: "testNamespace1",
					Name:      "testIndex1",
				},
			},
		},
	}

	mockIndexClient.
		EXPECT().
		List(gomock.Any(), gomock.Any()).
		Return(expectedIndexDefs, nil).
		Do(func(ctx context.Context, in *protos.IndexListRequest, opts ...grpc.CallOption) {
			assert.Equal(t, expectedIndexListRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()

	indexDefs, err := client.IndexList(ctx, true)

	assert.NoError(t, err)
	assert.EqualExportedValues(t, expectedIndexDefs, indexDefs)
}

func TestIndexList_FailGetConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	_, err = client.IndexList(ctx, true)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to get indexes", fmt.Errorf("foo")))
}

func TestIndexList_FailDropCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	mockIndexClient.
		EXPECT().
		List(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("bar"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()

	_, err = client.IndexList(ctx, true)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to get indexes", fmt.Errorf("bar")))
}

func TestIndexGet_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	expectedIndexListRequest := &protos.IndexGetRequest{
		IndexId: &protos.IndexId{
			Namespace: "testNamespace",
			Name:      "testIndex",
		},
		ApplyDefaults: GetBoolPtr(true),
	}

	expectedIndexDefs := &protos.IndexDefinition{
		Id: &protos.IndexId{
			Namespace: "testNamespace",
			Name:      "testIndex",
		},
	}

	mockIndexClient.
		EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(expectedIndexDefs, nil).
		Do(func(ctx context.Context, in *protos.IndexGetRequest, opts ...grpc.CallOption) {
			assert.Equal(t, expectedIndexListRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testNamespace := "testNamespace"
	testIndex := "testIndex"

	indexDefs, err := client.IndexGet(ctx, testNamespace, testIndex, true)

	assert.NoError(t, err)
	assert.EqualExportedValues(t, expectedIndexDefs, indexDefs)
}

func TestIndexGet_FailGetConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testNamespace := "testNamespace"
	testIndex := "testIndex"

	_, err = client.IndexGet(ctx, testNamespace, testIndex, true)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to get index", fmt.Errorf("foo")))
}

func TestIndexGet_FailDropCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	mockIndexClient.
		EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("bar"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testNamespace := "testNamespace"
	testIndex := "testIndex"

	_, err = client.IndexGet(ctx, testNamespace, testIndex, true)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to get index", fmt.Errorf("bar")))
}

func TestIndexGetStatus_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	expectedIndexStatusRequest := &protos.IndexStatusRequest{
		IndexId: &protos.IndexId{
			Namespace: "testNamespace",
			Name:      "testIndex",
		},
	}

	expectedIndexStatusResp := &protos.IndexStatusResponse{
		UnmergedRecordCount: 9,
	}

	mockIndexClient.
		EXPECT().
		GetStatus(gomock.Any(), gomock.Any()).
		Return(expectedIndexStatusResp, nil).
		Do(func(ctx context.Context, in *protos.IndexStatusRequest, opts ...grpc.CallOption) {
			assert.Equal(t, expectedIndexStatusRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testNamespace := "testNamespace"
	testIndex := "testIndex"

	indexDefs, err := client.IndexGetStatus(ctx, testNamespace, testIndex)

	assert.NoError(t, err)
	assert.EqualExportedValues(t, expectedIndexStatusResp, indexDefs)
}

func TestIndexGetStatus_FailGetConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testNamespace := "testNamespace"
	testIndex := "testIndex"

	_, err = client.IndexGetStatus(ctx, testNamespace, testIndex)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to get index status", fmt.Errorf("foo")))
}

func TestIndexGetStatus_FailDropCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	mockIndexClient.
		EXPECT().
		GetStatus(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("bar"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testNamespace := "testNamespace"
	testIndex := "testIndex"

	_, err = client.IndexGetStatus(ctx, testNamespace, testIndex)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to get index status", fmt.Errorf("bar")))
}

func TestIndexGcInvalidVertices_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	cutoffTime := time.Now()
	expectedIndexStatusRequest := &protos.GcInvalidVerticesRequest{
		IndexId: &protos.IndexId{
			Namespace: "testNamespace",
			Name:      "testIndex",
		},
		CutoffTimestamp: cutoffTime.Unix(),
	}

	mockIndexClient.
		EXPECT().
		GcInvalidVertices(gomock.Any(), gomock.Any()).
		Return(&emptypb.Empty{}, nil).
		Do(func(ctx context.Context, in *protos.GcInvalidVerticesRequest, opts ...grpc.CallOption) {
			assert.Equal(t, expectedIndexStatusRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testNamespace := "testNamespace"
	testIndex := "testIndex"

	err = client.GcInvalidVertices(ctx, testNamespace, testIndex, cutoffTime)

	assert.NoError(t, err)
}

func TestIndexGcInvalidVertices_FailGetConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testNamespace := "testNamespace"
	testIndex := "testIndex"
	cutoffTime := time.Now()

	err = client.GcInvalidVertices(ctx, testNamespace, testIndex, cutoffTime)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to garbage collect invalid vertices", fmt.Errorf("foo")))
}

func TestIndexGcInvalidVertices_FailDropCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockIndexClient := protos.NewMockIndexServiceClient(ctrl)
	mockConn := &connection{
		indexClient: mockIndexClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	mockIndexClient.
		EXPECT().
		GcInvalidVertices(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("bar"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testNamespace := "testNamespace"
	testIndex := "testIndex"
	cutoffTime := time.Now()

	err = client.GcInvalidVertices(ctx, testNamespace, testIndex, cutoffTime)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to garbage collect invalid vertices", fmt.Errorf("bar")))
}

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockUserAdminClient := protos.NewMockUserAdminServiceClient(ctrl)
	mockConn := &connection{
		userAdminClient: mockUserAdminClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	expectedRequest := &protos.AddUserRequest{
		Credentials: &protos.Credentials{
			Username: "testUser",
			Credentials: &protos.Credentials_PasswordCredentials{
				PasswordCredentials: &protos.PasswordCredentials{
					Password: "testPass",
				},
			},
		},
		Roles: []string{
			"testRole",
		},
	}

	mockUserAdminClient.
		EXPECT().
		AddUser(gomock.Any(), gomock.Any()).
		Return(&emptypb.Empty{}, nil).
		Do(func(ctx context.Context, in *protos.AddUserRequest, opts ...grpc.CallOption) {
			assert.Equal(t, expectedRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testUser := "testUser"
	testPass := "testPass"
	testRoles := []string{"testRole"}

	err = client.CreateUser(ctx, testUser, testPass, testRoles)

	assert.NoError(t, err)
}

func TestCreateUser_FailGetConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockUserAdminClient := protos.NewMockUserAdminServiceClient(ctrl)
	mockConn := &connection{
		userAdminClient: mockUserAdminClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testUser := "testUser"
	testPass := "testPass"
	testRoles := []string{"testRole"}

	err = client.CreateUser(ctx, testUser, testPass, testRoles)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to create user", fmt.Errorf("foo")))
}

func TestCreateUser_FailCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockUserClient := protos.NewMockUserAdminServiceClient(ctrl)
	mockConn := &connection{
		userAdminClient: mockUserClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	mockUserClient.
		EXPECT().
		AddUser(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("bar"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testUser := "testUser"
	testPass := "testPass"
	testRoles := []string{"testRole"}

	err = client.CreateUser(ctx, testUser, testPass, testRoles)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to create user", fmt.Errorf("bar")))
}

func TestUpdateCredentials_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockUserAdminClient := protos.NewMockUserAdminServiceClient(ctrl)
	mockConn := &connection{
		userAdminClient: mockUserAdminClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	expectedRequest := &protos.UpdateCredentialsRequest{
		Credentials: &protos.Credentials{
			Username: "testUser",
			Credentials: &protos.Credentials_PasswordCredentials{
				PasswordCredentials: &protos.PasswordCredentials{
					Password: "testPass",
				},
			},
		},
	}

	mockUserAdminClient.
		EXPECT().
		UpdateCredentials(gomock.Any(), gomock.Any()).
		Return(&emptypb.Empty{}, nil).
		Do(func(ctx context.Context, in *protos.UpdateCredentialsRequest, opts ...grpc.CallOption) {
			assert.Equal(t, expectedRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testUser := "testUser"
	testPass := "testPass"

	err = client.UpdateCredentials(ctx, testUser, testPass)

	assert.NoError(t, err)
}

func TestUpdateCredentials_FailGetConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockUserAdminClient := protos.NewMockUserAdminServiceClient(ctrl)
	mockConn := &connection{
		userAdminClient: mockUserAdminClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testUser := "testUser"
	testPass := "testPass"

	err = client.UpdateCredentials(ctx, testUser, testPass)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to update user credentials", fmt.Errorf("foo")))
}

func TestUpdateCredentials_FailCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockUserClient := protos.NewMockUserAdminServiceClient(ctrl)
	mockConn := &connection{
		userAdminClient: mockUserClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	mockUserClient.
		EXPECT().
		UpdateCredentials(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("bar"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testUser := "testUser"
	testPass := "testPass"

	err = client.UpdateCredentials(ctx, testUser, testPass)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to update user credentials", fmt.Errorf("bar")))
}

func TestDropUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockUserAdminClient := protos.NewMockUserAdminServiceClient(ctrl)
	mockConn := &connection{
		userAdminClient: mockUserAdminClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	expectedRequest := &protos.DropUserRequest{
		Username: "testUser",
	}

	mockUserAdminClient.
		EXPECT().
		DropUser(gomock.Any(), gomock.Any()).
		Return(&emptypb.Empty{}, nil).
		Do(func(ctx context.Context, in *protos.DropUserRequest, opts ...grpc.CallOption) {
			assert.Equal(t, expectedRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testUser := "testUser"

	err = client.DropUser(ctx, testUser)

	assert.NoError(t, err)
}

func TestDropUser_FailGetConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockUserAdminClient := protos.NewMockUserAdminServiceClient(ctrl)
	mockConn := &connection{
		userAdminClient: mockUserAdminClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testUser := "testUser"

	err = client.DropUser(ctx, testUser)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to drop user", fmt.Errorf("foo")))
}

func TestDropUser_FailCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockUserClient := protos.NewMockUserAdminServiceClient(ctrl)
	mockConn := &connection{
		userAdminClient: mockUserClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	mockUserClient.
		EXPECT().
		DropUser(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("bar"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testUser := "testUser"

	err = client.DropUser(ctx, testUser)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to drop user", fmt.Errorf("bar")))
}

func TestGetUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockUserAdminClient := protos.NewMockUserAdminServiceClient(ctrl)
	mockConn := &connection{
		userAdminClient: mockUserAdminClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	expectedRequest := &protos.GetUserRequest{
		Username: "testUser",
	}

	expectedUser := &protos.User{
		Username: "testUser",
		Roles: []string{
			"testRole",
		},
	}

	mockUserAdminClient.
		EXPECT().
		GetUser(gomock.Any(), gomock.Any()).
		Return(expectedUser, nil).
		Do(func(ctx context.Context, in *protos.GetUserRequest, opts ...grpc.CallOption) {
			assert.Equal(t, expectedRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testUser := "testUser"

	user, err := client.GetUser(ctx, testUser)

	assert.NoError(t, err)
	assert.EqualExportedValues(t, expectedUser, user)
}

func TestGetUser_FailGetConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockUserAdminClient := protos.NewMockUserAdminServiceClient(ctrl)
	mockConn := &connection{
		userAdminClient: mockUserAdminClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testUser := "testUser"

	_, err = client.GetUser(ctx, testUser)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to get user", fmt.Errorf("foo")))
}

func TestGetUser_FailCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockUserClient := protos.NewMockUserAdminServiceClient(ctrl)
	mockConn := &connection{
		userAdminClient: mockUserClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	mockUserClient.
		EXPECT().
		GetUser(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("bar"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testUser := "testUser"

	_, err = client.GetUser(ctx, testUser)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to get user", fmt.Errorf("bar")))
}

func TestListUsers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockUserAdminClient := protos.NewMockUserAdminServiceClient(ctrl)
	mockConn := &connection{
		userAdminClient: mockUserAdminClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	expectedRequest := &emptypb.Empty{}

	expectedUsers := &protos.ListUsersResponse{
		Users: []*protos.User{
			{

				Username: "testUser",
				Roles: []string{
					"testRole",
				},
			},
		},
	}

	mockUserAdminClient.
		EXPECT().
		ListUsers(gomock.Any(), gomock.Any()).
		Return(expectedUsers, nil).
		Do(func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) {
			assert.Equal(t, expectedRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()

	user, err := client.ListUsers(ctx)

	assert.NoError(t, err)
	assert.EqualExportedValues(t, expectedUsers, user)
}

func TestListUsers_FailGetConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockUserAdminClient := protos.NewMockUserAdminServiceClient(ctrl)
	mockConn := &connection{
		userAdminClient: mockUserAdminClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()

	_, err = client.ListUsers(ctx)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to list users", fmt.Errorf("foo")))
}

func TestListUsers_FailCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockUserClient := protos.NewMockUserAdminServiceClient(ctrl)
	mockConn := &connection{
		userAdminClient: mockUserClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	mockUserClient.
		EXPECT().
		ListUsers(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("bar"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()

	_, err = client.ListUsers(ctx)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to list users", fmt.Errorf("bar")))
}

func TestRevokeRoles_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockUserAdminClient := protos.NewMockUserAdminServiceClient(ctrl)
	mockConn := &connection{
		userAdminClient: mockUserAdminClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	expectedRequest := &protos.RevokeRolesRequest{
		Username: "testUser",
		Roles:    []string{"testRole"},
	}

	mockUserAdminClient.
		EXPECT().
		RevokeRoles(gomock.Any(), gomock.Any()).
		Return(&emptypb.Empty{}, nil).
		Do(func(ctx context.Context, in *protos.RevokeRolesRequest, opts ...grpc.CallOption) {
			assert.Equal(t, expectedRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testUser := "testUser"
	testRoles := []string{"testRole"}

	err = client.RevokeRoles(ctx, testUser, testRoles)

	assert.NoError(t, err)
}

func TestRevokeRoles_FailGetConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockUserAdminClient := protos.NewMockUserAdminServiceClient(ctrl)
	mockConn := &connection{
		userAdminClient: mockUserAdminClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testUser := "testUser"
	testRoles := []string{"testRole"}

	err = client.RevokeRoles(ctx, testUser, testRoles)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to revoke user roles", fmt.Errorf("foo")))
}

func TestRevokeRoles_FailCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockUserClient := protos.NewMockUserAdminServiceClient(ctrl)
	mockConn := &connection{
		userAdminClient: mockUserClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	mockUserClient.
		EXPECT().
		RevokeRoles(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("bar"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	testUser := "testUser"
	testRoles := []string{"testRole"}

	err = client.RevokeRoles(ctx, testUser, testRoles)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to revoke user roles", fmt.Errorf("bar")))
}

func TestListRoles_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockUserAdminClient := protos.NewMockUserAdminServiceClient(ctrl)
	mockConn := &connection{
		userAdminClient: mockUserAdminClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	expectedRequest := &emptypb.Empty{}

	expectedRoles := &protos.ListRolesResponse{
		Roles: []*protos.Role{
			{
				Id: "testRole",
			},
		},
	}

	mockUserAdminClient.
		EXPECT().
		ListRoles(gomock.Any(), gomock.Any()).
		Return(expectedRoles, nil).
		Do(func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) {
			assert.Equal(t, expectedRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()

	user, err := client.ListRoles(ctx)

	assert.NoError(t, err)
	assert.EqualExportedValues(t, expectedRoles, user)
}

func TestListRoles_FailGetConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockUserAdminClient := protos.NewMockUserAdminServiceClient(ctrl)
	mockConn := &connection{
		userAdminClient: mockUserAdminClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()

	_, err = client.ListRoles(ctx)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to list roles", fmt.Errorf("foo")))
}

func TestListRoles_FailCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockUserClient := protos.NewMockUserAdminServiceClient(ctrl)
	mockConn := &connection{
		userAdminClient: mockUserClient,
	}

	mockConnProvider.
		EXPECT().
		GetRandomConn().
		Return(mockConn, nil)

	mockUserClient.
		EXPECT().
		ListRoles(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("bar"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()

	_, err = client.ListRoles(ctx)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to list roles", fmt.Errorf("bar")))
}

func TestConnectedNodeEndpoint_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockGrpcConn := NewMockgrpcClientConn(ctrl)
	mockConn := &connection{
		grpcConn: mockGrpcConn,
	}

	mockConnProvider.
		EXPECT().
		GetSeedConn().
		Return(mockConn, nil)

	expectedEndpoint := &protos.ServerEndpoint{
		Address: "1.1.1.1",
		Port:    3000,
	}

	mockGrpcConn.
		EXPECT().
		Target().
		Return("1.1.1.1:3000")

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()

	endpoint, err := client.ConnectedNodeEndpoint(ctx, nil)

	assert.NoError(t, err)
	assert.EqualExportedValues(t, expectedEndpoint, endpoint)
}

func TestConnectedNodeEndpoint_FailedGetConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockGrpcConn := NewMockgrpcClientConn(ctrl)
	mockConn := &connection{
		grpcConn: mockGrpcConn,
	}

	mockConnProvider.
		EXPECT().
		GetSeedConn().
		Return(mockConn, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()

	_, err = client.ConnectedNodeEndpoint(ctx, nil)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to get connected endpoint", fmt.Errorf("foo")))
}

func TestConnectedNodeEndpoint_FailParsePort(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockGrpcConn := NewMockgrpcClientConn(ctrl)
	mockConn := &connection{
		grpcConn: mockGrpcConn,
	}

	mockConnProvider.
		EXPECT().
		GetSeedConn().
		Return(mockConn, nil)

	mockGrpcConn.
		EXPECT().
		Target().
		Return("1.1.1.1:aaaa")

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()

	_, err = client.ConnectedNodeEndpoint(ctx, nil)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Contains(t, avsError.Error(), "failed to parse port")
}

func TestClusteringState_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockClusterInfoClient := protos.NewMockClusterInfoServiceClient(ctrl)
	mockConn := &connection{
		clusterInfoClient: mockClusterInfoClient,
	}

	mockConnProvider.
		EXPECT().
		GetSeedConn().
		Return(mockConn, nil)

	expectedRequest := &emptypb.Empty{}

	expectedRoles := &protos.ClusteringState{
		IsInCluster: true,
	}

	mockClusterInfoClient.
		EXPECT().
		GetClusteringState(gomock.Any(), gomock.Any()).
		Return(expectedRoles, nil).
		Do(func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) {
			assert.Equal(t, expectedRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()

	user, err := client.ClusteringState(ctx, nil)

	assert.NoError(t, err)
	assert.EqualExportedValues(t, expectedRoles, user)
}

func TestClusteringState_FailGetConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockClusterInfoClient := protos.NewMockClusterInfoServiceClient(ctrl)
	mockConn := &connection{
		clusterInfoClient: mockClusterInfoClient,
	}

	mockConnProvider.
		EXPECT().
		GetSeedConn().
		Return(mockConn, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()

	_, err = client.ClusteringState(ctx, nil)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to get clustering state", fmt.Errorf("foo")))
}

func TestClusteringState_FailCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockClusterInfoClient := protos.NewMockClusterInfoServiceClient(ctrl)
	mockConn := &connection{
		clusterInfoClient: mockClusterInfoClient,
	}

	mockConnProvider.
		EXPECT().
		GetSeedConn().
		Return(mockConn, nil)

	mockClusterInfoClient.
		EXPECT().
		GetClusteringState(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("bar"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()

	_, err = client.ClusteringState(ctx, nil)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to get clustering state", fmt.Errorf("bar")))
}

func TestClusterEndpoints_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockClusterInfoClient := protos.NewMockClusterInfoServiceClient(ctrl)
	mockConn := &connection{
		clusterInfoClient: mockClusterInfoClient,
	}

	mockConnProvider.
		EXPECT().
		GetSeedConn().
		Return(mockConn, nil)

	listenerName := "test-listener"
	expectedRequest := &protos.ClusterNodeEndpointsRequest{
		ListenerName: &listenerName,
	}

	expectedResp := &protos.ClusterNodeEndpoints{}

	mockClusterInfoClient.
		EXPECT().
		GetClusterEndpoints(gomock.Any(), gomock.Any()).
		Return(expectedResp, nil).
		Do(func(ctx context.Context, in *protos.ClusterNodeEndpointsRequest, opts ...grpc.CallOption) {
			assert.Equal(t, expectedRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()

	user, err := client.ClusterEndpoints(ctx, nil, &listenerName)

	assert.NoError(t, err)
	assert.EqualExportedValues(t, expectedResp, user)
}

func TestClusterEndpoints_FailGetConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockClusterInfoClient := protos.NewMockClusterInfoServiceClient(ctrl)
	mockConn := &connection{
		clusterInfoClient: mockClusterInfoClient,
	}

	mockConnProvider.
		EXPECT().
		GetSeedConn().
		Return(mockConn, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	listenerName := "test-listener"

	_, err = client.ClusterEndpoints(ctx, nil, &listenerName)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to get cluster endpoints", fmt.Errorf("foo")))
}

func TestClusterEndpoints_FailCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockClusterInfoClient := protos.NewMockClusterInfoServiceClient(ctrl)
	mockConn := &connection{
		clusterInfoClient: mockClusterInfoClient,
	}

	mockConnProvider.
		EXPECT().
		GetSeedConn().
		Return(mockConn, nil)

	mockClusterInfoClient.
		EXPECT().
		GetClusterEndpoints(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("bar"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	listenerName := "test-listener"
	_, err = client.ClusterEndpoints(ctx, nil, &listenerName)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to get cluster endpoints", fmt.Errorf("bar")))
}

func TestAbout_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockAboutClient := protos.NewMockAboutServiceClient(ctrl)
	mockConn := &connection{
		aboutClient: mockAboutClient,
	}

	mockConnProvider.
		EXPECT().
		GetSeedConn().
		Return(mockConn, nil)

	expectedRequest := &protos.AboutRequest{}
	expectedResp := &protos.AboutResponse{}

	mockAboutClient.
		EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(expectedResp, nil).
		Do(func(ctx context.Context, in *protos.AboutRequest, opts ...grpc.CallOption) {
			assert.Equal(t, expectedRequest, in)
		})

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()

	user, err := client.About(ctx, nil)

	assert.NoError(t, err)
	assert.EqualExportedValues(t, expectedResp, user)
}

func TestAbout_FailGetConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockAboutClient := protos.NewMockAboutServiceClient(ctrl)
	mockConn := &connection{
		aboutClient: mockAboutClient,
	}

	mockConnProvider.
		EXPECT().
		GetSeedConn().
		Return(mockConn, fmt.Errorf("foo"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()

	_, err = client.About(ctx, nil)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to make about request", fmt.Errorf("foo")))
}

func TestAbout_FailCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnProvider := NewMockconnProvider(ctrl)
	mockAboutClient := protos.NewMockAboutServiceClient(ctrl)
	mockConn := &connection{
		aboutClient: mockAboutClient,
	}

	mockConnProvider.
		EXPECT().
		GetSeedConn().
		Return(mockConn, nil)

	mockAboutClient.
		EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("bar"))

	// Create the client with the mock connProvider
	client, err := newClient(mockConnProvider, slog.Default())
	assert.NoError(t, err)

	// Prepare input parameters
	ctx := context.Background()
	_, err = client.About(ctx, nil)

	var avsError *Error
	assert.ErrorAs(t, err, &avsError)
	assert.Equal(t, avsError, NewAVSError("failed to make about request", fmt.Errorf("bar")))
}
