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

---

### Client (tablet):
Copy the client and the config.yaml to the device

