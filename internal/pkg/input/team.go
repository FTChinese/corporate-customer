package input

import (
	"github.com/FTChinese/ftacademy/pkg/validator"
	"github.com/FTChinese/go-rest/render"
	"github.com/guregu/null"
	"strings"
)

type TeamParams struct {
	OrgName      string      `json:"orgName" db:"org_name"`
	InvoiceTitle null.String `json:"invoiceTitle" db:"invoice_title"`
}

func (t *TeamParams) Validate() *render.ValidationError {
	t.OrgName = strings.TrimSpace(t.OrgName)
	title := strings.TrimSpace(t.InvoiceTitle.String)
	t.InvoiceTitle = null.NewString(title, title != "")

	ve := validator.New("orgName").Required().MaxLen(128).Validate(t.OrgName)
	if ve != nil {
		return ve
	}

	return validator.New("invoiceTitle").MaxLen(512).Validate(title)
}

func (t *TeamParams) IsEqual(newVal TeamParams) bool {
	return t.OrgName == newVal.OrgName && t.InvoiceTitle == newVal.InvoiceTitle
}
