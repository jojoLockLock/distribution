package manager

import (
	"context"
	"fmt"
	"log"

	pb "github.com/docker/distribution/registry/storage/manager/storage-path"
	"google.golang.org/grpc"
)

func GetDockerStoragePath(address string, host string, subpath string) (string, error) {

	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)

		finalError := fmt.Errorf("[ERROR] getDockerStoragePath did not connect: %v", err)

		fmt.Println(finalError)

		return "", finalError
	}

	defer conn.Close()

	c := pb.NewStoragePathClient(conn)

	r, err := c.GetDockerStoragePath(context.Background(), &pb.DockerStoragePathRequest{Host: host, SubPath: subpath})

	if err != nil {

		fmt.Println(err)

		return "", err
	}

	return r.Path, nil
}
