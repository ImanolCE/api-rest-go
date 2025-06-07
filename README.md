# API REST en Go con Fiber y MongoDB

Este proyecto es una API REST  que se desarrolla con Go usando el framework Fiber y con la utilizacion de  MongoDB como base de datos y autenticación JWT.

## Características
- Gestión de usuarios (CRUD)
- Gestión de tareas (Tasks CRUD)
- Autenticación JWT (validez 10 minutos)
- Conexión con MongoDB
- Estructura que fomente buenas practicas

## Rutas principales
- `POST /api/users/register` - Registro de usuarios
- `POST /api/users/login` - Login de usuario (retorna JWT)
- `PUT /api/users/:id`- Actualizar usuario 
- `DELETE /api/users/:id`- Borrar un usuario 
- `GET /api/tasks/user/:id`- Tareas de un usuario


## Estructura del proyecto

    📁config
    └── db.go
    📁handlers
    └── task_handler.go
    └── user_handler.go
    📁middleware
    └── jwt_middleware.go
    📁models
    └── task.go
    └── user.go
    📁routes
    └── routes.go
    📁test
    📁utils
    └── jwt.go
    CHANGELOG.md
    go.mod
    go.sum
    main.go
    README.md

## Requisitos
- Go 1.21+
- Fiber v2.52.8 
- MongoDB de forma local o MongoAtlas 
- Thunder Client para pruebas o Postman 

Repositorio: (https://github.com/ImanolCE/api-rest-go)