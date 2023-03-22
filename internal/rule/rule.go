package rule

import (
	"github.com/aws/amazon-cloudwatch-agent-test/filesystem"
	"log"
)

type Rule[T any] struct {
	Conditions []ICondition[T]
}

func (r *Rule[T]) Evaluate(target T) (bool, error) {
	for _, c := range r.Conditions {
		log.Printf("Evaluating condition: %v", c.Name())
		success, err := c.Evaluate(target)
		if err != nil {
			log.Printf("Evaluation did not run. Error was %v", err)
			return false, err
		}
		log.Printf("Evaluate result: %v", success)
		if !success {
			return false, nil
		}
	}
	return true, nil
}

type ICondition[T any] interface {
	Name() string
	Evaluate(T) (bool, error)
}

type FilePermissionExpected struct {
	PermissionCompared filesystem.FilePermission
	ShouldExist        bool
}

func (e *FilePermissionExpected) Name() string {
	return "FilePermissionExpected"
}

func (e *FilePermissionExpected) Evaluate(target string) (bool, error) {
	has, err := filesystem.FileHasPermission(target, e.PermissionCompared)
	if err != nil {
		return false, err
	}

	if e.ShouldExist {
		return has, nil
	} else {
		return !has, nil
	}
}

var _ ICondition[string] = (*FilePermissionExpected)(nil)

type PermittedEntityMatch struct {
	ExpectedOwner *string
	ExpectedGroup *string
}

func (e *PermittedEntityMatch) Name() string {
	return "PermittedEntityMatch"
}

func (e *PermittedEntityMatch) Evaluate(target string) (bool, error) {
	if e.ExpectedOwner != nil {
		name, err := filesystem.GetFileOwnerUserName(target)
		log.Printf("FileOwnerUsername is: %v", name)
		if err != nil {
			return false, err
		} else if name != *e.ExpectedOwner {
			return false, nil
		}
	}

	if e.ExpectedGroup != nil {
		name, err := filesystem.GetFileGroupName(target)
		log.Printf("FileGroupName is: %v", name)
		if err != nil {
			return false, err
		} else if name != *e.ExpectedGroup {
			return false, nil
		}
	}

	return true, nil
}

var _ ICondition[string] = (*PermittedEntityMatch)(nil)
