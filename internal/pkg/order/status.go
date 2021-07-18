package order

type Status string

const (
	StatusPending    Status = "pending"
	StatusPaid       Status = "paid"
	StatusProcessing Status = "processing"
	StatusCancelled  Status = "Cancelled"
)
