// models/task.go
package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

// coleccion de task
type Task struct {
    ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    Titulo    string             `json:"titulo" bson:"titulo"`
    Descripcion string           `json:"descripcion" bson:"descripcion"`
    FechaInicio  time.Time       `json:"fecha_inicio" bson:"fecha_inicio"`
    FechaFinal     time.Time       `json:"fecha_final" bson:"fecha_final"`
    UsuarioID    primitive.ObjectID `json:"usuario_id" bson:"usuario_id"`
}
