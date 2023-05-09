package spark

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	client "github.com/azarc-io/vth-faas-sdk-go/cmd/spark-web-wrapper/internal/gen"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
	"net/http"
	"time"
)

type Spark struct {
	baseUrl string
	client  client.ClientWithResponsesInterface
	spec    *client.Spec
}

func (s *Spark) Init(ctx sparkv1.InitContext) error {
	var err error
	cfg, err := ctx.Config().Raw()
	if err != nil {
		return err
	}

	resp, err := s.client.PostInitWithBodyWithResponse(context.Background(), string(client.Applicationjson), bytes.NewReader(cfg))
	if err != nil {
		return err
	}

	if resp.JSON500 != nil {
		return fmt.Errorf("%s: %s", resp.JSON500.Message, resp.JSON500.Code)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("error initialising proxy: %s", string(resp.Body))
	}

	return nil
}

func (s *Spark) Stop() {
}

func (s *Spark) initClientAndSpec() error {
	var err error
	s.client, err = client.NewClientWithResponses(s.baseUrl)
	if err != nil {
		return err
	}

	resp, err := s.client.GetSpecWithResponse(context.Background())
	if err != nil {
		return err
	}

	if resp.JSON500 != nil {
		return fmt.Errorf("%s: %s", resp.JSON500.Message, resp.JSON500.Code)
	}

	s.spec = resp.JSON200
	if s.spec == nil || len(s.spec.Stages) == 0 {
		return fmt.Errorf("unable to find any stages in response specification")
	}

	return nil
}

func (s *Spark) BuildChain(b sparkv1.Builder) sparkv1.Chain {
	err := s.initClientAndSpec()
	if err != nil {
		// TODO: can only panic here. Should look at allowing returning an error in the SDK
		panic(err)
	}

	chain := b.NewChain("main")
	stg := chain.Stage(s.spec.Stages[0].Name, s.getStage(&s.spec.Stages[0]))

	for _, stage := range s.spec.Stages[1:] {
		stg = stg.Stage(stage.Name, s.getStage(&stage))
	}

	return stg.Complete(s.getCompleteStage(&s.spec.Complete))
}

func (s *Spark) getStage(stage *client.Stage) sparkv1.StageDefinitionFn {
	return func(ctx sparkv1.StageContext) (any, sparkv1.StageError) {
		// get all inputs
		inputs, err := s.getStageInputs(ctx, stage.Inputs)
		if err != nil {
			return nil, sparkv1.NewStageError(err)
		}

		// get all previous stages
		psr, err := s.getPreviousStageResults(ctx, stage.Name)
		if err != nil {
			return nil, sparkv1.NewStageError(err)
		}

		req := client.StageRequest{
			JobKey:         ctx.JobKey(),
			CorrelationID:  ctx.CorrelationID(),
			TransactionID:  ctx.TransactionID(),
			PreviousStages: psr,
			Inputs:         &inputs,
		}
		reqData, err := json.MarshalIndent(req, "", "   ")
		if err != nil {
			return nil, sparkv1.NewStageError(err)
		}

		// make request to wrapped services stage
		resp, err := s.client.PostStagesNameWithBodyWithResponse(context.Background(), stage.Name, string(client.Applicationjson), bytes.NewReader(reqData))
		if err != nil {
			return nil, sparkv1.NewStageError(err)
		}

		if resp.JSON500 != nil {
			e := resp.JSON500
			opts := []sparkv1.ErrorOption{
				sparkv1.WithMetadata(e.Metadata),
				sparkv1.WithErrorCode(sparkv1.ErrorCode(e.Code)),
			}
			if e.Retry != nil {
				opts = append(opts, sparkv1.WithRetry(uint(e.Retry.Times), uint(e.Retry.BackoffMultiplier), time.Duration(e.Retry.FirstBackoffWait)))
			}
			return nil, sparkv1.NewStageError(errors.New(e.Message), opts...)
		}

		if resp.JSON200 != nil {
			res := resp.JSON200
			s.renderLogs(ctx.Log(), res.Logs)
			return res.Value, nil
		}

		return nil, sparkv1.NewStageError(fmt.Errorf("unable to process complete stage response: http code %d: %s", resp.HTTPResponse.StatusCode, string(resp.Body)))
	}
}

