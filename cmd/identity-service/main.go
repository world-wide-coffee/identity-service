package main

import (
	"log"

	"github.com/world-wide-coffee/identity-service/boundaries/grpc"
	"github.com/world-wide-coffee/identity-service/logic"
	"github.com/world-wide-coffee/identity-service/persistence/memdb"
)

func main() {
	persister := memdb.NewPersister()
	logic := logic.NewLogic().SetPersister(persister)
	server := grpc.NewServer().SetLogic(logic)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to listen and serve: %v", err)
	}
}
