[![codecov](https://codecov.io/gh/aerospike/avs-client-go/graph/badge.svg?token=811TWWPW6S)](https://codecov.io/gh/aerospike/avs-client-go)

# Aerospike Vector Search Go Client

> :warning: The go client is currently in development. APIs will break in the future!

## Overview

This repository contains the Go client for interacting with Aerospike Vector
Search (AVS) capabilities. It is designed to provide an easy-to-use interface for
performing vector searches with an Aerospike cluster.

## Current Functionality
- Connection to Aerospike clusters
- Create, List, Drop Indexes

## Installation
To install the Aerospike Vector Search Go Client, use the following command:

```bash
go get github.com/aerospike/avs-client-go
```

## Example Usage
Here's a quick example to get you started:
```go
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	avs "github.com/aerospike/avs-client-go"
	"github.com/aerospike/avs-client-go/protos"
)

func main() {
	connectCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	hostPort := NewHostPort("127.0.0.1", 5000)
	loadBalancer := true
	credentials := NewCredentialsFromUserPass("admin", "admin")
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	var (
		listenerName *string
		tlsConfig *tls.Config
	)

	client, err := avs.NewClient(
		connectCtx,
		HostPortSlice{hostPort},
		listenerName,
		loadBalancer,
		credentials
		tlsConfig,
		logger,
	)
	cancel()

	if err != nil {
		logger.Error("failed to create admin client", slog.Any("error", err))
		return
	}

	logger.Info("successfully connected to server")
	defer client.Close()

	namespace := "test"
	indexName := "bookIndex"
	vectorField := "vector"
	dimensions := uint32(10)
	distanceMetric := protos.VectorDistanceMetric_MANHATTAN
	indexOpts := &IndexCreateOpts{
		Set: []string{"testset"}
	}

	err = client.IndexCreate(
		context.Background(),
		namespace,
		indexName,
		vectorField,
		dimensions,
		distanceMetric,
		indexOpts,
	)
	if err != nil {
		logger.Error("failed to create index", slog.Any("error", err))
		return
	}

	logger.Info("successfully created index")

	indexes, err := client.IndexList(context.Background(), true)
	if err != nil {
		logger.Error("failed to list indexes", slog.Any("error", err))
		return
	}

	for _, index := range indexes.GetIndices() {
		fmt.Println(index.String())
	}

	err = client.IndexDrop(context.Background(), namespace, indexName)
	if err != nil {
		logger.Error("failed to drop index", slog.Any("error", err))
		return
	}

	logger.Info("successfully dropped index")
	logger.Info("done!")
}
```

## License
This project is licensed under the [Apache License 2.0](./LICENSE)
