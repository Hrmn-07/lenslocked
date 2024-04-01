package main

import (
	"fmt"
	"net/http"

	"lenslocked/controllers"
	"lenslocked/models"
	"lenslocked/templates"
	"lenslocked/views"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", controllers.StaticHandler(views.Must(
		views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(views.Must(
		views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))

	r.Get("/faq", controllers.FAQ(views.Must(
		views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	// setup a database connection
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// setup a model service
	userService := models.UserService{
		DB: db,
	}

	usersC := controllers.Users{
		UserService: &userService,
	}

	usersC.Templates.New = (views.Must(
		views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml")))
	usersC.Templates.SignIn = (views.Must(
		views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml")))

	r.Get("/signup", usersC.New)
	r.Get("/signin", usersC.SignIn)
	r.Get("/users/me", usersC.CurrentUser)
	r.Post("/signup", usersC.Create)
	r.Post("/signin", usersC.ProcessSignIn)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
