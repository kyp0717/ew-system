package docker

import (
"context"
"fmt"

"github.com/docker/docker/api/types"
"github.com/docker/docker/api/types/container"
"github.com/docker/docker/client"
"github.com/docker/go-connections/nat"
)

func CreateNewContainer(image string) (string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}

	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: "8000",
	}
	containerPort, err := nat.NewPort("tcp", "80")
	if err != nil {
		panic("Unable to get the port")
	}

	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: image,
		},
		&container.HostConfig{
			PortBindings: portBinding,
		}, nil, "")
	if err != nil {
		panic(err)
	}
	
	cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	fmt.Printf("Container %s is started", cont.ID)
	return cont.ID, nil
}

func ListContainer() error {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
  
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
  
	if len(containers) > 0 {
		for _, container := range containers {
			fmt.Printf("Container ID: %s", container.ID)
		}
	} else {
		fmt.Println("There are no containers running")
	}
	return nil
}

func StopContainer(containerID string) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	
	err = cli.ContainerStop(context.Background(), containerID, nil)
	if err != nil {
		panic(err)
	}
	return err
}

func StartPostgres() {
    cli, err := client.NewEnvClient()
    if err != nil {
        panic(err)
    }

    ctx := context.Background()
    resp, err := cli.ContainerCreate(ctx, &container.Config{
        Image:        "postgres",
        ExposedPorts: nat.PortSet{"5432": struct{}{}},
    }, &container.HostConfig{
        PortBindings: map[nat.Port][]nat.PortBinding{nat.Port("5432"): {{HostIP: "127.0.0.1", HostPort: "5432"}}},
    }, nil, "mongo-go-cli")
    if err != nil {
        panic(err)
    }

    if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
        panic(err)
    }
}
