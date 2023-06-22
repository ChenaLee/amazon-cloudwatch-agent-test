// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

//go:build !windows

package assume_role

const (
	credsDir = "/tmp/.aws"
)

func getCommands(roleArn string) []string {
	return []string{
		"mkdir -p " + credsDir,
		"printf '[default]\naws_access_key_id=%s\naws_secret_access_key=%s\naws_session_token=%s' $(aws sts assume-role --role-arn " + roleArn + " --role-session-name test --query 'Credentials.[AccessKeyId,SecretAccessKey,SessionToken]' --output text) | tee " + credsDir + "/credentials>/dev/null",
		"printf '[default]\nregion = us-west-2' > " + credsDir + "/config",
		"printf '[credentials]\n  shared_credential_profile = \"default\"\n  shared_credential_file = \"" + credsDir + "/credentials\"' | sudo tee /opt/aws/amazon-cloudwatch-agent/etc/common-config.toml>/dev/null",
	}
}
