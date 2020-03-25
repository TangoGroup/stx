package stx

import (
	"errors"

	"cuelang.org/go/cue"
	"github.com/TangoGroup/stx/logger"
)

// Stack represents the decoded value of stacks[stackname]
type Stack struct {
	Name, Profile, SopsProfile, Region, Environment, RegionCode string
	DependsOn                                                   []string
	Tags                                                        map[string]string
	Hooks                                                       struct {
		PostSave []string
	}
}

// StacksIterator is a wrapper around cue.Iterator that allows for filtering based on stack fields
type StacksIterator struct {
	cueIter cue.Iterator
	flags   Flags
	log     *logger.Logger
}

// NewStacksIterator returns *StacksIterator
func NewStacksIterator(cueInstance *cue.Instance, flags Flags, log *logger.Logger) (*StacksIterator, error) {
	log.Debug("Getting stacks...")
	stacks := cueInstance.Value().Lookup("Stacks")
	if !stacks.Exists() {
		return nil, errors.New("Stacks is undefined")
	}

	fields, fieldsErr := stacks.Fields()
	if fieldsErr != nil {
		return nil, fieldsErr
	}

	return &StacksIterator{cueIter: fields, flags: flags, log: log}, nil
}

// Next moves the index forward and applies global filters. returns true if there is a value that passes the filters
func (it *StacksIterator) Next() bool {
	if !it.cueIter.Next() {
		return false
	}

	currentValue := it.cueIter.Value()
	// currentLabel, _ := currentValue.Label()

	// apply filters to the current value
	if it.flags.Environment != "" {
		environmentValue := currentValue.Lookup("Environment")
		if !environmentValue.Exists() {
			return it.Next()
		}
		environment, environmentErr := environmentValue.String()
		if environmentErr != nil {
			it.log.Error(environmentErr)
			return it.Next()
		}
		if it.flags.Environment != environment {
			return it.Next()
		}
	}

	if it.flags.RegionCode != "" {
		regionCodeValue := currentValue.Lookup("RegionCode")
		if !regionCodeValue.Exists() {
			return it.Next()
		}
		regionCode, regionCodeErr := regionCodeValue.String()
		if regionCodeErr != nil {
			it.log.Error(regionCodeErr)
			return it.Next()
		}
		if it.flags.RegionCode != regionCode {
			return it.Next()
		}
	}

	if it.flags.Profile != "" {
		profileValue := currentValue.Lookup("Profile")
		if !profileValue.Exists() {
			return it.Next()
		}
		profile, profileErr := profileValue.String()
		if profileErr != nil {
			it.log.Error(profileErr)
			return it.Next()
		}
		if it.flags.Profile != profile {
			return it.Next()
		}
	}

	return true
}

// Value returns the value from the cue.Iterator
func (it *StacksIterator) Value() cue.Value {
	return it.cueIter.Value()
}
