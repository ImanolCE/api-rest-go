
package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

// coleccion del usuario
type User struct {
    ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    Nombre      string             `json:"nombre" bson:"nombre"`
    Apellidos    string             `json:"apellidos" bson:"apellidos"`
    Email         string             `json:"email" bson:"email"`
    Password        string             `json:"password" bson:"password"`
    FechaNacimiento  time.Time          `json:"fecha_nacimiento" bson:"fecha_nacimiento"`
    PreguntaSecreta  string             `json:"pregunta_secreta" bson:"pregunta_secreta"`
    RespuestaSecreta string             `json:"respuesta_secreta" bson:"respuesta_secreta"`
}
