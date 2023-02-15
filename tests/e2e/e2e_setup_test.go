package e2e

import (
	"testing"

	"github.com/ojo-network/ojo/client"
	"github.com/ojo-network/ojo/tests/e2e/orchestrator"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite

	orchestrator orchestrator.Orchestrator
	ojoClient    *client.OjoClient
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.orchestrator = orchestrator.Orchestrator{}
	s.T().Log("---> initializing docker resources")
	s.Require().NoError(s.orchestrator.InitDockerResources(s.T()))
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("---> tearing down")
	s.Require().NoError(s.orchestrator.TearDownDockerResources())
}
