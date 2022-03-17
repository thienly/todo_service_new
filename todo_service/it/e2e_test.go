package it

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"net/http"
	"new_todo_project/server"
	"new_todo_project/server/config"
	"new_todo_project/server/handlers"
	"testing"
)

var configDir = flag.String("cfgDir", "", "configuration directory")

type e2eTestSuite struct {
	suite.Suite
}

func (suite *e2eTestSuite) SetupTest() {
}

func (suite *e2eTestSuite) SetupSuite() {
	loadConfig, err := config.LoadConfig(*configDir)
	if err != nil {
		panic(fmt.Sprintf("panic test suite: %v", err))
	}
	r := make(chan bool)
	newService := &server.Server{AppConfig: loadConfig,
		ReadyChan: r}
	go newService.Start()
	<-newService.ReadyChan
}

func (suite *e2eTestSuite) TearDownSuite() {
	ctx := context.Background()
	loadConfig, _ := config.LoadConfig(*configDir)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(loadConfig.Db.Conn))
	client.Database("todo").Drop(context.Background())
}

//TestUserRegistration is user registration endpoint to register new user.
func (suite *e2eTestSuite) TestUserRegistration() {
	client := http.DefaultClient
	body := handlers.UserRegisterRequest{
		Name:     "test1",
		Email:    "test1@gmail.com",
		Password: "abc",
	}
	marshal, _ := json.Marshal(body)
	resBody, err := client.Post("http://localhost:8080/users/register", "application/json", bytes.NewReader(marshal))
	if err != nil || resBody.Status != "201 Created" {
		suite.Fail("failed to create user", err)
	}
	defer resBody.Body.Close()
}

//TestGenerateToken is token generation endpoint
func (suite *e2eTestSuite) TestGenerateToken() {
	client := http.DefaultClient
	body := handlers.UserRegisterRequest{
		Name:     "test",
		Email:    "test@gmail.com",
		Password: "abc",
	}
	marshal, _ := json.Marshal(body)
	resBody, err := client.Post("http://localhost:8080/users/register", "application/json", bytes.NewReader(marshal))
	if err != nil || resBody.Status != "201 Created" {
		suite.Fail("failed to create user", err)
	}
	defer resBody.Body.Close()
	token:= handlers.TokenRequest{
		Email:    "test@gmail.com",
		Password: "abc",
	}
	tokenMarshal, _ := json.Marshal(token)
	resp, err := client.Post("http://localhost:8080/token", "application/json", bytes.NewReader(tokenMarshal))
	if err != nil {
		suite.Fail("can not get token")
	}
	defer resp.Body.Close()
	tokenData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		suite.Fail("can not get token")
	}
	tokenResponse:= &handlers.TokenResponse{}
	_= json.Unmarshal(tokenData, tokenResponse)
	if tokenResponse.Token == "" {
		suite.Fail("can not get token")
	}
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(e2eTestSuite))
}
