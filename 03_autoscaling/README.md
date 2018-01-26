# Overview
## This Demo
This demo is a proof of technology I did for a customer. It´s a simple Kubernetes pod based autoscaler which scales the backend pods based on the queue depth of the message broker. The following components will be deployed here:

- Kubernetes service for RabbitMQ (service discovery, ext. access for mgmt. UI)
- Kubernetes deployments for RabbitMQ, sender, receiver and autoscaler

The sender just generates random integers in random intervals, publishs them to RabbitMQ, where they´re pulled by *n* (autoscaled) workers which calculate the fibonacci sequence to simulate CPU load.  

<img src=./img/rmq_qdepth_example.png width="800">

---
**NOTE** The code works and demonstrates the capabilities. But it also needs a big overhaul (the [Issue#2](https://github.com/embano1/vmworld2017/issues/2)). Bear with me please ,) 

---

## Why not the HPA (Horizontal Pod Autoscaler)? 
I wanted to use a custom metric, i.e. the queue depth of RabbitMQ as a trigger for scaling. This was before the Kubernetes HPA v2, which added support for autoscaling based on custom metrics. 

## High-level Architecture

The following picture describes the components involved and their interaction. 

<img src=./img/high-level_architecture.png width="800">

---

**NOTE** The autoscaler can also be run outside the cluster. Therefore you have to compile the source code located in `autoscaler_ext`.

--- 

# Quick Start
## A Note on Kubernetes Role-Based Access Control (RBAC)
Depending on which version of Kubernetes you're running, you might use a cluster with Kubernetes RBAC turned on by default. The Kubernetes manifest files in this repo are prepared for environments with/without RBAC enabled. Keep this in mind when following along with the deployment steps described below.

You can read more about Kubernetes RBAC in the following resources:

- http://blog.kubernetes.io/2017/04/rbac-support-in-kubernetes.html
- https://kubernetes.io/docs/admin/authorization/rbac/
- https://kubernetes.io/docs/admin/authorization/

## Software Requirements to run the Demo
- kubectl (must match supported Kubernetes version requirements for the tested Kubernetes cluster versions mentioned below)
- Have a running Kubernetes environment
- Tested Kubernetes environments:
  - minikube (tested with v0.20.0 and Kubernetes v1.6.4)
  - kubeadm (tested with Kubernetes v1.9.1)

## Run

```bash
# Clone the repo
$ git clone https://github.com/embano1/vmworld2017

# Change into the autoscaler directory
$ cd vmworld2017/03_autoscaling

# Depending on whether your cluster runs with RBAC enabled/ disabled use <norbac> or <rbac>
$ kubectl create -f kubernetes/rbac # or norbac

# Wait for images to be pulled and pods being started in namespace default, then
# Access RabbitMQ UI through NodePort and NodeIP (see "Accessing RabbitMQ statistics" below)
# Scale sender deployment and see how the system (RabbitMQ queue details) adapts: 
$ kubectl scale deploy sender --replicas=20
```

---

**NOTE** You might see sender/receiver/autoscaler pods to fail initially. This is due to the rabbitmq pod still being pending, e.g. during the image pull phase. The application is designed to stabilize, leveraging Kubernetes backoff logic. 

---
  
## Accessing RabbitMQ Statistics

RabbitMQ is exposed using Kubernetes `NodePort`. This demo does not use an external load-balancer. Username/ password: guest/guest.

```bash
# Get the nodeport where the RabitMQ UI is exposed
$ kubectl get svc rabbitmq -o jsonpath='{.spec.ports[1].nodePort}'
```

---

**NOTE** If you work with minikube, you can get the port of the RabbitMQ management UI with "minikube service list" (look for target port **15672**). Then you can access the management portal through your browser, e.g. `http://$(minikube ip):30383` (as a `nodeport` example).

---

## Get autoscaler Metrics
    - kubectl get po | grep autoscaler  
    - kubectl logs -f <autoscaler_pod>

# Cleanup

```bash
# Depending on whether your cluster runs with RBAC enabled/ disabled use <norbac> or <rbac>
$ kubectl delete -f kubernetes/rbac # or norbac
```
