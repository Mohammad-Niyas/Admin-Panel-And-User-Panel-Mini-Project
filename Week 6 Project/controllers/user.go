package controllers

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"main.go/auth"
	"main.go/database"
	"main.go/model"
)

var Error string
var Fetch model.UserModel
const RoleUser = "user"


// Landing Page (home page)

func Landingpage(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.HTML(http.StatusOK, "index.html", nil)
}


// user login page open handle

func Handler(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	session := sessions.Default(c)
	check := session.Get(RoleUser)
	if check == nil {
		c.HTML(http.StatusSeeOther, "user_login.html",Error)
		Error = ""
	} else {
		c.Redirect(http.StatusSeeOther, "/home")
	}
}

// user login handler

func LoginHandler(c *gin.Context) {
	database.DB.First(&Fetch, "email=?", c.Request.FormValue("email"))
	plainpassword := c.Request.FormValue("password")
	err := bcrypt.CompareHashAndPassword([]byte(Fetch.Password), []byte(plainpassword))
	if err != nil {
		Error = "Invalid Email ID or Password"
		c.Redirect(http.StatusFound, "/user/login")
	} else {
		if Fetch.Status == "Blocked" {
			Error = "Blocked User"
			c.Redirect(http.StatusFound, "/user/login")
		} else {
			auth.JwtToken(c, Fetch.Email, RoleUser)
			c.Redirect(http.StatusFound, "/home")
		}
	}
}

// User Home page open handler

func HomeHandler(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	session := sessions.Default(c)
	check := session.Get(RoleUser)

	if check != nil {
		c.HTML(http.StatusOK, "user_home.html", Fetch)
	} else {
		c.Redirect(http.StatusSeeOther, "/user/login")
	}
}

// User Logout Handler

func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(RoleUser)
	err := session.Save()
	if err != nil {
		log.Println("Error saving session:", err)
		c.String(http.StatusInternalServerError, "Logout failed.")
		return
	}
	c.Redirect(http.StatusSeeOther, "/user/login")
}

// User Sign Up Page Open Handler

func SignupHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "user_signup.html", nil)
}

// user signup post handler

func SignupPostHandler(c *gin.Context) {
	hashedPassword, error := bcrypt.GenerateFromPassword([]byte(c.Request.PostFormValue("password")), 10)
	if error != nil {
		log.Fatal("======= Password Hashing Failed =======")
	}
	err := database.DB.Create(&model.UserModel{
		Name:     c.Request.PostFormValue("username"),
		Email:    c.Request.PostFormValue("email"),
		Password: string(hashedPassword),
		Status:   "Active",
	}).Error
	if err != nil {
		Error = "Email ID Already Exists"
		c.Redirect(http.StatusSeeOther, "/user/signup")
	} else {
		Error = "User Successfully Created"
		c.Redirect(http.StatusSeeOther, "/user/login")
	}
}
