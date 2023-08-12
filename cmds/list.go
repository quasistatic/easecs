package cmds

import (
	"context"
	"fmt"
)

type List BaseCmd

func NewListCommand(args []string) *List {
	return &List{
		name: ListCommand,
		args: args,
	}
}

func (d *List) Execute() {
	tree := GenerateTree(context.Background())

	fmt.Println("[clusters]")
	for _, c := range tree.Clusters {
		fmt.Printf("\t%s\n", c.Name)
		if len(c.Services) == 0 {
			continue
		}
		fmt.Printf("\t[services]\n")
		for _, s := range c.Services {
			fmt.Printf("\t\t%s\n", s.Name)
			if len(s.Tasks) == 0 {
				continue
			}
			fmt.Printf("\t\t[tasks]\n")
			for _, t := range s.Tasks {
				fmt.Printf("\t\t\t%s\n", t.ARN)
				fmt.Printf("\t\t\t[containers]\n")
				for _, ctr := range t.Containers {
					fmt.Printf("\t\t\t\t%s\n", ctr.Name)
				}
			}
		}
	}
}
