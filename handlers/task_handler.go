
package handlers

import (
    "context"
    "time"

    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"

    "github.com/tu_usuario/fiber-tasks-api/config"
    "github.com/tu_usuario/fiber-tasks-api/models"
)

// getCollectionTasks devuelve la colección "tasks"
func getCollectionTasks() *mongo.Collection {
    return config.ClientMongo.Database(config.DBName).Collection("tasks")
}

// CreateTask permite a un usuario autenticado crear una nueva task
func CreateTask(c *fiber.Ctx) error {
    type Request struct {
        Titulo      string `json:"titulo"`
        Descripcion string `json:"descripcion"`
        FechaInicio string `json:"fecha_inicio"` // "2006-01-02T15:04:05Z07:00"
        FechaFinal    string `json:"fecha_final"`
    }
    var body Request
    if err := c.BodyParser(&body); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
    }

    // Obtener userID de locals (puesto por el middleware JWT)
    userIDHex := c.Locals("userID").(string)
    userObjID, err := primitive.ObjectIDFromHex(userIDHex)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "UserID inválido"})
    }

    // Parsear fechas
    fechaInicio, err := time.Parse(time.RFC3339, body.FechaInicio)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Fecha de inicio inválida"})
    }
    FechaFinal, err := time.Parse(time.RFC3339, body.FechaFinal)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "FechaFinal inválido"})
    }

    newTask := models.Task{
        ID:           primitive.NewObjectID(),
        Titulo:       body.Titulo,
        Descripcion:  body.Descripcion,
        FechaInicio:  fechaInicio,
        FechaFinal:     FechaFinal,
        UsuarioID:    userObjID,
    }

    col := getCollectionTasks()
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err = col.InsertOne(ctx, newTask)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al crear task"})
    }
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Task creada exitosamente"})
}

// GetTasks retorna las tasks del usuario autenticado
func GetTasks(c *fiber.Ctx) error {
    userIDHex := c.Locals("userID").(string)
    userObjID, err := primitive.ObjectIDFromHex(userIDHex)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "UserID inválido"})
    }

    col := getCollectionTasks()
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := col.Find(ctx, bson.M{"usuario_id": userObjID})
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al listar tasks"})
    }
    defer cursor.Close(ctx)

    var tasks []models.Task
    if err := cursor.All(ctx, &tasks); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer tasks"})
    }
    return c.JSON(tasks)
}

// GetTask obtiene una task específica (solo si pertenece al usuario)
func GetTask(c *fiber.Ctx) error {
    idParam := c.Params("id")
    taskID, err := primitive.ObjectIDFromHex(idParam)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
    }
    userIDHex := c.Locals("userID").(string)
    userObjID, _ := primitive.ObjectIDFromHex(userIDHex)

    col := getCollectionTasks()
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var task models.Task
    filter := bson.M{"_id": taskID, "usuario_id": userObjID}
    err = col.FindOne(ctx, filter).Decode(&task)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task no encontrada"})
    }
    return c.JSON(task)
}

// UpdateTask actualiza campos de una task (solo si es del usuario)
func UpdateTask(c *fiber.Ctx) error {
    idParam := c.Params("id")
    taskID, err := primitive.ObjectIDFromHex(idParam)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
    }
    userIDHex := c.Locals("userID").(string)
    userObjID, _ := primitive.ObjectIDFromHex(userIDHex)

    var updates map[string]interface{}
    if err := c.BodyParser(&updates); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
    }

    // Si vienen fechas en string, parsearlas
    if val, ok := updates["fecha_inicio"].(string); ok {
        t, err := time.Parse(time.RFC3339, val)
        if err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Fecha de inicio inválida"})
        }
        updates["fecha_inicio"] = t
    }
    if val, ok := updates["fecha_final"].(string); ok {
        t, err := time.Parse(time.RFC3339, val)
        if err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "FechaFinal inválido"})
        }
        updates["fecha_final"] = t
    }

    col := getCollectionTasks()
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    filter := bson.M{"_id": taskID, "usuario_id": userObjID}
    result, err := col.UpdateOne(ctx, filter, bson.M{"$set": updates})
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "No se pudo actualizar la task"})
    }
    if result.MatchedCount == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task no encontrada o no autorizada"})
    }
    return c.JSON(fiber.Map{"message": "Task actualizada exitosamente"})
}

// DeleteTask elimina una task (si pertenece al usuario)
func DeleteTask(c *fiber.Ctx) error {
    idParam := c.Params("id")
    taskID, err := primitive.ObjectIDFromHex(idParam)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
    }
    userIDHex := c.Locals("userID").(string)
    userObjID, _ := primitive.ObjectIDFromHex(userIDHex)

    col := getCollectionTasks()
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    filter := bson.M{"_id": taskID, "usuario_id": userObjID}
    result, err := col.DeleteOne(ctx, filter)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "No se pudo eliminar la task"})
    }
    if result.DeletedCount == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task no encontrada o no autorizada"})
    }
    return c.JSON(fiber.Map{"message": "Task eliminada exitosamente"})
}

