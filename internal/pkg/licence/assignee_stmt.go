package licence

// Retrieve Assignee as JSON object so that we don't need to
// create extra types to convert string to JSON.
const colAssigneeJSON = `
JSON_OBJECT(
	"ftcId", a.user_id,
	"email", a.email,
	"userName", a.user_name
) AS assignee`

const selectAssignee = `
SELECT user_id 	AS ftc_id,
	email 	AS user_email,
	user_name AS user_name
FROM cmstmp01.userinfo
`

const StmtAssigneeByID = selectAssignee + `
WHERE user_id = ?
LIMIT 1`

const StmtAssigneeByEmail = selectAssignee + `
WHERE email = ?
LIMIT 1`
