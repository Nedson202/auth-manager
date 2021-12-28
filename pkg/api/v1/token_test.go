package v1

import (
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/smartystreets/goconvey/convey"
)

var user = &UserWithoutPassword{
	Username: "Convey",
}

func Test_generateToken(t *testing.T) {
	convey.Convey("Generate token", t, func() {
		convey.Convey("should generate token", func() {
			cfg := &Config{
				JwtSecret: "jwt_secret",
			}
			tokenService := NewTokenService(*cfg)
			token, err := tokenService.GenerateToken(user)

			convey.So(err, convey.ShouldBeNil)
			convey.So(token, convey.ShouldNotBeNil)
		})
	})
}

func Test_verifyToken(t *testing.T) {
	convey.Convey("Verify token", t, func() {
		convey.Convey("should verify and decode token", func() {
			cfg := &Config{
				JwtSecret: "jwt_secret",
			}
			tokenService := NewTokenService(*cfg)
			token, _ := tokenService.GenerateToken(user)

			claims, err := tokenService.VerifyToken(token)
			convey.So(err, convey.ShouldBeNil)
			convey.So(claims.Valid(), convey.ShouldBeNil)

			decodedUser := &UserWithoutPassword{}
			mapstructure.Decode(claims, &decodedUser)

			convey.So(decodedUser.Username, convey.ShouldEqual, user.Username)
		})
	})
}
