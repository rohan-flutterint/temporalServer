package workerdeployment

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"sync"

	deploymentpb "go.temporal.io/api/deployment/v1"
	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/api/serviceerror"
	taskqueuepb "go.temporal.io/api/taskqueue/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
	deploymentspb "go.temporal.io/server/api/deployment/v1"
	"go.temporal.io/server/api/matchingservice/v1"
	"go.temporal.io/server/common/namespace"
	"go.temporal.io/server/common/resource"
	"go.temporal.io/server/common/worker_versioning"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	VersionActivities struct {
		namespace        *namespace.Namespace
		deploymentClient Client
		matchingClient   resource.MatchingClient
	}
)

func (a *VersionActivities) StartWorkerDeploymentWorkflow(
	ctx context.Context,
	input *deploymentspb.StartWorkerDeploymentRequest,
) error {
	logger := activity.GetLogger(ctx)
	logger.Info("starting worker-deployment workflow", "deploymentName", input.DeploymentName)
	identity := "deployment-version workflow " + activity.GetInfo(ctx).WorkflowExecution.ID
	err := a.deploymentClient.StartWorkerDeployment(ctx, a.namespace, input.DeploymentName, identity, input.RequestId)
	var precond *serviceerror.FailedPrecondition
	if errors.As(err, &precond) {
		return temporal.NewNonRetryableApplicationError("failed to create deployment", errTooManyDeployments, err)
	}
	return err
}

func (a *VersionActivities) SyncDeploymentVersionUserData(
	ctx context.Context,
	input *deploymentspb.SyncDeploymentVersionUserDataRequest,
) (*deploymentspb.SyncDeploymentVersionUserDataResponse, error) {
	logger := activity.GetLogger(ctx)

	errs := make(chan error)

	var lock sync.Mutex
	maxVersionByName := make(map[string]int64)

	for _, e := range input.Sync {
		go func(syncData *deploymentspb.SyncDeploymentVersionUserDataRequest_SyncUserData) {
			logger.Info("syncing task queue userdata for deployment version", "taskQueue", syncData.Name, "types", syncData.Types)

			var res *matchingservice.SyncDeploymentUserDataResponse
			var err error

			if input.ForgetVersion {
				res, err = a.matchingClient.SyncDeploymentUserData(ctx, &matchingservice.SyncDeploymentUserDataRequest{
					NamespaceId:    a.namespace.ID().String(),
					TaskQueue:      syncData.Name,
					TaskQueueTypes: syncData.Types,
					Operation: &matchingservice.SyncDeploymentUserDataRequest_ForgetVersion{
						ForgetVersion: input.Version,
					},
				})
			} else {
				res, err = a.matchingClient.SyncDeploymentUserData(ctx, &matchingservice.SyncDeploymentUserDataRequest{
					NamespaceId:    a.namespace.ID().String(),
					TaskQueue:      syncData.Name,
					TaskQueueTypes: syncData.Types,
					Operation: &matchingservice.SyncDeploymentUserDataRequest_UpdateVersionData{
						UpdateVersionData: syncData.Data,
					},
				})
			}

			if err != nil {
				logger.Error("syncing task queue userdata", "taskQueue", syncData.Name, "types", syncData.Types, "error", err)
			} else {
				lock.Lock()
				maxVersionByName[syncData.Name] = max(maxVersionByName[syncData.Name], res.Version)
				lock.Unlock()
			}
			errs <- err
		}(e)
	}

	var err error
	for range input.Sync {
		err = cmp.Or(err, <-errs)
	}
	if err != nil {
		return nil, err
	}
	return &deploymentspb.SyncDeploymentVersionUserDataResponse{TaskQueueMaxVersions: maxVersionByName}, nil
}

func (a *VersionActivities) CheckWorkerDeploymentUserDataPropagation(ctx context.Context, input *deploymentspb.CheckWorkerDeploymentUserDataPropagationRequest) error {
	logger := activity.GetLogger(ctx)

	errs := make(chan error)

	for n, v := range input.TaskQueueMaxVersions {
		go func(name string, version int64) {
			logger.Info("waiting for userdata propagation", "taskQueue", name, "version", version)
			_, err := a.matchingClient.CheckTaskQueueUserDataPropagation(ctx, &matchingservice.CheckTaskQueueUserDataPropagationRequest{
				NamespaceId: a.namespace.ID().String(),
				TaskQueue:   name,
				Version:     version,
			})
			if err != nil {
				logger.Error("waiting for userdata", "taskQueue", name, "type", version, "error", err)
			}
			errs <- err
		}(n, v)
	}

	var err error
	for range input.TaskQueueMaxVersions {
		err = cmp.Or(err, <-errs)
	}
	return err
}

