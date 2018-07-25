package funding

type Fund struct {
	// balance unexported (private), because it's lowercase
	balance int
}

// A regular function returning a pointer to a fund
func NewFund(initialBalance int) *Fund {
	return &Fund{
		balance: initialBalance,
	}
}

// Methods start with a reciever
func (f *Fund) Balance() int {
	return f.balance
}

func (f *Fund) Withdraw(amount int) {
	f.balance -= amount
}
