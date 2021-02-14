package rate_limit_service

import (
	"context"
	"fmt"
	"github.com/erkanzileli/rate-limiter/pkg/model"
	"github.com/erkanzileli/rate-limiter/pkg/repository"
	"github.com/erkanzileli/rate-limiter/pkg/repository/rate-limit-rule-repository"
	"log"
	"time"
)

const (
	requestHashFormat  = "%s %s"
	windowedRuleFormat = "%s_%s"
)

type RateLimitService interface {
	CanProceed(ctx context.Context, method, uri string) (bool, error)
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
func (s *service) CanProceed(ctx context.Context, method, uri string) (bool, error) {
	requestHash := fmt.Sprintf(requestHashFormat, method, uri)
	matchedMinimumLimitRule := findMatchedMinimumLimitRule(s.ruleRepository.GetRules(), requestHash)

	if matchedMinimumLimitRule == nil {
		return true, nil
	}

	windowedPatternKey := fmt.Sprintf(windowedRuleFormat, matchedMinimumLimitRule.Pattern, getTimeWindow())
	actualUsage, err := s.cacheRepository.Increment(ctx, windowedPatternKey)

	if err != nil {
		log.Printf("Request can't be limited due to cache repository error: %+v\n", err)
		return true, err
	}

	if actualUsage > matchedMinimumLimitRule.Limit {
		log.Printf("Key %s cannot be processed due to pattern %s with limit %d was exceeded and actual is %d\n",
			requestHash, matchedMinimumLimitRule.Pattern, matchedMinimumLimitRule.Limit, actualUsage)
		return false, err
	}

	return true, nil
}

// findMatchedMinimumLimitRule loops over whole rules and returns a rule that has matched with the requestHash and its limit is lowest
func findMatchedMinimumLimitRule(rules []*model.Rule, requestHash string) (matchedMinimumLimitRule *model.Rule) {
	for _, rule := range rules {
		if matched := rule.Regex.MatchString(requestHash); matched {
			if matchedMinimumLimitRule == nil {
				matchedMinimumLimitRule = rule
				continue
			}
			if rule.Limit < matchedMinimumLimitRule.Limit {
				matchedMinimumLimitRule = rule
			}
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
