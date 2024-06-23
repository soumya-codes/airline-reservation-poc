package booking

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/soumya-codes/airline-reservation-poc/config"
	"github.com/soumya-codes/airline-reservation-poc/internal/booking/seat"
	pgtx "github.com/soumya-codes/airline-reservation-poc/internal/postgres/transaction"
)

func TestBookSeats(t *testing.T) {
	isolationLevels := []pgtx.IsolationLevel{
		pgtx.ReadCommitted,
		pgtx.RepeatableRead,
		pgtx.Serializable,
	}

	lockStrategies := []struct {
		strategy     seat.LockStrategy
		strategyName string
	}{
		{strategy: seat.GetSeatWithNoLock, strategyName: "GetSeatWithNoLock"},
		{strategy: seat.GetSeatWithSharedLock, strategyName: "GetSeatWithSharedLock"},
		{strategy: seat.GetSeatWithSharedLockSkipped, strategyName: "GetSeatWithSharedLockSkipped"},
		{strategy: seat.GetSeatWithExclusiveLock, strategyName: "GetSeatWithExclusiveLock"},
		{strategy: seat.GetSeatWithExclusiveLockSkipped, strategyName: "GetSeatWithExclusiveLockSkipped"},
	}

	poolSizes := []int{1, 5, 50, 180}

	retries := 3

	for _, isolationLevel := range isolationLevels {
		for _, strategy := range lockStrategies {
			for _, poolSize := range poolSizes {
				t.Run(fmt.Sprintf("IsolationLevel=%v_LockStrategy=%s_PoolSize=%d_Retries=%d",
					isolationLevel, strategy.strategyName, poolSize, retries),
					func(t *testing.T) {
						cfg := config.NewConfig(config.WithMaxConn(poolSize),
							config.WithTxIsolation(isolationLevel),
							config.WithLockStrategy(strategy.strategy),
							config.WithMaxRetries(retries),
						)

						// Adjust the timeout based on your infra setup.
						ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
						defer cancel()

						err := BookSeats(ctx, cfg)
						if err != nil {
							t.Logf("error booking seats: %v", err)
						}
					})

				// Sleep for 3 seconds to allow the connections to be released
				time.Sleep(3 * time.Second)
			}
		}
	}
}
