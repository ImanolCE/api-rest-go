/* package main

import (

	"context"
    "fmt"
    "log"
    "net/http"

	"firebase.google.com/go"
    "firebase.google.com/go/db"
    "github.com/gorilla/mux"
    "google.golang.org/api/option"

	"api-rest-go/hanflers"
	"api-rest-go/middleware"

	"github.com/gofiber/fiber/v2"

)
 */


package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/tu_usuario/api-rest-go/config"
    "github.com/tu_usuario/api-rest-go/routes"
)

func main() {
    // 1. Conectar a BD
    config.ConnectDB()

    // 2. Crear instancia de Fiber
    app := fiber.New()

    // 3. Registrar rutas
    routes.Setup(app)

    // 4. Iniciar servidor en puerto 3000
    app.Listen(":3000")
}


