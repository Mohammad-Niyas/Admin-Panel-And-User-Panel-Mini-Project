package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"main.go/controllers"
	"main.go/database"
)

func init() {
	database.DBconnect()
}

func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// load templates

	r.LoadHTMLGlob("template/*")

	// Landing Page Route

	r.GET("/", controllers.Landingpage)

	// User Routes
	
	r.GET("/user/signup", controllers.SignupHandler)
	r.POST("/signup", controllers.SignupPostHandler)
	r.GET("/user/login", controllers.Handler)
	r.POST("/login", controllers.LoginHandler)
	r.GET("/home", controllers.HomeHandler)
	r.GET("/logout", controllers.LogoutHandler)

	// Admin Routes

	r.GET("/admin", controllers.AdminLoginHandler)
	r.POST("/admin/login", controllers.AdminPostHandler)
	r.GET("/admin/home", controllers.AdminHomeHandle)
	r.GET("/delete/:ID", controllers.DeleteUserHandler)
	r.GET("/block/:ID",controllers.BlockUserHandler)
	r.GET("/edit/:ID",controllers.UserEditHandler)
	r.POST("/User/edit/:ID",controllers.UserPostEditHandler)
	r.GET("/admin/logout",controllers.AdminLogoutHandler)
	r.POST("/add/user",controllers.AdminAddUserHandler)

	r.Run(":9999")
}
