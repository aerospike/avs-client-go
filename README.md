# Aerospike Vector Search Go Client

> :warning: The go client is currently in development. APIs will break!

## Overview

This repository contains the Go client for interacting with Aerospike's vector search capabilities. It is designed to provide an easy-to-use interface for performing vector searches within an Aerospike cluster.

## Example Usage

```go
// Example code will be provided once the API is stable.
```

## Current Functionality
- Connection to Aerospike clusters
- Create, List, Drop Indexes

## Installation
To install the Aerospike Vector Search Go Client, use the following command:

```bash
go get github.com/aerospike/aerospike-vector-search-go-client
```

## Quick Start
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
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	hostPort := NewHostPort("127.0.0.1", 5000)

	adminClient, err := avs.NewAdminClient(
		connectCtx,
		HostPortSlice{hostPort},
		nil,
		true,
		logger,
	)
	cancel()

	if err != nil {
		logger.Error("failed to create admin client", slog.Any("error", err))
		return
	}

	logger.Info("successfully connected to server")
	defer adminClient.Close()

	namespace := "test"
	set := "testset"
	indexName := "bookIndex"
	vectorField := "vector"
	dimensions := 10
	distanceMetric := protos.VectorDistanceMetric_MANHATTAN

	err = adminClient.IndexCreate(
		context.Background(),
		namespace,
		[]string{set},
		indexName,
		vectorField,
		uint32(dimensions),
		distanceMetric,
		nil,
		nil,
		nil,
	)
	if err != nil {
		logger.Error("failed to create index", slog.Any("error", err))
		return
	}

	logger.Info("successfully created index")

	indexes, err := adminClient.IndexList(context.Background())
	if err != nil {
		logger.Error("failed to list indexes", slog.Any("error", err))
		return
	}

	for _, index := range indexes.GetIndices() {
		fmt.Println(index.String())
	}

	err = adminClient.IndexDrop(context.Background(), namespace, indexName)
	if err != nil {
		logger.Error("failed to drop index", slog.Any("error", err))
		return
	}

	logger.Info("successfully dropped index")
	logger.Info("done!")
}
```

## License
This project is licensed under the MIT License - see the LICENSE file for details.
