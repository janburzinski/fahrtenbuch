package handlers

import "fahrtenbuch/pkg/models"

// these are basic structs
// some responses maybe need some more data
// in that case these are irrelevant and this is just dumb

/** LOGIN **/
type SuccessResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

type SuccessUserMeResponse struct {
	OK      bool        `json:"ok"`
	Message string      `json:"message"`
	User    models.User `json:"data"`
}

/** ERROR **/
type ErrorResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}
