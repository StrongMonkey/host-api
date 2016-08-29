package events

import (
	//"fmt"
	"github.com/fsouza/go-dockerclient"
	"os"
	"path"
	"fmt"
	"runtime"
	engineapi "github.com/docker/engine-api/client"
)

const (
	defaultUnixSocket = "unix:///var/run/docker.sock"
	defaultApiVersion = "1.18"
)

func NewDockerClient() (*docker.Client, error) {
	apiVersion := getenv("DOCKER_API_VERSION", defaultApiVersion)
	var endpoint string
	if runtime.GOOS == "windows" {
		apiVersion = "1.24"
		endpoint = fmt.Sprintf("tcp://%v:2375", os.Getenv("CATTLE_AGENT_IP"))
	} else {
		endpoint = defaultUnixSocket
	}

	if os.Getenv("CATTLE_DOCKER_USE_BOOT2DOCKER") == "true" {
		endpoint = os.Getenv("DOCKER_HOST")
		certPath := os.Getenv("DOCKER_CERT_PATH")
		tlsVerify := os.Getenv("DOCKER_TLS_VERIFY") != ""

		if tlsVerify && certPath != "" {
			cert := path.Join(certPath, "cert.pem")
			key := path.Join(certPath, "key.pem")
			ca := path.Join(certPath, "ca.pem")
			return docker.NewVersionedTLSClient(endpoint, cert, key, ca, apiVersion)
		}
	}

	return docker.NewVersionedClient(endpoint, apiVersion)
}

func DockerClient() (*engineapi.Client, error) {
	apiVersion := getenv("DOCKER_API_VERSION", defaultApiVersion)
	var endpoint string
	if runtime.GOOS == "windows" {
		apiVersion = "1.24"
		endpoint = fmt.Sprintf("tcp://%v:2375", os.Getenv("CATTLE_AGENT_IP"))
	} else {
		endpoint = defaultUnixSocket
	}
	return engineapi.NewClient(endpoint, apiVersion, nil, nil)
}

func getenv(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultVal
	}
	return val
}
