package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	pb "github.com/lassenordahl/disaggui/obs/proto"
	"github.com/lassenordahl/disaggui/obs/uihandler"
	"github.com/rs/cors"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCRDBServiceServer
	db *sql.DB
}

func (s *server) ProcessFingerprint(ctx context.Context, req *pb.Fingerprint) (*pb.Ack, error) {
	timestamp := time.Now().Format(time.RFC3339)
	req.Timestamp = timestamp

	err := storeFingerprint(s.db, req.GetInput(), timestamp)
	if err != nil {
		return nil, err
	}

	log.Printf("Stored fingerprint: %s at %s", req.GetInput(), timestamp)
	return &pb.Ack{Message: "Fingerprint processed"}, nil
}

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		s := grpc.NewServer()
		pb.RegisterCRDBServiceServer(s, &server{db: db})
		log.Println("gRPC server is running on port :50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	s := &server{db: db}
	r := mux.NewRouter()

	// Add CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow all origins
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/fingerprints", s.listFingerprints).Methods("GET")
	apiRouter.HandleFunc("/fingerprints/count", s.listFingerprintCounts).Methods("GET")
	apiRouter.HandleFunc("/health", s.health).Methods("GET")

	// Serve the latest UI bundle
	uihandler.Serve("v1.0", r)

	handler := c.Handler(r)

	log.Println("HTTP server is running on port :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
