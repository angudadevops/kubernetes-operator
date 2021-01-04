# Service Deployment Operator

This Operator help you to deploy Application Deployment and Service on Kubernetes/Openshift Cluster.

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

Once Operator is up and running, use the below snippet and update your parameters to deploy your code on your cluster.

```
apiVersion: servicedeployment.com/v1alpha1
kind: SD
metadata:
  name: nginx-sd
spec:
  # Add fields here
  replicas: 3
  image: nginx
  containerPort: 80
  nodePort: 31001
```

You can apply the below yaml for your testing 
```
kubectl apply -f deploy/crds/servicedeployment.com_v1alpha1_sd_cr.yaml
```

