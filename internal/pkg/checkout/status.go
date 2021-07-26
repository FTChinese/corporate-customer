package checkout

type Status string

const (
	StatusPending    Status = "pending_payment"
	StatusPaid       Status = "paid"
	StatusProcessing Status = "processing"
	StatusCancelled  Status = "cancelled"
)
