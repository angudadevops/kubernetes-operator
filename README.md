# Service Deployment Operator

This Operator help you to deploy Application Deployment and Service on Kubernetes/Openshift Cluster. And this Operator built using [Operator Framework](https://github.com/operator-framework)

## Prerequisites
  - [Operator Framework](https://github.com/operator-framework)
  - [Golang](https://golang.org/dl/)
  - Kubernetes/Openshift Cluster

## Build an Operator

Run the below commands to build an operator on your machine

```
operator-sdk generate k8s  ## if you change the types of varaibles
operator-sdk generate crds 

operator-sdk build <Docker-Image>
```

## Deploy on your cluster

Create a CRD on your cluster with below command 

```
kubectl apply -f deploy/crds/servicedeployment_v1alpha1_sd_crd.yaml
```

Now Create service Account and role bindings for SD Operator with below commands
```
kubectl apply -f deploy/service_account.yaml
kubectl apply -f deploy/role.yaml
kubectl apply -f deploy/role_binding.yaml
```
Deploy Service Deployment Operator with below command
```
kubectl apply -f deploy/operator.yaml
```

## Create a Custom Resource



