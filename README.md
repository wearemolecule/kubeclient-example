![Molecule Software](https://avatars1.githubusercontent.com/u/2736908?v=3&s=100 "Molecule Software")
# Kubeclient Example App

A CLI application that uses [Molecule's kubeclient library](https://github.com/wearemolecule/kubeclient). __kubeclient__ is a minimal dependency library to easily interact with a kubernetes cluster. The purpose of this application is to demonstrate how to use the library in a golang project. The CLI is a simple [kubectl](http://kubernetes.io/docs/user-guide/kubectl-overview/) clone.

### Commands
* __create__ *(create --help)*
  * Create a pod or replication controller from the sample JSON config (pod/replication-controller.json)
* __delete__ *(delete --help)*
  * Delete a pod or replication controller given a name
* __list__ *(list --help)*
  * List pods or replication controllers information
* __logs__ *(logs --help)*
  * Print the logs of the first container for a given pod

### Required env vars

* __CERTS_PATH__
  * Path to your kubernetes certs (should be the same dir that kubectl uses)
* __KUBERNETES_SERVICE_HOST__
  * Ex: "kubernetes-master.your-company.com"
* __KUBERNETES_SERVICE_PORT__
  * Ex: "001"

