package main

import (
	"log"
	"os"
	"net"
	"os/signal"
	"syscall"

	grpc "google.golang.org/grpc"
	pb "service/config/api"
	def "service/definitions"
	"service/utils"
)

const port = ":55045"

func main() {
	defer log.Println("Server shutdown successful")
	defer func() {
		if err := recover(); err != nil {
			log.Fatalln(err)
		} else {
			log.Println("User Server shutting down...")
		}
	}()

	channel := make(chan os.Signal, 1)
	defer close(channel)
	errChan := make(chan error, 1)
	defer close(errChan)

	signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM)
	
	lis, err := net.Listen("tcp", port)
	defer lis.Close()
	if err != nil {
		log.Fatalln(err)
	}	

	s := grpc.NewServer(grpc.UnaryInterceptor(utils.RecoveryInterceptor))
	pb.RegisterUserServer(s, &def.Server{})
	
	go func (){
		if err := s.Serve(lis); err != nil {
			errChan <- err
		}
	}()
	defer s.Stop()

	select {
		case <-channel:
			return
		case err := <-errChan:
			panic(err)
	}
}
