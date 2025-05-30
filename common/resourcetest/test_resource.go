package resourcetest

import (
	"net"

	"github.com/uber-go/tally/v4"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/server/api/adminservice/v1"
	"go.temporal.io/server/api/adminservicemock/v1"
	"go.temporal.io/server/api/historyservice/v1"
	"go.temporal.io/server/api/historyservicemock/v1"
	"go.temporal.io/server/api/matchingservice/v1"
	"go.temporal.io/server/api/matchingservicemock/v1"
	"go.temporal.io/server/client"
	"go.temporal.io/server/common/archiver"
	"go.temporal.io/server/common/archiver/provider"
	"go.temporal.io/server/common/clock"
	"go.temporal.io/server/common/cluster"
	"go.temporal.io/server/common/log"
	"go.temporal.io/server/common/membership"
	"go.temporal.io/server/common/metrics"
	"go.temporal.io/server/common/namespace"
	"go.temporal.io/server/common/persistence"
	"go.temporal.io/server/common/persistence/serialization"
	"go.temporal.io/server/common/persistence/visibility/manager"
	esclient "go.temporal.io/server/common/persistence/visibility/store/elasticsearch/client"
	"go.temporal.io/server/common/primitives"
	"go.temporal.io/server/common/sdk"
	"go.temporal.io/server/common/searchattribute"
	"go.temporal.io/server/common/testing/mockapi/workflowservicemock/v1"
	"go.uber.org/mock/gomock"
)

// TODO: replace with test specific Fx

type (
	// Test is the test implementation used for testing
	Test struct {
		MetricsScope                   tally.Scope
		ClusterMetadata                *cluster.MockMetadata
		SearchAttributesProvider       *searchattribute.MockProvider
		SearchAttributesManager        *searchattribute.MockManager
		SearchAttributesMapperProvider *searchattribute.MockMapperProvider

		// other common resources

		NamespaceCache    *namespace.MockRegistry
		TimeSource        clock.TimeSource
		PayloadSerializer serialization.Serializer
		MetricsHandler    metrics.Handler
		ArchivalMetadata  archiver.MetadataMock
		ArchiverProvider  *provider.MockArchiverProvider

		// membership infos

		MembershipMonitor       *membership.MockMonitor
		HostInfoProvider        *membership.MockHostInfoProvider
		FrontendServiceResolver *membership.MockServiceResolver
		MatchingServiceResolver *membership.MockServiceResolver
		HistoryServiceResolver  *membership.MockServiceResolver
		WorkerServiceResolver   *membership.MockServiceResolver

		// internal services clients

		SDKClientFactory     *sdk.MockClientFactory
		FrontendClient       *workflowservicemock.MockWorkflowServiceClient
		MatchingClient       *matchingservicemock.MockMatchingServiceClient
		HistoryClient        *historyservicemock.MockHistoryServiceClient
		RemoteAdminClient    *adminservicemock.MockAdminServiceClient
		RemoteFrontendClient *workflowservicemock.MockWorkflowServiceClient
		ClientBean           *client.MockBean
		ClientFactory        *client.MockFactory
		ESClient             *esclient.MockClient
		VisibilityManager    *manager.MockVisibilityManager

		// persistence clients

		MetadataMgr               *persistence.MockMetadataManager
		ClusterMetadataMgr        *persistence.MockClusterMetadataManager
		TaskMgr                   *persistence.MockTaskManager
		NamespaceReplicationQueue persistence.NamespaceReplicationQueue
		ShardMgr                  *persistence.MockShardManager
		ExecutionMgr              *persistence.MockExecutionManager
		NexusEndpointManager      *persistence.MockNexusEndpointManager

		Logger log.Logger
	}
)

const (
	testHostName = "test_host"
)

var testHostInfo = membership.NewHostInfoFromAddress(testHostName)

