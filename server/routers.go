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
	// diagnose url
	r.GET("/health", home.Health)

	r.GET("/signup", auth.SignupGet)
	r.POST("/signup",
		middlewares.BindForm(forms.SignUpForm{}),
		auth.SignupPost(signupService),
	)
	// auth urls
	r.POST("/login", auth.LoginPost)
	r.GET("/login", auth.LoginGet)
	r.Any("/logout", auth.Logout)
	r.GET("/reset", auth.ResetPasswordGet)
	r.POST("/reset", middlewares.BindForm(forms.ResetForm{}), auth.ResetPasswordPost)

	authRouters := r.Group("")
	{
		authRouters.Use(middlewares.AuthRequired)
		authRouters.Any("/", home.Home)
		authRouters.Any("/home", home.Home)
		authRouters.GET("/favourite", home.Favourite)
		authRouters.GET("/archived", home.Archived)
		authRouters.GET("/trash", home.Trash)

		authRouters.POST("/new", home.AddWeblink)

		// setting handler
		settingsRoute := authRouters.Group("/setting")
		{
			settingsRoute.GET("", setting.Setting)
			settingsRoute.DELETE("/deletemyaccount", setting.DeleteMyAccount)
			settingsRoute.PUT("/disablemyaccount", setting.DisableMyAccount)
		}

		// navbar handler
		//authRouters.
	}

	// public routes
	_ = r.Group("")
	{
		//publicRouters.GET("/", public.StartGet)
		//publicRouters.Any("/start", public.StartGet)
	}

	// this route will accept all the params
	r.NoRoute(auth.LoginGet)
}