// CheckIfTaskQueuesHavePollers returns true if any of the given task queues has any pollers
func (a *VersionActivities) CheckIfTaskQueuesHavePollers(ctx context.Context, args *deploymentspb.CheckTaskQueuesHavePollersActivityArgs) (bool, error) {
	versionStr := worker_versioning.ExternalWorkerDeploymentVersionToString(worker_versioning.ExternalWorkerDeploymentVersionFromVersion(args.WorkerDeploymentVersion))
	for tqName, tqTypes := range args.TaskQueuesAndTypes {
		res, err := a.matchingClient.DescribeTaskQueue(ctx, &matchingservice.DescribeTaskQueueRequest{
			NamespaceId: a.namespace.ID().String(),
			DescRequest: &workflowservice.DescribeTaskQueueRequest{
				Namespace:      a.namespace.Name().String(),
				TaskQueue:      &taskqueuepb.TaskQueue{Name: tqName, Kind: enumspb.TASK_QUEUE_KIND_NORMAL},
				ApiMode:        enumspb.DESCRIBE_TASK_QUEUE_MODE_ENHANCED,
				Versions:       &taskqueuepb.TaskQueueVersionSelection{BuildIds: []string{versionStr}},
				ReportPollers:  true,
				TaskQueueType:  enumspb.TASK_QUEUE_TYPE_WORKFLOW,
				TaskQueueTypes: tqTypes.Types,
			},
		})
		if err != nil {
			return false, fmt.Errorf("error describing task queue with name %s: %s", tqName, err)
		}
		typesInfo := res.GetDescResponse().GetVersionsInfo()[versionStr].GetTypesInfo()
		if len(typesInfo[int32(enumspb.TASK_QUEUE_TYPE_WORKFLOW)].GetPollers()) > 0 {
			return true, nil
		}
		if len(typesInfo[int32(enumspb.TASK_QUEUE_TYPE_ACTIVITY)].GetPollers()) > 0 {
			return true, nil
		}
		if len(typesInfo[int32(enumspb.TASK_QUEUE_TYPE_NEXUS)].GetPollers()) > 0 {
			return true, nil
		}
	}
	return false, nil
}

func (a *VersionActivities) AddVersionToWorkerDeployment(ctx context.Context, input *deploymentspb.AddVersionToWorkerDeploymentRequest) (*deploymentspb.AddVersionToWorkerDeploymentResponse, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("adding version to worker-deployment", "deploymentName", input.DeploymentName, "version", input.UpdateArgs.Version)
	identity := "deployment-version workflow " + activity.GetInfo(ctx).WorkflowExecution.ID
	resp, err := a.deploymentClient.AddVersionToWorkerDeployment(ctx, a.namespace, input.DeploymentName, input.UpdateArgs, identity, input.RequestId)
	var precond *serviceerror.FailedPrecondition
	if errors.As(err, &precond) {
		return nil, temporal.NewNonRetryableApplicationError("failed to add version to deployment", errTooManyVersions, err)
	}
	return resp, err
}

func (a *VersionActivities) GetVersionDrainageStatus(ctx context.Context, version *deploymentspb.WorkerDeploymentVersion) (*deploymentpb.VersionDrainageInfo, error) {
	logger := activity.GetLogger(ctx)
	response, err := a.deploymentClient.GetVersionDrainageStatus(ctx, a.namespace, worker_versioning.WorkerDeploymentVersionToStringV31(version))
	if err != nil {
		logger.Error("error counting workflows for drainage status", "error", err)
		return nil, err
	}
	return &deploymentpb.VersionDrainageInfo{
		Status:          response,
		LastChangedTime: nil, // ignored; whether Status changed will be evaluated by the receiver
		LastCheckedTime: timestamppb.Now(),
	}, nil
}
