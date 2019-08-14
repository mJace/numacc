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
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// cidCmd represents the cid command
var cidCmd = &cobra.Command{
	Use:   "cid <container_id>",
	Short: "check container's numa configuration by container ID",
	Long: "check container's numa configuration by container ID \n" +
		"Eg. \n" +
		"numacc cid 1234ABCD",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		numacc(args[0])
	},
}

func init() {
	RootCmd.AddCommand(cidCmd)
}

func numacc(cid string) error {
	pidCpuMap := initPidMapByContainerID(cid)
	fillCpuIdByPidMap(pidCpuMap)
	fmt.Println("Process and CPU for Container ", cid)
	fmt.Println("PID\tCurrentCpu\tCpuAffinity")
	tasksetSlice := checkPidMapTaskset(pidCpuMap)
	for k, v := range pidCpuMap {
		fmt.Println(k, "\t", v, "\t", tasksetSlice[k])
	}

	fmt.Println("The NIC NUMA for container ", cid)
	for k, v := range getNicNumaByContainerId(cid) {
		fmt.Print(k, "\t")
		if v == "" {
			fmt.Println("N/A")
		} else {
			fmt.Println(v)
		}
	}

	return nil
}

func getCpuIDByPid(id string) string {
	cmd := exec.Command("ps", "-o", "psr", "-p", id)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	cpuID := strings.Split(string(out), "\n")[1]
	cpuID = strings.TrimPrefix(cpuID, "  ")
	cpuID = strings.TrimSuffix(cpuID, " ")
	return cpuID
}

func initPidMapByContainerID(id string) map[string]string {
	pidCpuMap := make(map[string]string)
	cmd := exec.Command("docker", "top", id)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf(string(out))
	}
	for counter := 1; counter < len(strings.Split(string(out), "\n"))-1; counter++ {
		tmpId := strings.Split(string(out), "\n")[counter]
		re, _ := regexp.Compile("\\s(\\w+)\\S")

		pId := re.FindStringSubmatch(tmpId)
		pId[0] = strings.TrimPrefix(pId[0], " ")
		pId[0] = strings.TrimSuffix(pId[0], " ")
		pidCpuMap[pId[0]] = "0"
	}
	return pidCpuMap
}

func fillCpuIdByPidMap(inputMap map[string]string) {
	for k := range inputMap {
		inputMap[k] = getCpuIDByPid(k)
	}
}

func checkPidMapTaskset(inputMap map[string]string) map[string]string {
	tasksetMap := make(map[string]string)
	for k := range inputMap {
		cmd := exec.Command("taskset", "-cp", k)
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		tmp := strings.Split(string(out), " ")
		affinity := tmp[len(tmp)-1]
		tasksetMap[k] = strings.TrimSuffix(string(affinity), "\n")
	}
	return tasksetMap
}

func getNicNumaByContainerId(cid string) map[string]string {
	//Get all nic in container
	cmd := exec.Command("docker", "exec", cid, "ls", "/sys/class/net")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(out))
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	allNic := strings.Split(string(out), "\n")
	if len(allNic) > 0 && allNic[len(allNic)-1] == "" {
		allNic = allNic[:len(allNic)-1]
	}

	nicCpuMap := make(map[string]string)
	for i := range allNic {
		cmd := exec.Command("docker", "exec", cid, "cat", "/sys/class/net/"+allNic[i]+"/device/numa_node")
		out, _ := cmd.CombinedOutput()
		if strings.Contains(string(out), "No such file or directory") {
			nicCpuMap[allNic[i]] = ""
		} else {
			nicCpuMap[allNic[i]] = strings.TrimSuffix(string(out), "\n")
		}
	}
	return nicCpuMap
}
