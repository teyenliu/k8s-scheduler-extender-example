package main

import (
	"k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
	"log"
)

type Predicate struct {
	Name string
	Func func(pod v1.Pod, node v1.Node, mem int,gputype string) (bool, error)
}

func (p Predicate) Handler(args schedulerapi.ExtenderArgs) *schedulerapi.ExtenderFilterResult {
	pod := args.Pod
	canSchedule := make([]v1.Node, 0, len(args.Nodes.Items))
	canNotSchedule := make(map[string]string)

	podmem := pod.Spec
        container := podmem.Containers
        var gputype string
        gputype = pod.ObjectMeta.Annotations["gputype"]

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


	for _, node := range args.Nodes.Items {
		result, err := p.Func(*pod, node,total,gputype)
		if err != nil {
			canNotSchedule[node.Name] = err.Error()
		} else {
			if result {
				canSchedule = append(canSchedule, node)
				canSchedule = append(canSchedule, node)
			}
		}
	}

	result := schedulerapi.ExtenderFilterResult{
		Nodes: &v1.NodeList{
			Items: canSchedule,
		},
		FailedNodes: canNotSchedule,
		Error:       "",
	}

	return &result
}
