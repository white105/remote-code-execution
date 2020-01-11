package services

import (
	"context"
	// "fmt"
	"github.com/docker/docker/client"
	"log"
	"os"
	dockerManager "remote-code-execution/rce-engine"
)

func CompileCode(language, fileName string) {

	dockerClient, err := client.NewEnvClient()
	if err != nil {
		log.Fatal("Cannot create Docker client ", err)
	}
	docker := dockerManager.DockerManager{
		Ctx:    context.Background(),
		Client: dockerClient,
	}
	// docker.PullImage(language)
	_, _, containerID := docker.CreateContainer("59acf2b3028c")
	defer func() {
		deleteFile(fileName)
		removeContainer(docker, containerID)
	}()
	success, respID, result := docker.ExecuteCommand(containerID, []string{"python", fileName})
	log.Println(success, respID, result)
}

func removeContainer(docker dockerManager.DockerManager, containerID string) error {
	err := docker.RemoveContainer(containerID)
	if err != nil {
		log.Fatal("Cannot remove container")
	}
	return err
}

func deleteFile(fileName string) error {
	err := os.Remove(fileName)
	if err != nil {
		log.Fatalf("Cannot delete file name %s", fileName)
	}
	return err
}
