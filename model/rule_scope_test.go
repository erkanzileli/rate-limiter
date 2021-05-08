package model_test

import (
	"github.com/erkanzileli/rate-limiter/model"
	testifyAssert "github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewRuleScope_when_scope_is_valid(t *testing.T) {
 	assert := testifyAssert.New(t)

	// Given
	scope := "pattern"

	// When
	ruleScope := model.NewRuleScope(scope)

	// Then
	assert.Equal(model.PatternScope, ruleScope)
}

func Test_NewRuleScope_when_scope_is_invalid(t *testing.T) {
 	assert := testifyAssert.New(t)

	// Given
	scope := ""

	// When
	ruleScope := model.NewRuleScope(scope)

	// Then
	assert.Equal(model.PathScope, ruleScope)
}
