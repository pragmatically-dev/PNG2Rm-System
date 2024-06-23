# PNG2RM CONVERSION SYSTEM
⚠️ Project under development ⚠️

## Overview
The System is implemented using Go and leverages gRPC for communication. It provides a server that can receive PNG files in chunks, save them, convert them using an external tool (drawj2d), and stream the resulting Remarkable document back to the client.

### Special Acknowledgment to:
>*This wouldn't have been possible without your help and incredible developments.*

> [<img src="https://github.com/Eeems.png" alt="Eeems" width="60"/>](https://github.com/Eeems)
> [<img src="https://a.fsdn.com/con/images/sandiego/icons/default-avatar.png" alt="A.Vontobel" width="60"/>](https://sourceforge.net/u/qwert2003/profile/)
> [<img src="https://github.com/rM-self-serve.png" alt="reMiss" width="60"/>](https://github.com/rM-self-serve)
> [<img src="https://github.com/mb1986.png" alt="mb1986" width="60"/>](https://github.com/mb1986)
> [<img src="https://github.com/atngames.png" alt="atngames" width="60"/>](https://github.com/atngames)

--- 

### How it works:

![alt text](doc/5af92af8-47d3-43cf-aa18-f74750ed8da5.jpeg)


## Requirements:
- Golang: https://go.dev/dl
### Server side:
- Protoc: https://grpc.io/docs/protoc-installation
- Drawj2d: https://sourceforge.net/projects/drawj2d/

### Client side (tablet):
- rm-hacks (enables screenshot feature): https://github.com/mb1986/rm-hacks

- webinterface-onboot: https://github.com/rM-self-serve/webinterface-onboot 


## How to setup (locally):
### Server:
1. **Install Go and necessary dependencies**:
   - Ensure you have Go installed on your system. You can download and install it from [golang.org](https://golang.org/).
   - Set up your development environment to work with Go. Configure `$GOPATH` and add `$GOPATH/bin` to your `$PATH`.

2. **Create the `server-config.yaml` file**:
   - This file should contain the necessary configuration for the server. Create a file named `server-config.yaml` in the same directory as the `main.go` file with the following content:

     ```yaml
     image_folder: "/path/to/image/folder"
     run_path: "/path/to/run/path"
     server_address: "localhost:4040"
     ```

   - Adjust the values of `image_folder`, `run_path`, and `server_address` according to your needs.


---

### Client (tablet):
1. **Transfrer the client binary and the config yaml file to the remarkable**:
   - After transfering the files you have to edit the config file with nano acording to your server config
   you can find the ip address of the server by running `ifconfig` or `ipconfig` 
   
   then:
   ```bash
   $remarkable: ~/ ./client &
   ```
now you should be able to convert your screenshots to rmlines in less than 4 sec

