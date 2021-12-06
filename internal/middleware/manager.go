package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"joborder/internal/joborder"
	"joborder/internal/model"
	"net/http"
	"time"
)

type MiddleWare struct {
	service joborder.Service
}

func NewMiddleWareManager(service joborder.Service) *MiddleWare {
	return &MiddleWare{
		service: service,
	}
}

func (middleWare *MiddleWare) MiddleWareHandler()  *jwt.GinJWTMiddleware{
	var identityKey = "id"
	var passCode = "passCode"
	var role = "role"
	return &jwt.GinJWTMiddleware{
		Realm:       "user",
		Key:         []byte("12345"),
		Timeout:     time.Minute * 3,
		MaxRefresh:  time.Minute * 3,
		IdentityKey: "code",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					identityKey: v.PhoneNumber,
					passCode: v.Passcode,
					role: v.Role,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(context *gin.Context) interface{} {
			claims := jwt.ExtractClaims(context)
			return &model.User{
				PhoneNumber: claims[identityKey].(string),
				Passcode: claims[passCode].(string),
				Role: claims[role].(string),
			}
		},
		LoginResponse: func(context *gin.Context, i int, token string, expire time.Time) {
			claims := jwt.ExtractClaims(context)

			context.JSON(http.StatusOK, gin.H{
				"message": claims[identityKey],
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},
		RefreshResponse: func(context *gin.Context, i int, token string, expire time.Time) {
			claims := jwt.ExtractClaims(context)

			context.JSON(http.StatusOK, gin.H{
				"message": claims[identityKey],
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},
		Authenticator: func(context *gin.Context) (interface{}, error) {

			var loginReq model.Login
			if err := context.ShouldBind(&loginReq); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			response, err := middleWare.service.CheckUserByPhoneNumber(context.Request.Context(), &loginReq)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return &response,nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v,ok := data.(*model.User);ok && v.Role == "admin" || v.Role == "employer" {
				return true
			}
			return false
		},
		Unauthorized: func(context *gin.Context, code int, message string) {
			context.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: token",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
		SendCookie:       true,
		SecureCookie:     true, //non HTTPS dev environments
		CookieHTTPOnly:   true,  // JS can't modify
		CookieDomain:     "localhost:8080",
		CookieName:       "token",
		CookieSameSite:   http.SameSiteDefaultMode, //SameSiteDefaultMode, SameSiteLaxMode, SameSiteStrictMode, SameSiteNoneMode
	}
}
