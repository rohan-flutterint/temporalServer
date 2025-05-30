package tests

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/server/common/persistence/sql/sqlplugin"
	"go.temporal.io/server/common/primitives"
	"go.temporal.io/server/common/shuffle"
	"go.temporal.io/server/common/util"
)

type (
	visibilitySuite struct {
		suite.Suite
		*require.Assertions

		store sqlplugin.Visibility

		version int64 // used to create increasing version numbers for `_version` column
	}
)

const (
	testVisibilityEncoding         = "random encoding"
	testVisibilityWorkflowTypeName = "random workflow type name"
	testVisibilityWorkflowID       = "random workflow ID"
)

var (
	testVisibilityData = []byte("random history execution activity data")
)

var testVisibilityCloseStatus = []enumspb.WorkflowExecutionStatus{
	enumspb.WORKFLOW_EXECUTION_STATUS_COMPLETED,
	enumspb.WORKFLOW_EXECUTION_STATUS_FAILED,
	enumspb.WORKFLOW_EXECUTION_STATUS_CANCELED,
	enumspb.WORKFLOW_EXECUTION_STATUS_TERMINATED,
	enumspb.WORKFLOW_EXECUTION_STATUS_CONTINUED_AS_NEW,
	enumspb.WORKFLOW_EXECUTION_STATUS_TIMED_OUT,
}

func NewVisibilitySuite(
	t *testing.T,
	store sqlplugin.Visibility,
) *visibilitySuite {
	return &visibilitySuite{
		Assertions: require.New(t),
		store:      store,
	}
}

func (s *visibilitySuite) SetupSuite() {

}

func (s *visibilitySuite) TearDownSuite() {

}

func (s *visibilitySuite) SetupTest() {
	s.Assertions = require.New(s.T())
}

func (s *visibilitySuite) TearDownTest() {

}

func (s *visibilitySuite) TestInsertSelect_NonExists() {
	namespaceID := primitives.NewUUID()
	runID := primitives.NewUUID()
	workflowTypeName := shuffle.String(testVisibilityWorkflowTypeName)
	workflowID := shuffle.String(testVisibilityWorkflowID)
	startTime := s.now()
	executionTime := startTime.Add(time.Second)
	status := int32(0)
	closeTime := (*time.Time)(nil)
	historyLength := (*int64)(nil)

	visibility := s.newRandomVisibilityRow(
		namespaceID,
		runID,
		workflowTypeName,
		workflowID,
		startTime,
		executionTime,
		status,
		closeTime,
		historyLength,
	)
	result, err := s.store.InsertIntoVisibility(newVisibilityContext(), &visibility)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	selectFilter := sqlplugin.VisibilitySelectFilter{
		NamespaceID: namespaceID.String(),
		RunID:       util.Ptr(runID.String()),
	}
	rows, err := s.store.SelectFromVisibility(newVisibilityContext(), selectFilter)
	s.NoError(err)
	for index := range rows {
		rows[index].NamespaceID = namespaceID.String()
	}
	s.Equal([]sqlplugin.VisibilityRow{visibility}, rows)
}

func (s *visibilitySuite) TestInsertSelect_Exists() {
	namespaceID := primitives.NewUUID()
	runID := primitives.NewUUID()
	workflowTypeName := shuffle.String(testVisibilityWorkflowTypeName)
	workflowID := shuffle.String(testVisibilityWorkflowID)
	startTime := s.now()
	executionTime := startTime.Add(time.Second)
	status := int32(0)
	closeTime := (*time.Time)(nil)
	historyLength := (*int64)(nil)

	visibility1 := s.newRandomVisibilityRow(
		namespaceID,
		runID,
		workflowTypeName,
		workflowID,
		startTime,
		executionTime,
		status,
		closeTime,
		historyLength,
	)
	result, err := s.store.InsertIntoVisibility(newVisibilityContext(), &visibility1)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	visibility2 := s.newRandomVisibilityRow(
		namespaceID,
		runID,
		workflowTypeName,
		workflowID,
		startTime,
		executionTime,
		status,
		closeTime,
		historyLength,
	)
	_, err = s.store.InsertIntoVisibility(newVisibilityContext(), &visibility2)
	s.NoError(err)
	// NOTE: cannot do assertion on affected rows
	//  PostgreSQL will return 0
	//  MySQL will return 1: ref https://dev.mysql.com/doc/c-api/5.7/en/mysql-affected-rows.html

	selectFilter := sqlplugin.VisibilitySelectFilter{
		NamespaceID: namespaceID.String(),
		RunID:       util.Ptr(runID.String()),
	}
	rows, err := s.store.SelectFromVisibility(newVisibilityContext(), selectFilter)
	s.NoError(err)
	for index := range rows {
		rows[index].NamespaceID = namespaceID.String()
	}
	s.Equal([]sqlplugin.VisibilityRow{visibility1}, rows)
}

func (s *visibilitySuite) TestReplaceSelect_NonExists() {
	namespaceID := primitives.NewUUID()
	runID := primitives.NewUUID()
	workflowTypeName := shuffle.String(testVisibilityWorkflowTypeName)
	workflowID := shuffle.String(testVisibilityWorkflowID)
	startTime := s.now()
	executionTime := startTime.Add(time.Second)
	status := int32(0)
	closeTime := executionTime.Add(time.Second)
	historyLength := rand.Int63()

	visibility := s.newRandomVisibilityRow(
		namespaceID,
		runID,
		workflowTypeName,
		workflowID,
		startTime,
		executionTime,
		status,
		timePtr(closeTime),
		util.Ptr(historyLength),
	)
	result, err := s.store.ReplaceIntoVisibility(newVisibilityContext(), &visibility)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	selectFilter := sqlplugin.VisibilitySelectFilter{
		NamespaceID: namespaceID.String(),
		RunID:       util.Ptr(runID.String()),
	}
	rows, err := s.store.SelectFromVisibility(newVisibilityContext(), selectFilter)
	s.NoError(err)
	for index := range rows {
		rows[index].NamespaceID = namespaceID.String()
	}
	s.Equal([]sqlplugin.VisibilityRow{visibility}, rows)
}

