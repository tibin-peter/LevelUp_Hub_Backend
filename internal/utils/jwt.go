package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

//func for access token
func GenerateAccessToken(userID uint,email,role,jwtkey string)(string,time.Time,error){
	exp:=time.Now().Add(15 *time.Minute)
	claims:=&Claims{
		UserID: userID,
		Email: email,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	signed,err:=token.SignedString([]byte(jwtkey))
	return signed,exp,err
}

//func for refresh token
func GenerateRefreshToken(userID uint,email,role,jwtkey string)(string,time.Time,error){
	exp:=time.Now().Add(7*24 *time.Hour)
	claims:=&Claims{
		UserID: userID,
		Email: email,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	signed,err:=token.SignedString([]byte(jwtkey))
	return signed,exp,err
}

//func for validate token
func ValidateToken(tokenStr,jwtkey string)(*Claims,error){
	claims:=&Claims{}
		token,err:=jwt.ParseWithClaims(tokenStr,claims,func(t *jwt.Token) (any, error) {
       
			 //ensure correct siging method
			 if _,ok:=t.Method.(*jwt.SigningMethodHMAC);!ok{
				return nil,jwt.ErrSignatureInvalid
			 }

		return []byte(jwtkey),nil
		})
		if err!=nil||!token.Valid{
			return nil,err
		}
		return claims,nil
}
