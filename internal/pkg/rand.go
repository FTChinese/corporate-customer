package pkg

import "github.com/FTChinese/go-rest/rand"

func OrderID() string {
	return "ord_" + rand.String(12)
}

func OrderItemID() string {
	return "ori_" + rand.String(12)
}
func LicenceID() string {
	return "lic_" + rand.String(12)
}

func InvitationID() string {
	return "invite_" + rand.String(12)
}

func ReceiptID() string {
	return "rcpt_" + rand.String(12)
}

func ReceiptItemID() string {
	return "rci_" + rand.String(12)
}

func SnapshotID() string {
	return "snp_" + rand.String(12)
}

// InvoiceID creates an id for addon invoice.
func InvoiceID() string {
	return "inv_" + rand.String(12)
}
