package service_test

import (
	"context"
	"net"
	"pcbook/pb"
	"pcbook/sample"
	"pcbook/service"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

func startTestLaptopServer(t *testing.T) (*service.LaptopServer, string) {
	store := service.NewInMemoryLaptopStore()
	laptopServer := service.NewLaptopServer(store)
	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	go grpcServer.Serve(listener)

	return laptopServer, listener.Addr().String()
}

func newTestLaptopClient(t *testing.T, serverAddress string) pb.LaptopServiceClient {
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	return pb.NewLaptopServiceClient(conn)
}

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopServer, serverAddress := startTestLaptopServer(t)
	laptopClient := newTestLaptopClient(t, serverAddress)

	laptop := sample.NewLaptop()
	expectedID := laptop.Id
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	res, err := laptopClient.CreateLaptop(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedID, res.Id)

	// Verify the laptop is stored in the server
	foundLaptop, err := laptopServer.Store.Get(expectedID)
	require.NoError(t, err)
	require.NotNil(t, foundLaptop)

	// Check the same
	requireSameLaptop(t, laptop, foundLaptop)
}

func requireSameLaptop(t *testing.T, laptop1 *pb.Laptop, laptop2 *pb.Laptop) {
	json1, err1 := protojson.Marshal(laptop1)
	require.NoError(t, err1)

	json2, err2 := protojson.Marshal(laptop2)
	require.NoError(t, err2)

	require.Equal(t, json1, json2)
}
