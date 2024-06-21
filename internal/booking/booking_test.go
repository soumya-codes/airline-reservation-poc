package booking

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/soumya-codes/airline-reservation-poc/config"
	bookingseat "github.com/soumya-codes/airline-reservation-poc/internal/booking/seat"
	postgrestransaction "github.com/soumya-codes/airline-reservation-poc/internal/postgres/transaction"
)

func TestBookSeats(t *testing.T) {
	lockStrategies := []struct {
		strategy     bookingseat.LockStrategy
		strategyName string
	}{
		{strategy: bookingseat.GetSeatWithNoLock, strategyName: "GetSeatWithNoLock"},
		{strategy: bookingseat.GetSeatWithSharedLock, strategyName: "GetSeatWithSharedLock"},
		{strategy: bookingseat.GetSeatWithSharedLockSkipped, strategyName: "GetSeatWithSharedLockSkipped"},
		{strategy: bookingseat.GetSeatWithExclusiveLock, strategyName: "GetSeatWithExclusiveLock"},
		{strategy: bookingseat.GetSeatWithExclusiveLockSkipped, strategyName: "GetSeatWithExclusiveLockSkipped"},
	}

	poolSizes := []int{1, 5, 50, 180}

	isolationLevels := []postgrestransaction.IsolationLevel{
		postgrestransaction.ReadCommitted,
		postgrestransaction.RepeatableRead,
		postgrestransaction.Serializable,
	}

	for _, isolationLevel := range isolationLevels {
		for _, strategy := range lockStrategies {
			for _, poolSize := range poolSizes {
				t.Run(fmt.Sprintf("IsolationLevel=%v_LockStrategy=%s_PoolSize=%d",
					isolationLevel, strategy.strategyName, poolSize),
					func(t *testing.T) {
						cfg := config.NewConfig(config.WithMaxConn(poolSize),
							config.WithTxIsolation(isolationLevel),
							config.WithLockStrategy(strategy.strategy))

						// Adjust the timeout based on your infra setup.
						ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
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
