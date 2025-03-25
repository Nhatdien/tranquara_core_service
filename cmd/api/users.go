package main

import (
	"errors"
	"net/http"
	"time"

	"tranquara.net/internal/data"
	"tranquara.net/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		FullName string `json:"full_name"`
		Age      int8   `json:"age"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJson(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		FullName:  input.FullName,
		Age:       input.Age,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.User.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			app.badRequestResponse(w, r, data.ErrDuplicateEmail)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	token, err := app.models.Token.New(user.UserID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusCreated, user, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		data := map[string]interface{}{
			"activationToken": token.PlainText,
			"UserID":          user.UserID,
		}

		err := app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})

}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}
	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	if data.ValidateTokenPlainText(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.User.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	user.Activated = true

	err = app.models.User.Update(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Token.DeleteAllForUser(data.ScopeActivation, user.UserID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, user, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
