package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"pcbook/pb"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxImageSize = 1 << 20 // 1MB

// LaptopServer is the server for the LaptopService.
type LaptopServer struct {
	pb.UnimplementedLaptopServiceServer
	laptopStore LaptopStore
	imageStore  ImageStore
}

// NewLaptopServer creates a new LaptopServer.
func NewLaptopServer(laptopStore LaptopStore, imageStore ImageStore) *LaptopServer {
	return &LaptopServer{
		laptopStore: laptopStore,
		imageStore:  imageStore,
	}
}

// CreateLaptop implements the Unary CreateLaptop RPC method.
func (server *LaptopServer) CreateLaptop(
	ctx context.Context,
	req *pb.CreateLaptopRequest,
) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	log.Printf("Receive CreateLaptop request with ID: %s", laptop.Id)

	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid laptop ID: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate new ID: %v", err)
		}
		laptop.Id = id.String()
	}

	// heavy task for timeout simulation
	// time.Sleep(6 * time.Second)

	err := server.laptopStore.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save laptop: %v", err)
	}

	log.Printf("Laptop saved with ID: %s", laptop.Id)

	return &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}, nil
}

func (server *LaptopServer) SearchLaptop(
	req *pb.SearchLaptopRequest,
	stream pb.LaptopService_SearchLaptopServer,
) error {
	filter := req.GetFilter()
	log.Printf("receive a laptop filter")

	err := server.laptopStore.Search(
		stream.Context(),
		filter,
		func(laptop *pb.Laptop) error {
			res := &pb.SearchLaptopResponse{Laptop: laptop}

			err := stream.Send(res)
			if err != nil {
				return err
			}

			log.Printf("sent laptop with id: %s", laptop.GetId())
			return nil
		},
	)

	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}
	return nil
}

func (server *LaptopServer) UploadImage(
	stream pb.LaptopService_UploadImageServer,
) error {
	req, err := stream.Recv()
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot receive image info: %v", err))
	}

	laptopID := req.GetInfo().GetLaptopId()
	imageType := req.GetInfo().GetImageType()
	log.Printf("receive an upload-image request for laptop %s with image type %s", laptopID, imageType)

	laptop, err := server.laptopStore.Get(laptopID)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot get laptop: %v", err))
	}
	if laptop == nil {
		return logError(status.Errorf(codes.NotFound, "laptop not found: %s", laptopID))
	}

	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		// check context error
		if err := contextError(stream.Context()); err != nil {
			return logError(err)
		}

		log.Print("waiting to receive more data")

		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "cannot receive image chunk: %v", err))
		}

		chunkData := req.GetChunkData()
		size := len(chunkData)

		log.Printf("received a chunk with size: %d", size)

		imageSize += size
		if imageSize > maxImageSize {
			return logError(status.Errorf(codes.InvalidArgument, "image size is too large: %d", imageSize))
		}

		// Write slowly
		// time.Sleep(2 * time.Second)

		_, err = imageData.Write(chunkData)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "cannot write image data: %v", err))
		}
	}

	imageID, err := server.imageStore.Save(laptopID, imageType, imageData)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot save image: %v", err))
	}

	res := &pb.UploadImageResponse{
		ImageId: imageID,
		Size:    uint32(imageSize),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot send response: %v", err))
	}

	log.Printf("image saved with id: %s", imageID)
	return nil
}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return status.Errorf(codes.Canceled, "request cancelled")
	case context.DeadlineExceeded:
		return status.Errorf(codes.DeadlineExceeded, "deadline exceeded")
	default:
		return nil
	}
}
