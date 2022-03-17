package it

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"new_todo_project/server"
	"new_todo_project/server/config"
	"testing"
)

type e2eTestSuite struct {
	suite.Suite
}

func (suite *e2eTestSuite) SetupTest() {
}
func (suite *e2eTestSuite) SetupSuite() {
	loadConfig, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("panic test suite: %v", err))
	}
	server.Start()

}
func (suite *e2eTestSuite) TearDownSuite() {

}
func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(e2eTestSuite))
}
