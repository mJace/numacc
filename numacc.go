package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/exec"
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

	cmd := exec.Command("ps", "-o", "psr" ,"-p", config.containerID)
	fmt.Println(cmd)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))

	cpuID := strings.Split(string(out),"\n")[1]
	fmt.Println("cpu ID : ", cpuID)

	return nil
}