package rate_limit_service

import (
	"context"
	"fmt"
	"github.com/erkanzileli/rate-limiter/model"
	"github.com/erkanzileli/rate-limiter/repository"
	"github.com/erkanzileli/rate-limiter/repository/rate-limit-rule-repository"
	"log"
	"time"
)

const (
	requestHashFormat  = "%s %s"
	incrementKeyFormat = "%s_%s"
)

type RateLimitService interface {
	CanProceed(ctx context.Context, method, path string) (bool, error)
}

type service struct {
	cacheRepository repository.CacheRepository
	ruleRepository  rate_limit_rule_repository.RateLimitRuleRepository
}

func New(
	cacheRepository repository.CacheRepository,
	ruleRepository rate_limit_rule_repository.RateLimitRuleRepository) RateLimitService {
	return &service{cacheRepository, ruleRepository}
}

// CanProceed takes key and tries to found a pattern that matching with given key.
// When it found then it compares pattern value with actualUsage.
func (s *service) CanProceed(ctx context.Context, method, path string) (canProceed bool, err error) {
	requestHash := fmt.Sprintf(requestHashFormat, method, path)
	matchedRule, anyMatch := findMatchedMinimumLimitRule(s.ruleRepository.GetRules(), requestHash)

	if !anyMatch {
		return true, nil
	}

	incrementKey := getIncrementKey(matchedRule, requestHash)
	actualUsage, err := s.cacheRepository.Increment(ctx, incrementKey)

	if err != nil {
		log.Printf("Request can't be limited due to cache repository error: %+v\n", err)
		return true, err
	}

	if actualUsage > matchedRule.Limit {
		log.Printf("Key %s cannot be processed due to pattern %s with limit %d was exceeded and actual is %d\n",
			requestHash, matchedRule.Pattern, matchedRule.Limit, actualUsage)
		return false, err
	}

	return true, nil
}

// getIncrementKey decides to which key will be incremented by the rule's scope
func getIncrementKey(rule model.Rule, requestHash string) string {
	incrementScope := requestHash
	if rule.IsPatternScope() {
		incrementScope = rule.Pattern
	}
	return fmt.Sprintf(incrementKeyFormat, incrementScope, getTimeWindow())
}

// findMatchedMinimumLimitRule loops over whole rules and returns a rule that has matched with the requestHash and its limit is lowest
func findMatchedMinimumLimitRule(rules []model.Rule, requestHash string) (result model.Rule, anyMatch bool) {
	for _, rule := range rules {
		if matched := rule.Regex.MatchString(requestHash); !matched {
			continue
		}
		if !anyMatch {
			anyMatch = true
			result = rule
			continue
		}
		if rule.Limit < result.Limit {
			result = rule
		}
	}
	return
}

// getTimeWindow returns actual time as hhmm format
// For example assume that current time is 15:04:05 then getTimeWindow returns 1504
func getTimeWindow() string {
	now := time.Now()
	return now.Format("1504")
}
