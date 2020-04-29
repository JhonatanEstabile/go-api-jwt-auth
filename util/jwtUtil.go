package util

import (
	"api-jwt-auth/model"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

//CreateJwtToken generate a jwt token that will expire in 15 minutes
func CreateJwtToken(userid uint64) (*model.TokenDetails, error) {
	var err error

	td := &model.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	td, err = createAuthenticationToken(userid, td)
	if err != nil {
		return nil, err
	}

	return createRefreshToken(userid, td)
}

func createAuthenticationToken(userid uint64, td *model.TokenDetails) (*model.TokenDetails, error) {
	var err error

	os.Setenv("ACCESS_SECRET", "tknjdpaimxalsmdrk") //this should be in an env file

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func createRefreshToken(userid uint64, td *model.TokenDetails) (*model.TokenDetails, error) {
	var err error

	os.Setenv("REFRESH_SECRET", "pasijdybelaceidjuiybbak") //this should be in an env file

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
}
