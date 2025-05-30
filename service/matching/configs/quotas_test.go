package configs

import (
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/server/api/matchingservice/v1"
	"go.temporal.io/server/common/headers"
	"go.temporal.io/server/common/quotas"
	"go.temporal.io/server/common/testing/temporalapi"
)

type (
	quotasSuite struct {
		suite.Suite
		*require.Assertions
	}
)

func TestQuotasSuite(t *testing.T) {
	s := new(quotasSuite)
	suite.Run(t, s)
}

func (s *quotasSuite) SetupSuite() {
}

func (s *quotasSuite) TearDownSuite() {
}

func (s *quotasSuite) SetupTest() {
	s.Assertions = require.New(s.T())
}

func (s *quotasSuite) TearDownTest() {
}

func (s *quotasSuite) TestAPIToPriorityMapping() {
	for _, priority := range APIToPriority {
		index := slices.Index(APIPrioritiesOrdered, priority)
		s.NotEqual(-1, index)
	}
}

func (s *quotasSuite) TestAPIPrioritiesOrdered() {
	for idx := range APIPrioritiesOrdered[1:] {
		s.True(APIPrioritiesOrdered[idx] < APIPrioritiesOrdered[idx+1])
	}
}

func (s *quotasSuite) TestAPIs() {
	var service matchingservice.MatchingServiceServer
	apiToPriority := make(map[string]int)
	temporalapi.WalkExportedMethods(&service, func(m reflect.Method) {
		fullName := "/temporal.server.api.matchingservice.v1.MatchingService/" + m.Name
		apiToPriority[fullName] = APIToPriority[fullName]
	})
	s.Equal(apiToPriority, APIToPriority)
}

func (s *quotasSuite) TestOperatorPrioritized() {
	rateFn := func() float64 { return 5 }
	operatorRPSRatioFn := func() float64 { return 0.2 }
	limiter := NewPriorityRateLimiter(rateFn, operatorRPSRatioFn)

	operatorRequest := quotas.NewRequest(
		"/temporal.server.api.matchingservice.v1.MatchingService/QueryWorkflow",
		1,
		"",
		headers.CallerTypeOperator,
		-1,
		"")

	apiRequest := quotas.NewRequest(
		"/temporal.server.api.matchingservice.v1.MatchingService/QueryWorkflow",
		1,
		"",
		headers.CallerTypeAPI,
		-1,
		"")

	requestTime := time.Now()
	limitCount := 0

	for i := 0; i < 12; i++ {
		if !limiter.Allow(requestTime, apiRequest) {
			limitCount++
			s.True(limiter.Allow(requestTime, operatorRequest))
		}
	}
	s.Equal(2, limitCount)
}
