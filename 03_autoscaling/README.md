# Overview
## This Demo
This demo is a proof of technology I did for a customer. It´s a simple Kubernetes pod based autoscaler which scales the backend pods based on the queue depth of the message broker. The following components will be deployed here:

- Kubernetes service for RabbitMQ (service discovery, ext. access for mgmt. UI)
- Kubernetes deployments for RabbitMQ, sender, receiver and autoscaler

The sender just generates random integers in random intervals, publishs them to RabbitMQ, where they´re pulled by up to *n* (default: 5, see `Flags to configure Autoscaler` below) autoscaled workers which calculate the fibonacci sequence to simulate CPU load.  

<img src=./img/rmq_qdepth_example.png width="800">


## Why not the HPA (Horizontal Pod Autoscaler)? 
I wanted to use a custom metric, i.e. the queue depth of RabbitMQ as a trigger for scaling. Kubernetes HPA v1 has no support for custom metrics. HPA v2 added this capability but this demo shows how easy it is adding auto-scaling capablities with a custom controller - sometimes easier than using HPA v2 custom metrics. Also, if you run an older Kubernetes version without HPA v2 support, the controller based approach still works.

## High-level Architecture

The following picture describes the components involved and their interaction. 

<img src=./img/high-level_architecture.png width="800">

---

**NOTE** By default, the autoscaler is configured to run outside the Kubernetes cluster. This makes playing with the code easier. Use `-i` flag to start in-cluster mode. The Kubernetes deployment manifests are correctly prepared.

--- 

# Quick Start
## A Note on Kubernetes Role-Based Access Control (RBAC)
Depending on which version of Kubernetes you're running, you might use a cluster with Kubernetes RBAC turned on by default. The Kubernetes manifest files in this repo are prepared for environments with RBAC enabled. Keep this in mind when following along with the deployment steps described below.

You can read more about Kubernetes RBAC in the following resources:

- http://blog.kubernetes.io/2017/04/rbac-support-in-kubernetes.html
- https://kubernetes.io/docs/admin/authorization/rbac/
- https://kubernetes.io/docs/admin/authorization/

## Software Requirements to run the Demo
- kubectl
- Obviously, have a running Kubernetes environment :)
- Tested Kubernetes environments:
  - minikube (tested with v0.25.0 and Kubernetes v1.9.4)
  - kubeadm Cluster (tested with Kubernetes v1.9.x)
  - RabbitMQ services are exposed using `type: NodePort` - please change if you want to use `type: LoadBalancer`

## Run

```bash
# Clone the repo
$ git clone https://github.com/embano1/vmworld2017

# Change into the autoscaler directory
$ cd vmworld2017/03_autoscaling

# Create a custom Kubernetes namespace where you want to deploy this demo, e.g. "pubsub"
$ NAMESPACE=pubsub
$ kubectl create ns ${NAMESPACE}

# Use the namespace you created above
$ kubectl -n ${NAMESPACE} create -f kubernetes/

# Wait for images to be pulled and pods being started in namespace default, then
# Access RabbitMQ UI through NodePort and NodeIP (see "Accessing RabbitMQ statistics" below)
# Scale sender deployment and see how the system (RabbitMQ queue details) adapts: 
$ kubectl -n ${NAMESPACE} scale deploy sender --replicas=20
```

---

**NOTE** You might see sender/receiver/autoscaler pods to fail initially. This is due to the rabbitmq pod still being pending, e.g. during the image pull phase. The application is designed to stabilize, leveraging Kubernetes backoff logic. 

---
  
## Accessing RabbitMQ Statistics

RabbitMQ is exposed using Kubernetes `NodePort`. This demo does not use an external load-balancer. Username/ password: guest/guest.

```bash
# Get the nodeport where the RabitMQ Management UI is exposed
$ kubectl -n ${NAMESPACE} get svc rabbitmq -o jsonpath='{.spec.ports[1].nodePort}'
```

---

**NOTE** If you work with minikube, you can get the port of the RabbitMQ management UI with "minikube service list" (look for target port **15672**). Then you can access the management portal through your browser, e.g. `http://$(minikube ip):30383` (as a `NodePort` example; port will be different in your environment!).

---

## Get autoscaler logs

```bash
$ SCALER=$(k -n ${NAMESPACE} get po -l app=autoscaler -o jsonpath='{.items[].metadata.name}')
$ kubectl -n ${NAMESPACE} logs -f ${SCALER}
```

## Flags to configure Autoscaler

Removed `glog` flags for clarity. 

```bash
$ ./autoscaler -h
Usage of ./autoscaler:

  -a          Fully qualified address of the AMQP broker (default "amqp://guest:guest@rabbitmq:5672/")
  -d          Deployment to scale (default "receiver")
  -i          Run autoscaller inside Kubernetes cluster.
  -kubeconfig (optional) absolute path to the kubeconfig file (default "/Users/somename/.kube/config")
  -n          Namespace of deployment (default "pubsub")
  -w          Max workers (Kubernetes Pods) to scale up (default 5)
```

# Cleanup

```bash
# Delete all resources, incl. RBAC related objects
$ kubectl -n ${NAMESPACE} delete -f kubernetes/
$ kubectl delete ns ${NAMESPACE} 
```
