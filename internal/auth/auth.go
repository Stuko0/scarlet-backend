package auth

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/Stuko0/scarlet-backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	accessTokenLife time.Duration
}

func NewJWTManager(privateKey []byte, publicKey []byte, accessTokenLife time.Duration) (*JWTManager,error) {
	privKey, err:= jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {return nil, err}
	pubKey, err:= jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {return nil, err}
	return &JWTManager{
		privateKey: privKey,
		publicKey:  pubKey,
		accessTokenLife: accessTokenLife,
	},nil
}

func (m *JWTManager) Generate(user *models.User)(string, error){
	claims :=Claims{
		UserId: user.UserId,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessTokenLife)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token:=jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(m.privateKey)
}

func (m *JWTManager) Verify(tokenString string)(*Claims, error){
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return m.publicKey, nil
	},
	)
	if err != nil {return nil, fmt.Errorf("invalid token: %v", err)}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims,nil
}