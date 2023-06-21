// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

//go:build !windows

package metric_value_benchmark

import (
	"github.com/aws/amazon-cloudwatch-agent-test/internal/common"
	"github.com/aws/amazon-cloudwatch-agent-test/test/metric"
	"github.com/aws/amazon-cloudwatch-agent-test/test/metric/dimension"
	"github.com/aws/amazon-cloudwatch-agent-test/test/status"
	"github.com/aws/amazon-cloudwatch-agent-test/test/test_runner"
	"github.com/aws/aws-sdk-go-v2/aws"
	"log"
	"os"
)

type ProcessesTestRunner struct {
	test_runner.BaseTestRunner
}

var _ test_runner.ITestRunner = (*ProcessesTestRunner)(nil)

func (m *ProcessesTestRunner) Validate() status.TestGroupResult {
	metricsToFetch := m.GetMeasuredMetrics()
	testResults := make([]status.TestResult, len(metricsToFetch))
	for i, name := range metricsToFetch {
		testResults[i] = m.validateProcessesMetric(name)
	}

	return status.TestGroupResult{
		Name:        m.GetTestName(),
		TestResults: testResults,
	}
}

func (m *ProcessesTestRunner) GetTestName() string {
	return "Processes"
}

func (m *ProcessesTestRunner) GetAgentConfigFileName() string {
	return "processes_config.json"
}

func (m *ProcessesTestRunner) GetMeasuredMetrics() []string {
	return []string{
		"processes_blocked", "processes_dead", "processes_idle", "processes_paging", "processes_running", "processes_sleeping", "processes_stopped",
		"processes_total", "processes_total_threads", "processes_zombies"}
}

func (m *ProcessesTestRunner) validateProcessesMetric(metricName string) status.TestResult {
	testResult := status.TestResult{
		Name:   metricName,
		Status: status.FAILED,
	}

	hostName, err := os.Hostname()
	if err != nil {
		log.Printf("Hostname was not found")

		m.Fatalf("Can't get hostname")
	}
	log.Printf("Hostname found %s", hostName)

	dims, failed := m.DimensionFactory.GetDimensions([]dimension.Instruction{
		{
			Key:   common.Host,
			Value: dimension.ExpectedDimensionValue{Value: aws.String(hostName)},
		},
	})
	if len(failed) > 0 {
		return testResult
	}

	fetcher := metric.MetricValueFetcher{}
	values, err := fetcher.Fetch(namespace, metricName, dims, metric.AVERAGE, metric.HighResolutionStatPeriod)
	if err != nil {
		return testResult
	}

	if !metric.IsAllValuesGreaterThanOrEqualToExpectedValue(metricName, values, 0) {
		return testResult
	}

	testResult.Status = status.SUCCESSFUL
	return testResult
}
