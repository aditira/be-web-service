package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"tira.com/src/db"
	"tira.com/src/helper"
	"tira.com/src/model"
	"tira.com/src/module"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/lib/pq"

	jwt "github.com/dgrijalva/jwt-go"
)

type M map[string]interface{}

type MyClaims struct {
	jwt.StandardClaims
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

var APPLICATION_NAME = "Web Service APP"
var LOGIN_EXPIRATION_DURATION = time.Duration(1) * time.Hour
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
var JWT_SIGNATURE_KEY = []byte("the secret of ID")

func main() {
	router := new(CustomMux)

	// corsMiddleware := cors.New(cors.Options{
	// 	AllowedOrigins: []string{"*"},
	// 	// AllowedMethods: []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
	// 	// AllowedHeaders: []string{"append", "delete", "entries", "foreach", "get", "has", "keys", "set", "values", "Authorization"},
	// 	AllowedHeaders: []string{"Authorization"},
	// 	MaxAge:         3600,
	// 	// Debug:          true,
	// })

	// Handle CORS:
	// router.RegisterMiddleware(corsMiddleware.Handler)
	// Middleware:
	router.RegisterMiddleware(MiddlewareJWTAuthorization)

	fmt.Println("|===============================================================|")
	fmt.Println("|=====================>> " + APPLICATION_NAME + " <<=====================|")
	fmt.Println("|===============================================================|")

	router.HandleFunc("/api/public/signin", HandlerLogin)
	router.HandleFunc("/api/public/signup", HandlerRegister)
	router.HandleFunc("/api/auth/index", HandlerIndex)

	// API test:
	router.HandleFunc("/api/public/all", module.UserPublic)
	router.HandleFunc("/api/auth/user", module.UserGuest)
	router.HandleFunc("/api/auth/mod", module.UserModerator)
	router.HandleFunc("/api/auth/admin", module.UserAdmin)

	// API books:
	router.HandleFunc("/api/auth/get-books", module.GetBooks)
	router.HandleFunc("/api/auth/post-books", module.CreateBook)
	router.HandleFunc("/api/auth/delete-books/{bookid}", module.DeleteBook)
	router.HandleFunc("/api/auth/delete-books", module.DeleteBooks)

	// API collegers:
	router.HandleFunc("/api/auth/get-mahasiswa", module.GetCollegers)
	router.HandleFunc("/api/auth/post-mahasiswa", module.CreateColleger)
	router.HandleFunc("/api/auth/delete-mahasiswa/{nim}", module.DeleteColleger)
	router.HandleFunc("/api/auth/delete-mahasiswa", module.DeleteCollegers)

	// API colleger teacher:
	router.HandleFunc("/api/auth/get-dosen", module.GetCollegers)
	router.HandleFunc("/api/auth/post-dosen", module.CreateColleger)
	router.HandleFunc("/api/auth/delete-dosen/{nim}", module.DeleteColleger)
	router.HandleFunc("/api/auth/delete-dosen", module.DeleteCollegers)

	// Listen Port:
	log.Fatal(http.ListenAndServe(":8089", router))
}

func HandlerIndex(w http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	fmt.Println("Context: ", userInfo)
	message := fmt.Sprintf("Context info: username => %s | email => %s | role => %s", userInfo["username"], userInfo["email"], userInfo["role"])
	w.Write([]byte(message))
}

func HandlerRegister(w http.ResponseWriter, r *http.Request) {
	fmt.Println("|=====================>> Register")
	helper.CheckMethod(w, r, "POST")
	db := db.SetupDBPostgres()

	var p model.Register

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := p.Username
	email := p.Email
	password := p.Password
	role := p.Role

	var response = model.JsonResponse{}

	// Check Existing username:
	check_username, err := db.Query("SELECT username FROM user_demo WHERE username = $1", username)
	// Check Existing email:
	check_email, err := db.Query("SELECT email FROM user_demo WHERE email = $1", email)
	helper.CheckErr(err)

	if check_username.Next() {
		fmt.Println("|=====================>> Username exist")
		response = model.JsonResponse{Type: "failed", Message: "Username is already taken!"}
	} else if check_email.Next() {
		fmt.Println("|=====================>> Email exist")
		response = model.JsonResponse{Type: "failed", Message: "Email is already in use!"}
	} else if username == "" || email == "" || password == "" || role == "" {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
	} else {

		helper.PrintMessage("Inserting user into DB")

		var lastInsertID int
		err := db.QueryRow("INSERT INTO user_demo(username, email, password, role) VALUES($1, $2, SHA256($3), $4) returning id;", username, email, password, role).Scan(&lastInsertID)
		helper.CheckErr(err)

		response = model.JsonResponse{Type: "success", Message: "User registered successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

func HandlerLogin(w http.ResponseWriter, r *http.Request) {
	helper.CheckMethod(w, r, "POST")
	fmt.Println("|=====================>> Login")

	var p model.Login

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := p.Username
	password := p.Password

	// Check on DB:
	ok, userInfo := authenticateUserPg(username, password)
	if !ok {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    APPLICATION_NAME,
			ExpiresAt: time.Now().Add(LOGIN_EXPIRATION_DURATION).Unix(),
		},
		Id:       userInfo.Id,
		Username: userInfo.Username,
		Email:    userInfo.Email,
		Role:     userInfo.Role,
	}

	token := jwt.NewWithClaims(
		JWT_SIGNING_METHOD,
		claims,
	)

	signedToken, err := token.SignedString(JWT_SIGNATURE_KEY)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update DB Token:
	updateDBtoken(userInfo.Id, signedToken)

	tokenString, _ := json.Marshal(M{"token": signedToken})
	fmt.Println("|=>> Token: ", tokenString)
	// w.Write([]byte(tokenString))

	response := model.ResUser{
		Id:       userInfo.Id,
		Username: userInfo.Username,
		Email:    userInfo.Email,
		Roles:    userInfo.Role,
		Type:     "Bearer",
		Token:    signedToken,
	}
	json.NewEncoder(w).Encode(response)

}

func updateDBtoken(id int, token string) {
	db := db.SetupDBPostgres()
	_, err := db.Exec("UPDATE user_demo SET token = $1 WHERE id = $2", token, id)
	helper.CheckErr(err)

	fmt.Println("Token DB Updated!")
}

func authenticateUserPg(username, password string) (bool, model.AuthUser) {
	fmt.Println("User:", username)
	fmt.Println("Password:", password)
	db := db.SetupDBPostgres()

	rows, err := db.Query("SELECT id, username, email, role FROM user_demo WHERE username = $1 and password=sha256($2)", username, password)
	helper.CheckErr(err)

	var res model.AuthUser

	for rows.Next() {
		var id int
		var username string
		var email string
		var role string

		err = rows.Scan(&id, &username, &email, &role)
		helper.CheckErr(err)

		res = model.AuthUser{Id: id, Username: username, Email: email, Role: role}
	}

	if res.Username != "" {
		return true, res
	}

	return false, res
}

func MiddlewareJWTAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Miidleware path request: ", r.URL.Path)
		fmt.Println("Check Allow: ", strings.Contains(r.URL.Path, "public"))

		// Allow url start with /piblic/
		// if strings.HasPrefix(r.URL.Path, "/api/public/") {
		// 	next.ServeHTTP(w, r)
		// 	return
		// }

		// Allow url start with find text /piblic/
		if strings.Contains(r.URL.Path, "public") {
			next.ServeHTTP(w, r)
			return
		}

		authorizationHeader := r.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Signing method invalid")
			} else if method != JWT_SIGNING_METHOD {
				return nil, fmt.Errorf("Signing method invalid")
			}

			return JWT_SIGNATURE_KEY, nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(context.Background(), "userInfo", claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
