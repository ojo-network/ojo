//go:build norace
// +build norace

package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	ojoapp "github.com/ojo-network/ojo/app"
)

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
