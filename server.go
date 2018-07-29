package funding

import (
	"fmt"
)

type FundServer struct {
	commands chan interface{} // Lowercase name, unexported
	fund     Fund
}

// Typedef the callback for readability
type Transactor func(fund Fund)

// Add a new command type with a callback and a semaphore channel
type TransactionCommand struct {
	Transactor Transactor
	Done       chan bool
}

func (s *FundServer) Balance() int {
	responseChan := make(chan int)
	s.commands <- BalanceCommand{Response: responseChan}
	return <-responseChan
}

func (s *FundServer) Withdraw(amount int) {
	s.commands <- WithdrawCommand{Amount: amount}
}

func NewFundServer(initialBalance int) *FundServer {

	server := &FundServer{
		// make() creates builtins like channels, maps, and slices
		commands: make(chan interface{}),
		fund:     *NewFund(initialBalance),
	}

	// Spawn off the server's main loop immediately
	go server.loop()
	return server
}

func (s *FundServer) loop() {
	// The built-in "range" clause can iterate over channels, amongst other things
	for command := range s.commands {
		// Handle command
		switch command.(type) {

		case WithdrawCommand:
			// And then use a "type assertion" to convert it
			withdrawal := command.(WithdrawCommand)
			s.fund.Withdraw(withdrawal.Amount)

		case BalanceCommand:
			getBalance := command.(BalanceCommand)
			balance := s.fund.Balance()
			getBalance.Response <- balance

		case TransactionCommand:
			transaction := command.(TransactionCommand)
			transaction.Transactor(s.fund)
			transaction.Done <- true

		default:
			panic(fmt.Sprintf("Unrecognized command: %v", command))
		}

	}
}

type WithdrawCommand struct {
	Amount int
}

type BalanceCommand struct {
	Response chan int
}

func (s *FundServer) Transact(transactor Transactor) {
	command := TransactionCommand{
		Transactor: transactor,
		Done:       make(chan bool),
	}

	s.commands <- command
	<-command.Done
}
