### About this sample:
This sample demonstrates how to structure a Command-Query Responsibility Segregation (CQRS) + Event Sourcing actor workflow that listens for commands on a signal channel. 

The workflow runs an activity to process the commands into events (which are both persisted in the workflow history). This activity represents an Aggregate Root from the Domain Driven Design (DDD) lexicon. 

The events and projections of the events are exposed through the workflow query interface.


### Steps to run this sample:
1) You need a Temporal service running. See details in README.md
2) Run the following command to start the worker
```
go run actor/worker/main.go
```
3) Run the following command to start the example
```
go run actor/starter/main.go
```
