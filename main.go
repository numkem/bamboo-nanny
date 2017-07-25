package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/numkem/bamboo-nanny/agent"
	"github.com/numkem/bamboo-nanny/task"
)

var tasks = map[string]task.Task{
	"cleanup job folder": task.XMLDataCleanup{},
}

type parameters struct {
	Interval       int
	BambooUsername string
	BambooPassword string
	Statuses       chan string
	Errors         chan error
}

func agentWorker(a *agent.Agent, param *parameters) {
	for {
		status, err := agent.GetAgentStatus(a, param.BambooUsername, param.BambooPassword)
		if err != nil {
			param.Errors <- fmt.Errorf("%s: Couldn't get the status of agent: %v", a.Name, err)
		}

		if status == nil {
			param.Errors <- fmt.Errorf("%s: Couldn't find agent as running on the server", a.Name)
		} else { // TODO: Fix this, this is pretty ugly
			if !status.Busy {
				for name, t := range tasks {
					param.Statuses <- fmt.Sprintf("%s: Starting task %s...", a.Name, name)
					err := t.Run(a)
					if err != nil {
						param.Errors <- fmt.Errorf("%s: Error while running task %s: %v", a.Name, name, err)
					}
					param.Statuses <- fmt.Sprintf("%s: Done with task %s...", a.Name, name)
				}
			}
		}

		time.Sleep(time.Duration(param.Interval) * time.Minute)
	}
}

func printStatuses(param *parameters) {
	for {
		select {
		case status := <-param.Statuses:
			log.Info(status)
		}
	}
}

func printError(param *parameters) {
	for {
		select {
		case err := <-param.Errors:
			log.Error(err)
		}
	}
}

func main() {
	p := parameters{}
	p.Statuses = make(chan string, 10)
	p.Errors = make(chan error, 10)

	flag.IntVar(&p.Interval, "interval", 5, "time in minute to wait before each run")
	flag.StringVar(&p.BambooUsername, "username", "", "Bamboo username to authenticate with API")
	flag.StringVar(&p.BambooPassword, "password", "", "Bamboo password to authenticate with API")
	flag.Parse()

	if len(os.Args) == 1 {
		log.Fatal("One or more agent directory needs to be provided")
	}

	if p.BambooUsername == "" || p.BambooPassword == "" {
		log.Fatal("Bamboo username/password are required")
	}

	go printError(&p)
	go printStatuses(&p)
	for _, dir := range flag.Args() {
		a, err := agent.NewAgentFromDir(dir)
		if err != nil {
			log.Errorf("Error while parsing agent directory %s: %v", dir, err)
		}

		log.Infof("Starting worker for agent %s", a.Name)
		go agentWorker(a, &p)
	}

	wait := make(chan bool)
	<-wait
}
