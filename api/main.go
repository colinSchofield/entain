package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/colinSchofield/entain/api/proto/racing"
	"github.com/colinSchofield/entain/api/proto/sporting"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

var (
	apiEndpoint          = flag.String("api-endpoint", "localhost:8000", "API endpoint")
	grpcRacingEndpoint   = flag.String("grpc-racing-endpoint", "localhost:9000", "gRPC server racing endpoint")
	grpcSportingEndpoint = flag.String("grpc-sporting-endpoint", "localhost:9001", "gRPC server sporting endpoint")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Printf("failed running api server: %s\n", err)
	}
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	if err := racing.RegisterRacingHandlerFromEndpoint(
		ctx,
		mux,
		*grpcRacingEndpoint,
		[]grpc.DialOption{grpc.WithInsecure()},
	); err != nil {
		return err
	}

	if err := sporting.RegisterSportingHandlerFromEndpoint(
		ctx,
		mux,
		*grpcSportingEndpoint,
		[]grpc.DialOption{grpc.WithInsecure()},
	); err != nil {
		return err
	}

	log.Printf("API server listening on: %s\n", *apiEndpoint)

	return http.ListenAndServe(*apiEndpoint, mux)
}
