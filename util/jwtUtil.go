package util

import (
	"api-jwt-auth/model"
	"fmt"
	"net/http"
	"os"
	"strconv"
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

//VerifyToken check if a token is valid
func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

//TokenValid check if token is expired
func TokenValid(r *http.Request) error {
	tokenString := ExtractToken(r)

	token, err := VerifyToken(tokenString)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

//ExtractTokenMetadata decrypt the token content
func ExtractTokenMetadata(tokenString string) (*model.AccessDetails, error) {
	token, err := VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}

		userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}

		return &model.AccessDetails{
			AccessUuid: accessUUID,
			UserId:     userID,
		}, nil
	}

	return nil, err
}
