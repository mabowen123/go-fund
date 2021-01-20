package middleware

import (
	"crypto/md5"
	"fmt"
	"fund/handlers"
	"fund/mysql"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var (
	secret      = "AOtUF2JXQ8nrudTw" //salt
	expireTime  = time.Hour * 24     //token expire time
	identityKey = "id"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func AuthMiddleware() *jwt.GinJWTMiddleware {
	authMiddleware, _ := jwt.New(&jwt.GinJWTMiddleware{
		Key:            []byte(secret),
		Timeout:        expireTime,
		MaxRefresh:     expireTime,
		IdentityKey:    identityKey,
		PayloadFunc:    payloadFunc(),
		Authenticator:  authenticator(),
		LoginResponse:  loginResponse(),
		LogoutResponse: logoutResponse(),
		Authorizator:   authorizator(),
		Unauthorized:   unauthorized(),
		TokenLookup:    "header: Authorization,  cookie: jwt",
		TokenHeadName:  "Bearer",
		TimeFunc:       time.Now,
	})

	return authMiddleware
}

func authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if id, ok := data.(float64); ok {
			user := &mysql.User{}
			result := mysql.Db.Select("nickname, username").First(user, uint64(id))
			if result.RecordNotFound() {
				return false
			}
			claims := jwt.ExtractClaims(c)
			return claims["nickname"] == user.Nickname && claims["username"] == user.Username
		}
		return false
	}
}
func loginResponse() func(*gin.Context, int, string, time.Time) {
	return func(c *gin.Context, httpCode int, token string, expire time.Time) {
		c.JSON(httpCode, handlers.Success.WithData(map[string]interface{}{
			"token":  token,
			"expire": expire.Format("2006-01-02 15:04:05"),
		}))
	}
}

func logoutResponse() func(ctx *gin.Context, code int) {
	return func(c *gin.Context, httpCode int, ) {
		c.JSON(httpCode, handlers.Success)
	}
}

func unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, handlers.Fail.WithMsg("unauthorized"))
	}
}

func authenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		loginVal := login{}
		if err := c.ShouldBind(&loginVal); err != nil {
			return "", jwt.ErrMissingLoginValues
		}

		user := &mysql.User{}
		result := mysql.Db.Where(&mysql.User{Username: loginVal.Username, Password: md5Password(loginVal.Password)}).First(user)
		if !result.RecordNotFound() {
			return user, nil
		}

		return nil, jwt.ErrFailedAuthentication
	}
}

func payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if user, ok := data.(*mysql.User); ok {
			return jwt.MapClaims{
				"id":       user.ID,
				"username": user.Username,
				"nickname": user.Nickname,
			}
		}
		return jwt.MapClaims{}
	}
}

type user struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Nickname string `form:"nickname" json:"nickname" binding:"required"`
}

func Register(c *gin.Context) {
	var params user
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, handlers.Fail.WithMsg(err.Error()))
		return
	}

	user := &mysql.User{}
	res := mysql.Db.Where(&mysql.User{Username: params.Username}).First(&user)

	if !res.RecordNotFound() {
		c.JSON(http.StatusOK, handlers.Fail.WithMsg("用户名已存在"))
		return
	}

	createData := &mysql.User{
		Username: params.Username,
		Password: md5Password(params.Password),
		Nickname: params.Nickname,
	}
	res = mysql.Db.Create(createData)
	if res.RecordNotFound() {
		c.JSON(http.StatusOK, handlers.Fail.WithMsg("创建失败"))
		return
	}

	tokenString, _, _ := AuthMiddleware().TokenGenerator(createData)
	c.JSON(http.StatusOK, handlers.Success.WithData(map[string]interface{}{
		"token": tokenString,
	}))
}

func md5Password(password string) string {
	salt := "Dw1cjXYANM6eo5mL"
	str := []byte(password + salt)
	return fmt.Sprintf("%x", md5.Sum(str))
}
