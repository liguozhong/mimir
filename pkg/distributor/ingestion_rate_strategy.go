// SPDX-License-Identifier: AGPL-3.0-only
// Provenance-includes-location: https://github.com/cortexproject/cortex/blob/master/pkg/distributor/ingestion_rate_strategy.go
// Provenance-includes-license: Apache-2.0
// Provenance-includes-copyright: The Cortex Authors.

package distributor

import (
	"github.com/grafana/dskit/limiter"
	"golang.org/x/time/rate"

	"github.com/grafana/mimir/pkg/util/validation"
)

// ReadLifecycler represents the read interface to the lifecycler.
type ReadLifecycler interface {
	HealthyInstancesCount() int
}

type globalStrategy struct {
	limits *validation.Overrides
	ring   ReadLifecycler
}

func newGlobalIngestionRateStrategy(limits *validation.Overrides, ring ReadLifecycler) limiter.RateLimiterStrategy {
	return &globalStrategy{
		limits: limits,
		ring:   ring,
	}
}

func (s *globalStrategy) Limit(tenantID string) float64 {
	numDistributors := s.ring.HealthyInstancesCount()

	if numDistributors == 0 {
		return s.limits.IngestionRate(tenantID)
	}

	return s.limits.IngestionRate(tenantID) / float64(numDistributors)
}

func (s *globalStrategy) Burst(tenantID string) int {
	// The meaning of burst doesn't change for the global strategy, in order
	// to keep it easier to understand for users / operators.
	return s.limits.IngestionBurstSize(tenantID)
}

type infiniteStrategy struct{}

func newInfiniteIngestionRateStrategy() limiter.RateLimiterStrategy {
	return &infiniteStrategy{}
}

func (s *infiniteStrategy) Limit(tenantID string) float64 {
	return float64(rate.Inf)
}

func (s *infiniteStrategy) Burst(tenantID string) int {
	// Burst is ignored when limit = rate.Inf
	return 0
}
