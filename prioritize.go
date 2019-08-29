package main

import (
	"k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
	"log"
	//"strings"
	//"strconv"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/client-go/rest"
	//"k8s.io/client-go/kubernetes"
	//"k8s.io/apimachinery/pkg/fields"
	//"k8s.io/apimachinery/pkg/labels"
	//"k8s.io/apimachinery/pkg/api/resource"
)

type Prioritize struct {
	Name string
	Func func(pod v1.Pod, nodes []v1.Node, mem int,gputype string) (*schedulerapi.HostPriorityList, error)
}


func (p Prioritize) Handler(args schedulerapi.ExtenderArgs) (*schedulerapi.HostPriorityList, error) {
	log.Printf("Entry Prioritize\n")
	/*
	a := args.Nodes.Items
	for _,b := range(a){
		//log.Printf("****************************\n")
		log.Printf("NodeGPU:%+v\n",b.ObjectMeta.Annotations["GPUInfo"])
		log.Printf("****************************\n")

		//log.Printf("NodeStatus:\n Config:%+v\n",b.Status.Config)

		//log.Printf("****************************\n")
	}*/

	//log.Printf("Get Pod Info:\n")

	pod := args.Pod.Spec
	container := pod.Containers
	//log.Printf("Container info :%+v\n",container)
	var gputype string
	gputype = args.Pod.ObjectMeta.Annotations["gputype"]

	var total int
	for _,i := range(container){
		resourcereq := i.Resources
		//log.Printf("Req:%+v\n",resourcereq.Limits)
		if val,ok := resourcereq.Limits["test/gpu"];ok{
			total += int(val.Value())
			log.Printf("Need Memory:%d\n",int(val.Value()))
		}else{
			log.Printf("Have no memory need\n")
		}

	}

	return p.Func(*args.Pod, args.Nodes.Items, total, gputype)
}
