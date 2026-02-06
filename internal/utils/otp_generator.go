package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

//func for genetate random otp
func GenerateOTP()(string,error){
	max:=big.NewInt(100000)
	n,err:=rand.Int(rand.Reader,max)
	if err!=nil{
		return "",err
	}
	return fmt.Sprintf("%06d",n.Int64()),nil
}