func (s *visibilitySuite) TestReplaceSelect_Exists() {
	namespaceID := primitives.NewUUID()
	runID := primitives.NewUUID()
	workflowTypeName := shuffle.String(testVisibilityWorkflowTypeName)
	workflowID := shuffle.String(testVisibilityWorkflowID)
	startTime := s.now()
	executionTime := startTime.Add(time.Second)
	status := int32(0)
	closeTime := executionTime.Add(time.Second)
	historyLength := rand.Int63()

	visibility := s.newRandomVisibilityRow(
		namespaceID,
		runID,
		workflowTypeName,
		workflowID,
		startTime,
		executionTime,
		status,
		timePtr(closeTime),
		util.Ptr(historyLength),
	)
	result, err := s.store.ReplaceIntoVisibility(newVisibilityContext(), &visibility)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	visibility = s.newRandomVisibilityRow(
		namespaceID,
		runID,
		workflowTypeName,
		workflowID,
		startTime,
		executionTime,
		status,
		timePtr(closeTime),
		util.Ptr(historyLength),
	)
	_, err = s.store.ReplaceIntoVisibility(newVisibilityContext(), &visibility)
	s.NoError(err)
	// NOTE: cannot do assertion on affected rows
	//  PostgreSQL will return 1
	//  MySQL will return 2: ref https://dev.mysql.com/doc/c-api/5.7/en/mysql-affected-rows.html

	selectFilter := sqlplugin.VisibilitySelectFilter{
		NamespaceID: namespaceID.String(),
		RunID:       util.Ptr(runID.String()),
	}
	rows, err := s.store.SelectFromVisibility(newVisibilityContext(), selectFilter)
	s.NoError(err)
	for index := range rows {
		rows[index].NamespaceID = namespaceID.String()
	}
	s.Equal([]sqlplugin.VisibilityRow{visibility}, rows)
}

func (s *visibilitySuite) TestDeleteGet() {
	namespaceID := primitives.NewUUID()
	runID := primitives.NewUUID()

	deleteFilter := sqlplugin.VisibilityDeleteFilter{
		NamespaceID: namespaceID.String(),
		RunID:       runID.String(),
	}
	result, err := s.store.DeleteFromVisibility(newVisibilityContext(), deleteFilter)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(0, int(rowsAffected))

	getFilter := sqlplugin.VisibilityGetFilter{
		NamespaceID: namespaceID.String(),
		RunID:       runID.String(),
	}
	_, err = s.store.GetFromVisibility(newVisibilityContext(), getFilter)
	s.Error(err) // TODO persistence layer should do proper error translation
}

func (s *visibilitySuite) TestInsertDeleteGet() {
	namespaceID := primitives.NewUUID()
	runID := primitives.NewUUID()
	workflowTypeName := shuffle.String(testVisibilityWorkflowTypeName)
	workflowID := shuffle.String(testVisibilityWorkflowID)
	startTime := s.now()
	executionTime := startTime.Add(time.Second)
	status := int32(0)
	closeTime := (*time.Time)(nil)
	historyLength := (*int64)(nil)

	visibility := s.newRandomVisibilityRow(
		namespaceID,
		runID,
		workflowTypeName,
		workflowID,
		startTime,
		executionTime,
		status,
		closeTime,
		historyLength,
	)
	result, err := s.store.InsertIntoVisibility(newVisibilityContext(), &visibility)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	deleteFilter := sqlplugin.VisibilityDeleteFilter{
		NamespaceID: namespaceID.String(),
		RunID:       runID.String(),
	}
	result, err = s.store.DeleteFromVisibility(newVisibilityContext(), deleteFilter)
	s.NoError(err)
	rowsAffected, err = result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	getFilter := sqlplugin.VisibilityGetFilter{
		NamespaceID: namespaceID.String(),
		RunID:       runID.String(),
	}
	_, err = s.store.GetFromVisibility(newVisibilityContext(), getFilter)
	s.Error(err) // TODO persistence layer should do proper error translation
}

func (s *visibilitySuite) TestReplaceDeleteGet() {
	namespaceID := primitives.NewUUID()
	runID := primitives.NewUUID()
	workflowTypeName := shuffle.String(testVisibilityWorkflowTypeName)
	workflowID := shuffle.String(testVisibilityWorkflowID)
	startTime := s.now()
	executionTime := startTime.Add(time.Second)
	status := int32(0)
	closeTime := executionTime.Add(time.Second)
	historyLength := rand.Int63()

	visibility := s.newRandomVisibilityRow(
		namespaceID,
		runID,
		workflowTypeName,
		workflowID,
		startTime,
		executionTime,
		status,
		timePtr(closeTime),
		util.Ptr(historyLength),
	)
	result, err := s.store.ReplaceIntoVisibility(newVisibilityContext(), &visibility)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	deleteFilter := sqlplugin.VisibilityDeleteFilter{
		NamespaceID: namespaceID.String(),
		RunID:       runID.String(),
	}
	result, err = s.store.DeleteFromVisibility(newVisibilityContext(), deleteFilter)
	s.NoError(err)
	rowsAffected, err = result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	getFilter := sqlplugin.VisibilityGetFilter{
		NamespaceID: namespaceID.String(),
		RunID:       runID.String(),
	}
	_, err = s.store.GetFromVisibility(newVisibilityContext(), getFilter)
	s.Error(err) // TODO persistence layer should do proper error translation
}

