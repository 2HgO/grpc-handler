package main

import (
	"log"
	"net"

	grpc "google.golang.org/grpc"
	pb "service/config/api"
	def "service/definitions"
	"service/utils"
)

const port = ":55045"

func main() {
	defer func() interface{} {
		if err := recover(); err != nil {
			println(err)
			return err
		}
		panic(0)
	}()
	
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(utils.RecoveryInterceptor))
	pb.RegisterUserServer(s, &def.Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
