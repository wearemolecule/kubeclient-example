package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	kube "github.com/wearemolecule/kubeclient"
	"golang.org/x/build/kubernetes/api"
	"golang.org/x/net/context"
)

func main() {
	kubeClient, err := kube.GetKubeClientFromEnv()
	if err != nil {
		panic(fmt.Errorf("Failed to connect to kubernetes.", err))
	}

	commandPtr := flag.String("command", "noop", "kubeclient command")
	namespacePtr := flag.String("namespace", "", "namespace to work in")
	namePtr := flag.String("name", "", "name of pod/rc")
	imagePtr := flag.String("image", "", "image of rc")
	versionPtr := flag.String("version", "", "version of rc image")
	flag.Parse()

	pod := &Pod{kubeClient}
	rc := &ReplicationController{kubeClient}

	switch *commandPtr {
	case "create-pod":
		pod.create(*namespacePtr)
	case "create-rc":
		rc.create(*namespacePtr)
	case "list-pods":
		pod.list(*namespacePtr)
	case "list-rcs":
		rc.list(*namespacePtr)
	case "delete-pod":
		pod.delete(*namespacePtr, *namePtr)
	case "delete-rc":
		rc.delete(*namespacePtr, *namePtr)
	case "update-rc":
		rc.update(*namespacePtr, *namePtr, *imagePtr, *versionPtr)
	default:
		fmt.Println("Can't understand command.")
	}
}

type Pod struct {
	kubeClient *kube.Client
}

func (p *Pod) create(namespace string) {
	ctx := context.TODO()

	podData, err := ioutil.ReadFile("pod.json")
	if err != nil {
		fmt.Printf("Error reading pod.\n%v", err)
		return
	}

	pod := api.Pod{}
	err = json.Unmarshal(podData, &pod)
	if err != nil {
		fmt.Printf("Error parsing pod.\n%v", err)
		return
	}

	pod.ObjectMeta.Namespace = namespace

	newPod, err := p.kubeClient.CreatePod(ctx, &pod)
	if err != nil {
		fmt.Printf("Error creating pod.\n%v", err)
		return
	}

	statuses, err := p.kubeClient.WatchPod(ctx, newPod.Namespace, newPod.Name, newPod.ResourceVersion)
	if err != nil {
		fmt.Printf("Error watching task pod.\n%v", err)
		return
	}

	for status := range statuses {
		podStatus := status.Pod.Status
		if podStatus.Phase == "Failed" {
			_ = p.kubeClient.DeletePod(ctx, newPod.Namespace, newPod.Name)
			fmt.Printf("Task pod failed.\n")
			return
		}
		if podStatus.Phase == "Succeeded" {
			fmt.Printf("Task pod succeeded.\n")
			return
		}
	}

	return
}

func (p *Pod) delete(namespace, name string) {
	err := p.kubeClient.DeletePod(context.TODO(), namespace, name)
	if err != nil {
		fmt.Printf("Error deleting pod: %v", err)
	}
}

func (p *Pod) list(namespace string) {
	pods, err := p.kubeClient.PodList(context.TODO(), namespace, "")
	if err != nil {
		fmt.Printf("Error listing pods: %v", err)
	}

	for _, pod := range pods {
		fmt.Println(pod.Name)
	}

}

type ReplicationController struct {
	kubeClient *kube.Client
}

func (rc *ReplicationController) create(namespace string) {
	ctx := context.TODO()

	rcData, err := ioutil.ReadFile("replication-controller.json")
	if err != nil {
		fmt.Printf("Error reading rc.\n%v", err)
		return
	}

	repc := api.ReplicationController{}
	err = json.Unmarshal(rcData, &repc)
	if err != nil {
		fmt.Printf("Error parsing rc.\n%v", err)
		return
	}

	repc.ObjectMeta.Namespace = namespace

	_, err = rc.kubeClient.CreateReplicationController(ctx, &repc)
	if err != nil {
		fmt.Printf("Error creating rc.\n%v", err)
		return
	}

	return
}

func (rc *ReplicationController) update(namespace, name, image, version string) {
	err := rc.kubeClient.UpdateReplicationControllerImage(context.TODO(), namespace, name, image, version)
	if err != nil {
		fmt.Printf("Error deleting rc: %v", err)
	}
}

func (rc *ReplicationController) delete(namespace, name string) {
	err := rc.kubeClient.DeleteReplicationController(context.TODO(), namespace, name)
	if err != nil {
		fmt.Printf("Error deleting rc: %v", err)
	}
}

func (rc *ReplicationController) list(namespace string) {
	rcs, err := rc.kubeClient.ReplicationControllerList(context.TODO(), namespace, "")
	if err != nil {
		fmt.Printf("Error listing pods: %v", err)
	}

	for _, rc := range rcs {
		fmt.Println(rc.Name)
	}

}
