package auth

import (
	"fmt"
	"net/http"

	"github.com/Moe-Ajam/ldr-sync-server/internal/helpers"
)

func VerifyToken(w http.ResponseWriter, r *http.Request) (string, error) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			helpers.RespondWithError(w, http.StatusUnauthorized, "No cookie present")
			return "", err
		}
		fmt.Printf("Something went wrong while hanlding the request: %v\n", err)
		helpers.RespondWithError(w, http.StatusBadRequest, "Something went wrong")
		return "", err
	}

	tknString := c.Value
	return tknString, nil
}