func (s *visibilitySuite) TestSelect_MinStartTime_MaxStartTime_WorkflowID_StatusOpen_Single() {
	pageSize := 1

	namespaceID := primitives.NewUUID()
	runID := primitives.NewUUID()
	workflowTypeName := shuffle.String(testVisibilityWorkflowTypeName)
	workflowID := shuffle.String(testVisibilityWorkflowID)
	startTime := s.now()
	executionTime := startTime.Add(time.Second)
	status := int32(enumspb.WORKFLOW_EXECUTION_STATUS_RUNNING)
	closeTime := (*time.Time)(nil)
	historyLength := (*int64)(nil)

	visibility := s.newRandomVisibilityRow(
		namespaceID,
		runID,
		workflowTypeName,
		workflowID,
		startTime,
		executionTime,
		status,
		closeTime,
		historyLength,
	)
	result, err := s.store.InsertIntoVisibility(newVisibilityContext(), &visibility)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	minStartTime := startTime
	maxStartTime := startTime
	selectFilter := sqlplugin.VisibilitySelectFilter{
		NamespaceID:      namespaceID.String(),
		WorkflowID:       util.Ptr(workflowID),
		RunID:            util.Ptr(""),
		WorkflowTypeName: nil,
		MinTime:          timePtr(minStartTime),
		MaxTime:          timePtr(maxStartTime),
		Status:           int32(enumspb.WORKFLOW_EXECUTION_STATUS_RUNNING),
		PageSize:         util.Ptr(pageSize),
	}
	rows, err := s.store.SelectFromVisibility(newVisibilityContext(), selectFilter)
	s.NoError(err)
	for index := range rows {
		rows[index].NamespaceID = namespaceID.String()
	}
	s.Equal([]sqlplugin.VisibilityRow{visibility}, rows)
}

func (s *visibilitySuite) TestSelect_MinStartTime_MaxStartTime_WorkflowID_StatusOpen_Multiple() {
	numStartTime := 20
	visibilityPerStartTime := 4
	pageSize := 5

	var visibilities []sqlplugin.VisibilityRow

	namespaceID := primitives.NewUUID()
	workflowID := shuffle.String(testVisibilityWorkflowID)
	startTime := s.now()
	minStartTime := startTime
	maxStartTime := startTime.Add(time.Duration(numStartTime) * time.Second)
	executionTime := startTime.Add(time.Second)
	status := int32(enumspb.WORKFLOW_EXECUTION_STATUS_RUNNING)
	closeTime := (*time.Time)(nil)
	historyLength := (*int64)(nil)
	for i := 0; i < numStartTime; i++ {
		for j := 0; j < visibilityPerStartTime; j++ {
			runID := primitives.NewUUID()
			workflowTypeName := shuffle.String(testVisibilityWorkflowTypeName)
			visibility := s.newRandomVisibilityRow(
				namespaceID,
				runID,
				workflowTypeName,
				workflowID,
				startTime,
				executionTime,
				status,
				closeTime,
				historyLength,
			)
			result, err := s.store.InsertIntoVisibility(newVisibilityContext(), &visibility)
			s.NoError(err)
			rowsAffected, err := result.RowsAffected()
			s.NoError(err)
			s.Equal(1, int(rowsAffected))

			visibilities = append(visibilities, visibility)
		}

		startTime = startTime.Add(time.Second)
	}

	selectFilter := sqlplugin.VisibilitySelectFilter{
		NamespaceID:      namespaceID.String(),
		WorkflowID:       util.Ptr(workflowID),
		RunID:            util.Ptr(""),
		WorkflowTypeName: nil,
		MinTime:          timePtr(minStartTime),
		MaxTime:          timePtr(maxStartTime),
		Status:           int32(enumspb.WORKFLOW_EXECUTION_STATUS_RUNNING),
		PageSize:         util.Ptr(pageSize),
	}
	var rows []sqlplugin.VisibilityRow
	for {
		rowsPerPage, err := s.store.SelectFromVisibility(newVisibilityContext(), selectFilter)
		s.NoError(err)
		rows = append(rows, rowsPerPage...)

		if len(rowsPerPage) > 0 {
			lastVisibility := rowsPerPage[len(rowsPerPage)-1]
			selectFilter.MaxTime = timePtr(lastVisibility.StartTime)
			selectFilter.RunID = util.Ptr(lastVisibility.RunID)
		} else {
			break
		}
	}
	s.Len(rows, len(visibilities))
	s.sortByStartTimeDescRunIDAsc(visibilities)
	for index := range rows {
		rows[index].NamespaceID = namespaceID.String()
	}
	s.Equal(visibilities, rows)
}

func (s *visibilitySuite) TestSelect_MinStartTime_MaxStartTime_WorkflowID_StatusClose_Single() {
	pageSize := 1

	namespaceID := primitives.NewUUID()
	runID := primitives.NewUUID()
	workflowTypeName := shuffle.String(testVisibilityWorkflowTypeName)
	workflowID := shuffle.String(testVisibilityWorkflowID)
	startTime := s.now()
	executionTime := startTime.Add(time.Second)
	status := int32(enumspb.WORKFLOW_EXECUTION_STATUS_COMPLETED)
	closeTime := executionTime.Add(time.Second)
	historyLength := rand.Int63()

	visibility := s.newRandomVisibilityRow(
		namespaceID,
		runID,
		workflowTypeName,
		workflowID,
		startTime,
		executionTime,
		status,
		timePtr(closeTime),
		util.Ptr(historyLength),
	)
	result, err := s.store.ReplaceIntoVisibility(newVisibilityContext(), &visibility)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	minStartTime := closeTime
	maxStartTime := closeTime
	selectFilter := sqlplugin.VisibilitySelectFilter{
		NamespaceID:      namespaceID.String(),
		WorkflowID:       util.Ptr(workflowID),
		RunID:            util.Ptr(""),
		WorkflowTypeName: nil,
		MinTime:          timePtr(minStartTime),
		MaxTime:          timePtr(maxStartTime),
		Status:           int32(enumspb.WORKFLOW_EXECUTION_STATUS_UNSPECIFIED),
		PageSize:         util.Ptr(pageSize),
	}
	rows, err := s.store.SelectFromVisibility(newVisibilityContext(), selectFilter)
	s.NoError(err)
	for index := range rows {
		rows[index].NamespaceID = namespaceID.String()
	}
	s.Equal([]sqlplugin.VisibilityRow{visibility}, rows)
}

