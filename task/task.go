package task

import "github.com/numkem/bamboo-nanny/agent"

type Task interface {
	Run(*agent.Agent) error
}
