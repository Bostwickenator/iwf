package integ

import (
	"context"
	"fmt"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/indeedeng/iwf/service/common/timeparser"
	"strconv"
	"testing"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/persistence"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
)

func TestPersistenceWorkflowTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestPersistenceWorkflow(t, service.BackendTypeTemporal)
		time.Sleep(time.Millisecond * time.Duration(*repeatInterval))
	}
}

func TestPersistenceWorkflowCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestPersistenceWorkflow(t, service.BackendTypeCadence)
		time.Sleep(time.Millisecond * time.Duration(*repeatInterval))
	}
}

func doTestPersistenceWorkflow(t *testing.T, backendType service.BackendType) {
	assertions := assert.New(t)
	wfHandler := persistence.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler)
	defer closeFunc1()

	closeFunc2 := startIwfService(backendType)
	defer closeFunc2()

	// start a workflow
	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: "http://localhost:" + testIwfServerPort,
			},
		},
	})
	wfId := persistence.WorkflowType + strconv.Itoa(int(time.Now().Unix()))
	nowTime := time.Now()
	nowTimeStr := nowTime.Format(timeparser.DateTimeFormat)

	reqStart := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	wfReq := iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        persistence.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           persistence.State1,
		StateOptions: &iwfidl.WorkflowStateOptions{
			SearchAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
				PersistenceLoadingType: ptr.Any(iwfidl.ALL_WITHOUT_LOCKING),
			},
			DataObjectsLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
				PersistenceLoadingType: ptr.Any(iwfidl.ALL_WITHOUT_LOCKING),
			},
		},
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			SearchAttributes: []iwfidl.SearchAttribute{
				{
					Key:         iwfidl.PtrString("CustomDatetimeField"),
					ValueType:   ptr.Any(iwfidl.DATETIME),
					StringValue: iwfidl.PtrString(nowTimeStr),
				},
			},
		},
	}
	_, httpResp, err := reqStart.WorkflowStartRequest(wfReq).Execute()
	panicAtHttpError(err, httpResp)

	reqWait := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	wfResponse, httpResp, err := reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp)

	reqQry := apiClient.DefaultApi.ApiV1WorkflowDataobjectsGetPost(context.Background())
	queryResult1, httpResp, err := reqQry.WorkflowGetDataObjectsRequest(iwfidl.WorkflowGetDataObjectsRequest{
		WorkflowId: wfId,
		Keys: []string{
			persistence.TestDataObjectKey,
		},
	}).Execute()
	panicAtHttpError(err, httpResp)

	queryResult2, httpResp, err := reqQry.WorkflowGetDataObjectsRequest(iwfidl.WorkflowGetDataObjectsRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp)

	reqSearch := apiClient.DefaultApi.ApiV1WorkflowSearchattributesGetPost(context.Background())
	searchResult1, httpResp, err := reqSearch.WorkflowGetSearchAttributesRequest(iwfidl.WorkflowGetSearchAttributesRequest{
		WorkflowId:    wfId,
		WorkflowRunId: &wfResponse.WorkflowRunId,
		Keys: []iwfidl.SearchAttributeKeyAndType{
			{
				Key:       iwfidl.PtrString(persistence.TestSearchAttributeKeywordKey),
				ValueType: ptr.Any(iwfidl.KEYWORD),
			},
		},
	}).Execute()
	panicAtHttpError(err, httpResp)

	searchResult2, httpResp, err := reqSearch.WorkflowGetSearchAttributesRequest(iwfidl.WorkflowGetSearchAttributesRequest{
		WorkflowId:    wfId,
		WorkflowRunId: &wfResponse.WorkflowRunId,
		Keys: []iwfidl.SearchAttributeKeyAndType{
			{
				Key:       iwfidl.PtrString(persistence.TestSearchAttributeIntKey),
				ValueType: ptr.Any(iwfidl.INT),
			},
		},
	}).Execute()
	panicAtHttpError(err, httpResp)

	// assertion
	history, data := wfHandler.GetTestResult()
	assertions.Equalf(map[string]int64{
		"S1_start":  1,
		"S1_decide": 1,
		"S2_start":  1,
		"S2_decide": 1,
		"S3_start":  1,
		"S3_decide": 1,
	}, history, "attribute test fail, %v", history)

	assertions.Equal(map[string]interface{}{
		"S1_decide_intSaFounds": 1,
		"S1_decide_kwSaFounds":  1,
		"S2_decide_intSaFounds": 1,
		"S2_decide_kwSaFounds":  1,
		"S2_start_intSaFounds":  1,
		"S2_start_kwSaFounds":   1,

		"S1_decide_localAttFound": true,
		"S1_decide_queryAttFound": 2,
		"S2_decide_queryAttFound": true,
		"S2_start_queryAttFound":  true,
	}, data)

	expected1 := []iwfidl.KeyValue{
		{
			Key:   iwfidl.PtrString(persistence.TestDataObjectKey),
			Value: &persistence.TestDataObjectVal2,
		},
	}
	expected2 := []iwfidl.KeyValue{
		{
			Key:   iwfidl.PtrString(persistence.TestDataObjectKey),
			Value: &persistence.TestDataObjectVal2,
		},
		{
			Key:   iwfidl.PtrString(persistence.TestDataObjectKey2),
			Value: &persistence.TestDataObjectVal1,
		},
	}
	assertions.ElementsMatch(expected1, queryResult1.GetObjects())
	assertions.ElementsMatch(expected2, queryResult2.GetObjects())

	expectedSearchAttributeInt := iwfidl.SearchAttribute{
		Key:          iwfidl.PtrString(persistence.TestSearchAttributeIntKey),
		ValueType:    ptr.Any(iwfidl.INT),
		IntegerValue: iwfidl.PtrInt64(persistence.TestSearchAttributeIntValue2),
	}

	expectedSearchAttributeKeyword := iwfidl.SearchAttribute{
		Key:         iwfidl.PtrString(persistence.TestSearchAttributeKeywordKey),
		ValueType:   ptr.Any(iwfidl.KEYWORD),
		StringValue: iwfidl.PtrString(persistence.TestSearchAttributeKeywordValue2),
	}

	assertions.Equal([]iwfidl.SearchAttribute{expectedSearchAttributeKeyword}, searchResult1.GetSearchAttributes())
	assertions.Equal([]iwfidl.SearchAttribute{expectedSearchAttributeInt}, searchResult2.GetSearchAttributes())

	if *testSearchIntegTest {
		// start more WFs in order to test pagination
		firstWfId := wfReq.WorkflowId
		wfReq.WorkflowId = firstWfId + "-1"
		newSa := iwfidl.SearchAttribute{
			Key:       iwfidl.PtrString("CustomBoolField"),
			ValueType: ptr.Any(iwfidl.BOOL),
			BoolValue: ptr.Any(true),
		}
		wfReq.WorkflowStartOptions.SearchAttributes = append(wfReq.WorkflowStartOptions.SearchAttributes, newSa)

		_, httpResp, err := reqStart.WorkflowStartRequest(wfReq).Execute()
		panicAtHttpError(err, httpResp)

		wfReq.WorkflowId = firstWfId + "-2"
		newSa = iwfidl.SearchAttribute{
			Key:         iwfidl.PtrString("CustomDoubleField"),
			ValueType:   ptr.Any(iwfidl.DOUBLE),
			DoubleValue: ptr.Any(0.01),
		}
		wfReq.WorkflowStartOptions.SearchAttributes = append(wfReq.WorkflowStartOptions.SearchAttributes, newSa)

		_, httpResp, err = reqStart.WorkflowStartRequest(wfReq).Execute()
		panicAtHttpError(err, httpResp)

		wfReq.WorkflowId = firstWfId + "-3"
		newSa = iwfidl.SearchAttribute{
			Key:              iwfidl.PtrString("CustomKeywordField"),
			ValueType:        ptr.Any(iwfidl.KEYWORD_ARRAY),
			StringArrayValue: []string{"keyword1", "keyword2"},
		}
		wfReq.WorkflowStartOptions.SearchAttributes = append(wfReq.WorkflowStartOptions.SearchAttributes, newSa)
		_, httpResp, err = reqStart.WorkflowStartRequest(wfReq).Execute()
		panicAtHttpError(err, httpResp)

		wfReq.WorkflowId = firstWfId + "-4"
		newSa = iwfidl.SearchAttribute{
			Key:         iwfidl.PtrString("CustomStringField"),
			ValueType:   ptr.Any(iwfidl.TEXT),
			StringValue: iwfidl.PtrString("My name is Quanzheng Long"),
		}
		wfReq.WorkflowStartOptions.SearchAttributes = append(wfReq.WorkflowStartOptions.SearchAttributes, newSa)
		_, httpResp, err = reqStart.WorkflowStartRequest(wfReq).Execute()
		panicAtHttpError(err, httpResp)

		// wait for all completed
		_, httpResp, err = reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
			WorkflowId: firstWfId + "-1",
		}).Execute()
		panicAtHttpError(err, httpResp)
		_, httpResp, err = reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
			WorkflowId: firstWfId + "-2",
		}).Execute()
		panicAtHttpError(err, httpResp)
		_, httpResp, err = reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
			WorkflowId: firstWfId + "-3",
		}).Execute()
		panicAtHttpError(err, httpResp)
		_, httpResp, err = reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
			WorkflowId: firstWfId + "-4",
		}).Execute()
		panicAtHttpError(err, httpResp)

		// wait for the search attribute index to be ready in ElasticSearch
		time.Sleep(time.Duration(*searchWaitTimeIntegTest) * time.Millisecond)

		assertSearch(fmt.Sprintf("CustomDatetimeField='%v'", nowTimeStr), 5, apiClient, assertions)
		assertSearch(fmt.Sprintf("CustomDatetimeField='%v' AND CustomStringField='%v'", nowTimeStr, "Quanzheng"), 1, apiClient, assertions)
		assertSearch(fmt.Sprintf("CustomDatetimeField='%v' AND CustomDoubleField='%v'", nowTimeStr, "0.01"), 3, apiClient, assertions)
		assertSearch(fmt.Sprintf("CustomDatetimeField='%v' AND CustomBoolField='%v'", nowTimeStr, "true"), 0, apiClient, assertions) // this got changed during WF execution

		// TODO?? research how to use text
		//assertSearch(fmt.Sprintf("CustomDatetimeField='%v' AND CustomKeywordField='%v'", nowTimeStrForSearch, "keyword-value1"), 5, apiClient, assertions)
	}
}