func (s *visibilitySuite) TestSelect_MinStartTime_MaxStartTime_WorkflowID_StatusClose_Multiple() {
	numStartTime := 20
	visibilityPerStartTime := 4
	pageSize := 5

	var visibilities []sqlplugin.VisibilityRow

	namespaceID := primitives.NewUUID()
	workflowID := shuffle.String(testVisibilityWorkflowID)
	startTime := s.now()
	executionTime := startTime.Add(time.Second)
	status := int32(enumspb.WORKFLOW_EXECUTION_STATUS_COMPLETED)
	closeTime := executionTime.Add(time.Second)
	historyLength := rand.Int63()
	minStartTime := closeTime
	maxStartTime := closeTime.Add(time.Duration(numStartTime) * time.Second)
	for i := 0; i < numStartTime; i++ {
		for j := 0; j < visibilityPerStartTime; j++ {
			runID := primitives.NewUUID()
			workflowTypeName := shuffle.String(testVisibilityWorkflowTypeName)
			visibility := s.newRandomVisibilityRow(
				namespaceID,
				runID,
				workflowTypeName,
				workflowID,
				startTime,
				executionTime,
				status,
				timePtr(closeTime),
				util.Ptr(historyLength),
			)
			result, err := s.store.ReplaceIntoVisibility(newVisibilityContext(), &visibility)
			s.NoError(err)
			rowsAffected, err := result.RowsAffected()
			s.NoError(err)
			s.Equal(1, int(rowsAffected))

			visibilities = append(visibilities, visibility)
		}
		closeTime = closeTime.Add(time.Second)
	}

	selectFilter := sqlplugin.VisibilitySelectFilter{
		NamespaceID:      namespaceID.String(),
		WorkflowID:       util.Ptr(workflowID),
		RunID:            util.Ptr(""),
		WorkflowTypeName: nil,
		MinTime:          timePtr(minStartTime),
		MaxTime:          timePtr(maxStartTime),
		Status:           int32(enumspb.WORKFLOW_EXECUTION_STATUS_UNSPECIFIED),
		PageSize:         util.Ptr(pageSize),
	}
	var rows []sqlplugin.VisibilityRow
	for {
		rowsPerPage, err := s.store.SelectFromVisibility(newVisibilityContext(), selectFilter)
		s.NoError(err)
		rows = append(rows, rowsPerPage...)

		if len(rowsPerPage) > 0 {
			lastVisibility := rowsPerPage[len(rowsPerPage)-1]
			selectFilter.MaxTime = lastVisibility.CloseTime
			selectFilter.RunID = util.Ptr(lastVisibility.RunID)
		} else {
			break
		}
	}
	s.Len(rows, len(visibilities))
	s.sortByCloseTimeDescRunIDAsc(visibilities)
	for index := range rows {
		rows[index].NamespaceID = namespaceID.String()
	}
	s.Equal(visibilities, rows)
}

func (s *visibilitySuite) TestSelect_MinStartTime_MaxStartTime_WorkflowTypeName_StatusOpen_Single() {
	pageSize := 1

	namespaceID := primitives.NewUUID()
	runID := primitives.NewUUID()
	workflowTypeName := shuffle.String(testVisibilityWorkflowTypeName)
	workflowID := shuffle.String(testVisibilityWorkflowID)
	startTime := s.now()
	executionTime := startTime.Add(time.Second)
	status := int32(enumspb.WORKFLOW_EXECUTION_STATUS_RUNNING)
	closeTime := (*time.Time)(nil)
	historyLength := (*int64)(nil)

	visibility := s.newRandomVisibilityRow(
		namespaceID,
		runID,
		workflowTypeName,
		workflowID,
		startTime,
		executionTime,
		status,
		closeTime,
		historyLength,
	)
	result, err := s.store.InsertIntoVisibility(newVisibilityContext(), &visibility)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	minStartTime := startTime
	maxStartTime := startTime
	selectFilter := sqlplugin.VisibilitySelectFilter{
		NamespaceID:      namespaceID.String(),
		WorkflowID:       nil,
		RunID:            util.Ptr(""),
		WorkflowTypeName: util.Ptr(workflowTypeName),
		MinTime:          timePtr(minStartTime),
		MaxTime:          timePtr(maxStartTime),
		Status:           int32(enumspb.WORKFLOW_EXECUTION_STATUS_RUNNING),
		PageSize:         util.Ptr(pageSize),
	}
	rows, err := s.store.SelectFromVisibility(newVisibilityContext(), selectFilter)
	s.NoError(err)
	for index := range rows {
		rows[index].NamespaceID = namespaceID.String()
	}
	s.Equal([]sqlplugin.VisibilityRow{visibility}, rows)
}

