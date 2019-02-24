package main

import (
	"fmt"
	"net/http"
	"context"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/salihkemaloglu/DemMain-beta-001/proto"
	"github.com/salihkemaloglu/DemMain-beta-001/middleware"
	"github.com/salihkemaloglu/DemMain-beta-001/proxy"
	"github.com/salihkemaloglu/DemMain-beta-001/ipfs"
	"github.com/salihkemaloglu/DemMain-beta-001/rabbitmq"
)

type server struct {
}
func (s *server) SayHello(ctx context.Context, req *demMN.HelloRequest) (*demMN.HelloResponse, error) {
	fmt.Printf("Main service is working...Received rpc from client, message=%s\n", req.GetName())
	resultHash,err:=ipfs.UploadToIpfs([]byte(req.GetName()))
	if err!=nil{
		return nil,err
	}

	result,err:=rabbitmq.PublishDataToQueue(resultHash)
	if err!=nil{
		return nil,err
	}
	return &demMN.HelloResponse{Message: result+" :"+resultHash }, nil
}

func main(){
	fmt.Println("Main Service Started")
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	demMN.RegisterDemServiceServer(grpcServer, &server{})

	wrappedGrpc := grpcweb.WrapServer(grpcServer)

	router := chi.NewRouter()
	router.Use(
		chiMiddleware.Logger,
		chiMiddleware.Recoverer,
		middleware. NewGrpcWebMiddleware(wrappedGrpc).Handler,// Must come before general CORS handling
		cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}).Handler,
	)

	router.Get("/article-proxy", proxy.Article)

	if err := http.ListenAndServe(":8900", router); err != nil {
		grpclog.Fatalf("failed starting http2 server: %v", err)
	}
}