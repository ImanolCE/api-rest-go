// middleware/jwt_middleware.go
package middleware

import (
    "github.com/gofiber/fiber/v2"
   // "github.com/golang-jwt/jwt/v4"
    "github.com/ImanolCE/api-rest-go/utils"
)

// ExtractToken, es el que extrae el token del header "Authorization: Bearer <token>"
func ExtractToken(c *fiber.Ctx) string {
    authHeader := c.Get("Authorization")
    // El header viene como "Bearer <token>"
    if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
        return authHeader[7:]
    }
    return ""
}

// JWTMiddleware verifica la validez del token,	 antes de permitir el acceso a rutas protegidas
func JWTMiddleware(c *fiber.Ctx) error {
    tokenString := ExtractToken(c)
    if tokenString == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Falta token de autorización"})
    }
    claims, err := utils.ValidarToken(tokenString)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token inválido: " + err.Error()})
    }
    // Puedes pasar el UserID en locals para usarlo en el handler
    c.Locals("userID", claims.UserID)
    return c.Next()
}