func (s *visibilitySuite) TestSelect_MinStartTime_MaxStartTime_WorkflowTypeName_StatusOpen_Multiple() {
	numStartTime := 20
	visibilityPerStartTime := 4
	pageSize := 5

	var visibilities []sqlplugin.VisibilityRow

	namespaceID := primitives.NewUUID()
	workflowTypeName := shuffle.String(testVisibilityWorkflowTypeName)
	startTime := s.now()
	minStartTime := startTime
	maxStartTime := startTime.Add(time.Duration(numStartTime) * time.Second)
	executionTime := startTime.Add(time.Second)
	status := int32(enumspb.WORKFLOW_EXECUTION_STATUS_RUNNING)
	closeTime := (*time.Time)(nil)
	historyLength := (*int64)(nil)
	for i := 0; i < numStartTime; i++ {
		for j := 0; j < visibilityPerStartTime; j++ {
			workflowID := shuffle.String(testVisibilityWorkflowID)
			runID := primitives.NewUUID()
			visibility := s.newRandomVisibilityRow(
				namespaceID,
				runID,
				workflowTypeName,
				workflowID,
				startTime,
				executionTime,
				status,
				closeTime,
				historyLength,
			)
			result, err := s.store.InsertIntoVisibility(newVisibilityContext(), &visibility)
			s.NoError(err)
			rowsAffected, err := result.RowsAffected()
			s.NoError(err)
			s.Equal(1, int(rowsAffected))

			visibilities = append(visibilities, visibility)
		}

		startTime = startTime.Add(time.Second)
	}

	selectFilter := sqlplugin.VisibilitySelectFilter{
		NamespaceID:      namespaceID.String(),
		WorkflowID:       nil,
		RunID:            util.Ptr(""),
		WorkflowTypeName: util.Ptr(workflowTypeName),
		MinTime:          timePtr(minStartTime),
		MaxTime:          timePtr(maxStartTime),
		Status:           int32(enumspb.WORKFLOW_EXECUTION_STATUS_RUNNING),
		PageSize:         util.Ptr(pageSize),
	}
	var rows []sqlplugin.VisibilityRow
	for {
		rowsPerPage, err := s.store.SelectFromVisibility(newVisibilityContext(), selectFilter)
		s.NoError(err)
		rows = append(rows, rowsPerPage...)

		if len(rowsPerPage) > 0 {
			lastVisibility := rowsPerPage[len(rowsPerPage)-1]
			selectFilter.MaxTime = timePtr(lastVisibility.StartTime)
			selectFilter.RunID = util.Ptr(lastVisibility.RunID)
		} else {
			break
		}
	}
	s.Len(rows, len(visibilities))
	s.sortByStartTimeDescRunIDAsc(visibilities)
	for index := range rows {
		rows[index].NamespaceID = namespaceID.String()
	}
	s.Equal(visibilities, rows)
}

func (s *visibilitySuite) TestSelect_MinStartTime_MaxStartTime_WorkflowTypeName_StatusClose_Single() {
	pageSize := 1

	namespaceID := primitives.NewUUID()
	runID := primitives.NewUUID()
	workflowTypeName := shuffle.String(testVisibilityWorkflowTypeName)
	workflowID := shuffle.String(testVisibilityWorkflowID)
	startTime := s.now()
	executionTime := startTime.Add(time.Second)
	status := int32(enumspb.WORKFLOW_EXECUTION_STATUS_COMPLETED)
	closeTime := executionTime.Add(time.Second)
	historyLength := rand.Int63()

	visibility := s.newRandomVisibilityRow(
		namespaceID,
		runID,
		workflowTypeName,
		workflowID,
		startTime,
		executionTime,
		status,
		timePtr(closeTime),
		util.Ptr(historyLength),
	)
	result, err := s.store.ReplaceIntoVisibility(newVisibilityContext(), &visibility)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	minStartTime := closeTime
	maxStartTime := closeTime
	selectFilter := sqlplugin.VisibilitySelectFilter{
		NamespaceID:      namespaceID.String(),
		WorkflowID:       nil,
		RunID:            util.Ptr(""),
		WorkflowTypeName: util.Ptr(workflowTypeName),
		MinTime:          timePtr(minStartTime),
		MaxTime:          timePtr(maxStartTime),
		Status:           int32(enumspb.WORKFLOW_EXECUTION_STATUS_UNSPECIFIED),
		PageSize:         util.Ptr(pageSize),
	}
	rows, err := s.store.SelectFromVisibility(newVisibilityContext(), selectFilter)
	s.NoError(err)
	for index := range rows {
		rows[index].NamespaceID = namespaceID.String()
	}
	s.Equal([]sqlplugin.VisibilityRow{visibility}, rows)
}

func (s *visibilitySuite) TestSelect_MinStartTime_MaxStartTime_WorkflowTypeName_StatusClose_Multiple() {
	numStartTime := 20
	visibilityPerStartTime := 4
	pageSize := 5

	var visibilities []sqlplugin.VisibilityRow

	namespaceID := primitives.NewUUID()
	workflowTypeName := shuffle.String(testVisibilityWorkflowTypeName)
	startTime := s.now()
	executionTime := startTime.Add(time.Second)
	status := int32(enumspb.WORKFLOW_EXECUTION_STATUS_COMPLETED)
	closeTime := executionTime.Add(time.Second)
	historyLength := rand.Int63()
	minStartTime := closeTime
	maxStartTime := closeTime.Add(time.Duration(numStartTime) * time.Second)
	for i := 0; i < numStartTime; i++ {
		for j := 0; j < visibilityPerStartTime; j++ {
			workflowID := shuffle.String(testVisibilityWorkflowID)
			runID := primitives.NewUUID()
			visibility := s.newRandomVisibilityRow(
				namespaceID,
				runID,
				workflowTypeName,
				workflowID,
				startTime,
				executionTime,
				status,
				timePtr(closeTime),
				util.Ptr(historyLength),
			)
			result, err := s.store.ReplaceIntoVisibility(newVisibilityContext(), &visibility)
			s.NoError(err)
			rowsAffected, err := result.RowsAffected()
			s.NoError(err)
			s.Equal(1, int(rowsAffected))

			visibilities = append(visibilities, visibility)
		}
		closeTime = closeTime.Add(time.Second)
	}

	selectFilter := sqlplugin.VisibilitySelectFilter{
		NamespaceID:      namespaceID.String(),
		WorkflowID:       nil,
		RunID:            util.Ptr(""),
		WorkflowTypeName: util.Ptr(workflowTypeName),
		MinTime:          timePtr(minStartTime),
		MaxTime:          timePtr(maxStartTime),
		Status:           int32(enumspb.WORKFLOW_EXECUTION_STATUS_UNSPECIFIED),
		PageSize:         util.Ptr(pageSize),
	}
	var rows []sqlplugin.VisibilityRow
	for {
		rowsPerPage, err := s.store.SelectFromVisibility(newVisibilityContext(), selectFilter)
		s.NoError(err)
		rows = append(rows, rowsPerPage...)

		if len(rowsPerPage) > 0 {
			lastVisibility := rowsPerPage[len(rowsPerPage)-1]
			selectFilter.MaxTime = lastVisibility.CloseTime
			selectFilter.RunID = util.Ptr(lastVisibility.RunID)
		} else {
			break
		}
	}
	s.Len(rows, len(visibilities))
	s.sortByCloseTimeDescRunIDAsc(visibilities)
	for index := range rows {
		rows[index].NamespaceID = namespaceID.String()
	}
	s.Equal(visibilities, rows)
}

