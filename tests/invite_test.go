package tests

import (
	"fmt"
	"sync"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"
)

const (
	testCode = "test-code"
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

func (s *InviteTestSuite) TestInviteRegistrationSuccess() {
	testEmail := "test@test.com"
	apitest.New().
		EnableNetworking().
		Post(s.BaseURL() + "/invite").
		ContentType("application/json").
		Body(fmt.Sprintf(`{"email":"%s","code":"%s"}`, testEmail, testCode)).
		Expect(s.T()).
		Status(200).
		Body("").
		End()

	//todo check code was written to the db
	// add test check wrong emails wasn't written
}
func (s *InviteTestSuite) TestMultipleParalleInviteRegistrationsSuccess() {
	testEmail := "test@test.com"

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			apitest.New().
				EnableNetworking().
				Post(s.BaseURL() + "/invite").
				ContentType("application/json").
				Body(fmt.Sprintf(`{"email":"%s", "code":"%s"}`, testEmail, testCode)).
				Expect(s.T()).
				Status(200).
				Body("").
				End()

			wg.Done()
		}()
	}

}
