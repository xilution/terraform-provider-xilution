package provider

import (
	"errors"
	"fmt"
	"strings"
	"time"

	xc "github.com/xilution/xilution-client-go"
)

const CREATE_COMPLETE = "CREATE_COMPLETE"
const UPDATE_COMPLETE = "UPDATE_COMPLETE"
const UPDATE_ROLLBACK_COMPLETE = "UPDATE_ROLLBACK_COMPLETE"
const SUCCEEDED = "SUCCEEDED"
const FAILED = "FAILED"
const NOT_FOUND = "NOT_FOUND"

func getIdFromLocationUrl(location *string) *string {
	index := strings.LastIndex(*location, "/")
	id := string((*location)[(index + 1):])

	return &id
}

func waitForPipelineEventToComplete(
	eventType string,
	timeout time.Duration,
	waitIncrement time.Duration,
	getPipelineStatusFunc func() (*xc.PipelineStatus, error),
) error {
	if eventType == "PROVISION" || eventType == "RUN_NOW" {
		err := waitForPipelineUpToSucceeded(10*time.Minute, 5*time.Second, getPipelineStatusFunc)
		if err != nil {
			return err
		}
	} else if eventType == "REPROVISION" {
		err := waitForPipelineInfrastructureUpdateComplete(10*time.Minute, 5*time.Second, getPipelineStatusFunc)
		if err != nil {
			return err
		}
	} else if eventType == "DEPROVISION" {
		err := waitForPipelineInfrastructureNotFound(10*time.Minute, 5*time.Second, getPipelineStatusFunc)
		if err != nil {
			return err
		}
	}

	return nil
}

func waitForPipelineUpToSucceeded(
	timeout time.Duration,
	waitIncrement time.Duration,
	getPipelineStatusFunc func() (*xc.PipelineStatus, error),
) error {
	if waitIncrement < 5*time.Second {
		return errors.New("wait increment must be greater than 5 seconds")
	}
	
	done := false
	start := time.Now()
	for !done {
		status, err := getPipelineStatusFunc()
		if err != nil {
			return err
		}
		if status != nil {
			infrastructureStatus := status.InfrastructureStatus
			if infrastructureStatus == CREATE_COMPLETE {
				continuousIntegrationStatus := status.ContinuousIntegrationStatus
				if continuousIntegrationStatus != nil {
					latestUpExecutionStatus := continuousIntegrationStatus.LatestUpExecutionStatus
					if latestUpExecutionStatus == SUCCEEDED {
						done = true
					} else if strings.HasSuffix(latestUpExecutionStatus, FAILED) {
						return fmt.Errorf("pipeline up status is %s", latestUpExecutionStatus)
					}
				}
			} else if strings.HasSuffix(infrastructureStatus, FAILED) {
				return fmt.Errorf("pipeline infrastructure status is %s", infrastructureStatus)
			}
		} else {
			if time.Since(start) > timeout {
				return errors.New("timeout waiting for pipeline to succeed")
			}
			time.Sleep(waitIncrement)
		}
	}

	return nil
}

func waitForPipelineInfrastructureUpdateComplete(
	timeout time.Duration,
	waitIncrement time.Duration,
	getPipelineStatusFunc func() (*xc.PipelineStatus, error),
) error {
	if waitIncrement < 5*time.Second {
		return errors.New("wait increment must be greater than 5 seconds")
	}

	done := false
	start := time.Now()
	for !done {
		status, err := getPipelineStatusFunc()
		if err != nil {
			return err
		}
		if status != nil {
			infrastructureStatus := status.InfrastructureStatus
			if infrastructureStatus == UPDATE_COMPLETE {
				done = true
			} else if infrastructureStatus == UPDATE_ROLLBACK_COMPLETE ||
				strings.HasSuffix(infrastructureStatus, FAILED) {
				return fmt.Errorf("pipeline infrastructure status is %s", infrastructureStatus)
			}
		} else {
			if time.Since(start) > timeout {
				return errors.New("timeout waiting for pipeline to succeed")
			}
			time.Sleep(waitIncrement)
		}
	}

	return nil
}

func waitForPipelineInfrastructureNotFound(
	timeout time.Duration,
	waitIncrement time.Duration,
	getPipelineStatusFunc func() (*xc.PipelineStatus, error),
) error {
	if waitIncrement < 5*time.Second {
		return errors.New("wait increment must be greater than 5 seconds")
	}

	done := false
	start := time.Now()
	notFoundCount := 0
	for !done {
		status, err := getPipelineStatusFunc()
		if err != nil {
			return err
		}
		if status != nil {
			infrastructureStatus := status.InfrastructureStatus
			if infrastructureStatus == NOT_FOUND {
				if notFoundCount == 3 {
					done = true
				}
				notFoundCount = notFoundCount + 1
			} else if strings.HasSuffix(infrastructureStatus, FAILED) {
				return fmt.Errorf("pipeline infrastructure status is %s", infrastructureStatus)
			}
		} else {
			if time.Since(start) > timeout {
				return errors.New("timeout waiting for pipeline to succeed")
			}
			time.Sleep(waitIncrement)
		}
	}

	return nil
}
