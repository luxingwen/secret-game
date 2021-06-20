package controller

import (
	jwt "github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"

	"fmt"
	"time"

	"net/http"
)

var jwtSecret = []byte("secret-key")

type WxClaims struct {
	Id int `json:"id"`
	jwt.StandardClaims
}

func GenerateWxToken(id int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := WxClaims{
		id,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "secret-game",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseWxToken(token string) (*WxClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &WxClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*WxClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

func WxJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		token := c.Query("token")
		if token == "" {
			token = c.GetHeader("X-Token")
			if token == "" {
				code = 999
			}
		}
		if code == 0 && token != "" {
			claims, err := ParseWxToken(token)
			if err != nil {
				code = 1
				fmt.Println("code 11111..")
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = 2
				fmt.Println("code 222222..")
			} else {
				c.Set("wxUserId", claims.Id)
			}
		}
		if code != 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  "无效的token",
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
