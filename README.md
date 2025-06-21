# Simple chat application deployed in Kubernetes

This is a simple chat server that allows clients to connect to it and chat over TCP.

## Server side
`/server`

This is where the server that will run the chat app is stored. If you want to change to the way the server operates you will need to make your changes there.

## Client/User side
`/client`

This is the client side of things where users send messages to the server and display the received messages.
If a user wants to connect to the server to start sending and receiving messages the need to run the `client/main.go` file.

## K8s config
`/k8s`

The configurations needed to set up the K8s side of things are stored here. These are needed for the k8s to run the server.

## Run locally/Demo
You can follow the instructions in  [minikube_demo.md](./minikube_demo.md) to run this locally with Minikube.