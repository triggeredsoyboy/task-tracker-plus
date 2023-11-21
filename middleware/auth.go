package middleware

import (
	"first-project/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		// get cookie session_token from req
		sessionToken, err := ctx.Cookie("session_token")
		if err != nil {
			// no cookie
			if ctx.GetHeader("Content-Type") == "application/json" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			} else {
				// redirect to login page if there's no cookie
				ctx.Redirect(http.StatusSeeOther, "/user/login")
			}
			// stop req
			ctx.Abort()
			return
		}

		// parsing jwt token
		claims := &model.Claims{}
		token, err := jwt.ParseWithClaims(sessionToken, claims, func(t *jwt.Token) (interface{}, error) {
			return model.JwtKey, nil
		})
		
		// parsing token fail
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(err.Error()))
			ctx.Abort()
			return
		}

		// token not valid
		if !token.Valid {
			ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(err.Error()))
			ctx.Abort()
			return
		}

		// set email from claims to context
		ctx.Set("email", claims.Email)

		// continue to the handler
		ctx.Next()
	})
}