func (s *visibilitySuite) TestSelect_MinStartTime_MaxStartTime_StatusOpen_Single() {
	pageSize := 1

	namespaceID := primitives.NewUUID()
	runID := primitives.NewUUID()
	workflowTypeName := shuffle.String(testVisibilityWorkflowTypeName)
	workflowID := shuffle.String(testVisibilityWorkflowID)
	startTime := s.now()
	executionTime := startTime.Add(time.Second)
	status := int32(enumspb.WORKFLOW_EXECUTION_STATUS_RUNNING)
	closeTime := (*time.Time)(nil)
	historyLength := (*int64)(nil)

	visibility := s.newRandomVisibilityRow(
		namespaceID,
		runID,
		workflowTypeName,
		workflowID,
		startTime,
		executionTime,
		status,
		closeTime,
		historyLength,
	)
	result, err := s.store.InsertIntoVisibility(newVisibilityContext(), &visibility)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	minStartTime := startTime
	maxStartTime := startTime
	selectFilter := sqlplugin.VisibilitySelectFilter{
		NamespaceID:      namespaceID.String(),
		WorkflowID:       nil,
		RunID:            util.Ptr(""),
		WorkflowTypeName: nil,
		MinTime:          timePtr(minStartTime),
		MaxTime:          timePtr(maxStartTime),
		Status:           int32(enumspb.WORKFLOW_EXECUTION_STATUS_RUNNING),
		PageSize:         util.Ptr(pageSize),
	}
	rows, err := s.store.SelectFromVisibility(newVisibilityContext(), selectFilter)
	s.NoError(err)
	for index := range rows {
		rows[index].NamespaceID = namespaceID.String()
	}
	s.Equal([]sqlplugin.VisibilityRow{visibility}, rows)
}

func (s *visibilitySuite) TestSelect_MinStartTime_MaxStartTime_StatusOpen_Multiple() {
	numStartTime := 20
	visibilityPerStartTime := 4
	pageSize := 5

	var visibilities []sqlplugin.VisibilityRow

	namespaceID := primitives.NewUUID()
	startTime := s.now()
	minStartTime := startTime
	maxStartTime := startTime.Add(time.Duration(numStartTime) * time.Second)
	executionTime := startTime.Add(time.Second)
	status := int32(enumspb.WORKFLOW_EXECUTION_STATUS_RUNNING)
	closeTime := (*time.Time)(nil)
	historyLength := (*int64)(nil)
	for i := 0; i < numStartTime; i++ {
		for j := 0; j < visibilityPerStartTime; j++ {
			workflowID := shuffle.String(testVisibilityWorkflowID)
			runID := primitives.NewUUID()
			workflowTypeName := shuffle.String(testVisibilityWorkflowTypeName)
			visibility := s.newRandomVisibilityRow(
				namespaceID,
				runID,
				workflowTypeName,
				workflowID,
				startTime,
				executionTime,
				status,
				closeTime,
				historyLength,
			)
			result, err := s.store.InsertIntoVisibility(newVisibilityContext(), &visibility)
			s.NoError(err)
			rowsAffected, err := result.RowsAffected()
			s.NoError(err)
			s.Equal(1, int(rowsAffected))

			visibilities = append(visibilities, visibility)
		}

		startTime = startTime.Add(time.Second)
	}

	selectFilter := sqlplugin.VisibilitySelectFilter{
		NamespaceID:      namespaceID.String(),
		WorkflowID:       nil,
		RunID:            util.Ptr(""),
		WorkflowTypeName: nil,
		MinTime:          timePtr(minStartTime),
		MaxTime:          timePtr(maxStartTime),
		Status:           int32(enumspb.WORKFLOW_EXECUTION_STATUS_RUNNING),
		PageSize:         util.Ptr(pageSize),
	}
	var rows []sqlplugin.VisibilityRow
	for {
		rowsPerPage, err := s.store.SelectFromVisibility(newVisibilityContext(), selectFilter)
		s.NoError(err)
		rows = append(rows, rowsPerPage...)

		if len(rowsPerPage) > 0 {
			lastVisibility := rowsPerPage[len(rowsPerPage)-1]
			selectFilter.MaxTime = timePtr(lastVisibility.StartTime)
			selectFilter.RunID = util.Ptr(lastVisibility.RunID)
		} else {
			break
		}
	}
	s.Len(rows, len(visibilities))
	s.sortByStartTimeDescRunIDAsc(visibilities)
	for index := range rows {
		rows[index].NamespaceID = namespaceID.String()
	}
	s.Equal(visibilities, rows)
}

func (s *visibilitySuite) TestSelect_MinStartTime_MaxStartTime_StatusClose_Single() {
	for _, status := range testVisibilityCloseStatus {
		s.testSelectMinStartTimeMaxStartTimeStatusCloseSingle(status)
	}
}

