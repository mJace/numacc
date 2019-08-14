// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"flag"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// podCmd represents the pod command
var podCmd = &cobra.Command{
	Use:   "pod <podName>",
	Short: "check pod's container numa configuration by pod name",
	Long: "check pod's container numa configuration by pod name \n" +
		"Eg. \n" +
		"numacc pod podtest",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		podNumacc(args[0])
	},
}

func podNumacc(podName string) {
	containerID := getContainerIdByPodName(podName)
	numacc(containerID)
}

func getContainerIdByPodName(podName string) string {
	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", "/root/.kube/config", "absolute path to the kubeconfig file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pod, _ := clientset.CoreV1().Pods("kube-system").Get("etcd-minikube", metav1.GetOptions{})
	containerID := getFieldString(&pod.Status.ContainerStatuses[0], "ContainerID")
	containerID = strings.Replace(containerID, "docker://", "", -1)
	return containerID

}

func getFieldString(e *v1.ContainerStatus, field string) string {
	r := reflect.ValueOf(e)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}

func init() {
	RootCmd.AddCommand(podCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// podCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// podCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
