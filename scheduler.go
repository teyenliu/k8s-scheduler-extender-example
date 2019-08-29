package main

import (
        //"k8s.io/api/core/v1"
        //schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
        "log"
        "strings"
        "strconv"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        "k8s.io/client-go/rest"
        "k8s.io/client-go/kubernetes"
        "k8s.io/apimachinery/pkg/fields"
        "k8s.io/apimachinery/pkg/labels"
        //"k8s.io/apimachinery/pkg/api/resource"
)

const INT_MAX = int(^uint(0) >> 1)

type Podfig struct{
        state string
        annotation string
        mem uint
}

func nodescheduler2(nodegpu map[string]uint,mem int)int {
	log.Printf("NodeInfo:%+v\nNeedMeM:%d\n",nodegpu,mem)
	flag := 0
	var lowcapacity uint
	var score int
	for _,j := range(nodegpu){
		if flag==0{
			lowcapacity = j
			flag = 1

		}
		if j > uint(mem) && j <= lowcapacity && j >= 0{
			score = int( uint(mem) - j )

		}
	}
	return score

}


func nodescheduler(originalgpuflag map[string]uint,nodegpu map[string]uint,mem int,gputype string)int {
        //log.Printf("NodeInfo:%+v\nNeedMeM:%d\n",nodegpu,mem)
        var lowcapacity uint
	score := INT_MAX
	flag := 0
	if gputype == "memory"{
		for _,j := range(nodegpu){
			if flag==0 && j >= uint(mem) {
				lowcapacity = j
				flag = 1

			}
			if flag == 1 &&  j >= uint(mem) && j <= lowcapacity && j >= 0{
				score = int( uint(mem) - j )
			}
		}
		log.Printf("Memory Score:%d\n",score)
		if flag!= 0{
			return score
		}else{
			return 1
		}
	}else{
		score = len(originalgpuflag)
		log.Printf("Count Score:%d\n",score)
		for i,j := range(nodegpu){
			//log.Printf("Original:%d Final:%d\n",originalgpuflag[i],j)
			if(originalgpuflag[i]) != j{
				score -= 1
			}
		}
		if score >= mem{
			return score
		}else{
			return -1
		}
	}
}

func nodegputable(nodegpuinfo string)map[string]uint{
	gpuinfo := make(map[string]uint)
	idcount := strings.Split(nodegpuinfo,",")
	for _,i := range(idcount){
		h := strings.Split(i,":")
		allmem,_ := strconv.Atoi(h[1])
		gpuinfo[h[0]] = uint(allmem)
		//log.Printf("%s\n",i)
	}
	return gpuinfo
}


func Podinfo(nodename string,podmem int,nodegpuinfo string,gputype string)int{

	log.Printf("NodeGPUInfo:%s\n",nodegpuinfo)
        //gpuinfo := make(map[string]uint)
        //infoflag := int(0)

        config, err := rest.InClusterConfig()
        if err != nil {
                panic(err.Error())
        }

        clientset, err := kubernetes.NewForConfig(config)
        if err != nil {
                panic(err.Error())
        }
        selector := fields.SelectorFromSet(fields.Set{"spec.nodeName": nodename})
        pods, err := clientset.CoreV1().Pods("default").List(metav1.ListOptions{
                FieldSelector: selector.String(),
                LabelSelector: labels.Everything().String(),
        })

         if err != nil {
                        panic(err.Error())
                }

                log.Printf("There are %d pods in the cluster\n", len(pods.Items))
/*
        for _, pod := range pods.Items {
                log.Printf("Name: %s, Status: %s\n", pod.ObjectMeta.Name, pod.Status.Phase )

                if len(pod.ObjectMeta.Annotations["GPUAllInfo"]) != 0 && infoflag == 0 {
                        idcount := strings.Split(pod.ObjectMeta.Annotations["GPUAllInfo"],",")
                        for _,i := range(idcount){
                                h := strings.Split(i,":")
                                allmem,_ := strconv.Atoi(h[1])
                                gpuinfo[h[0]] = uint(allmem)
                                //log.Printf("%s\n",i)
                        }
                infoflag = 1
                }
		if  pod.Status.Phase == "Running" && len(pod.ObjectMeta.Annotations["GPUID"]) != 0 &&  len(pod.ObjectMeta.Annotations["GPUMEM"]) != 0 {
                        log.Printf("Pod State:%s\n",pod.Status.Phase)
                        usedmem,_ := strconv.Atoi(pod.ObjectMeta.Annotations["GPUMEM"])
                        gpuinfo[pod.ObjectMeta.Annotations["GPUID"]] -= uint(usedmem)
                }

        }
	for k,m := range(gpuinfo){
                log.Printf("GPUID:%s Mem:%d\n",k,m)
        }

        finalscore := nodescheduler(gpuinfo,podmem)
        return finalscore

*/
//***************************************************************
	//gpuflag := make(map[string]uint)
	//gpuflag := gputable(nodegpuinfo)
	originalgpuflag := nodegputable(nodegpuinfo)
	gpuflag := nodegputable(nodegpuinfo)
	for _, pod := range pods.Items {
		if len(pod.ObjectMeta.Annotations["GPUID"]) != 0 && len(pod.ObjectMeta.Annotations["GPUMEM"]) != 0 && pod.Status.Phase == "Running"{
			reqtype := pod.ObjectMeta.Annotations["gputype"]
			useid := pod.ObjectMeta.Annotations["GPUID"]
			usemem := pod.ObjectMeta.Annotations["GPUMEM"]
			if reqtype == "memory"{
				allmem,_ := strconv.Atoi(usemem)
				gpuflag[useid] -= uint(allmem)

			}else{
				idlist := strings.Split(useid,",")
				for _,id := range(idlist){
					gpuflag[id] = uint(0)
				}
			}
		}
	}
	/*
        for k,m := range(gpuflag){
                log.Printf("GPUID:%s Mem:%d\n",k,m)
        }*/


	finalscore := nodescheduler(originalgpuflag,gpuflag,podmem,gputype)
        return finalscore
//***************************************************************
}

