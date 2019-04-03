package main

import (
	"github.com/spf13/pflag"
	"google.golang.org/grpc/metadata"

	"io"
	"os"
	"log"
	"fmt"
	"bytes"
	"context"
	"net/http"
	"path/filepath"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/grpclog"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/salihkemaloglu/gignox-main-beta-001/proto"
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
var (
	
	flagAllowAllOrigins = pflag.Bool("allow_all_origins", true, "allow requests from any origin.")
	flagAllowedOrigins  = pflag.StringSlice("allowed_origins", nil, "comma-separated list of origin URLs which are allowed to make cross-origin requests.")

	// useWebsockets = pflag.Bool("use_websockets", false, "whether to use beta websocket transport layer")
	enableTls       = pflag.Bool("enable_tls", true, "Use TLS - required for HTTP2.")
	tlsCertFilePath = pflag.String("tls_cert_file", "ssl/fullchain.pem", "Path to the CRT/PEM file.")
	tlsKeyFilePath  = pflag.String("tls_key_file", "ssl/privkey.pem", "Path to the private key file.")
	// flagHttpMaxWriteTimeout = pflag.Duration("server_http_max_write_timeout", 10*time.Second, "HTTP server config, max write duration.")
	// flagHttpMaxReadTimeout  = pflag.Duration("server_http_max_read_timeout", 10*time.Second, "HTTP server config, max read duration.")
)
func main(){

	pflag.Parse()

	port :=8902
	if *enableTls {
		port = 8903
	}

	fmt.Println("Main Service Started")

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	gigxMN.RegisterGigxMNServiceServer(grpcServer, &server{})

	allowedOrigins := makeAllowedOrigins(*flagAllowedOrigins)

	options := []grpcweb.Option{
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
		grpcweb.WithOriginFunc(makeHttpOriginFunc(allowedOrigins)),
	}

	wrappedGrpc := grpcweb.WrapServer(grpcServer, options...)
	handler := func(resp http.ResponseWriter, req *http.Request) {
		wrappedGrpc.ServeHTTP(resp, req)
	}

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: http.HandlerFunc(handler),
	}

	grpclog.Printf("Starting server. http port: %d, with TLS: %v", port, *enableTls)

	if *enableTls {
		fmt.Printf("server started as  https and listen to port: %v \n",port)
		if err := httpServer.ListenAndServeTLS(*tlsCertFilePath, *tlsKeyFilePath); err != nil {
			grpclog.Fatalf("failed starting http2 server: %v", err)
		}
	} else {
		fmt.Printf("server started as http and listen to port: %v \n",port)
		if err := httpServer.ListenAndServe(); err != nil {
			grpclog.Fatalf("failed starting http server: %v", err)
		}
	}
}
func makeHttpOriginFunc(allowedOrigins *allowedOrigins) func(origin string) bool {
	if *flagAllowAllOrigins {
		return func(origin string) bool {
			return true
		}
	}
	return allowedOrigins.IsAllowed
}
func makeAllowedOrigins(origins []string) *allowedOrigins {
	o := map[string]struct{}{}
	for _, allowedOrigin := range origins {
		o[allowedOrigin] = struct{}{}
	}
	return &allowedOrigins{
		origins: o,
	}
}

type allowedOrigins struct {
	origins map[string]struct{}
}
func (a *allowedOrigins) IsAllowed(origin string) bool {
	_, ok := a.origins[origin]
	return ok
}