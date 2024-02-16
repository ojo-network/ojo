package e2e

import (
	// "testing"

	"github.com/stretchr/testify/suite"

	"github.com/ojo-network/ojo/tests/e2e/orchestrator"
)

type IntegrationTestSuite struct {
	suite.Suite

	orchestrator *orchestrator.Orchestrator
}

// TODO: Make e2e work with rollkit
// func TestIntegrationTestSuite(t *testing.T) {
// 	suite.Run(t, new(IntegrationTestSuite))
// }

// func (s *IntegrationTestSuite) SetupSuite() {
// 	s.orchestrator = &orchestrator.Orchestrator{}
// 	s.orchestrator.InitResources(s.T())
// }

// func (s *IntegrationTestSuite) TearDownSuite() {
// 	s.orchestrator.TearDownResources(s.T())
// }
