# CommandResults

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ActivityResults** | Pointer to [**[]ActivityResult**](ActivityResult.md) |  | [optional] 
**SignalResults** | Pointer to [**[]SignalResult**](SignalResult.md) |  | [optional] 
**InterStateChannelResults** | Pointer to [**[]InterStateChannelResult**](InterStateChannelResult.md) |  | [optional] 
**TimerResults** | Pointer to [**[]TimerResult**](TimerResult.md) |  | [optional] 
**WaitForQueryAttributeChangeResults** | Pointer to [**[]WaitForQueryAttributeChangeResult**](WaitForQueryAttributeChangeResult.md) |  | [optional] 

## Methods

### NewCommandResults

`func NewCommandResults() *CommandResults`

NewCommandResults instantiates a new CommandResults object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommandResultsWithDefaults

`func NewCommandResultsWithDefaults() *CommandResults`

NewCommandResultsWithDefaults instantiates a new CommandResults object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetActivityResults

`func (o *CommandResults) GetActivityResults() []ActivityResult`

GetActivityResults returns the ActivityResults field if non-nil, zero value otherwise.

### GetActivityResultsOk

`func (o *CommandResults) GetActivityResultsOk() (*[]ActivityResult, bool)`

GetActivityResultsOk returns a tuple with the ActivityResults field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActivityResults

`func (o *CommandResults) SetActivityResults(v []ActivityResult)`

SetActivityResults sets ActivityResults field to given value.

### HasActivityResults

`func (o *CommandResults) HasActivityResults() bool`

HasActivityResults returns a boolean if a field has been set.

### GetSignalResults

`func (o *CommandResults) GetSignalResults() []SignalResult`

GetSignalResults returns the SignalResults field if non-nil, zero value otherwise.

### GetSignalResultsOk

`func (o *CommandResults) GetSignalResultsOk() (*[]SignalResult, bool)`

GetSignalResultsOk returns a tuple with the SignalResults field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSignalResults

`func (o *CommandResults) SetSignalResults(v []SignalResult)`

SetSignalResults sets SignalResults field to given value.

### HasSignalResults

`func (o *CommandResults) HasSignalResults() bool`

HasSignalResults returns a boolean if a field has been set.

### GetInterStateChannelResults

`func (o *CommandResults) GetInterStateChannelResults() []InterStateChannelResult`

GetInterStateChannelResults returns the InterStateChannelResults field if non-nil, zero value otherwise.

### GetInterStateChannelResultsOk

`func (o *CommandResults) GetInterStateChannelResultsOk() (*[]InterStateChannelResult, bool)`

GetInterStateChannelResultsOk returns a tuple with the InterStateChannelResults field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInterStateChannelResults

`func (o *CommandResults) SetInterStateChannelResults(v []InterStateChannelResult)`

SetInterStateChannelResults sets InterStateChannelResults field to given value.

### HasInterStateChannelResults

`func (o *CommandResults) HasInterStateChannelResults() bool`

HasInterStateChannelResults returns a boolean if a field has been set.

### GetTimerResults

`func (o *CommandResults) GetTimerResults() []TimerResult`

GetTimerResults returns the TimerResults field if non-nil, zero value otherwise.

### GetTimerResultsOk

`func (o *CommandResults) GetTimerResultsOk() (*[]TimerResult, bool)`

GetTimerResultsOk returns a tuple with the TimerResults field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimerResults

`func (o *CommandResults) SetTimerResults(v []TimerResult)`

SetTimerResults sets TimerResults field to given value.

### HasTimerResults

`func (o *CommandResults) HasTimerResults() bool`

HasTimerResults returns a boolean if a field has been set.

### GetWaitForQueryAttributeChangeResults

`func (o *CommandResults) GetWaitForQueryAttributeChangeResults() []WaitForQueryAttributeChangeResult`

GetWaitForQueryAttributeChangeResults returns the WaitForQueryAttributeChangeResults field if non-nil, zero value otherwise.

### GetWaitForQueryAttributeChangeResultsOk

`func (o *CommandResults) GetWaitForQueryAttributeChangeResultsOk() (*[]WaitForQueryAttributeChangeResult, bool)`

GetWaitForQueryAttributeChangeResultsOk returns a tuple with the WaitForQueryAttributeChangeResults field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWaitForQueryAttributeChangeResults

`func (o *CommandResults) SetWaitForQueryAttributeChangeResults(v []WaitForQueryAttributeChangeResult)`

SetWaitForQueryAttributeChangeResults sets WaitForQueryAttributeChangeResults field to given value.

### HasWaitForQueryAttributeChangeResults

`func (o *CommandResults) HasWaitForQueryAttributeChangeResults() bool`

HasWaitForQueryAttributeChangeResults returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


