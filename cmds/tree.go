package cmds

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type Tree struct {
	Clusters []*Cluster
}

type Cluster struct {
	ARN      string
	Name     string
	Tasks    []*Task
	Services []*Service
}

type Task struct {
	ARN        string
	FamilyName string
	Containers []*Container
}

type Service struct {
	ARN   string
	Name  string
	Tasks []*Task
}

type Container struct {
	ARN  string
	Name string
}

func GenerateTree(ctx context.Context) *Tree {
	awsConfig, configLoadErr := config.LoadDefaultConfig(context.Background())
	if configLoadErr != nil {
		panic(configLoadErr)
	}
	if len(awsConfig.Region) == 0 {
		panic("no aws region found")
	}

	// ECS Client
	client := ecs.NewFromConfig(awsConfig)

	// Tree instance
	tree := &Tree{}

	// List Clusters
	rawClusters, er1 := client.ListClusters(ctx, nil)
	if er1 != nil {
		return nil
	}
	// Describe all clusters
	detailedClusters, er2 := client.DescribeClusters(ctx, &ecs.DescribeClustersInput{
		Clusters: rawClusters.ClusterArns,
	})
	if er2 != nil {
		return nil
	}
	tree.Clusters = make([]*Cluster, len(detailedClusters.Clusters))
	for i, cluster := range detailedClusters.Clusters {
		tree.Clusters[i] = &Cluster{
			ARN:  *cluster.ClusterArn,
			Name: *cluster.ClusterName,
		}
		rawServices, er3 := client.ListServices(ctx, &ecs.ListServicesInput{
			Cluster: cluster.ClusterArn,
		})
		if er3 != nil {
			return nil
		}
		if len(rawServices.ServiceArns) == 0 {
			tree.Clusters[i].Tasks = getAllTasks(ctx, client, *cluster.ClusterArn, nil)
			continue
		}
		detailedServices, er4 := client.DescribeServices(ctx, &ecs.DescribeServicesInput{
			Services: rawServices.ServiceArns,
			Cluster:  cluster.ClusterArn,
		})
		if er4 != nil {
			return nil
		}
		services := make([]*Service, len(detailedServices.Services))
		for i, service := range detailedServices.Services {
			services[i] = &Service{
				ARN:   *service.ServiceArn,
				Name:  *service.ServiceName,
				Tasks: getAllTasks(ctx, client, *cluster.ClusterArn, service.ServiceName),
			}
		}
		tree.Clusters[i].Services = services
	}
	return tree
}

func getAllTasks(ctx context.Context, client *ecs.Client, clusterArn string, serviceName *string) []*Task {
	rawTasks, er1 := client.ListTasks(ctx, &ecs.ListTasksInput{
		Cluster:       &clusterArn,
		DesiredStatus: types.DesiredStatusRunning,
		ServiceName:   serviceName,
	})
	if er1 != nil {
		return nil
	}
	if len(rawTasks.TaskArns) == 0 {
		return nil
	}
	detailedTasks, er2 := client.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Tasks:   rawTasks.TaskArns,
		Cluster: &clusterArn,
	})
	if er2 != nil {
		return nil
	}
	tasks := make([]*Task, len(detailedTasks.Tasks))
	for i, task := range detailedTasks.Tasks {
		taskDefinition, er3 := client.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{
			TaskDefinition: task.TaskDefinitionArn,
		})
		if er3 != nil {
			return nil
		}
		tasks[i] = &Task{
			ARN:        *task.TaskArn,
			FamilyName: *taskDefinition.TaskDefinition.Family,
			Containers: getAllContainers(ctx, client, *task.TaskArn, clusterArn),
		}
	}
	return tasks
}

func getAllContainers(ctx context.Context, client *ecs.Client, taskArn, clusterArn string) []*Container {
	detailedTasks, er1 := client.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Tasks:   []string{taskArn},
		Cluster: &clusterArn,
	})
	if er1 != nil {
		return nil
	}
	containers := make([]*Container, len(detailedTasks.Tasks[0].Containers))
	for i, container := range detailedTasks.Tasks[0].Containers {
		containers[i] = &Container{
			ARN:  *container.ContainerArn,
			Name: *container.Name,
		}
	}
	return containers
}