func (s *visibilitySuite) testSelectMinStartTimeMaxStartTimeStatusCloseSingle(
	status enumspb.WorkflowExecutionStatus,
) {
	pageSize := 1

	namespaceID := primitives.NewUUID()
	runID := primitives.NewUUID()
	workflowTypeName := shuffle.String(testVisibilityWorkflowTypeName)
	workflowID := shuffle.String(testVisibilityWorkflowID)
	startTime := s.now()
	executionTime := startTime.Add(time.Second)
	closeTime := executionTime.Add(time.Second)
	historyLength := rand.Int63()

	visibility := s.newRandomVisibilityRow(
		namespaceID,
		runID,
		workflowTypeName,
		workflowID,
		startTime,
		executionTime,
		int32(status),
		timePtr(closeTime),
		util.Ptr(historyLength),
	)
	result, err := s.store.ReplaceIntoVisibility(newVisibilityContext(), &visibility)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	minStartTime := closeTime
	maxStartTime := closeTime
	selectFilter := sqlplugin.VisibilitySelectFilter{
		NamespaceID:      namespaceID.String(),
		WorkflowID:       nil,
		RunID:            util.Ptr(""),
		WorkflowTypeName: nil,
		MinTime:          timePtr(minStartTime),
		MaxTime:          timePtr(maxStartTime),
		Status:           int32(enumspb.WORKFLOW_EXECUTION_STATUS_UNSPECIFIED),
		PageSize:         util.Ptr(pageSize),
	}
	rows, err := s.store.SelectFromVisibility(newVisibilityContext(), selectFilter)
	s.NoError(err)
	for index := range rows {
		rows[index].NamespaceID = namespaceID.String()
	}
	s.Equal([]sqlplugin.VisibilityRow{visibility}, rows)
}

func (s *visibilitySuite) TestSelect_MinStartTime_MaxStartTime_StatusClose_Multiple() {
	numStartTime := 20
	visibilityPerStartTime := 4
	pageSize := 5

	var visibilities []sqlplugin.VisibilityRow

	namespaceID := primitives.NewUUID()
	startTime := s.now()
	executionTime := startTime.Add(time.Second)

	closeTime := executionTime.Add(time.Second)
	historyLength := rand.Int63()
	minStartTime := closeTime
	maxStartTime := closeTime.Add(time.Duration(numStartTime) * time.Second)
	for i := 0; i < numStartTime; i++ {
		for j := 0; j < visibilityPerStartTime; j++ {
			workflowID := shuffle.String(testVisibilityWorkflowID)
			runID := primitives.NewUUID()
			workflowTypeName := shuffle.String(testVisibilityWorkflowTypeName)
			status := int32(testVisibilityCloseStatus[rand.Intn(len(testVisibilityCloseStatus))])
			visibility := s.newRandomVisibilityRow(
				namespaceID,
				runID,
				workflowTypeName,
				workflowID,
				startTime,
				executionTime,
				status,
				timePtr(closeTime),
				util.Ptr(historyLength),
			)
			result, err := s.store.ReplaceIntoVisibility(newVisibilityContext(), &visibility)
			s.NoError(err)
			rowsAffected, err := result.RowsAffected()
			s.NoError(err)
			s.Equal(1, int(rowsAffected))

			visibilities = append(visibilities, visibility)
		}
		closeTime = closeTime.Add(time.Second)
	}

	selectFilter := sqlplugin.VisibilitySelectFilter{
		NamespaceID:      namespaceID.String(),
		WorkflowID:       nil,
		RunID:            util.Ptr(""),
		WorkflowTypeName: nil,
		MinTime:          timePtr(minStartTime),
		MaxTime:          timePtr(maxStartTime),
		Status:           int32(enumspb.WORKFLOW_EXECUTION_STATUS_UNSPECIFIED),
		PageSize:         util.Ptr(pageSize),
	}
	var rows []sqlplugin.VisibilityRow
	for {
		rowsPerPage, err := s.store.SelectFromVisibility(newVisibilityContext(), selectFilter)
		s.NoError(err)
		rows = append(rows, rowsPerPage...)

		if len(rowsPerPage) > 0 {
			lastVisibility := rowsPerPage[len(rowsPerPage)-1]
			selectFilter.MaxTime = lastVisibility.CloseTime
			selectFilter.RunID = util.Ptr(lastVisibility.RunID)
		} else {
			break
		}
	}
	s.Len(rows, len(visibilities))
	s.sortByCloseTimeDescRunIDAsc(visibilities)
	for index := range rows {
		rows[index].NamespaceID = namespaceID.String()
	}
	s.Equal(visibilities, rows)
}

func (s *visibilitySuite) TestSelect_MinStartTime_MaxStartTime_StatusCloseByType_Single() {
	for _, status := range testVisibilityCloseStatus {
		s.testSelectMinStartTimeMaxStartTimeStatusCloseByTypeSingle(status)
	}
}

func (s *visibilitySuite) testSelectMinStartTimeMaxStartTimeStatusCloseByTypeSingle(
	status enumspb.WorkflowExecutionStatus,
) {
	pageSize := 1

	namespaceID := primitives.NewUUID()
	runID := primitives.NewUUID()
	workflowTypeName := shuffle.String(testVisibilityWorkflowTypeName)
	workflowID := shuffle.String(testVisibilityWorkflowID)
	startTime := s.now()
	executionTime := startTime.Add(time.Second)
	closeTime := executionTime.Add(time.Second)
	historyLength := rand.Int63()

	visibility := s.newRandomVisibilityRow(
		namespaceID,
		runID,
		workflowTypeName,
		workflowID,
		startTime,
		executionTime,
		int32(status),
		timePtr(closeTime),
		util.Ptr(historyLength),
	)
	result, err := s.store.ReplaceIntoVisibility(newVisibilityContext(), &visibility)
	s.NoError(err)
	rowsAffected, err := result.RowsAffected()
	s.NoError(err)
	s.Equal(1, int(rowsAffected))

	minStartTime := closeTime
	maxStartTime := closeTime
	selectFilter := sqlplugin.VisibilitySelectFilter{
		NamespaceID:      namespaceID.String(),
		WorkflowID:       nil,
		RunID:            util.Ptr(""),
		WorkflowTypeName: util.Ptr(workflowTypeName),
		MinTime:          timePtr(minStartTime),
		MaxTime:          timePtr(maxStartTime),
		Status:           int32(status),
		PageSize:         util.Ptr(pageSize),
	}
	rows, err := s.store.SelectFromVisibility(newVisibilityContext(), selectFilter)
	s.NoError(err)
	for index := range rows {
		rows[index].NamespaceID = namespaceID.String()
	}
	s.Equal([]sqlplugin.VisibilityRow{visibility}, rows)
}