// NewTest returns a new test resource instance
func NewTest(controller *gomock.Controller, serviceName primitives.ServiceName) *Test {
	logger := log.NewTestLogger()

	frontendClient := workflowservicemock.NewMockWorkflowServiceClient(controller)
	matchingClient := matchingservicemock.NewMockMatchingServiceClient(controller)
	historyClient := historyservicemock.NewMockHistoryServiceClient(controller)
	remoteFrontendClient := workflowservicemock.NewMockWorkflowServiceClient(controller)
	remoteAdminClient := adminservicemock.NewMockAdminServiceClient(controller)
	clusterMetadataManager := persistence.NewMockClusterMetadataManager(controller)
	clientBean := client.NewMockBean(controller)
	clientBean.EXPECT().GetFrontendClient().Return(frontendClient).AnyTimes()
	clientBean.EXPECT().GetMatchingClient(gomock.Any()).Return(matchingClient, nil).AnyTimes()
	clientBean.EXPECT().GetHistoryClient().Return(historyClient).AnyTimes()
	clientBean.EXPECT().GetRemoteAdminClient(gomock.Any()).Return(remoteAdminClient, nil).AnyTimes()
	clientBean.EXPECT().GetRemoteFrontendClient(gomock.Any()).Return(nil, remoteFrontendClient, nil).AnyTimes()
	clientFactory := client.NewMockFactory(controller)

	metadataMgr := persistence.NewMockMetadataManager(controller)
	taskMgr := persistence.NewMockTaskManager(controller)
	shardMgr := persistence.NewMockShardManager(controller)
	executionMgr := persistence.NewMockExecutionManager(controller)
	executionMgr.EXPECT().GetHistoryBranchUtil().Return(&persistence.HistoryBranchUtilImpl{}).AnyTimes()
	namespaceReplicationQueue := persistence.NewMockNamespaceReplicationQueue(controller)
	nexusEndpointMgr := persistence.NewMockNexusEndpointManager(controller)

	membershipMonitor := membership.NewMockMonitor(controller)
	hostInfoProvider := membership.NewMockHostInfoProvider(controller)
	frontendServiceResolver := membership.NewMockServiceResolver(controller)
	matchingServiceResolver := membership.NewMockServiceResolver(controller)
	historyServiceResolver := membership.NewMockServiceResolver(controller)
	workerServiceResolver := membership.NewMockServiceResolver(controller)
	membershipMonitor.EXPECT().GetResolver(primitives.FrontendService).Return(frontendServiceResolver, nil).AnyTimes()
	membershipMonitor.EXPECT().GetResolver(primitives.InternalFrontendService).Return(nil, membership.ErrUnknownService).AnyTimes()
	membershipMonitor.EXPECT().GetResolver(primitives.MatchingService).Return(matchingServiceResolver, nil).AnyTimes()
	membershipMonitor.EXPECT().GetResolver(primitives.HistoryService).Return(historyServiceResolver, nil).AnyTimes()
	membershipMonitor.EXPECT().GetResolver(primitives.WorkerService).Return(workerServiceResolver, nil).AnyTimes()
	membershipMonitor.EXPECT().WaitUntilInitialized(gomock.Any()).Return(nil).AnyTimes()

	scope := tally.NewTestScope("test", nil)
	metricsHandler := metrics.NewTallyMetricsHandler(metrics.ClientConfig{}, scope).WithTags(
		metrics.ServiceNameTag(serviceName),
	)

	return &Test{
		MetricsScope:                   scope,
		ClusterMetadata:                cluster.NewMockMetadata(controller),
		SearchAttributesProvider:       searchattribute.NewMockProvider(controller),
		SearchAttributesManager:        searchattribute.NewMockManager(controller),
		SearchAttributesMapperProvider: searchattribute.NewMockMapperProvider(controller),

		// other common resources

		NamespaceCache:    namespace.NewMockRegistry(controller),
		TimeSource:        clock.NewRealTimeSource(),
		PayloadSerializer: serialization.NewSerializer(),
		MetricsHandler:    metricsHandler,
		ArchivalMetadata:  archiver.NewMetadataMock(controller),
		ArchiverProvider:  provider.NewMockArchiverProvider(controller),

		// membership infos

		MembershipMonitor:       membershipMonitor,
		HostInfoProvider:        hostInfoProvider,
		FrontendServiceResolver: frontendServiceResolver,
		MatchingServiceResolver: matchingServiceResolver,
		HistoryServiceResolver:  historyServiceResolver,
		WorkerServiceResolver:   workerServiceResolver,

		// internal services clients

		SDKClientFactory:     sdk.NewMockClientFactory(controller),
		FrontendClient:       frontendClient,
		MatchingClient:       matchingClient,
		HistoryClient:        historyClient,
		RemoteAdminClient:    remoteAdminClient,
		RemoteFrontendClient: remoteFrontendClient,
		ClientBean:           clientBean,
		ClientFactory:        clientFactory,
		ESClient:             esclient.NewMockClient(controller),
		VisibilityManager:    manager.NewMockVisibilityManager(controller),

		// persistence clients

		MetadataMgr:               metadataMgr,
		ClusterMetadataMgr:        clusterMetadataManager,
		TaskMgr:                   taskMgr,
		NamespaceReplicationQueue: namespaceReplicationQueue,
		ShardMgr:                  shardMgr,
		ExecutionMgr:              executionMgr,
		NexusEndpointManager:      nexusEndpointMgr,

		// logger

		Logger: logger,
	}
}

// Start for testing
func (t *Test) Start() {
}

// Stop for testing
func (t *Test) Stop() {
}

// static infos

// GetServiceName for testing
func (t *Test) GetServiceName() string {
	panic("user should implement this method for test")
}

// GetHostName for testing
func (t *Test) GetHostName() string {
	return testHostInfo.Identity()
}

// GetHostInfo for testing
func (t *Test) GetHostInfo() membership.HostInfo {
	return testHostInfo
}

// GetClusterMetadata for testing
func (t *Test) GetClusterMetadata() cluster.Metadata {
	return t.ClusterMetadata
}

// GetClusterMetadata for testing
func (t *Test) GetClusterMetadataManager() persistence.ClusterMetadataManager {
	return t.ClusterMetadataMgr
}

// other common resources

