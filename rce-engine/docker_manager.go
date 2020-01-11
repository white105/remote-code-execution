package dockermanager

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	// "github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var languageImages = map[string]string{
	"go": "docker.io/library/golang:1.13.5-alpine",
	"py": "docker.io/library/python:3.8.0-alpine",
	"js": "docker.io/library/node:13.6-alpine",
}

// Docker manager

type DockerManager struct {
	Ctx    context.Context
	Client *client.Client
}

// Get container name of user running it
func (d *DockerManager) GetContainerName() string {
	return uuid.New().String()
}

//Pull Image if Image is not exist
func (d *DockerManager) PullImage(language string) bool {
	imageName := languageImages[language]
	resp, err := d.Client.ImagePull(d.Ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		log.Fatal("Error pulling image: ", err)
		return false
	}
	io.Copy(os.Stdout, resp)
	return true
}

func (d *DockerManager) FilterImage(language string) {
	// imageName := languageImages[language]
	// initFilter :=
	// resp, err := d.Client.ImageList(d.Ctx, types.ImageListOptions{
	// 	Filters: params,
	// })
	// if err != nil {
	// 	log.Fatal("Cannot search image ", err)
	// }
	// log.Println(resp)

}

func (d *DockerManager) CreateContainer(imageID string) (bool, string, string) {
	currentWorkDir, _ := os.Getwd()
	resp, err := d.Client.ContainerCreate(d.Ctx, &container.Config{
		Image:       imageID,
		Tty:         true,
		AttachStdin: true,
		WorkingDir:  "/project",
	}, &container.HostConfig{
		Binds: []string{fmt.Sprintf("%s/project:/project", currentWorkDir)},
	}, nil, d.GetContainerName())

	if err != nil {
		log.Fatal("Error create container ", err)
		return false, "", ""
	}
	containerID := resp.ID
	if errStartContainer := d.Client.ContainerStart(d.Ctx, containerID, types.ContainerStartOptions{}); errStartContainer != nil {
		d.RemoveContainer(containerID)

		log.Fatal("Error start container", errStartContainer)
		return false, d.GetContainerName(), containerID
	}
	return true, d.GetContainerName(), containerID
}

func (d *DockerManager) ExecuteCommand(containerName string, commands []string) (bool, string, string) {
	execResp, err := d.Client.ContainerExecCreate(d.Ctx, containerName, types.ExecConfig{
		Cmd:          commands,
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		log.Fatal("Error create execute the commands", err)
	}
	stream, errExecCmd := d.Client.ContainerExecAttach(d.Ctx, execResp.ID, types.ExecConfig{})
	if errExecCmd != nil {
		log.Fatal("Error execute commands", err)
		return false, execResp.ID, ""
	}
	message, errRead := ioutil.ReadAll(stream.Reader)
	defer stream.Close()
	if errRead != nil && errRead != io.EOF {
		log.Fatal("Error to read the result", err)
	}
	result := string(message)
	defer stream.CloseWrite()
	return true, execResp.ID, result
}

func (d *DockerManager) GetContainerLogs(containerID string) {
	logs, err := d.Client.ContainerLogs(d.Ctx, containerID, types.ContainerLogsOptions{
		Timestamps: true,
		ShowStderr: true,
		ShowStdout: true,
	})
	if err != nil {
		log.Fatalf("Get error %s when get log from %s", err, containerID)
	}
	defer logs.Close()
}

func (d *DockerManager) RemoveContainer(containerID string) error {
	err := d.Client.ContainerRemove(d.Ctx, containerID, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		log.Fatal("Cannot remove container", err)
		return err
	}
	log.Println("Removed container ... ")
	return nil
}

func main() {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		log.Fatal("Cannot start docker client", err)
		return
	}
	dockerManager := &DockerManager{
		Ctx:    context.Background(),
		Client: dockerClient,
	}
	// success := dockerManager.PullImage(languageImages["go"])
	// if !success {
	// 	log.Fatal("Cannot pull image")
	// }
	errRemove := dockerManager.RemoveContainer("2e09fbea3e65c43b1098fdfdae0fe324f3cf2485f0d385fca46455822e7bbd0e")
	if errRemove != nil {
		log.Fatal(errRemove)
	}
}
