package workers

import (
	"fmt"
	"sync"
	"time"

	"gitlab.cee.redhat.com/service/managed-services-api/pkg/metrics"

	"gitlab.cee.redhat.com/service/managed-services-api/pkg/api"
	"gitlab.cee.redhat.com/service/managed-services-api/pkg/ocm"

	"github.com/getsentry/sentry-go"
	"github.com/golang/glog"
	"gitlab.cee.redhat.com/service/managed-services-api/pkg/services"
)

// KafkaManager represents a kafka manager that periodically reconciles kafka requests
type KafkaManager struct {
	ocmClient      ocm.Client
	clusterService services.ClusterService
	kafkaService   services.KafkaService
	timer          *time.Timer
	imStop         chan struct{}
	syncTeardown   sync.WaitGroup
	reconciler     Reconciler
}

// NewKafkaManager creates a new kafka manager
func NewKafkaManager(kafkaService services.KafkaService, clusterService services.ClusterService, ocmClient ocm.Client) *KafkaManager {
	return &KafkaManager{
		ocmClient:      ocmClient,
		clusterService: clusterService,
		kafkaService:   kafkaService,
	}
}

func (k *KafkaManager) GetStopChan() *chan struct{} {
	return &k.imStop
}

func (k *KafkaManager) GetSyncGroup() *sync.WaitGroup {
	return &k.syncTeardown
}

// Start initializes the kafka manager to reconcile kafka requests
func (k *KafkaManager) Start() {
	k.reconciler.Start(k)
}

// Stop causes the process for reconciling kafka requests to stop.
func (k *KafkaManager) Stop() {
	k.reconciler.Stop(k)
}

func (k *KafkaManager) reconcile() {
	glog.V(5).Infoln("reconciling kafkas")

	// handle accepted kafkas
	acceptedKafkas, serviceErr := k.kafkaService.ListByStatus(services.KafkaRequestStatusAccepted)
	if serviceErr != nil {
		sentry.CaptureException(serviceErr)
		glog.Errorf("failed to list accepted kafkas: %s", serviceErr.Error())
	}

	for _, kafka := range acceptedKafkas {
		if err := k.reconcileAcceptedKafka(kafka); err != nil {
			sentry.CaptureException(err)
			glog.Errorf("failed to reconcile accepted kafka %s: %s", kafka.ID, err.Error())
			continue
		}
	}

	// handle provisioning kafkas
	provisioningKafkas, serviceErr := k.kafkaService.ListByStatus(services.KafkaRequestStatusProvisioning)
	if serviceErr != nil {
		sentry.CaptureException(serviceErr)
		glog.Errorf("failed to list accepted kafkas: %s", serviceErr.Error())
	}

	for _, kafka := range provisioningKafkas {
		if err := k.reconcileProvisionedKafka(kafka); err != nil {
			sentry.CaptureException(err)
			glog.Errorf("failed to reconcile accepted kafka %s: %s", kafka.ID, err.Error())
			continue
		}
	}
}

func (k *KafkaManager) reconcileAcceptedKafka(kafka *api.KafkaRequest) error {
	cluster, err := k.clusterService.FindCluster(services.FindClusterCriteria{
		Provider: kafka.CloudProvider,
		Region:   kafka.Region,
		MultiAZ:  kafka.MultiAZ,
		Status:   api.ClusterReady,
	})
	if err != nil {
		return fmt.Errorf("failed to find cluster for kafka request %s: %w", kafka.ID, err)
	}
	if cluster != nil {
		kafka.ClusterID = cluster.ClusterID
		kafka.Status = services.KafkaRequestStatusProvisioning.String()
		if err = k.kafkaService.Update(kafka); err != nil {
			return fmt.Errorf("failed to update kafka %s with cluster details: %w", kafka.ID, err)
		}
	}
	return nil
}

func (k *KafkaManager) reconcileProvisionedKafka(kafka *api.KafkaRequest) error {
	if err := k.kafkaService.Create(kafka); err != nil {
		return fmt.Errorf("failed to create kafka %s on cluster %s: %w", kafka.ID, kafka.ClusterID, err)
	}
	// consider the kafka in a complete state
	if err := k.kafkaService.UpdateStatus(kafka.ID, services.KafkaRequestStatusComplete); err != nil {
		return fmt.Errorf("failed to update kafka %s to status complete: %w", kafka.ID, err)
	}
	// line below (responsible for adding metrics) will need to be moved to a place
	// where successful kafka cluster creation check takes place (once this check is implemented)
	metrics.UpdateKafkaCreationDurationMetric(metrics.JobTypeKafkaCreate, time.Since(kafka.CreatedAt))
	return nil
}