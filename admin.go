package avs

import (
	"context"
	"log/slog"
	"time"

	"github.com/aerospike/aerospike-proximus-client-go/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminClient struct {
	logger          *slog.Logger
	channelProvider *ChannelProvider
}

func NewAdminClient(
	ctx context.Context,
	seeds HostPortSlice,
	listenerName *string,
	isLoadBalancer bool,
	logger *slog.Logger,
) (*AdminClient, error) {
	logger = logger.WithGroup("avs.admin")
	logger.Debug("creating new client")

	channelProvider, err := NewChannelProvider(ctx, seeds, listenerName, isLoadBalancer, logger)
	if err != nil {
		logger.Error("failed to create channel provider", slog.Any("error", err))
		return nil, NewAVSErrorFromGrpc("failed to connect to server", err)
	}

	return &AdminClient{
		logger:          logger,
		channelProvider: channelProvider,
	}, nil
}

func (c *AdminClient) Close() {
	c.logger.Info("Closing client")
	c.channelProvider.Close()
}

func (c *AdminClient) IndexCreate(
	ctx context.Context,
	namespace string,
	sets []string,
	name string,
	vectorField string,
	dimensions uint32,
	vectorDistanceMetric protos.VectorDistanceMetric,
	indexParams *protos.HnswParams,
	indexMetaData map[string]string,
	indexStorage *protos.IndexStorage,
) error {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("name", name))

	conn, err := c.channelProvider.GetConn()
	if err != nil {
		logger.Error("failed to create index", slog.Any("error", err))
		return NewAVSErrorFromGrpc("failed to create index", err)
	}

	var set *string

	if len(sets) > 0 {
		set = &sets[0]
		logger.Warn(
			"multiple sets not yet supported for index creation, only the first set will be used",
			slog.String("set", *set),
		)
	}

	indexDef := &protos.IndexDefinition{
		Id: &protos.IndexId{
			Namespace: namespace,
			Name:      name,
		},
		Dimensions:           dimensions,
		VectorDistanceMetric: vectorDistanceMetric,
		Field:                vectorField,
		SetFilter:            set,
		Params:               &protos.IndexDefinition_HnswParams{HnswParams: indexParams},
		Labels:               indexMetaData,
		Storage:              indexStorage,
	}

	client := protos.NewIndexServiceClient(conn)

	_, err = client.Create(ctx, indexDef)
	if err != nil {
		logger.Error("failed to create index", slog.Any("error", err))
		return NewAVSErrorFromGrpc("failed to create index", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*100_000)
	defer cancel()

	return c.waitForIndexCreation(ctx, namespace, name, time.Millisecond*100)
}

func (c *AdminClient) IndexDrop(ctx context.Context, namespace, name string) error {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("name", name))

	conn, err := c.channelProvider.GetConn()
	if err != nil {
		msg := "failed to drop index"
		logger.Error(msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	indexID := &protos.IndexId{
		Namespace: namespace,
		Name:      name,
	}

	client := protos.NewIndexServiceClient(conn)

	_, err = client.Drop(ctx, indexID)
	if err != nil {
		msg := "failed to drop index"

		logger.Error(msg, slog.Any("error", err))

		return NewAVSErrorFromGrpc(msg, err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*100_000)
	defer cancel()

	return c.waitForIndexDrop(ctx, namespace, name, time.Microsecond*100)
}

func (c *AdminClient) IndexList(ctx context.Context) (*protos.IndexDefinitionList, error) {
	conn, err := c.channelProvider.GetConn()
	if err != nil {
		msg := "failed to get indexes"
		c.logger.Error(msg, slog.Any("error", err))
		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	client := protos.NewIndexServiceClient(conn)

	indexList, err := client.List(ctx, nil)
	if err != nil {
		msg := "failed to get indexes"
		c.logger.Error(msg, slog.Any("error", err))
		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return indexList, nil
}

func (c *AdminClient) IndexGet(ctx context.Context, namespace, name string) (*protos.IndexDefinition, error) {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("name", name))

	conn, err := c.channelProvider.GetConn()
	if err != nil {
		msg := "failed to get index"
		logger.Error(msg, slog.Any("error", err))
		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	indexID := &protos.IndexId{
		Namespace: namespace,
		Name:      name,
	}
	client := protos.NewIndexServiceClient(conn)

	indexDef, err := client.Get(ctx, indexID)
	if err != nil {
		msg := "failed to get index"
		logger.Error(msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return indexDef, nil
}

func (c *AdminClient) IndexGetStatus(ctx context.Context, namespace, name string) (*protos.IndexStatusResponse, error) {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("name", name))

	conn, err := c.channelProvider.GetConn()
	if err != nil {
		msg := "failed to get index status"
		logger.Error(msg, slog.Any("error", err))
		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	indexID := &protos.IndexId{
		Namespace: namespace,
		Name:      name,
	}
	client := protos.NewIndexServiceClient(conn)

	indexStatus, err := client.GetStatus(ctx, indexID)
	if err != nil {
		msg := "failed to get index status"
		logger.Error(msg, slog.Any("error", err))

		return nil, NewAVSErrorFromGrpc(msg, err)
	}

	return indexStatus, nil
}

func (c *AdminClient) waitForIndexCreation(ctx context.Context,
	namespace,
	name string,
	waitInterval time.Duration,
) error {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("name", name))

	conn, err := c.channelProvider.GetConn()
	if err != nil {
		msg := "failed to wait for index creation"
		logger.Error(msg, slog.Any("error", err))
		return NewAVSErrorFromGrpc(msg, err)
	}

	indexID := &protos.IndexId{
		Namespace: namespace,
		Name:      name,
	}

	client := protos.NewIndexServiceClient(conn)

	for {
		_, err := client.GetStatus(ctx, indexID)
		if err != nil {
			code := status.Code(err)
			if code == codes.Unavailable || code == codes.NotFound {
				logger.Debug("index does not exist, waiting...")
				select {
				case <-time.After(waitInterval):
				case <-ctx.Done():
					logger.ErrorContext(ctx, "waiting for index creation canceled")
					return ctx.Err()
				}
			} else {
				msg := "unable to wait for index creation, an unexpected error occurred"
				logger.Error(msg, slog.Any("error", err))
				return NewAVSErrorFromGrpc(msg, err)
			}
		} else {
			logger.Info("index has been created")
			break
		}
	}

	return nil
}

func (c *AdminClient) waitForIndexDrop(ctx context.Context, namespace, name string, waitInterval time.Duration) error {
	logger := c.logger.With(slog.String("namespace", namespace), slog.String("name", name))

	conn, err := c.channelProvider.GetConn()
	if err != nil {
		msg := "failed to wait for index deletion"
		logger.Error(msg, slog.Any("error", err))
		return NewAVSErrorFromGrpc(msg, err)
	}

	indexID := &protos.IndexId{
		Namespace: namespace,
		Name:      name,
	}

	client := protos.NewIndexServiceClient(conn)

	for {
		_, err := client.GetStatus(ctx, indexID)
		if err != nil {
			code := status.Code(err)
			if code == codes.Unavailable || code == codes.NotFound {
				logger.Info("index is deleted")
				return nil
			}

			msg := "unable to wait for index deletion, an unexpected error occurred"
			logger.Error(msg, slog.Any("error", err))

			return NewAVSErrorFromGrpc(msg, err)
		}

		c.logger.Debug("index still exists, waiting...")
		select {
		case <-time.After(waitInterval):
		case <-ctx.Done():
			logger.ErrorContext(ctx, "waiting for index deletion canceled")
			return ctx.Err()
		}
	}
}
