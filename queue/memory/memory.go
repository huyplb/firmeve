package memory

import (
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/queue"
	"github.com/firmeve/firmeve/utils"
	"time"
)


type Memory struct {
	payload queue.Payloader
}

func NewMemory(config *config.Config) *Memory {
	return &Memory{
		payload: make(queue.Payloader),
	}
}

func (m *Memory) Push(jobName string, options ...utils.OptionFunc) {

	payload := queue.NewPayload(jobName,options...)

	if _, ok := m.payload[payload.QueueName]; !ok {
		m.payload[payload.QueueName] = make(chan *queue.Payload)
	}


	m.payload[payload.QueueName] <- payload
}

func (m *Memory) Pop(queueName string) <-chan *queue.Payload {
	return m.payload[queueName]
}

func (m *Memory) Later(delay time.Time, jobName string, options ...utils.OptionFunc) {
}
