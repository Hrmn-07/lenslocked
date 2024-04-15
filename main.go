package main

import (
	"fmt"
	"net/http"

	"lenslocked/controllers"
	"lenslocked/migrations"
	"lenslocked/models"
	"lenslocked/templates"
	"lenslocked/views"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func main() {
	// setup a database connection
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// run our goose migrations
	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	// setup model services
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}

	// setup middlewares
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}
	// setup csrf protection
	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		// TODO: fix before deploying (csrf needs https)
		csrf.Secure(false),
	)

	// setup controllers
	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	usersC.Templates.New = (views.Must(
		views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml")))
	usersC.Templates.SignIn = (views.Must(
		views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml")))

	// setup router
	r := chi.NewRouter()
	// these middlewares are used everywhere
	r.Use(csrfMw)
	r.Use(umw.SetUser)
	// now we setup routes
	r.Get("/", controllers.StaticHandler(views.Must(
		views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))
	r.Get("/contact", controllers.StaticHandler(views.Must(
		views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))
	r.Get("/faq", controllers.FAQ(views.Must(
		views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	r.Get("/signup", usersC.New)
	r.Get("/signin", usersC.SignIn)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})
	r.Post("/signup", usersC.Create)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// start the server
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