func (s *visibilitySuite) TestSelect_MinStartTime_MaxStartTime_StatusCloseByType_Multiple() {
	for _, status := range testVisibilityCloseStatus {
		s.testSelectMinStartTimeMaxStartTimeStatusCloseByTypeMultiple(status)
	}
}

func (s *visibilitySuite) testSelectMinStartTimeMaxStartTimeStatusCloseByTypeMultiple(
	status enumspb.WorkflowExecutionStatus,
) {
	numStartTime := 20
	visibilityPerStartTime := 4
	pageSize := 5

	var visibilities []sqlplugin.VisibilityRow

	namespaceID := primitives.NewUUID()
	startTime := s.now()
	executionTime := startTime.Add(time.Second)
	closeTime := executionTime.Add(time.Second)
	historyLength := rand.Int63()
	minStartTime := closeTime
	maxStartTime := closeTime.Add(time.Duration(numStartTime) * time.Second)
	for i := 0; i < numStartTime; i++ {
		for j := 0; j < visibilityPerStartTime; j++ {
			workflowID := shuffle.String(testVisibilityWorkflowID)
			runID := primitives.NewUUID()
			workflowTypeName := shuffle.String(testVisibilityWorkflowTypeName)
			visibility := s.newRandomVisibilityRow(
				namespaceID,
				runID,
				workflowTypeName,
				workflowID,
				startTime,
				executionTime,
				int32(status),
				timePtr(closeTime),
				util.Ptr(historyLength),
			)
			result, err := s.store.ReplaceIntoVisibility(newVisibilityContext(), &visibility)
			s.NoError(err)
			rowsAffected, err := result.RowsAffected()
			s.NoError(err)
			s.Equal(1, int(rowsAffected))

			visibilities = append(visibilities, visibility)
		}
		closeTime = closeTime.Add(time.Second)
	}

	selectFilter := sqlplugin.VisibilitySelectFilter{
		NamespaceID:      namespaceID.String(),
		WorkflowID:       nil,
		RunID:            util.Ptr(""),
		WorkflowTypeName: nil,
		MinTime:          timePtr(minStartTime),
		MaxTime:          timePtr(maxStartTime),
		Status:           int32(status),
		PageSize:         util.Ptr(pageSize),
	}
	var rows []sqlplugin.VisibilityRow
	for {
		rowsPerPage, err := s.store.SelectFromVisibility(newVisibilityContext(), selectFilter)
		s.NoError(err)
		rows = append(rows, rowsPerPage...)

		if len(rowsPerPage) > 0 {
			lastVisibility := rowsPerPage[len(rowsPerPage)-1]
			selectFilter.MaxTime = lastVisibility.CloseTime
			selectFilter.RunID = util.Ptr(lastVisibility.RunID)
		} else {
			break
		}
	}

	s.Len(rows, len(visibilities))
	s.sortByCloseTimeDescRunIDAsc(visibilities)
	for index := range rows {
		rows[index].NamespaceID = namespaceID.String()
	}
	s.Equal(visibilities, rows)
}

func (s *visibilitySuite) sortByStartTimeDescRunIDAsc(
	visibilities []sqlplugin.VisibilityRow,
) {
	sort.Slice(visibilities, func(i, j int) bool {
		this := visibilities[i]
		that := visibilities[j]

		// start time desc, run ID asc

		if this.StartTime.Before(that.StartTime) {
			return false
		} else if that.StartTime.Before(this.StartTime) {
			return true
		}

		if this.RunID < that.RunID {
			return true
		} else if this.RunID > that.RunID {
			return false
		}

		// same
		return true
	})
}

func (s *visibilitySuite) sortByCloseTimeDescRunIDAsc(
	visibilities []sqlplugin.VisibilityRow,
) {
	sort.Slice(visibilities, func(i, j int) bool {
		this := visibilities[i]
		that := visibilities[j]

		// close time desc, run ID asc

		if this.CloseTime.Before(*that.CloseTime) {
			return false
		} else if that.CloseTime.Before(*this.CloseTime) {
			return true
		}

		if this.RunID < that.RunID {
			return true
		} else if this.RunID > that.RunID {
			return false
		}

		// same
		return true
	})
}

func (s *visibilitySuite) now() time.Time {
	return time.Now().UTC().Truncate(time.Millisecond)
}

func (s *visibilitySuite) newRandomVisibilityRow(
	namespaceID primitives.UUID,
	runID primitives.UUID,
	workflowTypeName string,
	workflowID string,
	startTime time.Time,
	executionTime time.Time,
	status int32,
	closeTime *time.Time,
	historyLength *int64,
) sqlplugin.VisibilityRow {
	s.version++
	return sqlplugin.VisibilityRow{
		NamespaceID:      namespaceID.String(),
		RunID:            runID.String(),
		WorkflowTypeName: workflowTypeName,
		WorkflowID:       workflowID,
		StartTime:        startTime,
		ExecutionTime:    executionTime,
		Status:           status,
		CloseTime:        closeTime,
		HistoryLength:    historyLength,
		Memo:             shuffle.Bytes(testVisibilityData),
		Encoding:         testVisibilityEncoding,
		Version:          s.version,
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
