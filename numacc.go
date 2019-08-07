package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type (
	// Config information.
	Config struct {
		containerID string
	}
)

var config Config

func main() {
	app := cli.NewApp()
	app.Name = "NUMACC"
	app.Usage = "NUMA Checker for Containers"
	app.Author = "Jace Liang"
	app.Email = "b436412@gmail.com"
	app.Action = run
	app.Version = "0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "container-id,cid",
			Usage: "Container ID",
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	config = Config{
		containerID: c.String("cid"),
	}
	return numacc()
}

/*
Main Function of NUMACC
 */
func numacc() error {
	pidCpuMap := initPidMapByContainerID(config.containerID)
	fillCpuIdByPidMap(pidCpuMap)
	fmt.Println("Process and CPU for Container ",config.containerID,)
	fmt.Println("PID\tCurrentCpu\tCpuAffinity")
	tasksetSlice := checkPidMapTaskset(pidCpuMap)
	for k,v := range pidCpuMap {
		fmt.Println(k, "\t", v, "\t", tasksetSlice[k])
	}

	fmt.Println("The NIC NUMA for container ", config.containerID)
	for k,v := range getNicNumaByContainerId(config.containerID) {
		fmt.Print(k, "\t")
		if v == ""{
			fmt.Println("N/A")
		} else {
			fmt.Println(v)
		}
	}

	return nil
}

func getCpuIDByPid(id string) string {
	cmd := exec.Command("ps", "-o", "psr" ,"-p", id)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	cpuID := strings.Split(string(out),"\n")[1]
	cpuID = strings.TrimPrefix(cpuID, "  ")
	cpuID = strings.TrimSuffix(cpuID, " ")
	return cpuID
}

func initPidMapByContainerID(id string) map[string]string {
	pidCpuMap := make(map[string]string)
	cmd := exec.Command("docker", "top", id)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	for counter := 1; counter < len(strings.Split(string(out), "\n"))-1; counter++ {
		tmpId := strings.Split(string(out),"\n")[counter]
		re, _ := regexp.Compile("\\s(\\w+)\\S")

		pId := re.FindStringSubmatch(tmpId)
		pId[0] = strings.TrimPrefix(pId[0], " ")
		pId[0] = strings.TrimSuffix(pId[0], " ")
		pidCpuMap[pId[0]] = "0"
	}
	return pidCpuMap
}

func fillCpuIdByPidMap(inputMap map[string]string){
	for k := range inputMap {
		inputMap[k] = getCpuIDByPid(k)
	}
}

func checkPidMapTaskset(inputMap map[string]string) map[string]string{
	tasksetMap := make(map[string]string)
	for k := range inputMap {
		cmd := exec.Command("taskset", "-cp", k)
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		tmp := strings.Split(string(out)," ")
		affinity := tmp[len(tmp)-1]
		tasksetMap[k]  = strings.TrimSuffix(string(affinity), "\n")
	}
	return tasksetMap
}

func getNicNumaByContainerId(cid string) map[string]string {
	//Get all nic in container
	cmd := exec.Command("docker", "exec", cid, "ls" , "/sys/class/net")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(out))
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	allNic := strings.Split(string(out),"\n")
	if len(allNic) > 0 && allNic[len(allNic)-1] == ""{
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
