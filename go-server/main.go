package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func main() {
	var err error
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "user=postgres password=postgres dbname=bbb sslmode=disable"
	}
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("bbb-session", store))

	r.LoadHTMLGlob("templates/*")

	r.GET("/register", showRegister)
	r.POST("/register", register)
	r.GET("/login", showLogin)
	r.POST("/login", login)
	r.GET("/logout", logout)
	r.GET("/admin", adminPanel)

	log.Println("Server running at http://localhost:8080")
	r.Run(":8080")
}

func showRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "fill all fields"})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": err.Error()})
		return
	}
	_, err = db.Exec("insert into users(username, password_hash) values($1,$2)", username, string(hashed))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusSeeOther, "/login")
}

func showLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var hash string
	var isAdmin bool
	err := db.QueryRow("select password_hash, is_admin from users where username=$1", username).Scan(&hash, &isAdmin)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "invalid credentials"})
		return
	}
	session := sessions.Default(c)
	session.Set("user", username)
	session.Set("is_admin", isAdmin)
	session.Save()
	c.Redirect(http.StatusSeeOther, "/admin")
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusSeeOther, "/login")
}

func adminPanel(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	isAdmin := session.Get("is_admin")
	if user == nil || isAdmin != true {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	rows, err := db.Query("select id, username, is_admin from users order by id")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()
	type uinfo struct {
		ID       int
		Username string
		Admin    bool
	}
	var users []uinfo
	for rows.Next() {
		var u uinfo
		if err := rows.Scan(&u.ID, &u.Username, &u.Admin); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		users = append(users, u)
	}
	c.HTML(http.StatusOK, "admin.html", gin.H{"users": users, "user": user})
}
