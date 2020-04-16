package admin

import "github.com/FTChinese/go-rest/enum"

// Letter is the data passed to template to generate the content
// of an email.
type Letter struct {
	URL      string // The link for verification, or password reset.
	IsSignUp bool   // Determines greeting.
}

type InvitationLetter struct {
	AdminEmail string
	TeamName   string
	Tier       enum.Tier
	URL        string
}
