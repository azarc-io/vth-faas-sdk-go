package sparkv1

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

/************************************************************************/
// STAGE OPTIONS
/************************************************************************/

type stageOptionParams struct {
	stageName string
	ctx       Context
}

func (s stageOptionParams) StageName() string {
	return s.stageName
}

func (s stageOptionParams) Context() Context {
	return s.ctx
}

func newStageOptionParams(ctx Context, stageName string) StageOptionParams {
	return stageOptionParams{
		stageName: stageName,
		ctx:       ctx,
	}
}

/************************************************************************/
// SPARK OPTIONS
/************************************************************************/

type SparkOpts struct {
	log            Logger
	config         []byte
	configType     ConfigType
	configBasePath string
}

type Option = func(je *SparkOpts) *SparkOpts

func WithSparkConfig(cfg any) Option {
	d, err := json.Marshal(cfg)
	if err != nil {
		log.Fatal().Err(err).Msgf("unable to serialise config")
	}
	return func(je *SparkOpts) *SparkOpts {
		je.config = d
		return je
	}
}

/************************************************************************/
// WORKFLOW OPTIONS
/************************************************************************/
type workflowOpts struct {
	stageTracker InternalStageTracker
	config       *Config
	nc           *nats.Conn
}

type WorkflowOption = func(je *workflowOpts) *workflowOpts

func WithStageTracker(ist InternalStageTracker) WorkflowOption {
	return func(jw *workflowOpts) *workflowOpts {
		jw.stageTracker = ist
		return jw
	}
}

func WithConfig(cfg *Config) WorkflowOption {
	return func(jw *workflowOpts) *workflowOpts {
		jw.config = cfg
		return jw
	}
}

func WithNatsClient(nc *nats.Conn) WorkflowOption {
	return func(jw *workflowOpts) *workflowOpts {
		jw.nc = nc
		return jw
	}
}
