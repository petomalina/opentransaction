package opentransaction

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Serve(server TenantServer) error {
	lis, err := net.Listen("tcp", "50051")
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	RegisterTenantServer(s, server)
	reflection.Register(s)

	err = s.Serve(lis)
	if err != nil {
		log.Fatalln(err)
	}

	return err
}
