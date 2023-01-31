package rule

import "github.com/aws/amazon-cloudwatch-agent-test/filesystem"

type Rule[T any] struct {
	Conditions []ICondition[T]
}

func (r *Rule[T]) Evaluate(target T) (bool, error) {
	for _, c := range r.Conditions {
		success, err := c.Evaluate(target)
		if err != nil {
			return false, err
		}
		if !success {
			return false, nil
		}
	}
	return true, nil
}

type ICondition[T any] interface {
	Evaluate(T) (bool, error)
}

type FilePermissionExpected struct {
	PermissionCompared filesystem.FilePermission
	ShouldExist        bool
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

func (e *PermittedEntityMatch) Evaluate(target string) (bool, error) {
	if e.ExpectedOwner != nil {
		err := filesystem.CheckFileOwnerRights(target, *e.ExpectedOwner)
		if err != nil {
			return false, err
		}
	}

	if e.ExpectedGroup != nil {
		name, err := filesystem.GetFileGroupName(target)
		if err != nil {
			return false, err
		} else if name != *e.ExpectedGroup {
			return false, nil
		}
	}

	return true, nil
}

var _ ICondition[string] = (*PermittedEntityMatch)(nil)

