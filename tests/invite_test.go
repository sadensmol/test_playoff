package tests

import (
	"fmt"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"
)

type InviteTestSuite struct {
	suite.Suite
	IntegrationTest
}

func (s *InviteTestSuite) BaseURL() string {
	return fmt.Sprintf("http://localhost:%d", s.IntegrationTest.Cfg.HTTP.Port)
}

func (s *InviteTestSuite) SetupSuite() {
	s.IntegrationTest.Setup()
}

func (s *InviteTestSuite) TearDownSuite() {
	s.IntegrationTest.TearDown()
}

func TestIntAPITestSuite(t *testing.T) {
	suite.Run(t, new(InviteTestSuite))
}

func (s *InviteTestSuite) TestInviteSuccess() {
	testEmail := "test@test.com"
	apitest.New().
		EnableNetworking().
		Post(s.BaseURL() + "/invite").
		ContentType("application/json").
		Body(fmt.Sprintf(`{"email":"%s"}`, testEmail)).
		Expect(s.T()).
		Status(200).
		Body("").
		End()
}
