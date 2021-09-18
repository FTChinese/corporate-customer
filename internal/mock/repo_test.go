package mock

import (
	"testing"
)

func TestRepo_InsertOrderSchema(t *testing.T) {

	r := NewRepo()

	schema := NewAdmin().CartBuilder().
		AddNewStandardN(10).
		AddRenewalStandardN(10).
		AddNewPremiumN(5).
		AddRenewalPremiumN(5).
		BuildOrderSchema()

	r.InsertOrderSchema(schema)
}

func TestRepo_CreateGrantedLicence(t *testing.T) {

	granted := NewAdmin().StdLicenceBuilder().SetPersona(NewPersona()).BuildGranted()

	t.Logf("Assignee %v", granted.SignUp)
	t.Logf("Licence %v", granted.ExpLicence)
	t.Logf("AnteChange %v", granted.Membership)

	NewRepo().CreateGrantedLicence(granted)
}
