package rate_limit_service

import (
	"log"
	"rate-limiter/pkg/infra/repository"
	"rate-limiter/pkg/infra/repository/rate-limit-rule-repository"
	"regexp"
)

type RateLimitService interface {
	CanProceed(key string) bool
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
func (s *service) CanProceed(key string) bool {
	rules := s.ruleRepository.GetRules()
	actualUsage := s.cacheRepository.Increment(key)

	// todo: range will copy, solve it!
	for pattern, limit := range *rules {
		if matched, err := regexp.MatchString(pattern, key); err == nil {
			if matched && actualUsage > limit {
				log.Printf("Key %s cannot be processed due to pattern %s with limit %d exceeded, actual: %d\n",
					pattern, key, limit, actualUsage)
				return false
			}
		} else {
			log.Printf("Error matcing pattern: %s key: %s, err: %+v\n", pattern, key, err)
		}
	}

	return true
}
