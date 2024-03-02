package tests

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/sadensmol/test_playoff/internal/invite"
	"github.com/sadensmol/test_playoff/internal/utils"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	rndStr, err := utils.GenRandomString(10)
	s.NoError(err)

	startTime := time.Now()
	testEmail := fmt.Sprintf("test-%s@test.com", rndStr)
	apitest.New().
		EnableNetworking().
		Post(s.BaseURL() + "/invite").
		ContentType("application/json").
		Body(fmt.Sprintf(`{"email":"%s","code":"%s"}`, testEmail, testCode)).
		Expect(s.T()).
		Status(200).
		Body("").
		End()

	result := s.IntegrationTest.MongoClient.Database("invites").Collection(fmt.Sprintf("code-%s", testCode)).
		FindOne(*s.IntegrationTest.Ctx, bson.M{"email": testEmail})

	s.NoError(result.Err())
	var invite invite.InviteModel
	err = result.Decode(&invite)
	s.NoError(err)

	s.Equal(testEmail, invite.Email)
	s.True(invite.RegisteredAt.After(startTime) && invite.RegisteredAt.Before(time.Now()))
}

func (s *InviteTestSuite) TestInviteRegistrationWrongEmailSkipped() {
	rndStr, err := utils.GenRandomString(10)
	s.NoError(err)

	testEmail := fmt.Sprintf("test-%s@test.com", rndStr)
	apitest.New().
		EnableNetworking().
		Post(s.BaseURL() + "/invite").
		ContentType("application/json").
		Body(fmt.Sprintf(`{"email":"some_wrong_email","code":"%s"}`, testCode)).
		Expect(s.T()).
		Status(200).
		Body("").
		End()

	result := s.IntegrationTest.MongoClient.Database("invites").Collection(fmt.Sprintf("code-%s", testCode)).
		FindOne(*s.IntegrationTest.Ctx, bson.M{"email": testEmail})
	s.ErrorIs(result.Err(), mongo.ErrNoDocuments)
	_ = result
}
func (s *InviteTestSuite) TestInviteRegistrationBlankCodeSkipped() {
	rndStr, err := utils.GenRandomString(10)
	s.NoError(err)

	testEmail := fmt.Sprintf("test-%s@test.com", rndStr)
	emptyCode := "  "
	apitest.New().
		EnableNetworking().
		Post(s.BaseURL() + "/invite").
		ContentType("application/json").
		Body(fmt.Sprintf(`{"email":"%s","code":"%s"}`, testEmail, emptyCode)).
		Expect(s.T()).
		Status(200).
		Body("").
		End()

	result := s.IntegrationTest.MongoClient.Database("invites").Collection(fmt.Sprintf("code-%s", emptyCode)).
		FindOne(*s.IntegrationTest.Ctx, bson.M{"email": testEmail})
	s.ErrorIs(result.Err(), mongo.ErrNoDocuments)
	_ = result
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
