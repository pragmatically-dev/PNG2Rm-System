package service

import (
	"bufio" // Buffered I/O package
	"bytes" // Bytes manipulation package
	"context"

	// Context for managing deadlines and cancellation signals
	"fmt"     // Formatting package
	"io"      // Basic I/O interface
	"log"     // Logging package
	"os"      // Operating system package for file handling
	"os/exec" // Executing external commands

	"github.com/contiv/executor"
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

	// Create an HCL file necessary for the conversion process
	if err := createHCLFile(server.runPath, pngFilename); err != nil { // Handle error in creating HCL file
		return logError(status.Errorf(codes.Internal, "cannot create HCL file: %v", err))
	}

	// Prepare the command to convert PNG to Remarkable document using drawj2d
	cmd := exec.Command("./drawj2d", "-Trmdoc", fmt.Sprintf("%s/CaptureToConvert.hcl", server.runPath), "-o", fmt.Sprintf("%s/%s.rmdoc", server.runPath, pngFilename))
	exec := executor.New(cmd)
	exec.Start()
	er, err := exec.Wait(context.Background())

	fmt.Printf("Conversion Exit Code: %d\n", er.ExitStatus)
	if er.ExitStatus != 0 {
		return logError(status.Errorf(codes.Internal, "Drawj2d Error : %v", err))

	}
	// Define the path of the resulting Remarkable document
	rmdocPath := fmt.Sprintf("%s/%s.rmdoc", server.runPath, pngFilename)
	fmt.Println(rmdocPath)
	rmdoc, err := os.Open(rmdocPath) // Open the resulting Remarkable document file
	if err != nil {                  // Handle error in opening the document file
		return logError(status.Errorf(codes.Internal, "cannot open rmdoc file: %v", err))
	}
	defer rmdoc.Close() // Ensure the file is closed after processing

	// Send the document name as the first response to the client
	res := &png2rm.UploadPNGResponse{
		Data: &png2rm.UploadPNGResponse_Docname{
			Docname: pngFilename + ".rmdoc",
		},
	}
	if err := stream.Send(res); err != nil { // Handle error in sending the response
		return logError(status.Errorf(codes.Internal, "cannot send the first response: %v", err))
	}

	// Stream the Remarkable document back to the client in chunks
	reader := bufio.NewReader(rmdoc) // Create a buffered reader for the document
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

// createHCLFile creates an HCL file required for the PNG to Remarkable conversion
func createHCLFile(runPath string, filepath string) error {
	// Open a new file for writing, create it if it doesn't exist
	file, err := os.OpenFile(fmt.Sprintf("%s/CaptureToConvert.hcl", runPath), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil { // Handle error in opening the file
		return err
	}
	defer file.Close() // Ensure the file is closed after writing

	// Write the expression needed for conversion into the HCL file
	expression := fmt.Sprintf("image %s 300 0 0 1.32", runPath+"/ToConvert/"+filepath)
	if _, err := file.Write([]byte(expression)); err != nil { // Handle error in writing to the file
		return err
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
