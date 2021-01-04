# Service Deployment Operator

This Operator help you to deploy Application Deployment and Service on Kubernetes/Openshift Cluster. And this Operator built using [Operator Framework](https://github.com/operator-framework)

# Prerequisites
  - [Operator Framework](https://github.com/operator-framework)
  - [Golang](https://golang.org/dl/)
  - Kubernetes/Openshift Cluster

# Build an Operator

Run the below commands to build an operator on your machine

```
operator-sdk generate k8s  ## if you change the types of varaibles
operator-sdk generate crds 

operator-sdk build <Docker-Image>
```

# Deploy on your cluster




