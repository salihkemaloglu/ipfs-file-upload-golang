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
)

type server struct {
}
func (s *server) SayHello(ctx context.Context, in *demMN.HelloRequest) (*demMN.HelloResponse, error) {
	fmt.Printf("Main service is working...Received rpc from client, message=%s\n", in.GetName())
	return &demMN.HelloResponse{Message: "Hello Main service is working....!!!" }, nil
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