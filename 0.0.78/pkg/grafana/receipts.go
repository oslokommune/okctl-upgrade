package grafana

func NewReceiptsStack() Receipts {
	return Receipts{receipts: make([]Receipt, 0)}
}

// Push adds a receipt to the receipt stack
func (receiver *Receipts) Push(r Receipt) {
	receiver.receipts = append(receiver.receipts, r)
}
