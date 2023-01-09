package sparkv1

/************************************************************************/
// JOB CONTEXT
/************************************************************************/

type jobContext struct {
	metadata *sparkMetadata
	log      Logger
}

func (j *jobContext) JobKey() string {
	return j.metadata.jobKey
}

func (j *jobContext) CorrelationID() string {
	return j.metadata.correlationID
}

func (j *jobContext) TransactionID() string {
	return j.metadata.transactionID
}

func (j *jobContext) Log() Logger {
	return j.log
}

func NewJobContext(metadata Context, opts *SparkOpts) Context {
	m := sparkMetadata{
		jobKey:        metadata.JobKey(),
		correlationID: metadata.CorrelationID(),
		transactionID: metadata.TransactionID(),
	}
	return &jobContext{
		metadata: &m,
		log:      opts.log,
	}
}

/************************************************************************/
// METADATA
/************************************************************************/

type sparkMetadata struct {
	jobKey        string
	correlationID string
	transactionID string
	logger        Logger
}

func (j *sparkMetadata) Log() Logger {
	return j.logger
}

func (j *sparkMetadata) JobKey() string {
	return j.jobKey
}

func (j *sparkMetadata) CorrelationID() string {
	return j.correlationID
}

func (j *sparkMetadata) TransactionID() string {
	return j.transactionID
}

func NewSparkMetadata(jobKey, correlationID, transactionID string, logger Logger) Context {
	return &sparkMetadata{
		logger:        logger,
		jobKey:        jobKey,
		correlationID: correlationID,
		transactionID: transactionID,
	}
}

/************************************************************************/
// STAGE CONTEXT
/************************************************************************/

type stageContext struct {
	*ExecuteStageRequest
	workflowId  string
	runId       string
	logger      Logger
	name        string
	inputs      ExecuteSparkInputs
	sparkDataIO SparkDataIO
}

func NewCompleteContext(req *ExecuteStageRequest, sparkDataIO SparkDataIO, workflowId, runId, name string, logger Logger, inputs ExecuteSparkInputs) CompleteContext {
	return &completeContext{stageContext: stageContext{ExecuteStageRequest: req, name: name, logger: logger, inputs: inputs, sparkDataIO: sparkDataIO, workflowId: workflowId, runId: runId}}
}

func NewStageContext(req *ExecuteStageRequest, sparkDataIO SparkDataIO, workflowId, runId, name string, logger Logger, inputs ExecuteSparkInputs) StageContext {
	return stageContext{ExecuteStageRequest: req, sparkDataIO: sparkDataIO, name: name, logger: logger, inputs: inputs, workflowId: workflowId, runId: runId}
}

func (sc stageContext) Input(name string) Input {
	in, ok := sc.inputs[name]
	if !ok {
		return &bindable{}
	}

	return in
}

func (sc stageContext) StageResult(name string) Bindable {
	result, err := sc.sparkDataIO.GetStageResult(sc.workflowId, sc.runId, name)
	if err != nil {
		return NewBindableError(err)
	}
	return result
}

func (sc stageContext) Name() string {
	return sc.name
}

func (sc stageContext) Log() Logger {
	return sc.logger
}

type completeContext struct {
	stageContext
	outputs []*Var
}

func (cc *completeContext) Output(variables ...*Var) error {
	cc.outputs = append(cc.outputs, variables...)
	return nil
}

/************************************************************************/
// INIT CONTEXT
/************************************************************************/

type initContext struct {
	loader BindableConfig
	opts   *SparkOpts
}

func (i *initContext) Config() BindableConfig {
	if i.loader == nil {
		i.loader = newBindableConfig(i.opts)
	}
	return i.loader
}

func NewInitContext(opts *SparkOpts) InitContext {
	return &initContext{opts: opts}
}
