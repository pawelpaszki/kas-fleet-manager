package handlers

import (
	"context"
	"net/http"
	"testing"

	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/internal/kafka/internal/api/dbapi"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/internal/kafka/internal/config"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/internal/kafka/internal/services"
	mocks "github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/internal/kafka/test/mocks/kafkas"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/errors"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/services/authorization"
	"github.com/gorilla/mux"
	. "github.com/onsi/gomega"
)

var (
	id = "test-id"
)

func Test_KafkaHandler_Get(t *testing.T) {
	type fields struct {
		service        services.KafkaService
		providerConfig *config.ProviderConfig
		authService    authorization.Authorization
		kafkaConfig    *config.KafkaConfig
	}

	tests := []struct {
		name           string
		fields         fields
		wantStatusCode int
	}{
		{
			name: "should fail it kafkaService GET fails",
			fields: fields{
				service: &services.KafkaServiceMock{
					GetFunc: func(ctx context.Context, id string) (*dbapi.KafkaRequest, *errors.ServiceError) {
						return nil, errors.GeneralError("test")
					},
				},
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "should succeed if kafkaService GET succeeds",
			fields: fields{
				service: &services.KafkaServiceMock{
					GetFunc: func(ctx context.Context, id string) (*dbapi.KafkaRequest, *errors.ServiceError) {
						return mocks.BuildKafkaRequest(mocks.WithPredefinedTestValues()), nil
					},
				},
				kafkaConfig: &FullKafkaConfig,
			},
			wantStatusCode: http.StatusOK,
		},
	}

	RegisterTestingT(t)

	for _, testcase := range tests {
		tt := testcase
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := NewKafkaHandler(tt.fields.service, tt.fields.providerConfig, tt.fields.authService, tt.fields.kafkaConfig)

			req, rw := GetHandlerParams("GET", "/{id}", nil)
			req = mux.SetURLVars(req, map[string]string{"id": id})
			h.Get(rw, req)
			resp := rw.Result()
			Expect(resp.StatusCode).To(Equal(tt.wantStatusCode))
			resp.Body.Close()
		})
	}
}
