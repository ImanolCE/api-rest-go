// handlers/user_handler.go
package handlers

import (
    "context"
    "time"

    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"

    "github.com/ImanolCE/api-rest-go/config"
    "github.com/ImanolCE/api-rest-go/models"
    "github.com/ImanolCE/api-rest-go/utils"

    "golang.org/x/crypto/bcrypt"
)

// getCollectionUsers devuelve la colección "users"
func getCollectionUsers() *mongo.Collection {
    return config.ClientMongo.Database(config.DBName).Collection("users")
}

// RegisterUser crea un usuario nuevo donde hace el hash de contraseña + guardado en DB
func RegisterUser(c *fiber.Ctx) error {
    type Request struct {
        Nombre          string `json:"nombre"`
        Apellidos       string `json:"apellidos"`
        Email           string `json:"email"`
        Password        string `json:"password"`
        FechaNacimiento string `json:"fecha_nacimiento"` // ISO string: "2023-01-02"
        PreguntaSecreta string `json:"pregunta_secreta"`
        RespuestaSecreta string `json:"respuesta_secreta"`
    }

    var body Request
    if err := c.BodyParser(&body); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
    }

    // Verificar que el email no exista
    col := getCollectionUsers()
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    count, err := col.CountDocuments(ctx, bson.M{"email": body.Email})
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al verificar email"})
    }
    if count > 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email ya registrado"})
    }

    // Hashear contraseña
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al encriptar contraseña"})
    }

    // Parsear fecha
    fecha, err := time.Parse("2006-01-02", body.FechaNacimiento)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Fecha de nacimiento inválida (formato YYYY-MM-DD)"})
    }

    newUser := models.User{
        ID:               primitive.NewObjectID(),
        Nombre:           body.Nombre,
        Apellidos:        body.Apellidos,
        Email:            body.Email,
        Password:         string(hashedPassword),
        FechaNacimiento:  fecha,
        PreguntaSecreta:  body.PreguntaSecreta,
        RespuestaSecreta: body.RespuestaSecreta,
    }

    _, err = col.InsertOne(ctx, newUser)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al crear usuario"})
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Usuario registrado exitosamente"})
}

// LoginUser valida credenciales y retorna JWT
func LoginUser(c *fiber.Ctx) error {
    type Request struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    var body Request
    if err := c.BodyParser(&body); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
    }

    col := getCollectionUsers()
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var user models.User
    err := col.FindOne(ctx, bson.M{"email": body.Email}).Decode(&user)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email o contraseña incorrectos"})
    }

    // Comparar contraseña
    if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)) != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email o contraseña incorrectos"})
    }

    // Generar JWT
    token, err := utils.GenerarToken(user.ID.Hex())
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al generar token"})
    }

    return c.JSON(fiber.Map{"token": token})
}

// GetUsers lista todos los usuarios (solo para ejemplo; normalmente no expondrías datos sensibles)
func GetUsers(c *fiber.Ctx) error {
    col := getCollectionUsers()
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := col.Find(ctx, bson.M{})
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al listar usuarios"})
    }
    defer cursor.Close(ctx)

    var users []models.User
    if err := cursor.All(ctx, &users); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer usuarios"})
    }

    return c.JSON(users)
}

// GetUser obtiene un usuario por ID (parámetro URL)
func GetUser(c *fiber.Ctx) error {
    idParam := c.Params("id")
    objectID, err := primitive.ObjectIDFromHex(idParam)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
    }
    col := getCollectionUsers()
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var user models.User
    err = col.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Usuario no encontrado"})
    }
    return c.JSON(user)
}

// UpdateUser actualiza datos (excepto contraseña)
func UpdateUser(c *fiber.Ctx) error {
    idParam := c.Params("id")
    objectID, err := primitive.ObjectIDFromHex(idParam)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
    }

    var updates map[string]interface{}
    if err := c.BodyParser(&updates); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
    }

    // qui se prohibimos cambiar el campo "password" desde aquí
    delete(updates, "password")

    col := getCollectionUsers()
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := col.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updates})
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "No se pudo actualizar el usuario"})
    }
    if result.MatchedCount == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Usuario no encontrado"})
    }
    return c.JSON(fiber.Map{"message": "Usuario actualizado exitosamente"})
}

// DeleteUser elimina un usuario por ID
func DeleteUser(c *fiber.Ctx) error {
    idParam := c.Params("id")
    objectID, err := primitive.ObjectIDFromHex(idParam)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
    }

    col := getCollectionUsers()
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := col.DeleteOne(ctx, bson.M{"_id": objectID})
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "No se pudo eliminar el usuario"})
    }
    if result.DeletedCount == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Usuario no encontrado"})
    }
    return c.JSON(fiber.Map{"message": "Usuario eliminado exitosamente"})
}
