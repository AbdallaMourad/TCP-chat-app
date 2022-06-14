# TCP-chat-app

## Disclamer
This project was build to understand 
- How to build a TCP server in Go language.
- Understand Concurrency control for accepting multiple connections and sending messages in no pre-defined order.

## Build
``` bash
docker build . -t chat-app
docker run -dp <PORT_TO_BIND>:8080 chat-app
```

## Connect
``` bash
nc 127.0.0.1 <BOUND_PORT>
```
