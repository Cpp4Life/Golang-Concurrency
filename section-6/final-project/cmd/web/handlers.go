package main

import (
	"Concurrency/section-6/final-project/constants"
	"Concurrency/section-6/final-project/data"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *Config) HomePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.gohtml", nil)
}

func (app *Config) LoginPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", nil)
}

func (app *Config) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.RenewToken(r.Context())

	// parse form post
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	// get email and password from form post
	email := r.Form.Get(constants.EmailTag)
	password := r.Form.Get(constants.PasswordTag)

	user, err := app.Models.User.GetByEmail(email)
	if err != nil {
		app.Session.Put(r.Context(), constants.ErrorTag, "Invalid credentials.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// check password
	validPassword, err := user.PasswordMatches(password)
	if err != nil {
		app.Session.Put(r.Context(), constants.ErrorTag, "Invalid credentials.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !validPassword {
		msg := Message{
			To:      email,
			Subject: "Failed log in attempt",
			Data:    "Invalid login attempt!",
		}

		app.sendEmail(msg)

		app.Session.Put(r.Context(), constants.ErrorTag, "Invalid credentials.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// okay, so log user in
	app.Session.Put(r.Context(), constants.UserIdTag, user.ID)
	app.Session.Put(r.Context(), constants.UserTag, user)
	app.Session.Put(r.Context(), constants.FlashTag, "Successful login!")

	// redirect the user
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Config) Logout(w http.ResponseWriter, r *http.Request) {
	// clean up session
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) RegisterPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gohtml", nil)
}

func (app *Config) PostRegisterPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	// TODO - validate data

	// create a user
	u := data.User{
		Email:     r.Form.Get(constants.EmailTag),
		FirstName: r.Form.Get(constants.FirstNameTag),
		LastName:  r.Form.Get(constants.LastNameTag),
		Password:  r.Form.Get(constants.PasswordTag),
		Active:    0,
		IsAdmin:   0,
	}

	_, err = u.Insert(u)
	if err != nil {
		app.Session.Put(r.Context(), constants.ErrorTag, "Unable to create user.")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// send an activation email
	url := fmt.Sprintf("http://localhost:8080/activate?email=%s", u.Email)
	signedURL := GenerateTokenFromString(url)
	app.InfoLog.Println(signedURL)

	msg := Message{
		To:       u.Email,
		Subject:  "Activate your account",
		Template: "confirmation-email",
		Data:     template.HTML(signedURL),
	}

	app.sendEmail(msg)

	app.Session.Put(r.Context(), constants.FlashTag, "Confirmation email sent. Check your mail.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	// validate url
	url := r.RequestURI
	app.InfoLog.Println(url)
	testURL := fmt.Sprintf("http://localhost:8080%s", url)
	okay := VerifyToken(testURL)

	if !okay {
		app.Session.Put(r.Context(), constants.ErrorTag, "Invalid token.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// activate account
	u, err := app.Models.User.GetByEmail(r.URL.Query().Get(constants.EmailTag))
	if err != nil {
		app.Session.Put(r.Context(), constants.ErrorTag, "No user found.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	u.Active = 1
	err = u.Update()
	if err != nil {
		app.Session.Put(r.Context(), constants.ErrorTag, "Unable to update user.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), constants.FlashTag, "Account activated. You can now log in")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) ChooseSubscription(w http.ResponseWriter, r *http.Request) {
	plans, err := app.Models.Plan.GetAll()
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	dataMap := make(map[string]any)
	dataMap["plans"] = plans

	app.render(w, r, "plans.page.gohtml", &TemplateData{
		Data: dataMap,
	})
}

func (app *Config) SubscribeToPlan(w http.ResponseWriter, r *http.Request) {
	// get the id of the plan that is chosen
	id := r.URL.Query().Get(constants.IdTag)

	planId, _ := strconv.Atoi(id)

	// get the plan from the database
	plan, err := app.Models.Plan.GetOne(planId)
	if err != nil {
		app.Session.Put(r.Context(), constants.ErrorTag, "Unable to find plan.")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		return
	}

	// get the user from the session
	user, ok := app.Session.Get(r.Context(), constants.UserTag).(data.User)
	if !ok {
		app.Session.Put(r.Context(), constants.ErrorTag, "Log in first!")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// generate an invoice

	// send an email with the invoice attached

	// generate a manual

	// send an email with the manual attached

	// subscribe the user to an account

	// redirect
}
