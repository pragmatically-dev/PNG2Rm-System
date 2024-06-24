package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	fp "path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/fsnotify/fsnotify"
	png2rm "github.com/pragmatically-dev/png2rm/png2rm"

	"google.golang.org/grpc"
)

// Config structure to hold YAML configuration
type Config struct {
	DirToSearch   string `yaml:"dir_to_search"`
	DirToSave     string `yaml:"dir_to_save"`
	FilePrefix    string `yaml:"file_prefix"`
	ServerAddress string `yaml:"server_address"`
}

type PNG2RmServiceUploadAndConvertClient struct {
	service png2rm.PNG2RmServiceClient
}

func NewPNG2RmServiceClient(cc *grpc.ClientConn) *PNG2RmServiceUploadAndConvertClient {
	service := png2rm.NewPNG2RmServiceClient(cc)
	return &PNG2RmServiceUploadAndConvertClient{
		service: service,
	}
}

func uploadPNGAndConvert(png2rmClient png2rm.PNG2RmServiceClient, filepath string, dirToSave string) error {
	stream, err := png2rmClient.UploadAndConvert(context.Background())

	if err != nil {
		log.Fatalln("cannot upload image: ", err)
	}

	//START PNG UPLOAD PROCESS
	pngFile, err := os.Open(filepath)
	if err != nil {
		log.Fatalln("error when trying to open the screenshot: ", err)
	}
	defer pngFile.Close()

	res := &png2rm.UploadPNGRequest{
		Data: &png2rm.UploadPNGRequest_Filename{
			Filename: getFileName(filepath),
		},
	}

	if err := stream.Send(res); err != nil {
		log.Fatalln("couldn't send the first chunk of the screenshot: ", err)
		return err
	}

	reader := bufio.NewReader(pngFile)
	buff := make([]byte, 1024*16)

	for {
		n, err := reader.Read(buff)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		res := &png2rm.UploadPNGRequest{
			Data: &png2rm.UploadPNGRequest_DataChunck{
				DataChunck: buff[:n],
			},
		}
		if err := stream.Send(res); err != nil {
			return err
		}
	}

	if err := stream.CloseSend(); err != nil {
		return err
	}
	//END PNG UPLOAD PROCESS

	var docname string
	var imageData bytes.Buffer
	responseErrorChannel := make(chan error)

	go func() {
		defer close(responseErrorChannel)

		//Here we're receiving the rmDoc from the server
		//START RMDOC STREAMING
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				responseErrorChannel <- fmt.Errorf("cannot receive stream response: %v", err)
				return
			}

			if docname == "" {
				docname = res.GetDocname()

				//if docname != "" {
				//fmt.Printf("Received docname: %s\n", docname)
				//}
			}

			chunk := res.GetDataChunck()
			_, err = imageData.Write(chunk)

			if err != nil {
				responseErrorChannel <- err
				return
			}
		} //END RMDOC STREAMING

		//START SAVING RMDOC FILE
		file, err := os.Create(fmt.Sprintf("%s/%s", dirToSave, docname))
		if err != nil {
			responseErrorChannel <- err
			return
		}

		_, err = imageData.WriteTo(file)
		if err != nil {
			responseErrorChannel <- err
			return
		}
		//END SAVING RMDOC FILE

		responseErrorChannel <- nil
	}()

	//Checking possible error of the rmdoc streaming from the server
	if err := <-responseErrorChannel; err != nil {
		fmt.Printf("Error on server side streaming, %s\n", err)
		return err
	}

	go postRmDocToWebInterface(fmt.Sprintf("%s/%s", dirToSave, docname))

	return nil
}
func postRmDocToWebInterface(filepath string) {
	url := "http://10.11.99.1"

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create a form file field
	part, err := writer.CreateFormFile("file", fp.Base(file.Name()))
	if err != nil {
		fmt.Println("Error creating the form file field:", err)
		return
	}

	// Copy the file content to the form field
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error copying the file content:", err)
		return
	}

	// Close the writer to complete the multipart form
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing the writer:", err)
		return
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url+"/upload", &requestBody)
	if err != nil {
		fmt.Println("Error creating the request:", err)
		return
	}

	// Set the content type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending the request:", err)
		return
	}
	defer resp.Body.Close()

}

func watchForScreenshots(dirToSearch string, filePrefix string, client png2rm.PNG2RmServiceClient, dirToSave string) {
	//The watcher will give us two channels, one for Events and other for errors
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	err = watcher.Add(dirToSearch)
	if err != nil {
		log.Fatal(err)
	}

	go func() {

		for {
			// Here we receive data from two channels
			//A select blocks until one of its cases can run, then it executes that case.
			//more info on [https://blog.stackademic.com/go-concurrency-visually-explained-select-statement-b546596c8e6b]
			select {

			case event, ok := <-watcher.Events:

				if !ok {
					return
				}

				if (event.Has(fsnotify.Create)) && (strings.HasPrefix(filepath.Base(event.Name), filePrefix)) {

					//fmt.Printf("** File found: %s **\n", event.Name)
					time.Sleep(1200 * time.Millisecond)
					go uploadPNGAndConvert(client, event.Name, dirToSave)
					time.Sleep(5 * time.Second)
					if err := deleteFile(event.Name); err != nil {
						log.Printf("Error deleting file: %v\n", err)
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	//SYNC GO ROUTINE
	<-done
}

func getFileName(filepath string) string {
	filesplit := strings.Split(filepath, "/")
	filename := filesplit[len(filesplit)-1:]
	return filename[0]
}

func deleteFile(filepath string) error {
	return os.Remove(filepath)
}

func main() {
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	conn, err := grpc.NewClient(config.ServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Couldn't connect to the server")
	}
	png2rmClient := png2rm.NewPNG2RmServiceClient(conn)
	fmt.Println("<--- Looking for new Screenshots --->")
	watchForScreenshots(config.DirToSearch, config.FilePrefix, png2rmClient, config.DirToSave)
}