func (s *Spark) getCompleteStage(cs *client.CompleteStage) sparkv1.CompleteDefinitionFn {
	return func(ctx sparkv1.CompleteContext) sparkv1.StageError {
		// get all inputs
		inputs, err := s.getStageInputs(ctx, cs.Inputs)
		if err != nil {
			return sparkv1.NewStageError(err)
		}

		// get all previous stages
		psr, err := s.getPreviousStageResults(ctx, cs.Name)
		if err != nil {
			return sparkv1.NewStageError(err)
		}

		req := client.StageRequest{
			JobKey:         ctx.JobKey(),
			CorrelationID:  ctx.CorrelationID(),
			TransactionID:  ctx.TransactionID(),
			PreviousStages: psr,
			Inputs:         &inputs,
		}
		reqData, err := json.MarshalIndent(req, "", "   ")
		if err != nil {
			return sparkv1.NewStageError(err)
		}

		// make request to wrapped services stage
		resp, err := s.client.PostCompleteNameWithBodyWithResponse(context.Background(), cs.Name, string(client.Applicationjson), bytes.NewReader(reqData))
		if err != nil {
			return sparkv1.NewStageError(err)
		}

		if resp.JSON500 != nil {
			e := resp.JSON500
			opts := []sparkv1.ErrorOption{
				sparkv1.WithMetadata(e.Metadata),
				sparkv1.WithErrorCode(sparkv1.ErrorCode(e.Code)),
			}
			if e.Retry != nil {
				opts = append(opts, sparkv1.WithRetry(uint(e.Retry.Times), uint(e.Retry.BackoffMultiplier), time.Duration(e.Retry.FirstBackoffWait)))
			}
			return sparkv1.NewStageError(errors.New(e.Message), opts...)
		}

		if resp.JSON200 != nil {
			res := resp.JSON200
			s.renderLogs(ctx.Log(), res.Logs)

			for _, o := range res.Outputs {
				if o.Value != nil {
					val, err := json.Marshal(o.Value)
					if err != nil {
						return sparkv1.NewStageError(err)
					}

					if err := ctx.Output(sparkv1.NewVar(o.Name, codec.MimeType(o.Mimetype), val)); err != nil {
						return sparkv1.NewStageError(err)
					}
				}
			}

			return nil
		}

		return sparkv1.NewStageError(fmt.Errorf("unable to process complete stage response: http code %d: %s", resp.HTTPResponse.StatusCode, string(resp.Body)))
	}
}

func (s *Spark) getStageInputs(ctx sparkv1.StageContext, ins []string) ([]client.Input, error) {
	// get all inputs
	inputs := make([]client.Input, 0)
	for _, in := range ins {
		var val any
		err := ctx.Input(in).Bind(&val)
		if err != nil {
			return nil, fmt.Errorf("unable to get input '%s' value: %w", in, err)
		}

		inputs = append(inputs, client.Input{
			Value: val,
			Name:  in,
		})
	}

	return inputs, nil
}

func (s *Spark) getPreviousStageResults(ctx sparkv1.StageContext, currentStageName string) ([]client.PreviousStageResult, error) {
	var psr []client.PreviousStageResult
	for _, ps := range s.spec.Stages {
		if ps.Name == currentStageName {
			// this is the current stage so exit
			break
		}

		var val any
		err := ctx.StageResult(ps.Name).Bind(&val)
		if err != nil {
			return nil, fmt.Errorf("unable to get stage '%s' value: %w", ps.Name, err)
		}

		psr = append(psr, client.PreviousStageResult{
			Value: val,
			Name:  ps.Name,
		})
	}

	return psr, nil
}

func (s *Spark) renderLogs(logger sparkv1.Logger, logs *[]client.Log) {
	if logs != nil {
		for _, log := range *logs {
			switch log.Level {
			case client.LogLevelInfo:
				logger.Info(log.Message)
			case client.LogLevelDebug:
				logger.Debug(log.Message)
			case client.LogLevelError:
				logger.Error(errors.New(log.Message), "")
			}
		}
	}
}

// NewSpark creates a Spark
// shutdown is a reference to the context cancel function, it can be used to gracefully stop the worker if needed
func NewSpark(baseUrl string) sparkv1.Spark {
	return &Spark{baseUrl: baseUrl}
}
