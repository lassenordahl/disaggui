package main

import (
	"context"
	"log"
	"net"

	pb "github.com/lassenordahl/disaggui/obs/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCRDBServiceServer
}

func (s *server) ProcessFingerprint(ctx context.Context, req *pb.Fingerprint) (*pb.Ack, error) {
	log.Printf("Received fingerprint: %s", req.GetInput())
	// Here you can process the fingerprint as needed
	return &pb.Ack{Message: "Fingerprint processed"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCRDBServiceServer(s, &server{})
	log.Println("Server is running on port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
