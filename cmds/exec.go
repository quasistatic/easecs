package cmds

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Exec BaseCmd

func NewExecCommand(args []string) *Exec {
	return &Exec{
		name: ExecCommand,
		args: args,
	}
}

func (cmd *Exec) Execute() {
	if len(cmd.args) < 2 || len(cmd.args) > 3 {
		fmt.Println("[ERROR] invalid args")
		return
	}
	cluster := cmd.args[0]
	service := cmd.args[1]
	container := ""
	if len(cmd.args) == 3 {
		container = cmd.args[2]
	}

	tree := GenerateTree(context.Background())

	var pickedC *Cluster
	for _, c := range tree.Clusters {
		if match(c.Name, cluster) {
			pickedC = c
			break
		}
	}
	if pickedC == nil {
		fmt.Println("[ERROR] could not find cluster")
		return
	}

	var pickedS *Service
	for _, s := range pickedC.Services {
		if match(s.Name, service) {
			pickedS = s
			break
		}
	}
	if pickedS == nil {
		fmt.Println("[ERROR] could not find service")
		return
	}

	if len(pickedS.Tasks) == 0 {
		fmt.Println("[ERROR] could not find tasks")
		return
	}
	var pickedT *Task = pickedS.Tasks[0]

	var pickedCtr *Container = pickedT.Containers[0]
	if len(container) > 0 {
		for _, ctr := range pickedT.Containers {
			if match(ctr.Name, container) {
				pickedCtr = ctr
				break
			}
		}
	}

	awsCommand := fmt.Sprintf(`ecs execute-command --cluster %s --task %s --container %s --interactive --command "/bin/bash"`, pickedC.Name, pickedT.ARN, pickedCtr.Name)

	command := exec.Command("aws", strings.Split(awsCommand, " ")...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		panic(err)
	}
}

func match(haystack, needle string) bool {
	matched, _ := regexp.MatchString(needle, haystack)
	return matched
}
