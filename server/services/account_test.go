package services

import (
	"log"
	"testing"
	"time"

	jwtutils "mlib.com/gofy/server/utils/jwt"
)

func TestJwt(t *testing.T) {
	jtokenstr, _, err := jwtutils.GenToken("ec7d12a0-f2a9-4f55-aa64-32140b7ef28f", "YqK3FapQWkWMlLiPCvqMWaTT822rfs8MVEnDbM4w0LjJbP/BOwH4jt9S", "SELF_HOSTED", 12*3600*time.Second)
	if err != nil {
		log.Printf("GenToken failed:%v\n", err)
	} else {
		cclaims, err := jwtutils.ParseToken(jtokenstr, "YqK3FapQWkWMlLiPCvqMWaTT822rfs8MVEnDbM4w0LjJbP/BOwH4jt9S")
		if err != nil {
			log.Printf("ParseToken failed:%v\n", err)
		} else {
			log.Printf("------cclaims=%#v\n", cclaims)
		}
	}
}
