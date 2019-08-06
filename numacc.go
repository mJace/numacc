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

func numacc() error {
	fmt.Println("Container ID:", config.containerID)

	log.Println("initPidMapByContainerId")
	pidCpuMap := initPidMapByContainerID(config.containerID)
	//fmt.Println(pidCpuMap)
	//log.Println("fill cpu id by pid map")
	fillCpuIdByPidMap(pidCpuMap)
	fmt.Println(pidCpuMap)
	return nil
}

func getCpuIDByPid(id string) string {
	cmd := exec.Command("ps", "-o", "psr" ,"-p", id)

	out, err := cmd.CombinedOutput()
	if err != nil {
		//fmt.Println(out)
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	//fmt.Printf("combined out:\n%s\n", string(out))

	cpuID := strings.Split(string(out),"\n")[1]
	cpuID = strings.TrimPrefix(cpuID, "  ")
	cpuID = strings.TrimSuffix(cpuID, " ")
	return cpuID
}

func initPidMapByContainerID(id string) map[string]string {
	log.Println("init Pid Map by container ID...")
	pidCpuMap := make(map[string]string)

	cmd := exec.Command("docker", "top", config.containerID)
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