package admin

// AccessRight is used to limit admin's right to data.
// Image what happens when data is only retrieved by row id?
// Admin could visit a specific row using url like:
// https://....../:id
// Just by simply change the id an admin could access the row,
// even if the row is not created by him/her.
// Ideally we should only allow a team access its data created by it.
type AccessRight struct {
	RowID  string
	TeamID string
}
