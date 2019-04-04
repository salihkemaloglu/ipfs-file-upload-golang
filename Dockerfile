FROM golang:1.9

RUN go get gopkg.in/mgo.v2/bson
RUN go get github.com/spf13/pflag
RUN go get golang.org/x/net/http2									
RUN go get golang.org/x/text/secure/bidirule	
RUN go get golang.org/x/text/unicode/bidi		
RUN go get golang.org/x/text/unicode/norm	
RUN go get github.com/golang/protobuf/proto	
RUN go get google.golang.org/grpc				
RUN go get google.golang.org/genproto/googleapis/rpc/status	
RUN go get golang.org/x/sys/unix		
RUN go get github.com/improbable-eng/grpc-web/go/grpcweb 
RUN go get github.com/dgrijalva/jwt-go
RUN go get golang.org/x/crypto/acme/autocert
# set environment path
ENV PATH /go/bin:$PATH

# cd into the api code directory
WORKDIR /go/src/github.com/salihkemaloglu/gignox-main-beta-001

# create ssh directory
RUN mkdir ~/.ssh
RUN touch ~/.ssh/known_hosts
RUN ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts

# allow private repo pull
RUN git config --global url."https://e4d5159cc774d99744024453431f00ddbb8d7b1d:x-oauth-basic@github.com/".insteadOf "https://github.com/"

# copy the local package files to the container's workspace
ADD . /go/src/github.com/salihkemaloglu/gignox-main-beta-001

# install the program
RUN go install github.com/salihkemaloglu/gignox-main-beta-001

# expose default port
EXPOSE 80 443
# start application
CMD ["go","run","main.go"] 