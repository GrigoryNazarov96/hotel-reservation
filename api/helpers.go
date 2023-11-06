package api

import (
	"context"
	"net/http"

	"github.com/GrigoryNazarov96/hotel-reservation.git/types"
)

func getAuthenticatedUser(c context.Context) (*types.User, error) {
	user, ok := c.Value("user").(*types.User)
	if !ok {
		return nil, NewError(http.StatusUnauthorized, "you should be logged in to perform this action")
	}
	return user, nil
}
