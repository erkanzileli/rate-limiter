package rate_limit_service

import (
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
	CanProceed(method, uri string) bool
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
func (s *service) CanProceed(method, uri string) bool {
	windowedRequestKey := fmt.Sprintf(windowedRequestKeyFormat, method, uri, getTimeWindow())
	actualUsage := s.cacheRepository.Increment(windowedRequestKey)

	rules := s.ruleRepository.GetRules()
	requestKey := fmt.Sprintf(requestKeyFormat, method, uri)

	for _, rule := range rules {
		if matched, err := regexp.MatchString(rule.Pattern, requestKey); err == nil {
			if matched && actualUsage > rule.Limit {
				log.Printf("Key %s cannot be processed due to pattern %s with limit %d exceeded, actual: %d\n",
					rule.Pattern, requestKey, rule.Limit, actualUsage)
				return false
			}
		} else {
			log.Printf("Error matcing pattern: %s key: %s, err: %+v\n", rule.Pattern, requestKey, err)
		}
	}

	return true
}

// getTimeWindow returns actual time as hhmm format
// For example assume that current time is 15:04:05
// getTimeWindow returns 1504
func getTimeWindow() string {
	now := time.Now()
	return now.Format("1504")
}
