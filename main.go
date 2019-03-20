package main

import (

	"io"
	"os"
	"log"
	"fmt"
	"bytes"
	"context"
	"net/http"
	"path/filepath"
	"github.com/rs/cors"
	"github.com/go-chi/chi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/salihkemaloglu/gignox-main-beta-001/proto"
	"github.com/salihkemaloglu/gignox-main-beta-001/middleware"
	"github.com/salihkemaloglu/gignox-main-beta-001/proxy"
	// "github.com/salihkemaloglu/DemMain-beta-001/ipfs"
	//"github.com/salihkemaloglu/DemMain-beta-001/rabbitmq"
)

type server struct {
}
func (s *server) SayHello(ctx context.Context, req *gigxMN.HelloRequest) (*gigxMN.HelloResponse, error) {
	fmt.Printf("Main service is working...Received rpc from client, message=%s\n", req.GetMessage())
	// resultHash,err:=ipfs.UploadToIpfs([]byte(req.GetName()))
	// if err!=nil{
	// 	return nil,err
	// }

	// result,err:=rabbitmq.PublishDataToQueue(req.GetMessage())
	// if err!=nil{
	// 	return nil,err
	// }
	return &gigxMN.HelloResponse{Message: "Hello RR service is working..."}, nil
}

func (*server) LongGreet(stream gigxMN.GigxMNService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked with a streaming request\n")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			return stream.SendAndClose(&gigxMN.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		//uploadToIpfs(req.GetGreeting())
		//fmt.Println(req.GetGreeting())
		file := req.GetGreeting().GetContent()
		// copy example
		absPath, _ := filepath.Abs("handlersss.jpg")
		f, err := os.OpenFile(absPath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			
		}
		defer f.Close() //if there is a bug. If the call to os.Create fails, the function will return without closing the source file. defer is closing it.
		reader := bytes.NewReader(file)
		io.Copy(f, reader)

		
		 result += "Upload is success ! "
	}
}

func (*server) UploadFile(stream gigxMN.GigxMNService_UploadFileServer) error {
	fmt.Printf("GreetEveryone function was invoked with a streaming request\n")
	headers, ok := metadata.FromIncomingContext(stream.Context())
	if ok!=true{
		return  status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Authorization token is need"),
		)
	}
	token := headers["authorization"]
	fmt.Println(token)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		file := req.GetFile().GetContent()
		// copy example
		absPath, _ := filepath.Abs("handler1.jpg")
		f, err := os.OpenFile(absPath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			
		}
		defer f.Close() //if there is a bug. If the call to os.Create fails, the function will return without closing the source file. defer is closing it.
		reader := bytes.NewReader(file)
		io.Copy(f, reader)

		result := "Hello "

		sendErr := stream.Send(&gigxMN.UploadFileResponse{
			Result: result,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", sendErr)
			return sendErr
		}
	}

}
func (s *server) InsertFile(ctx context.Context, req *gigxMN.InsertFileRequest) (*gigxMN.InsertFileResponse, error) {
	return nil,nil
}
func main(){
	fmt.Println("Main Service Started")
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	gigxMN.RegisterGigxMNServiceServer(grpcServer, &server{})

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