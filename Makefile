gen:
	protoc --proto_path=proto proto/*.proto --go_out=./ --go-grpc_out=./

clean:
	rm png2rm/*.go 

