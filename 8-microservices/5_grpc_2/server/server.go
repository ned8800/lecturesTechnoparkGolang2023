package main

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/lectures/8-microservices/4_grpc/session"
	"log"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/tap"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

func authInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	md, _ := metadata.FromIncomingContext(ctx)

	reply, err := handler(ctx, req)

	fmt.Printf(`--
	after incoming call=%v
	req=%#v
	reply=%#v
	time=%v
	md=%v
	err=%v
`, info.FullMethod, req, reply, time.Since(start), md, err)
	return reply, err
}

// -----

func rateLimiter(ctx context.Context, info *tap.Info) (context.Context, error) {
	fmt.Printf("--\ncheck ratelim for %s\n", info.FullMethodName)
	pass := rand.Int31n(2)
	if pass == 0 {
		return ctx, errors.New("you shall not pass")
	}
	return ctx, nil
}

// -----

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("cant listet port", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor),
		grpc.InTapHandle(rateLimiter),
	)

	session.RegisterAuthCheckerServer(server, NewSessionManager())

	fmt.Println("starting server at :8081")
	server.Serve(lis)
}
