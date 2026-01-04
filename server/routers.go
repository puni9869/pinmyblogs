package server

import (
	session "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/internal/signup"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/puni9869/pinmyblogs/server/auth"
	"github.com/puni9869/pinmyblogs/server/home"
	"github.com/puni9869/pinmyblogs/server/middlewares"
	"github.com/puni9869/pinmyblogs/server/public"
	"github.com/puni9869/pinmyblogs/server/search"
	"github.com/puni9869/pinmyblogs/server/setting"
	"github.com/puni9869/pinmyblogs/types/forms"
)

// RegisterRoutes configures the available Web server routes.
func RegisterRoutes(r *gin.Engine, sessionStore session.Store) {
	log := logger.NewLogger()

	db := database.Db()
	signupService := signup.NewSignupService(db, log)

	//r.Use(middlewares.Cors())
	r.Use(middlewares.Session(sessionStore))
	loginRoutes := r.Group("")
	{
		loginRoutes.GET("/signup", auth.SignupGet)
		loginRoutes.POST("/signup",
			middlewares.Bind(forms.SignUpForm{}),
			auth.SignupPost(signupService),
		)
		loginRoutes.POST("/login", auth.LoginPost)
		loginRoutes.GET("/login", auth.LoginGet)
		loginRoutes.Any("/logout", auth.Logout)
		loginRoutes.GET("/reset-password", auth.ResetPasswordGet)
		loginRoutes.GET("/reset-password/sent", auth.ResetPasswordSentGet)
		loginRoutes.POST("/reset-password", middlewares.Bind(forms.ResetForm{}), auth.ResetPasswordPost)
		loginRoutes.GET("/reset-password/:hash", auth.ResetPasswordSetGet)
		loginRoutes.POST("/set-password", middlewares.Bind(forms.ResetPasswordForm{}), auth.ResetPasswordSetPost)
		loginRoutes.GET("/enable-my-account/:hash", setting.EnableMyAccount)
	}

	authRouters := r.Group("")
	{
		authRouters.Use(middlewares.AuthRequired)
		authRouters.Any("/", home.Home)
		authRouters.Any("/home", home.Home)
		authRouters.GET("/favourite", home.Favourite)
		authRouters.GET("/archived", home.Archived)
		authRouters.GET("/trash", home.Trash)
		authRouters.GET("/share/:id", home.Share)
		authRouters.PUT("/actions", home.Actions)
		authRouters.POST("/new", middlewares.Bind(forms.WeblinkRequest{}), home.AddWeblink)

		//Search routes
		authRouters.GET("/search", search.Search)

		// setting handler
		settingsRoute := authRouters.Group("/setting")
		{
			// Base renderer route
			settingsRoute.GET("", setting.Setting)

			// Profile routes
			settingsRoute.GET("/profile/:action", setting.ProfileAction)

			// Data routes
			settingsRoute.GET("/download-my-data/:format", setting.DownloadMyData)

			// Account related routes
			settingsRoute.DELETE("/delete-my-account", setting.DeleteMyAccount)
			settingsRoute.PUT("/disable-my-account", setting.DisableMyAccount)
		}
	}
	// public routes
	publicRouters := r.Group("")
	{
		// Register the service worker at the root level
		publicRouters.GET("/service-worker.js", public.ServiceWorker)
		publicRouters.GET("/offline.html", public.OfflinePage)

		publicRouters.GET("/health", public.Health)
		publicRouters.GET("/policies", public.PrivacyPolicyGet)
		publicRouters.GET("/support", public.SupportGet)
		publicRouters.GET("/start", public.LandingPageGet)

		publicRouters.POST("/start", middlewares.Bind(forms.JoinWaitList{}), public.JoinWaitListPost)
		publicRouters.GET("/favicon.ico", public.FavIcon)
		publicRouters.GET("/500", public.Route500)

		//	New theme
		publicRouters.GET("/start2", public.LandingPage2Get)
	}
	// this route will accept all the params
	r.NoRoute(public.Route404)
}
