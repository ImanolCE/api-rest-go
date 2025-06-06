
package utils

import (
    "time"

    "github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("SECRETO_UTEQ2") 

// ClaimsPersonalizados hereda de jwt.RegisteredClaims para agregar campos extras
type ClaimsPersonalizados struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
}

// GenerarToken, es el qie genera un JWT para un userID, dado con expiración de 10 minutos
func GenerarToken(userID string) (string, error) {
    expira := time.Now().Add(10 * time.Minute)

    claims := &ClaimsPersonalizados{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expira),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// se valida el token (ValidarToken) parsea y valida un token JWT
func ValidarToken(tokenString string) (*ClaimsPersonalizados, error) {
    claims := &ClaimsPersonalizados{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    if err != nil || !token.Valid {
        return nil, err
    }
    return claims, nil
}

