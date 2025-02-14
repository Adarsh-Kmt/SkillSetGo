package helper

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	logger = log.New(os.Stdout, "SKILLSETGO AUTH >> ", 0)
)

func MakeAuthenticatedHandler(f http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Auth")

		if _, httpError := ValidateAccessToken(tokenString); httpError != nil {

			WriteJSON(w, httpError.StatusCode, map[string]any{"error": httpError.Error})
			return
		}

		f(w, r)
	}
}

func MakeAuthorizedHandler(f http.HandlerFunc, requiredRoles []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Auth")

		var httpError *HTTPError
		if _, httpError = ValidateAccessToken(tokenString); httpError != nil {
			WriteJSON(w, httpError.StatusCode, httpError.Error)
			return
		}

		if httpError = CheckAuthorization(tokenString, requiredRoles); httpError != nil {
			WriteJSON(w, httpError.StatusCode, httpError.Error)

			return
		}
		f(w, r)
	}
}

func IssueToken(uniqueId int32, roles []string) (tokenString string, httpError *HTTPError) {

	token := jwt.New(jwt.SigningMethodHS256)

	accessTokenClaims := token.Claims.(jwt.MapClaims)
	accessTokenClaims["id"] = uniqueId
	accessTokenClaims["exp"] = time.Now().Add(time.Hour * 2).Unix()
	accessTokenClaims["roles"] = roles

	secret := os.Getenv("JWT_PRIVATE_KEY")
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		logger.Printf("error : %s", err.Error())
		return "", &HTTPError{StatusCode: 500, Error: "internal server error."}
	}

	return tokenString, nil
}

func ValidateAccessToken(tokenString string) (id int, httpError *HTTPError) {

	secret := os.Getenv("JWT_PRIVATE_KEY")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("wrong signing key")
		}

		return []byte(secret), nil
	})

	if err != nil {
		logger.Printf("error : %s", err.Error())
		return 0, &HTTPError{StatusCode: 401, Error: err.Error()}
	}

	claims := token.Claims.(jwt.MapClaims)
	logger.Println(claims)

	idFloat, ok := claims["id"].(float64)

	id = int(idFloat)
	if !ok {
		logger.Printf("error : invalid id format")
		return 0, &HTTPError{StatusCode: 500, Error: "internal server error"}

	}

	return id, nil
}

func CheckAuthorization(tokenString string, requiredRoles []string) *HTTPError {

	secret := os.Getenv("JWT_PRIVATE_KEY")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("wrong signing key")
		}

		return []byte(secret), nil
	})

	if err != nil {
		logger.Printf("error : %s", err.Error())
		return &HTTPError{StatusCode: 401, Error: map[string]string{"error:": err.Error()}}
	}

	claims := token.Claims.(jwt.MapClaims)

	rolesInterface, exists := claims["roles"]
	if !exists {
		log.Println("Roles not found in token")
		return &HTTPError{StatusCode: 403, Error: map[string]string{"error": "missing permissions"}}
	}

	var existingRoles []string
	for _, roleVal := range rolesInterface.([]interface{}) {
		if strRole, ok := roleVal.(string); ok {
			existingRoles = append(existingRoles, strRole)
		}
	}

	for _, requiredRole := range requiredRoles {

		for _, exitingRole := range existingRoles {

			if requiredRole == exitingRole {
				return nil
			}
		}
	}

	return &HTTPError{StatusCode: 403, Error: map[string]string{"error": "missing permissions"}}
}
