package main

import (
	"bytes"
	"log"
	"net"
	"os"

	png2rm "github.com/pragmatically-dev/png2rm/png2rm"
	"github.com/pragmatically-dev/png2rm/service"
	_ "google.golang.org/grpc/encoding/gzip"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v2"
)

// Config structure to hold YAML configuration
type Config struct {
	ImageFolder   string `yaml:"image_folder"`
	RunPath       string `yaml:"run_path"`
	ServerAddress string `yaml:"server_address"`
}

// PNGStore interface defines a method to save a PNG file
type PNGStore interface {
	Save(filename string, data bytes.Buffer) (string, error)
}

func main() {

	configFile, err := os.ReadFile("server-config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// Create a new gRPC server instance
	grpcServer := grpc.NewServer()

	// Create an instance of PNGStore
	pngStore := service.NewPNGStore(config.ImageFolder)

	// Create an instance of PNG2RmServiceServer
	png2RmServiceServer := service.NewPNG2RmServer(pngStore, config.RunPath)

	// Register the PNG2RmServiceServer with the gRPC server
	png2rm.RegisterPNG2RmServiceServer(grpcServer, png2RmServiceServer)

	reflection.Register(grpcServer)

	// Listen on port 4040
	listener, err := net.Listen("tcp", config.ServerAddress)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Server is listening on port %s", config.ServerAddress)

	// Start the gRPC server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
