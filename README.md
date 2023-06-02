# iWF project - main & server repo

[![Go Reference](https://pkg.go.dev/badge/github.com/indeedeng/iwf.svg)](https://pkg.go.dev/github.com/indeedeng/iwf)
[![Go Report Card](https://goreportcard.com/badge/github.com/indeedeng/iwf)](https://goreportcard.com/report/github.com/indeedeng/iwf)
[![Coverage Status](https://codecov.io/github/indeedeng/iwf/coverage.svg?branch=main)](https://app.codecov.io/gh/indeedeng/iwf/branch/main)

[![Build status](https://github.com/indeedeng/iwf/actions/workflows/ci-cadence-integ-test.yml/badge.svg?branch=main)](https://github.com/indeedeng/iwf/actions/workflows/ci-cadence-integ-test.yml)
[![Build status](https://github.com/indeedeng/iwf/actions/workflows/ci-temporal-integ-test.yml/badge.svg?branch=main)](https://github.com/indeedeng/iwf/actions/workflows/ci-temporal-integ-test.yml)

**iWF will make you a 10x developer!**

iWF is a platform for developing resilient, fault-tolerant, scalable long-running applications. 
It offers a convenient abstraction for durable timers, background execution with backoff retry, 
persistence, indexing, message queues, RPC, and more. You will build long-running reliable processes faster than ever. 

iWF is built on top of [Cadence](https://github.com/uber/cadence)/[Temporal](https://github.com/temporalio/temporal).

Related projects:

* [iWF Java SDK](https://github.com/indeedeng/iwf-java-sdk)
* [iWF Java Samples](https://github.com/indeedeng/iwf-java-samples)
* [iWF Golang SDK](https://github.com/indeedeng/iwf-golang-sdk)
* [iWF Golang Samples](https://github.com/indeedeng/iwf-golang-samples)
* WIP [iWF Python SDK](https://github.com/indeedeng/iwf-python-sdk)
* WIP [iWF TypeScript SDK](https://github.com/indeedeng/iwf-ts-sdk)
# What is iWF


## Basic Concepts


The top level concept is **`ObjectWorkflow`** -- nearly any "object" can be an ObjectWorkflow, as long as it's long-lasting, at least a few seconds. 

User application creates ObjectWorkflow by implementing the Workflow interface, e.g. in
[Java](https://github.com/indeedeng/iwf-java-sdk/blob/main/src/main/java/io/iworkflow/core/ObjectWorkflow.java)
or [Golang](https://github.com/indeedeng/iwf-golang-sdk/blob/main/iwf/workflow.go).
An implementation of the interface is referred to as a `WorkflowDefinition`, consisting below components:

| Name             | Description                                                                                                                                          | 
|:-----------------|:-----------------------------------------------------------------------------------------------------------------------------------------------------| 
| Data Attribute   | Persistence field to storing data                                                                                                                    | 
| Search Attribute | "Searchable data attribute" -- attribute data is persisted and also indexed in search engine backed by ElasticSearch or OpenSearch                   | 
| Signal Channel   | Asynchronous message queue for the workflow object to receive message from external                                                                  |
| Internal Channel | "Internal Signal Channel" -- An internal message queue for workflow states/RPC                                                                       |
| Workflow State   | A background execution unit. State is super powerful like a small workflow of two steps: waitUntil(optional) and execute with default infinite retry |
| RPC              | Remote procedure call. Invoked by client, executed in worker, and interact with data/search attributes, internal channel and state execution         |

A workflow definition can be outlined like this:

![Example workflow diagram](https://user-images.githubusercontent.com/4523955/234424825-ff3673c0-af23-4eb7-887d-b1f421f3aaa4.png)

These are all the concepts that you need to build a super complicated workflow.
See this engagement workflow example in [Java](https://github.com/indeedeng/iwf-java-samples/tree/main/src/main/java/io/iworkflow/workflow/engagement)
or [Golang](https://github.com/indeedeng/iwf-golang-samples/tree/main/workflows/engagement)
for how it looks like!

Below are the detailed explanation of the concepts. 
They are powerful, also extremely simple to learn and use (as the philosophy of iWF).

## Persistence
Both data and search attributes are defined as "persistence schema". 
The schema just defined and maintained in the code along with other business logic.
Search attribute works like infinite indexes in traditional database. You
only need to specify which attributes should be indexed, without worrying about things in
a traditional database like the number of indexes, and the order of the fields in an index.

Logically, this workflow definition will have a persistence schema like below:

| Workflow Execution   | Search Attr A | Search Attr B | Data Attr C | Data Attr D |
|----------------------|---------------|:-------------:|------------:|------------:|
| Workflow Execution 1 | val 1         |     val 2     |       val 3 |       val 4 |
| Workflow Execution 2 | val 5         |     val 6     |       val 7 |       val 8 |
| ...                  | ...           |      ...      |         ... |         ... |

### Use memo for data attributes 
By default, data attributes is implemented with Cadence/Temporal [query API](https://docs.temporal.io/workflows#query), 
which is not optimized for very high volume reads on a single workflow execution(like 100 rps), because it could cause
too many replay with history, especially when workflows are closed.

However, you can enable the feature "useMemoForDataAttributes". This is currently only supported if the backend is Temporal, 
because [Cadence doesn't support mutable memo](https://github.com/uber/cadence/issues/3729).  

## Workflow State
A workflow state is like “a small workflow” of 1~2 steps:

**[ waitUntil ] → execute**

The `waitUntil` API can returns some commands to wait for. When the commands are completed, the `execute` API will be invoked.
The two APIs have access to read/write the persistence defined in the workflow.

The full detailed execution flow is like this:

![Workflow State diagram](https://user-images.githubusercontent.com/4523955/234921554-587d8ad4-84f5-4987-b838-959869293465.png)

The `waitUntil` API is optional. If not defined, then `execute` API will be invoked instead when the state started.

Note: the two APIs are invoked by iWF service with infinite backoff retry by default. See WorkflowStateOptions section for customization.  

The execute API will return a StateDecision:
* Single next state 
  * Go to to different state
  * Go to the same state as a loop
  * Go the the previous state as a loop
* Multiple next states, executing as multi threads in parallel
* Dead end -- Just stop the thread
* Graceful complete -- Stop the thread, and also will stop the workflow when all other threads are stopped
* Force complete -- Stop the workflow immediately
* Force fail  -- Stop the workflow immediately with failure

With decisions, the workflow definitions can have flows like these:

![decision flow1](https://user-images.githubusercontent.com/4523955/234919901-f327dfb6-5b38-4440-a2eb-5d1c832b694e.png)

or 

![decision flow2](https://user-images.githubusercontent.com/4523955/234919896-30db8628-daeb-4f1d-bd2b-7bf826989c75.png)

or even more complicated as needed.


### Commands for WorkflowState's WaitUntil API

iWF provides three types of commands:

* `SignalCommand`: will wait for a signal to be published to the workflow signal channel. External application can use
  SignalWorkflow API to signal a workflow.
* `TimerCommand`: will wait for a **durable timer** to fire.
* `InternalChannelCommand`: will wait for a message from InternalChannel.

The waitUntil API can return multiple commands along with a `CommandWaitingType`:

* `AllCommandCompleted`: This option waits for all commands to be completed.

* `AnyCommandCompleted`: This option waits for any of the commands to be completed.

* `AnyCommandCombinationCompleted`: This option waits for any combination of the commands in a specified list to be
  completed.

## RPC

RPC stands for "Remote Procedure Call". It's invoked by client, executed in workflow worker, and then respond back the results to client. 

RPC provides a simple and powerful mechansim to interact with external systems. With RPCs defined along with persistence, an ObjectWorkflow 
works like an durable object that provide methods to execute business logic. You can even uses iWF to implement a typical CRUD application like a 
blog post, **the pseudo code** looks like this:

```java
class BlogPost implements ObjectWorkflow{
    DataAttribute String title;
    DataAttribute String authorName;
    DataAttribute String body;

    @RPC 
    void updateTitle(String title){
         this.title = title
    } 
    
    @RPC 
    void updateBody(String body){
         this.body = body
    }     
    
    ...
    ... 
    
    @RPC 
    BlobPost get(){
         return new BlogPost(title, authorName, body);
    } 
}

```
This is just pseudo code. The real code will look similar depends on which SDK to use. 

See [this example of implementing an CRUD application](https://github.com/indeedeng/iwf-java-samples/tree/main/src/main/java/io/iworkflow/workflow/jobpost), with more capabilities like searching and background execution, just ~100 lines of code. 

### Atomicity of RPC APIs

It's important to note that in addition to read/write persistence fields, a RPC can **trigger new state executions, and publish message to InternalChannel, all atomically.**

Atomically sending internal channel, or triggering state executions is an important pattern to ensure consistency across dependencies for critical business – this 
solves a very common problem in many existing distributed system applications. Because most RPCs (like REST/gRPC/GraphQL) don't provide a way to invoke 
background execution when updating persistence. People sometimes have to use complicated design to acheive this. 

**But in iWF, it's all builtin, and user application just needs a few lines of code!** 

![flow with RPC](https://user-images.githubusercontent.com/4523955/234930263-40b98ca7-4401-44fa-af8a-32d5ae075438.png)

### Signal Channel vs RPC

There are two major ways for external clients to interact with workflows: Signal and RPC. So what are the difference? 

They are completely different:
* Signal is sent to iWF service without waiting for response of the processing
* RPC will wait for worker to process the RPC request synchronously
* Signal will be held in a signal channel until a workflow state consumes it
* RPC will be processed by worker immediately

![signals vs rpc](https://user-images.githubusercontent.com/4523955/234932674-b0d062b2-e5dd-4dbe-93b5-1b9863acc5e0.png)

So choose based on the situations/requirements

|                |        Availability        |                                        Latency |                                    Workflow Requirement |
|----------------|:-------------------------- |:----------------------------------------------- |:-------------------------------------------------------- |
| Signal Channel |            High            |                                            Low |                     Requires a WorkflowState to process |
| RPC            | Depends on workflow worker | Higher than signal, depends on workflow worker |                               No WorkflowState required |

## Advanced Customization

### WorkflowOptions

iWF let you deeply customize the workflow behaviors with the below options.

#### IdReusePolicy for WorkflowId

At any given time, there can be only one WorkflowExecution running for a specific workflowId.
A new WorkflowExecution can be initiated using the same workflowId by setting the appropriate `IdReusePolicy` in
WorkflowOptions.

* `ALLOW_IF_NO_RUNNING` 
    * Allow starting workflow if there is no execution running with the workflowId
    * This is the **default policy** if not specified in WorkflowOptions
* `ALLOW_IF_PREVIOUS_EXISTS_ABNORMALLY`
    * Allow starting workflow if a previous Workflow Execution with the same Workflow Id does not have a Completed
      status.
      Use this policy when there is a need to re-execute a Failed, Timed Out, Terminated or Cancelled workflow
      execution.
* `DISALLOW_REUSE` 
    * Not allow to start a new workflow execution with the same workflowId.
* `ALLOW_TERMINATE_IF_RUNNING`
    * Always allow starting workflow no matter what -- iWF server will terminate the current running one if it exists.

#### CRON Schedule

iWF allows you to start a workflow with a fixed cron schedule like below

```text
// CronSchedule - Optional cron schedule for workflow. If a cron schedule is specified, the workflow will run
// as a cron based on the schedule. The scheduling will be based on UTC time. The schedule for the next run only happens
// after the current run is completed/failed/timeout. If a RetryPolicy is also supplied, and the workflow failed
// or timed out, the workflow will be retried based on the retry policy. While the workflow is retrying, it won't
// schedule its next run. If the next schedule is due while the workflow is running (or retrying), then it will skip
that
// schedule. Cron workflow will not stop until it is terminated or cancelled (by returning cadence.CanceledError).
// The cron spec is as follows:
// ┌───────────── minute (0 - 59)
// │ ┌───────────── hour (0 - 23)
// │ │ ┌───────────── day of the month (1 - 31)
// │ │ │ ┌───────────── month (1 - 12)
// │ │ │ │ ┌───────────── day of the week (0 - 6) (Sunday to Saturday)
// │ │ │ │ │
// │ │ │ │ │
// * * * * *
```

NOTE:

* iWF also
  supports [more advanced cron expressions](https://pkg.go.dev/github.com/robfig/cron#hdr-CRON_Expression_Format)
* The [crontab guru](https://crontab.guru/) site is useful for testing your cron expressions.
* To cancel a cron schedule, use terminate of cancel type to stop the workflow execution.
* By default, there is no cron schedule.

#### RetryPolicy for workflow

Workflow execution can have a backoff retry policy which will retry on failed or timeout.

By default, there is no retry policy.

#### Initial Search Attributes

Client can specify some initial search attributes when starting the workflow.

By default, there is no initial search attributes.

### WorkflowStateOptions

Similarly, users can customize the WorkflowState

#### WorkflowState WaitUntil/Execute API timeout and retry policy

By default, the API timeout is 30s with infinite backoff retry. 
Users can customize the API timeout and retry policy:

- InitialIntervalSeconds: 1
- MaxInternalSeconds:100
- MaximumAttempts: 0
- MaximumAttemptsDurationSeconds: 0
- BackoffCoefficient: 2

Where zero means infinite attempts.

Both MaximumAttempts and MaximumAttemptsDurationSeconds are used for controlling the maximum attempts for the retry
policy.
MaximumAttempts is directly by number of attempts, where MaximumAttemptsDurationSeconds is by the total time duration of
all attempts including retries. It will be capped to the minimum if both are provided.

#### Persistence loading policy

When a workflowState/RPC API loads DataAttributes/SearchAttributes, by default it will use `LOAD_ALL_WITOUT_LOCKING` to load everything.

For WorkflowState, there is a 2MB limit by default to load data. User can use another loading policy `LOAD_PARTIAL_WITHOUT_LOCKING`
to specify certain DataAttributes/SearchAttributes only to load.

`WITHOUT_LOCKING` here means if multiple StateExecutions/RPC try to upsert the same DataAttribute/SearchAttribute, they can be
done in parallel without locking.

If racing conditions could be a problem, using`PARTIAL_WITH_EXCLUSIVE_LOCK` allows specifying some keys to be locked during the execution.

#### WaitUntil API failure policy

By default, the workflow execution will fail when API max out the retry attempts. In some cases that
workflow want to ignore the errors.

Using `PROCEED_ON_API_FAILURE` for `WaitUntilApiFailurePolicy` will let workflow continue to execute decide
API when the API fails with maxing out all the retry attempts.

Alternatively, WorkflowState can utilize `attempts` or `firstAttemptTime` from the context to decide ignore the
exception/error.

## Limitation

Though iWF can be used for a very wide range of use case even just CRUD, iWF is NOT for everything. It is not suitable for use cases like:

* High performance transaction( e.g. within 10ms)
* High frequent writes on a single workflow execution(like a single record in database) for hot partition issue
  * High frequent reads on a single workflow execution is okay if using memo for data attributes
* Join operation across different workflows
* Transaction for operation across multiple workflows


# Architecture

An iWF application is composed of several iWF workflow workers. These workers host REST APIs as "worker APIs" for server to call. This callback pattern similar to AWS Step Functions invoking Lambdas, if you are familiar with.

An application also perform actions on workflow executions, such as starting, stopping, signaling, and retrieving results 
by calling iWF service APIs as "service APIs".

The service APIs are provided by the "API service" in iWF server. Internally, this API service communicates with the Cadence/Temporal service as its backend.

In addition, the iWF server also runs the Cadence/Temporal workers as "worker service". The worker service
hosts [an interpreter workflow](https://github.com/indeedeng/iwf/blob/main/service/interpreter/workflowImpl.go).
This workflow implements all the core features as described above, and also things like "Auto ContinueAsNew" to let you use 
iWF without any scaling limitation. 

![architecture diagram](https://user-images.githubusercontent.com/4523955/234935630-e69c648e-7714-4672-beb2-d9867bedf940.png)

# How to use

## Using docker image & docker-compose

Checkout this repo, go to the docker-compose folder and run it:

```shell
cd docker-compose && docker-compose up
```

This by default will run Temporal server with it.
And it will also register a `default` Temporal namespace and required search attributes by iWF.
Link to the Temporal WebUI: http://localhost:8233/namespaces/default/workflows

By default, iWF server is serving port **8801**, server URL is http://localhost:8801/ )

NOTE:

Use `docker pull iworkflowio/iwf-server:latest` to update the latest image.Or update the docker-compose file to specify
the version tag.

## How to build & run locally

* Run `make bins` to build the binary `iwf-server`
* Make sure you have registered the system search attributes required by iWF server:
    * Keyword: IwfWorkflowType
    * Int: IwfGlobalWorkflowVersion
    * Keyword: IwfExecutingStateIds
    * See [Contribution](./CONTRIBUTING.md) for more detailed commands.
    * For Cadence without advancedVisibility enabled,
      set [disableSystemSearchAttributes](https://github.com/indeedeng/iwf/blob/main/config/development_cadence.yaml#L8)
      to true
* Then run  `./iwf-server start` to run the service . This defaults to serve workflows APIs with Temporal interpreter
  implementation. It requires to have local Temporal setup. See Run with local Temporal.
* Alternatively, run `./iwf-server --config config/development_cadence.yaml start` to run with local Cadence. See below
  instructions for setting up local Cadence.


## Troubleshooting

When something goes wrong in your applications, here are the tips:

* All the input/output to your workflow are stored in the activity input/output of history event. The input is
  in `ActivityTaskScheduledEvent`, output is in `ActivityTaskCompletedEvent` or in pending activity view if having
  errors.
* Use query handlers like (`GetDataObjects` or `GetCurrentTimerInfos`) in Cadence/Temporal WebUI to quickly understand
  the current status of the workflows.
    * DumpAllInternal will return all the internal status or the pending states
    * GetCurrentTimerInfos will return all the timers of the pending states
* Let your worker service return error stacktrace as the response body to iWF server. E.g.
  like [this example of Spring Boot using ExceptionHandler](https://github.com/indeedeng/iwf-java-samples/blob/2d500093e2aaecf2d728f78366fee776a73efd29/src/main/java/io/iworkflow/controller/IwfWorkerApiController.java#L51)
  .
* If you return the full stacktrace in response body, the pending activity view will show it to you! Then use
  Cadence/Temporal WebUI to debug your application.


## Operation

In additional of using Cadence/Temporal CLI, you can just
use [some HTTP script like this](./script/http/local/home.http) to operate on workflows to:

* Start a workflow
* Stop a workflow
* Reset a workflow
* Skip a timer
* etc, any APIs supported by the [iWF server API schema](https://github.com/indeedeng/iwf-idl/blob/main/iwf.yaml)

# Posts & Articles & Reference

* Temporal adopted
  as [the first community drive DSL framework/abstraction](https://github.com/temporalio/awesome-temporal) of Temporal
* Cadence adopted in its [README](https://github.com/uber/cadence#cadence)
  , [official documentation](https://cadenceworkflow.io/docs/get-started/#what-s-next)
  and [Cadence community spotlight](https://cadenceworkflow.io/blog/2023/01/31/community-spotlight-january-2023/)
