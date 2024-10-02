package auth

import (
	"NebuloGo/config"
	"NebuloGo/database"
	"NebuloGo/salt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var (
	identityKey = "id"
	userIdKey   = "user_id"
)

// User demo
type User struct {
	UserName string
	UserId   string
}

var JWTMiddleware *jwt.GinJWTMiddleware

func InitJWT() {
	var err error
	JWTMiddleware, err = jwt.New(initParams())
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
}

func HandlerMiddleWare(authMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return func(context *gin.Context) {
		errInit := authMiddleware.MiddlewareInit()
		if errInit != nil {
			log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
		}
	}
}

func initParams() *jwt.GinJWTMiddleware {

	return &jwt.GinJWTMiddleware{
		Realm:       "nebulogo",
		Key:         []byte(config.Configuration.JWT.Secret),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: payloadFunc(),

		IdentityHandler: identityHandler(),
		Authenticator:   authenticator(),
		Authorizator:    authorizator(),
		Unauthorized:    unauthorized(),
		TokenLookup:     "header: Authorization, query: token, cookie: token",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,

		SendCookie:     true,
		SecureCookie:   strings.HasPrefix(config.Configuration.Server.ServerURL, "https://"), //non HTTPS dev environments
		CookieHTTPOnly: true,                                                                 // JS can't modify
		CookieDomain:   strings.SplitN(config.Configuration.Server.ServerURL, "//", 2)[1],
		CookieName:     "token",
		CookieSameSite: http.SameSiteDefaultMode, //SameSiteDefaultMode, SameSiteLaxMode, SameSiteStrictMode, SameSiteNoneMode
	}
}

func payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*User); ok {
			return jwt.MapClaims{
				identityKey: v.UserName,
				userIdKey:   v.UserId,
			}
		}
		return jwt.MapClaims{}
	}
}

func identityHandler() func(c *gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		return &User{
			UserName: claims[identityKey].(string),
			UserId:   claims[userIdKey].(string),
		}
	}
}

func authenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var loginVals login
		var user *database.MongoUser
		var err error
		err = c.ShouldBind(&loginVals)
		if err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		loginID := loginVals.Username
		password := loginVals.Password
		user, err = database.ApplicationUserManager.GetUserByLoginID(loginID)
		if err != nil {
			return "", jwt.ErrFailedAuthentication
		}

		if salt.HashCompare(password, user.HashedPassword) {
			return &User{
				UserName: user.LoginID,
				UserId:   user.InternalID.String(),
			}, nil
		}
		return nil, jwt.ErrFailedAuthentication
	}
}

func authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if /*v*/ _, ok := data.(*User); ok /* && v.UserName == "neo.huyghe"*/ {
			return true
		}
		return false
	}
}

func unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		/*c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})*/
		c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("Expires", "0")
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
}

func HandleNoRoute() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("Expires", "0")
		c.Redirect(http.StatusMovedPermanently, "/drive")
		//claims := jwt.ExtractClaims(c)
		//log.Printf("NoRoute claims: %#v\n", claims)
		//c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	}
}

func HelloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)
	c.JSON(200, gin.H{
		"userID":   claims[identityKey],
		"userName": user.(*User).UserName,
		"text":     "Hello World.",
	})
}
