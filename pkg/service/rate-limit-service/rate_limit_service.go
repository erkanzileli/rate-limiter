package rate_limit_service

import (
	"context"
	"fmt"
	"log"
	"rate-limiter/pkg/repository"
	"rate-limiter/pkg/repository/rate-limit-rule-repository"
	"regexp"
	"time"
)

const (
	requestKeyFormat         = "%s_%s"
	windowedRequestKeyFormat = "%s_%s_%s"
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
	windowedRequestKey := fmt.Sprintf(windowedRequestKeyFormat, method, uri, getTimeWindow())

	actualUsage, err := s.cacheRepository.Increment(ctx, windowedRequestKey)
	if err != nil {
		return false, err
	}

	rules := s.ruleRepository.GetRules()
	requestKey := fmt.Sprintf(requestKeyFormat, method, uri)

	for _, rule := range rules {
		if matched, err := regexp.MatchString(rule.Pattern, requestKey); err == nil {
			if matched && actualUsage > rule.Limit {
				log.Printf("Key %s cannot be processed due to pattern %s with limit %d was exceeded and actual is %d\n",
					requestKey, rule.Pattern, rule.Limit, actualUsage)
				return false, nil
			}
		} else {
			log.Printf("Error matcing pattern: %s key: %s, err: %+v\n", rule.Pattern, requestKey, err)
		}
	}

	return true, nil
}

// getTimeWindow returns actual time as hhmm format
// For example assume that current time is 15:04:05 then getTimeWindow returns 1504
func getTimeWindow() string {
	now := time.Now()
	return now.Format("1504")
}
