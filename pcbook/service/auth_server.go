package service

import (
	"context"
	"pcbook/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	userStore  UserStore
	jwtManager *JWTManager
}

func NewAuthServer(userStore UserStore, jwtManager *JWTManager) *AuthServer {
	return &AuthServer{
		userStore:  userStore,
		jwtManager: jwtManager,
	}
}

func (server *AuthServer) Login(
	ctx context.Context,
	req *pb.LoginRequest,
) (*pb.LoginResponse, error) {
	user, err := server.userStore.Find(req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cant find user: %v", err)
	}

	if user == nil || !user.IsPasswordValid(req.GetPassword()) {
		return nil, status.Errorf(codes.NotFound, "username or password is wrong: %v", err)
	}

	token, err := server.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cant generate access token: %v", err)
	}

	res := &pb.LoginResponse{
		AccessToken: token,
	}

	return res, nil
}