// GetNamespaceRegistry for testing
func (t *Test) GetNamespaceRegistry() namespace.Registry {
	return t.NamespaceCache
}

// GetTimeSource for testing
func (t *Test) GetTimeSource() clock.TimeSource {
	return t.TimeSource
}

// GetPayloadSerializer for testing
func (t *Test) GetPayloadSerializer() serialization.Serializer {
	return t.PayloadSerializer
}

// GetMetricsHandler for testing
func (t *Test) GetMetricsHandler() metrics.Handler {
	return t.MetricsHandler
}

// GetArchivalMetadata for testing
func (t *Test) GetArchivalMetadata() archiver.ArchivalMetadata {
	return t.ArchivalMetadata
}

// GetArchiverProvider for testing
func (t *Test) GetArchiverProvider() provider.ArchiverProvider {
	return t.ArchiverProvider
}

// membership infos

// GetMembershipMonitor for testing
func (t *Test) GetMembershipMonitor() membership.Monitor {
	return t.MembershipMonitor
}

// GetHostInfoProvider for testing
func (t *Test) GetHostInfoProvider() membership.HostInfoProvider {
	return t.HostInfoProvider
}

// GetFrontendServiceResolver for testing
func (t *Test) GetFrontendServiceResolver() membership.ServiceResolver {
	return t.FrontendServiceResolver
}

// GetMatchingServiceResolver for testing
func (t *Test) GetMatchingServiceResolver() membership.ServiceResolver {
	return t.MatchingServiceResolver
}

// GetHistoryServiceResolver for testing
func (t *Test) GetHistoryServiceResolver() membership.ServiceResolver {
	return t.HistoryServiceResolver
}

// GetWorkerServiceResolver for testing
func (t *Test) GetWorkerServiceResolver() membership.ServiceResolver {
	return t.WorkerServiceResolver
}

// internal services clients

// GetSDKClientFactory for testing
func (t *Test) GetSDKClientFactory() sdk.ClientFactory {
	return t.SDKClientFactory
}

// GetFrontendClient for testing
func (t *Test) GetFrontendClient() workflowservice.WorkflowServiceClient {
	return t.FrontendClient
}

// GetMatchingRawClient for testing
func (t *Test) GetMatchingRawClient() matchingservice.MatchingServiceClient {
	return t.MatchingClient
}

// GetMatchingClient for testing
func (t *Test) GetMatchingClient() matchingservice.MatchingServiceClient {
	return t.MatchingClient
}

// GetHistoryRawClient for testing
func (t *Test) GetHistoryRawClient() historyservice.HistoryServiceClient {
	return t.HistoryClient
}

// GetHistoryClient for testing
func (t *Test) GetHistoryClient() historyservice.HistoryServiceClient {
	return t.HistoryClient
}

// GetRemoteAdminClient for testing
func (t *Test) GetRemoteAdminClient(
	cluster string,
) adminservice.AdminServiceClient {
	return t.RemoteAdminClient
}

// GetRemoteFrontendClient for testing
func (t *Test) GetRemoteFrontendClient(
	cluster string,
) workflowservice.WorkflowServiceClient {
	return t.RemoteFrontendClient
}

// GetClientBean for testing
func (t *Test) GetClientBean() client.Bean {
	return t.ClientBean
}

// GetClientFactory for testing
func (t *Test) GetClientFactory() client.Factory {
	return t.ClientFactory
}

// GetVisibilityManager for testing
func (t *Test) GetVisibilityManager() manager.VisibilityManager {
	return t.VisibilityManager
}

// persistence clients

// GetMetadataManager for testing
func (t *Test) GetMetadataManager() persistence.MetadataManager {
	return t.MetadataMgr
}

// GetTaskManager for testing
func (t *Test) GetTaskManager() persistence.TaskManager {
	return t.TaskMgr
}

// GetNamespaceReplicationQueue for testing
func (t *Test) GetNamespaceReplicationQueue() persistence.NamespaceReplicationQueue {
	// user should implement this method for test
	return t.NamespaceReplicationQueue
}

// GetShardManager for testing
func (t *Test) GetShardManager() persistence.ShardManager {
	return t.ShardMgr
}

// GetExecutionManager for testing
func (t *Test) GetExecutionManager() persistence.ExecutionManager {
	return t.ExecutionMgr
}

// loggers

// GetLogger for testing
func (t *Test) GetLogger() log.Logger {
	return t.Logger
}

// GetThrottledLogger for testing
func (t *Test) GetThrottledLogger() log.Logger {
	return t.Logger
}

// GetGRPCListener for testing
func (t *Test) GetGRPCListener() net.Listener {
	panic("user should implement this method for test")
}

func (t *Test) GetSearchAttributesProvider() searchattribute.Provider {
	return t.SearchAttributesProvider
}

func (t *Test) GetSearchAttributesManager() searchattribute.Manager {
	return t.SearchAttributesManager
}

func (t *Test) GetSearchAttributesMapperProvider() searchattribute.MapperProvider {
	return t.SearchAttributesMapperProvider
}
