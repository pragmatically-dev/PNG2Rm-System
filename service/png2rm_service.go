package service

import (
	"bufio" // Buffered I/O package
	"bytes" // Bytes manipulation package

	// Context for managing deadlines and cancellation signals
	"fmt" // Formatting package
	"io"  // Basic I/O interface
	"log" // Logging package

	// Operating system package for file handling
	// Executing external commands

	"github.com/pragmatically-dev/PoC-drawj2d-port-go/remarkablepage"
	"github.com/pragmatically-dev/png2rm/png2rm" // Custom package for PNG to Remarkable service
	"google.golang.org/grpc/codes"               // gRPC status codes
	"google.golang.org/grpc/status"              // gRPC error handling
)

// PNG2RmServiceServer struct to implement the server methods
type PNG2RmServiceServer struct {
	png2rm.UnimplementedPNG2RmServiceServer          // Embedding unimplemented server for forward compatibility
	pngStore                                PNGStore // Interface to handle PNG storage
	runPath                                 string   // Path where conversion operations run

}

// mustEmbedUnimplementedPNG2RmServiceServer required by gRPC to embed unimplemented server methods
func (server *PNG2RmServiceServer) mustEmbedUnimplementedPNG2RmServiceServer() {}

// NewPNG2RmServer constructs a new PNG2RmServiceServer
func NewPNG2RmServer(pngStore PNGStore, runPath string) *PNG2RmServiceServer {
	// Initialize the PNG2RmServiceServer with the provided paths and PNGStore
	return &PNG2RmServiceServer{
		pngStore: pngStore,
		runPath:  runPath,
	}
}

// UploadAndConvert handles streaming upload and conversion of PNG files
func (server *PNG2RmServiceServer) UploadAndConvert(stream png2rm.PNG2RmService_UploadAndConvertServer) error {
	var imageData bytes.Buffer // Buffer to store image data chunks
	var filename string        // Variable to store the filename

	// Receive data from the client in chunks
	for {
		req, err := stream.Recv() // Receive a chunk of data from the stream
		if err == io.EOF {        // End of file, no more data to receive
			log.Println("No more data")
			break
		}
		if err != nil { // Handle any other receiving errors
			return logError(status.Errorf(codes.Unknown, "cannot receive stream request: %v", err))
		}

		if filename == "" { // Set filename if not already set
			filename = req.GetFilename()
		}

		chunk := req.GetDataChunck()                      // Get the data chunk from the request
		if _, err := imageData.Write(chunk); err != nil { // Write chunk to the buffer
			return logError(status.Errorf(codes.Internal, "cannot write chunks: %v", err))
		}
	}

	if filename == "" { // Ensure a filename was provided
		return logError(status.Errorf(codes.InvalidArgument, "no filename provided"))
	}

	// Save the received PNG file using the PNGStore
	pngFilename, err := server.pngStore.Save(filename, imageData)
	if err != nil { // Handle error in saving the PNG file
		return logError(status.Errorf(codes.Internal, "cannot save image: %v", err))
	}

	decodedImg := remarkablepage.LaplacianEdgeDetection(server.runPath + "/ToConvert/" + pngFilename)
	if decodedImg == nil {
		return logError(status.Errorf(codes.Internal, "cannot decode to gray image: %v", err))
	}

	// Define the path of the resulting Remarkable document
	rmdocPath := fmt.Sprintf("%s/%s.rmdoc", server.runPath, pngFilename)
	fmt.Println(rmdocPath)
	rmzip, rmzipname := remarkablepage.CreateRmDoc(pngFilename, decodedImg)

	// Send the document name as the first response to the client
	res := &png2rm.UploadPNGResponse{
		Data: &png2rm.UploadPNGResponse_Docname{
			Docname: rmzipname,
		},
	}
	if err := stream.Send(res); err != nil { // Handle error in sending the response
		return logError(status.Errorf(codes.Internal, "cannot send the first response: %v", err))
	}

	// Stream the Remarkable document back to the client in chunks
	reader := bufio.NewReader(rmzip) // Create a buffered reader for the document
	buff := make([]byte, 1024*32)    // Buffer to hold file chunks

	for {
		n, err := reader.Read(buff) // Read a chunk of the document
		if err == io.EOF {          // End of file, no more data to read
			break
		}
		if err != nil { // Handle any other reading errors
			return logError(status.Errorf(codes.Internal, "error reading rmdoc file: %v", err))
		}
		res := &png2rm.UploadPNGResponse{ // Prepare the response with the data chunk
			Data: &png2rm.UploadPNGResponse_DataChunck{
				DataChunck: buff[:n],
			},
		}
		if err := stream.Send(res); err != nil { // Send the chunk to the client
			return logError(status.Errorf(codes.Internal, "cannot send chunk: %v", err))
		}
	}

	return nil // Indicate successful completion
}

// logError logs the provided error and returns it
func logError(err error) error {
	if err != nil { // If there is an error
		log.Print(err) // Log the error
	}
	return err // Return the error
}
