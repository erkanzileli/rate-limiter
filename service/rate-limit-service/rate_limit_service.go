package rate_limit_service

import (
	"context"
	"fmt"
	model2 "github.com/erkanzileli/rate-limiter/model"
	repository2 "github.com/erkanzileli/rate-limiter/repository"
	rate_limit_rule_repository2 "github.com/erkanzileli/rate-limiter/repository/rate-limit-rule-repository"
	"log"
	"time"
)

const (
	requestHashFormat  = "%s %s"
	windowedRuleFormat = "%s_%s"
)

type RateLimitService interface {
	CanProceed(ctx context.Context, method, path string) (bool, error)
}

type service struct {
	cacheRepository repository2.CacheRepository
	ruleRepository  rate_limit_rule_repository2.RateLimitRuleRepository
}

func New(
	cacheRepository repository2.CacheRepository,
	ruleRepository rate_limit_rule_repository2.RateLimitRuleRepository) RateLimitService {
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

	windowedPatternKey := fmt.Sprintf(windowedRuleFormat, matchedRule.Pattern, getTimeWindow())
	actualUsage, err := s.cacheRepository.Increment(ctx, windowedPatternKey)

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

// findMatchedMinimumLimitRule loops over whole rules and returns a rule that has matched with the requestHash and its limit is lowest
func findMatchedMinimumLimitRule(rules []model2.Rule, requestHash string) (result model2.Rule, anyMatch bool) {
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
