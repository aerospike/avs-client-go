package avs

import (
	"context"
	"fmt"
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
