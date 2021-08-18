package oauth

import (
	"errors"
	"net/http"
	"strings"
)

var errTokenRequired = errors.New("no access credentials provided")

// ParseBearer extracts Authorization header.
// Authorization: Bearer 19c7d9016b68221cc60f00afca7c498c36c361e3
func ParseBearer(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("empty authorization header")
	}

	s := strings.SplitN(authHeader, " ", 2)

	bearerExists := (len(s) == 2) && (strings.ToLower(s[0]) == "bearer")

	if !bearerExists {
		return "", errors.New("bearer not found")
	}

	return s[1], nil
}

// GetBearerAuth tries to extract the Bearer value from
// Authorization header.
func GetBearerAuth(header http.Header) (string, error) {
	authVal := header.Get("Authorization")

	if authVal == "" {
		return "", errTokenRequired
	}

	return ParseBearer(authVal)
}

func getTokenFromQuery(req *http.Request) (string, error) {
	token := req.Form.Get("access_token")
	if strings.TrimSpace(token) == "" {
		return "", errTokenRequired
	}

	return token, nil
}

// GetToken extract oauth token from http request.
// It first tries to find it from `Authorization: Bearer xxxxx`
// header, then fallback to url query parameter `access_token`
// field.
// If nothing is found, returns error.
func GetToken(req *http.Request) (string, error) {
	authHeader, err := GetBearerAuth(req.Header)
	if err == nil {
		return authHeader, nil
	}

	return getTokenFromQuery(req)
}
