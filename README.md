# Simple chat application deployed in Kubernetes

This is a very simple application created to get familiar with the basics of k8s and golang. It allows users to connect to it and chat.

## Server side
`/server`

This is where the server will be stored that runs the chat app.

## Client/User side
`/client`

This is the client side of things where users send messages to the server and display the received messages.
If a user wants to connect to the server to start sending and receiving messages the need to run the `client/main.go` file.

## K8s config
`/k8s`

The configurations needed to set up the K8s side of things are stored here. These are needed for the k8s to run the server.

## Run locally/Demo
You can follow the instructions in  [minikube_demo.md](./minikube_demo.md) to run this locally with Minikube.