# Minikube demo

If you don't have minikube installed you can do so from [here](https://minikube.sigs.k8s.io/docs/start/?arch=%2Fmacos%2Farm64%2Fstable%2Fbinary+download).

## 1. Start minikube
Spin up a minikube cluster if you don't have one already. The server side of the chat app will run there.
```sh
minikube start
```

## 2. Build the Docker image
Minikube runs its own Docker daemon. You must build the image inside Minikube’s Docker environment:
```sh
cd path/to/chat-app

# Switch to Minikube's Docker
eval $(minikube docker-env) 

# Build the image
docker build -t chat-server .
```
This builds an image named chat-server:latest and makes it available for Minikube to use.

Make sure image exists with 
```sh
docker images
```

SSH into minikube and check if docker image exists there - 
```sh
minikube ssh
```

If it doesn't, exit out of the minikube cluster terminal to go to your local one and load image into minikube:
```sh
minikube image load chat-server
```

## 3. Deploy to K8s
Apply the `deployment.yaml` and `service.yaml` files under the k8s folder.
```sh
kubectl apply -f ./k8s/
``` 
You can also apply each file individually:
```sh
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

## 4. Expose IP
Either use port forwarding or tunnel
```sh
# with port forwarding to pod directly (this will need to change when pod changes)
kubectl port-forward pod/<your-pod-name> <local-port>:<pod-port>

# with tunnel (will expose an external IP to service. This is possible becuase service is set to load balancer)
minikube tunnel
```

## Applying changes after making code changes
  1. Make sure you're using minikube docker daemon
      - `eval $(minikube docker-env)`
  2. Rebuild the image
      - `docker build -t chat-server .`
  3. Apply the changes in k8s
      - Rollout restart (terminate existing pods and create new ones): `kubectl rollout restart deployment chat-server`
      - If in k9s on the pods tab, just `ctrl + d` which will delete the existing pod and spin up a new one

## Cleanup
To clean up after setting this with minikube locally:
```sh
# Run the following commands if you want to cleanup but still keep your minikube cluster
kubectl delete -f k8s/deployment.yaml
kubectl delete -f k8s/service.yaml
minikube stop

# Run the following if you want to delete the cluster completely
minikube delete
```