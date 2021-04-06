package main

import (
	"context"
	"encoding/json"
	"go.temporal.io/sdk/client"
	"log"

	"github.com/temporalio/samples-go/actor"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	upvote(c, "dev", "some-entity-2")
}

func upvote(c client.Client, userId string, entityId string) {
	workflowId := entityId
	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowId,
		TaskQueue: actor.TaskQueueName,
	}
	command := actor.Command{Name: "upvote", UserId: userId}
	commandJson, err := json.Marshal(command)
	if err != nil {
		log.Fatalln("unable to serialize command", err)
	}
	log.Println(commandJson)
	c.SignalWithStartWorkflow(context.Background(), workflowId, actor.CommandSignalName, string(commandJson), workflowOptions, actor.Workflow)
}