func assertSearch(query string, expected int, apiClient *iwfidl.APIClient, assertions *assert.Assertions) {
	// search through all wfs using search API with pagination
	search := apiClient.DefaultApi.ApiV1WorkflowSearchPost(context.Background())

	var nextPageToken string
	current := 0
	for current < expected {
		searchResp, httpResp, err := search.WorkflowSearchRequest(iwfidl.WorkflowSearchRequest{
			Query:         query,
			PageSize:      iwfidl.PtrInt32(2),
			NextPageToken: &nextPageToken,
		}).Execute()
		panicAtHttpError(err, httpResp)

		current += len(searchResp.WorkflowExecutions)
		if current < expected {
			assertions.Equal(2, len(searchResp.WorkflowExecutions))
			assertions.True(len(searchResp.GetNextPageToken()) > 0)
			nextPageToken = *searchResp.NextPageToken
		} else if current == expected {
			if len(searchResp.GetNextPageToken()) > 0 {
				nextPageToken = *searchResp.NextPageToken
				// the next page must be empty
				searchResp, httpResp, err := search.WorkflowSearchRequest(iwfidl.WorkflowSearchRequest{
					Query:         query,
					PageSize:      iwfidl.PtrInt32(2),
					NextPageToken: &nextPageToken,
				}).Execute()
				panicAtHttpError(err, httpResp)
				assertions.Equal(0, len(searchResp.WorkflowExecutions))
				assertions.True(len(searchResp.GetNextPageToken()) == 0)
			}
		} else {
			assertions.Fail("cannot happen")
		}
	}

}
