package integration

import (
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/internal/kafka/internal/api/public"
	test2 "github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/internal/kafka/test"
	"net/http"
	"testing"

	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/config"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/test"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/test/mocks"
	. "github.com/onsi/gomega"
)

type TestEnv struct {
	helper   *test.Helper
	client   *public.APIClient
	teardown func()
}

func termsRequiredSetup(termsRequired bool, t *testing.T) TestEnv {
	ocmServerBuilder := mocks.NewMockConfigurableServerBuilder()
	termsReviewResponse, err := mocks.GetMockTermsReviewBuilder(nil).TermsRequired(termsRequired).Build()
	if err != nil {
		t.Fatalf(err.Error())
	}
	ocmServerBuilder.SetTermsReviewPostResponse(termsReviewResponse, nil)
	ocmServer := ocmServerBuilder.Build()

	// setup the test environment, if OCM_ENV=integration then the ocmServer provided will be used instead of actual
	// ocm
	h, client, tearDown := test2.NewKafkaHelperWithHooks(t, ocmServer, func(serverConfig *config.ServerConfig) {
		serverConfig.EnableTermsAcceptance = true
	})

	return TestEnv{
		helper: h,
		client: client,
		teardown: func() {
			ocmServer.Close()
			tearDown()
		},
	}
}

func TestTermsRequired_CreateKafkaTermsRequired(t *testing.T) {
	env := termsRequiredSetup(true, t)
	defer env.teardown()

	if test2.TestServices.OCMConfig.MockMode != config.MockModeEmulateServer {
		t.SkipNow()
	}

	// setup pre-requisites to performing requests
	account := env.helper.NewRandAccount()
	ctx := env.helper.NewAuthenticatedContext(account, nil)

	k := public.KafkaRequestPayload{
		Region:        mocks.MockCluster.Region().ID(),
		CloudProvider: mocks.MockCluster.CloudProvider().ID(),
		Name:          mockKafkaName,
		MultiAz:       testMultiAZ,
	}

	_, resp, err := env.client.DefaultApi.CreateKafka(ctx, true, k)

	Expect(err).To(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusForbidden))
}

func TestTermsRequired_CreateKafka_TermsNotRequired(t *testing.T) {
	env := termsRequiredSetup(false, t)
	defer env.teardown()

	if test2.TestServices.OCMConfig.MockMode != config.MockModeEmulateServer {
		t.SkipNow()
	}

	// setup pre-requisites to performing requests
	account := env.helper.NewRandAccount()
	ctx := env.helper.NewAuthenticatedContext(account, nil)

	k := public.KafkaRequestPayload{
		Region:        mocks.MockCluster.Region().ID(),
		CloudProvider: mocks.MockCluster.CloudProvider().ID(),
		Name:          mockKafkaName,
		MultiAz:       testMultiAZ,
	}

	_, resp, err := env.client.DefaultApi.CreateKafka(ctx, true, k)

	Expect(err).NotTo(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusAccepted))
}

func TestTermsRequired_ListKafkaTermsRequired(t *testing.T) {
	env := termsRequiredSetup(true, t)
	defer env.teardown()

	// setup pre-requisites to performing requests
	account := env.helper.NewRandAccount()
	ctx := env.helper.NewAuthenticatedContext(account, nil)

	_, resp, err := env.client.DefaultApi.GetKafkas(ctx, nil)

	Expect(err).NotTo(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))
}