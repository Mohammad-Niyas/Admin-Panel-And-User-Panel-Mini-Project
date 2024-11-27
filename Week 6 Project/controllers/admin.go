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

var AdFetch model.AdminModel
var AdError string
var UserTable []model.UserModel
var Updatefetch model.UserModel

const RoleAdmin = "admin"

// admin login handler

func AdminLoginHandler(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	session := sessions.Default(c)
	check := session.Get(RoleAdmin)
	if check == nil {
		c.HTML(http.StatusOK, "admin_login.html", AdError)
		AdError = ""
	} else {
		c.Redirect(http.StatusSeeOther, "/admin/home")
	}
}

// admin authentication 

func AdminPostHandler(c *gin.Context) {
	AdFetch = model.AdminModel{}
	database.DB.First(&AdFetch, "email=?", c.Request.PostFormValue("email"))
	if AdFetch.Password != c.Request.FormValue("password") {
		AdError = "Invalid Username or Password"
		c.Redirect(http.StatusSeeOther, "/admin")
	} else {
		auth.JwtToken(c, AdFetch.Email, RoleAdmin)
		c.Redirect(http.StatusSeeOther, "/admin/home")
	}
}

// admin home page
func AdminHomeHandle(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	session := sessions.Default(c)
	check := session.Get(RoleAdmin)
	if check != nil {
		database.DB.Find(&UserTable)
		c.HTML(http.StatusOK, "admin.html", gin.H{
			"Datas": UserTable,
			"Admin": AdFetch.Name,
			"Error": AdError,
		})
		AdError = ""
	} else {
		c.Redirect(http.StatusSeeOther, "/admin")
	}
}

// delete User

func DeleteUserHandler(c *gin.Context) {
	c.Header("Cache-control", "no-cache,no-store,musst-revalidate")
	session := sessions.Default(c)
	check := session.Get(RoleAdmin)
	if check != nil {
		User := c.Param("ID")
		database.DB.First(&Updatefetch, "ID=?", User)
		database.DB.Delete(&Updatefetch)
		Updatefetch = model.UserModel{}
		c.Redirect(http.StatusSeeOther, "/admin/home")
		AdError = "Deleted Successfully"
	} else {
		c.Redirect(http.StatusSeeOther, "/admin")
	}
}

// block User

func BlockUserHandler(c *gin.Context) {
	c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
	session := sessions.Default(c)
	check := session.Get(RoleAdmin)
	if check != nil {

		User := c.Param("ID")

		database.DB.First(&Updatefetch, "ID=?", User)

		if Updatefetch.Status == "Active" {
			Updatefetch.Status = "Blocked"
			database.DB.Save(&Updatefetch)
			Updatefetch = model.UserModel{}
			c.Redirect(http.StatusSeeOther, "/admin/home")
			AdError = "Blocked User!"

		} else {
			Updatefetch.Status = "Active"
			database.DB.Save(&Updatefetch)
			Updatefetch = model.UserModel{}
			c.Redirect(http.StatusSeeOther, "/admin/home")
			AdError = "Activated User"
		}
	} else {
		c.Redirect(http.StatusSeeOther, "/admin")
	}
}

// User Edit handler

func UserEditHandler(c *gin.Context) {
	c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
	session := sessions.Default(c)
	check := session.Get(RoleAdmin)
	if check != nil {
		User := c.Param("ID")
		c.HTML(http.StatusSeeOther, "admin_edit.html", User)
	} else {
		c.Redirect(http.StatusSeeOther, "/admin")
	}
}

// email and name update

func UserPostEditHandler(c *gin.Context) {
	User := c.Param("ID")
	database.DB.Find(&Updatefetch, "ID=?", User)
	Updatefetch.Name = c.Request.FormValue("name")
	Updatefetch.Email = c.Request.FormValue("email")
	database.DB.Save(&Updatefetch)
	AdError = "Updated Successfully"
	Updatefetch = model.UserModel{}
	c.Redirect(http.StatusSeeOther, "/admin/home")
}

// New User Add 

func AdminAddUserHandler(c *gin.Context) {
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
		AdError = "Email ID Already Exists"
		c.Redirect(http.StatusSeeOther, "/admin/home")
	} else {
		AdError = "User Successfully Created"
		c.Redirect(http.StatusSeeOther, "/admin/home")
	}
}
// admin logout handler

func AdminLogoutHandler(c *gin.Context) {
	c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
	session := sessions.Default(c)
	session.Delete(RoleAdmin)
	session.Save()
	c.Redirect(http.StatusSeeOther, "/admin")
}
