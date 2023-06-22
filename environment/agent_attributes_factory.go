package environment

const (
	DefaultEC2AgentStartCommand = "sudo /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl -a fetch-config -m ec2 -s -c "
)

type agentAttributesProvider struct {
	AgentStartCommand string
}

var agentProvider = &agentAttributesProvider{
	AgentStartCommand: DefaultEC2AgentStartCommand,
}

func SetAgentStartCommandAttribute(agentStartCommand string) {
	agentProvider.AgentStartCommand = agentStartCommand
}

func GetAgentAttributeProviderInstance() *agentAttributesProvider {
	return agentProvider
}
