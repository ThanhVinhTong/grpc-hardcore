package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"pcbook/pb"
	"pcbook/service"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// DEMO only
const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
)

func accessibleRoles() map[string][]string {
	const laptopServicePath = "/pcbook.LaptopService"
	return map[string][]string{
		laptopServicePath + "/CreateLaptop": {"admin"},
		// laptopServicePath + "/SearchLaptop": {"admin", "user"},
		laptopServicePath + "/UploadImage": {"admin"},
		laptopServicePath + "/RateLaptop":  {"admin", "user"},
	}
}

func seedUsers(userStore service.UserStore) error {
	err := createUser(userStore, "admin", "admin", "admin")
	if err != nil {
		return err
	}
	err = createUser(userStore, "user", "user", "user")
	if err != nil {
		return err
	}
	return nil
}

func createUser(userStore service.UserStore, username string, password string, role string) error {
	user, err := service.NewUser(username, password, role)
	if err != nil {
		return err
	}
	return userStore.Save(user)
}

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("start server at port %d", *port)

	userStore := service.NewInMemoryUserStore()
	err := seedUsers(userStore)
	if err != nil {
		log.Fatal("cannot seed users: ", err)
	}
	jwtManager := service.NewJWTManager(secretKey, tokenDuration)
	authServer := service.NewAuthServer(userStore, jwtManager)

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("img")
	ratingStore := service.NewInMemoryRatingStore()
	laptopServer := service.NewLaptopServer(laptopStore, imageStore, ratingStore)

	authInterceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	)

	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	reflection.Register(grpcServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
