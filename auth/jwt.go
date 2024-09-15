package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("your_secret_key")

// Claims defines the structure for JWT claims
type Claims struct {
    Email string `json:"email"`
    jwt.RegisteredClaims
}

// GenerateJWT generates a new JWT for a user
func GenerateJWT(email string) (string, error) {
    expirationTime := time.Now().Add(1 * time.Hour)

    claims := &Claims{
        Email: email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

// ValidateJWT validates the JWT token and extracts the claims
func ValidateJWT(tokenStr string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil || !token.Valid {
        return nil, err
    }
    return claims, nil
}
