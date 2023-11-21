package api

import (
	"first-project/model"
	"first-project/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserAPI interface {
	GetCurrentUser(c *gin.Context)
	Register(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

func (a *userAPI) GetCurrentUser(c *gin.Context) {
	currentUser := a.userService.GetCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse("user not authenticated"))
		return
	}

	c.JSON(http.StatusOK, currentUser)
}

func (a *userAPI) Register(c *gin.Context) {
	// get the email and pass off request body
	var req model.RegisterData

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid decode json"))
		return
	}

	if req.Fullname == "" || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("register data is empty"))
		return
	}

	// create the user
	var reqBody = model.User{
		Fullname: req.Fullname,
		Email:    req.Email,
		Password: req.Password,
	}

	reqBody, err := a.userService.Register(&reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("internal server error"))
		return
	}
	
	// send respond
	c.JSON(http.StatusCreated, gin.H{
		"message": "register success",
		"data": reqBody,
	})
}

func (a *userAPI) Login(c *gin.Context) {
	// get the email and pass of request body
	var req model.LoginData

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid decode json"))
		return
	}

	if req.Email == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("email is empty"))
		return
	}
	
	if req.Password == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("password is empty"))
		return
	}

	// look up requested user
	var reqBody = model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	tokenKey, err := a.userService.Login(&reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// generate a jwt token
	tokenString := *tokenKey
	claims := &model.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return model.JwtKey, nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("internal server error"))
		return
	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse("token invalid"))
		return
	}

	// create a cookie
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    tokenString,
		Path:     "/",
		Domain:   "localhost",
		Expires:  claims.ExpiresAt.Time,
		MaxAge:   int(claims.ExpiresAt.Unix()),
		Secure:   false,
		HttpOnly: true,
	}
	
	// send it back
	c.Writer.Header().Add("Set-Cookie", cookie.String())

	c.JSON(http.StatusOK, gin.H{
		"email": claims.Email,
		"message": "Login Success",
	})	
}

func (a *userAPI) Logout(c *gin.Context) {
	token, err := c.Cookie("session_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("No session token found"))
		return
	}

	err = a.userService.Logout(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("Error logging out"))
		return
	}

	c.SetCookie("session_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, model.NewSuccessResponse("logout success"))
}