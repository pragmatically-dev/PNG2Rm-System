# PNG2RM CONVERSION SYSTEM

[![rm1](https://img.shields.io/badge/rM1-supported-green)](https://remarkable.com/store/remarkable)
[![rm2](https://img.shields.io/badge/rM2-supported-green)](https://remarkable.com/store/remarkable-2)
[![Discord](https://img.shields.io/discord/385916768696139794.svg?label=reMarkable&logo=discord&logoColor=ffffff&color=7389D8&labelColor=6A7EC2)](https://discord.gg/ATqQGfu)
[![rM Hacks Discord](https://img.shields.io/discord/1153374327123759104.svg?label=rM%20Hacks&logo=discord&logoColor=ffffff&color=ffb759&labelColor=d99c4c)](https://discord.gg/bgVXW2bchN)


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

### Preview
<video width="320" height="240" controls>
  <source src="doc/de8466ca-0f8a-41e8-8f0f-f15885063855.mp4" type="video/mp4">
</video>


## Requirements:
- Golang: https://go.dev/dl
- Protoc: https://grpc.io/docs/protoc-installation/
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

---

⚠️ Please be sure to have rm-hacks and webinterface-onboot Installed⚠️

Now you should be able to convert your screenshots to rmlines in less than 4 sec


---

# Knowloadge base:
- The Wizzard: 
- https://www.youtube.com/playlist?list=PLy_6D98if3UJd5hxWNfAqKMr15HZqFnqf

- https://blog.stackademic.com/go-concurrency-visually-explained-select-statement-b546596c8e6b

- https://blog.stackademic.com/go-concurrency-visually-explained-channel-c6f88070aafa

- https://blog.stackademic.com/go-concurrency-visually-explained-select-statement-b546596c8e6bs