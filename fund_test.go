package funding

import (
	"sync"
	"testing"
)

const WORKERS = 10

func BenchmarkWithdrawals(b *testing.B) {
	// Skip N = 1
	if b.N < WORKERS {
		return
	}

	// Add as many dollars as we have iterations this run
	server := NewFundServer(b.N)

	// Casually assume b.N divides cleanly
	dollarsPerFounder := b.N / WORKERS

	// Waitgroup structs don't need to be initialized
	// (their "zero value" is ready to use)
	// So, we just declare one and then use it
	var wg sync.WaitGroup

	for i := 0; i < WORKERS; i++ {
		// Let waitgroup know we're adding a goroutine
		wg.Add(1)

		// Spawn off a founder worker as a closure
		go func() {
			// Mark this worker done when the function finishes
			defer wg.Done()

			for i := 0; i < dollarsPerFounder; i++ {
				server.Commands <- WithdrawCommand{Amount: 1}
			}
		}() // Remember to call the closure
	}

	// Wait for all the workers to finish
	wg.Wait()

	balanceResponseChan := make(chan int)
	server.Commands <- BalanceCommand{Response: balanceResponseChan}
	balance := <-balanceResponseChan

	if balance != 0 {
		b.Error("Balance wasn't zero:", balance)
	}
}

func BenchmarkFund(b *testing.B) {
	// Add as many dollars as we have iterations this run
	fund := NewFund(b.N)

	// Burn through them one at  a time until they are all gone
	for i := 0; i < b.N; i++ {
		fund.Withdraw(1)
	}

	if fund.Balance() != 0 {
		b.Error("Balance wasn't zero:", fund.Balance())
	}
}
