package main

import (
	"database/sql"
	"flag"
	"log"
	"net"

	"github.com/colinSchofield/entain/sporting/db"
	"github.com/colinSchofield/entain/sporting/proto/sporting"
	"github.com/colinSchofield/entain/sporting/service"
	"google.golang.org/grpc"
)

var (
	grpcEndpoint = flag.String("grpc-endpoint", "localhost:9001", "gRPC server endpoint")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("failed running grpc server: %s\n", err)
	}
}

func run() error {
	conn, err := net.Listen("tcp", ":9001")
	if err != nil {
		return err
	}

	sportingDB, err := sql.Open("sqlite3", "./db/sporting.db")
	if err != nil {
		return err
	}

	sportsRepo := db.NewSportsRepo(sportingDB)
	if err := sportsRepo.Init(); err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	sporting.RegisterSportingServer(
		grpcServer,
		service.NewSportingService(
			sportsRepo,
		),
	)

	log.Printf("gRPC server listening on: %s\n", *grpcEndpoint)

	if err := grpcServer.Serve(conn); err != nil {
		return err
	}

	return nil
}
