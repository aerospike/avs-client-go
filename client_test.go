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
