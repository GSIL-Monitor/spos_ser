// jwt
package umisc

import (
	//  "github.com/gin-gonic/gin"
	"errors"
	"github.com/dgrijalva/jwt-go"
	//	"log"
	//	"net/http"
	"fmt"
	"spos_auth/claims"
	"time"
)

//// midware, check token
//func JWTAuth() gin.HandlerFunc {
//    return func(c *gin.Context) {
//        token :=  c.Request.Header.Get("token")
//        if token == ""{
//            c.JSON(http.StatusOK,gin.H{
//                "status":-1,
//                "msg":"ÇëÇóÎ´Ð¯´øtoken£¬ÎÞÈ¨ÏÞ·ÃÎÊ",
//            })
//            c.Set("isPass", false)
//            return
//        }

//        log.Print("get token: ",token)

//        j := NewJWT()
//        // parseToken
//        claims, err := j.ParseToken(token)
//        if err != nil {
//            if err == TokenExpired {
//                c.JSON(http.StatusOK,gin.H{
//                    "status":-1,
//                    "msg":"ÊÚÈ¨ÒÑ¹ýÆÚ",
//                })
//                c.Set("isPass", false)
//                return
//            }
//            c.JSON(http.StatusOK, gin.H{
//                "status": -1,
//                "msg": err.Error(),
//                })
//            c.Set("isPass", false)
//            return
//        }
//        c.Set("isPass", true)
//        c.Set("claims",claims)
//    }
//}
// Sign
type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token not active yet")
	TokenMalformed   error  = errors.New("That's not even a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token:")
	SignKey          string = "unitone_spos"
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}
func GetSignKey() string {
	return SignKey
}
func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

func (j *JWT) CreateToken(claims claims.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenString string) (*claims.Claims, error) {

	fmt.Printf("ParseToken: tokenString=%s\n", tokenString)
	token, err := jwt.ParseWithClaims(tokenString, &claims.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	fmt.Printf("ParseToken: token=%v; err=%v\n", token, err)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}

	if claims, ok := token.Claims.(*claims.Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &claims.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*claims.Claims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
