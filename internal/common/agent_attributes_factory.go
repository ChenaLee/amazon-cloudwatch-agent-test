package common

import (
	"github.com/aws/amazon-cloudwatch-agent-test/environment"
)

type agentAttributesProvider struct {
	AgentStartCommand string
}

var agentprovider *agentAttributesProvider = &agentAttributesProvider{
	AgentStartCommand: DefaultEC2AgentStartCommand,
}

func SetAgentAttributesMetadata(testMetadata *environment.MetaData) {
	agentprovider.AgentStartCommand = testMetadata.AgentStartCommand
}

func GetAgentAttributeProviderInstance() *agentAttributesProvider {
	return agentprovider
}
