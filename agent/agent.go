package agent

import (
	"encoding/json"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/events"
	"github.com/docker/engine-api/types/filters"
	"golang.org/x/net/context"
)

func StartAgent(c *cli.Context) {
	cli, err := getDockerClient()
	if err != nil {
		logrus.Fatalf("Could not get Docker client: %s", err)
	}

	filterArgs := filters.NewArgs()
	filterArgs.Add("event", "start")

	eventOptions := types.EventsOptions{
		Filters: filterArgs,
	}

	eventsResp, err := cli.Events(context.Background(), eventOptions)
	if err != nil {
		logrus.Fatal(err)
	}
	defer eventsResp.Close()

	handler, err := NewMessageHandler(map[string]interface{}{
		"metadata-url": c.String("metadata-url"),
		"bridge-url":   c.String("bridge-url") + "/v1/message",
	})
	if err != nil {
		logrus.Fatalf("Error: %s", err)
	}

	logrus.Info("Entering event listening Loop")
	d := json.NewDecoder(eventsResp)
	var msg events.Message
	for {
		d.Decode(&msg)

		// For now... will need to add some throttling at some point.
		go handler.Handle(&msg)
	}

	os.Exit(0)
}

func getDockerClient() (*client.Client, error) {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	return client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
}