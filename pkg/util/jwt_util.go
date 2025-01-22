package util

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

//func MakeAuthorizedHandler(f http.HandlerFunc, role string) http.HandlerFunc {
//
//	return func(w http.ResponseWriter, r *http.Request) {
//
//		tokenString := r.Header.Get("Auth")
//		var userId int
//		var httpError *HTTPError
//		if userId, httpError = ValidateAccessToken(tokenString); httpError != nil {
//			WriteJSON(w, httpError.StatusCode, httpError.Error)
//			return
//		}
//
//		if httpError = CheckAuthorization(userId, role); httpError != nil {
//			WriteJSON(w, httpError.StatusCode, httpError.Error)
//
//			return
//		}
//		f(w, r)
//	}
//}

func IssueToken(uniqueId int32) (tokenString string, httpError *HTTPError) {

	token := jwt.New(jwt.SigningMethodHS256)

	accessTokenClaims := token.Claims.(jwt.MapClaims)
	accessTokenClaims["id"] = uniqueId
	accessTokenClaims["exp"] = time.Now().Add(time.Hour * 2).Unix()

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

//func CheckAuthorization(userId int, role string) *HTTPError {
//
//	if ok, err := db.Client.CheckAuthorization(context.Background(), sqlc.CheckAuthorizationParams{UserID: int32(userId), Role: role}); err != nil {
//		return &HTTPError{StatusCode: 500, Error: "internal server error"}
//
//	} else if !ok {
//		return &HTTPError{StatusCode: 403, Error: "missing permissions"}
//	}
//	return nil
//}
