// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

//go:build !windows

package metric_value_benchmark

import (
	"github.com/aws/amazon-cloudwatch-agent-test/internal/common"
	"log"
	"os"

	"github.com/aws/amazon-cloudwatch-agent-test/test/metric"
	"github.com/aws/amazon-cloudwatch-agent-test/test/metric/dimension"
	"github.com/aws/amazon-cloudwatch-agent-test/test/status"
	"github.com/aws/amazon-cloudwatch-agent-test/test/test_runner"
	"github.com/aws/aws-sdk-go-v2/aws"
)

type SwapTestRunner struct {
	test_runner.BaseTestRunner
}

var _ test_runner.ITestRunner = (*SwapTestRunner)(nil)

func (t *SwapTestRunner) Validate() status.TestGroupResult {
	metricsToFetch := t.GetMeasuredMetrics()
	testResults := make([]status.TestResult, len(metricsToFetch))
	for i, metricName := range metricsToFetch {
		testResults[i] = t.validateSwapMetric(metricName)
	}

	return status.TestGroupResult{
		Name:        t.GetTestName(),
		TestResults: testResults,
	}
}

func (t *SwapTestRunner) GetTestName() string {
	return "Swap"
}

func (t *SwapTestRunner) GetAgentConfigFileName() string {
	return "swap_config.json"
}

func (t *SwapTestRunner) GetMeasuredMetrics() []string {
	return []string{
		"swap_free",
		"swap_used",
		"swap_used_percent",
	}
}

func (t *SwapTestRunner) validateSwapMetric(metricName string) status.TestResult {
	testResult := status.TestResult{
		Name:   metricName,
		Status: status.FAILED,
	}

	hostName, err := os.Hostname()
	if err != nil {
		log.Printf("Hostname was not found")

		t.Fatalf("Can't get hostname")
	}
	log.Printf("Hostname found %s", hostName)

	dims, failed := t.DimensionFactory.GetDimensions([]dimension.Instruction{
		{
			Key:   aws.String(common.Host),
			Value: aws.String(hostName),
		},
	})

	if len(failed) > 0 {
		return testResult
	}

	fetcher := metric.MetricValueFetcher{}
	values, err := fetcher.Fetch(namespace, metricName, dims, metric.AVERAGE, metric.HighResolutionStatPeriod)
	log.Printf("metric values are %v", values)
	if err != nil {
		return testResult
	}

	if !metric.IsAllValuesGreaterThanOrEqualToExpectedValue(metricName, values, 0) {
		return testResult
	}

	testResult.Status = status.SUCCESSFUL
	return testResult
}
