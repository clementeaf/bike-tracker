# Bike Tracker API 🚴‍♂️

**Bike Tracker API** es un servicio backend escrito en Go, diseñado para gestionar viajes en bicicleta. Esta aplicación está desarrollada con las mejores prácticas, con un enfoque modular y escalable, utilizando el ecosistema Go para manejar datos, autenticación y la arquitectura de la aplicación.  

---

## **Características de la Aplicación**
- **Gestión de usuarios:**
  - Registro de nuevos usuarios.
  - Inicio de sesión con autenticación basada en email y contraseña.
  - Ingreso de saldo a la wallet del usuario.
  - Consulta de los datos actuales del usuario autenticado.
  - Eliminación de la cuenta del usuario autenticado.
- **Gestión de la wallet:**
  - Una wallet única por usuario.
  - Añadir transacciones (crédito o débito) a la wallet.
  - Obtener el saldo actual de la wallet del usuario autenticado.
  - Historial detallado de transacciones.
- **Gestión de viajes:**
  - Crear un nuevo viaje con origen, destino y estado.
  - Consultar la lista completa de viajes.
  - Obtener los detalles de un viaje específico por su ID.
- **Arquitectura modular y profesional:**
  - Separación de responsabilidades por paquetes (api, user, wallet, ride, etc.).
  - Mock de base de datos en memoria, preparado para integrar PostgreSQL.
- **Middlewares globales:**
  - Logging estructurado para todas las solicitudes.
  - Autenticación (base para futuros tokens JWT).
- **Estructura limpia y profesional:**
  - Diseño modular para facilitar la escalabilidad y el mantenimiento.
  - Mock de base de datos en memoria para desarrollo inicial.
- **Preparada para integración con PostgreSQL.**
- **Sistema de logs estructurados:**
  - Registro detallado de errores y eventos clave.
  - Fácil seguimiento del estado de la aplicación.
- **Pruebas unitarias planificadas para asegurar la calidad del código.**
  - Código estructurado para facilitar las pruebas unitarias y de integración.

---

## **Dependencias Instaladas**
| Dependencia                         | Descripción                                             |
|-------------------------------------|---------------------------------------------------------|
| `github.com/gin-gonic/gin`          | Framework web para construir APIs rápidas y livianas.   |
| `gorm.io/gorm`                      | ORM para el manejo de bases de datos.                   |
| `gorm.io/driver/postgres`           | Driver para conexión con PostgreSQL.                    |
| `github.com/dgrijalva/jwt-go`       | Manejo de autenticación basada en JWT.                  |
| `github.com/sirupsen/logrus`        | Logging estructurado y escalable.                       |
| `github.com/joho/godotenv`          | Carga de variables de entorno desde un archivo `.env`.  |

---

## **Estructura del Proyecto**

bike-tracker/
│── cmd/                     # Punto de entrada principal del proyecto
│   ├── main.go              # Configuración del servidor y registro de rutas
│
│── internal/                # Lógica de negocio y módulos internos
│   ├── api/                 # Configuración global de la API y middlewares
│   │   ├── api.go           # Inicialización principal de la API
│   │   ├── router.go        # Registro de rutas globales
│   │   ├── middleware.go    # Middlewares globales (autenticación, logging, etc.)
│   ├── ride/                # Módulo de viajes (CRUD de viajes)
│   │   ├── dto.go           # DTOs para transferir datos relacionados con viajes
│   │   ├── model.go         # Modelo de datos para viajes
│   │   ├── service.go       # Lógica de negocio para viajes
│   │   ├── handlers.go      # Controladores HTTP para viajes
│   │   ├── routes.go        # Registro de rutas específicas de viajes
│   ├── user/                # Módulo de usuarios (registro, login, datos personales)
│   │   ├── dto.go           # DTOs para transferir datos relacionados con usuarios
│   │   ├── model.go         # Modelo de datos para usuarios
│   │   ├── service.go       # Lógica de negocio para usuarios
│   │   ├── handlers.go      # Controladores HTTP para usuarios
│   │   ├── routes.go        # Registro de rutas específicas de usuarios
│   ├── wallet/              # Módulo de wallets
│   │   ├── model.go         # Modelos de datos para wallets
│   │   ├── service.go       # Lógica de negocio para wallets
│   │   ├── handlers.go      # Controladores HTTP para wallets
│   │   ├── routes.go        # Registro de rutas específicas de wallets
│
│── pkg/                     # Código reutilizable y configuraciones
│   ├── config/              # Manejo de variables de entorno
│   │   ├── config.go        # Lógica para cargar variables desde .env
│   ├── database/            # Configuración de la base de datos
│   │   ├── database.go      # Conexión con PostgreSQL (preparado)
│   ├── logger/              # Logger estructurado
│   │   ├── logger.go        # Configuración del sistema de logging
│   ├── mock/                # Mock de datos en memoria
│   │   ├── mock_user.go     # Simulación de usuarios en memoria
│   │   ├── mock_wallet.go   # Simulación de wallets en memoria
│   │   ├── mock_ride.go     # Simulación de viajes en memoria
│
│── .env                     # Archivo de configuración de entorno (excluido en .gitignore)
│── go.mod                   # Archivo de dependencias del proyecto
│── go.sum                   # Suma de verificación de dependencias
│── README.md                # Documentación del proyecto


## Rutas Disponibles

Método 	 | Endpoint	                  | Descripción
GET	     | /	                        | Verifica el estado del servidor
------------------------------------------------------------------------------------
GET	     | /rides	                    | Obtiene todos los viajes registrados.
POST	   | /rides/start	              | Crea un nuevo viaje.
POST	   | /rides/end	                | Finaliza un viaje.
GET	     | /rides/{id}	              | Obtiene un viaje específico por ID.
-------------------------------------------------------------------------------------
POST     | /users/register            | Registra un nuevo usuario.
POST     | /users/login               | Inicia sesión con email y contraseña.
POST     | /users/wallet              | Agrega saldo a la wallet de un usuario.
GET      | /users/me                  | Obtiene la información actual del usuario autenticado.
DELETE   | /users/me/delete           | Elimina la cuenta del usuario autenticado.
-------------------------------------------------------------------------------------
POST	   | /wallet/transactions/add	  | Añade una transacción a la wallet del usuario.
GET	     | /wallet/transactions	      | Obtiene el saldo actual de la wallet
GET	     | /wallet	                  | Obtiene el estado actual de la wallet.
-------------------------------------------------------------------------------------
POST     | /bikes/register            | Genera una nueva bicicleta
GET      | /bikes/available           | Obtiene arreglo de bicicletas disponibles
PUT      | /bikes/status              | Modifica el status de una bicicleta


## Cómo Ejecutar el Proyecto
Requisitos
 - Go (versión 1.20 o superior).
 - (Opcional) Docker para la base de datos PostgreSQL.
Instalación
 - Clona este repositorio:
    git clone https://github.com/clementeaf/bike-tracker.git
    cd bike-tracker
Instala las dependencias:
    go mod tidy
Crea un archivo .env en la raíz del proyecto con la siguiente configuración:
    DB_HOST=localhost
    DB_USER=postgres
    DB_PASSWORD=secret
    DB_NAME=biketracker
    DB_PORT=5432
    PORT=8080

## Ejecutar el Proyecto
Ejecuta el servidor:
    go run cmd/main.go
El servidor estará disponible en:
    http://localhost:8080


## Prueba de Rutas

SignUp POST /users/register
`  
{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "123456"
}
`
----------------------------------------

LogIn POST /users/login
`  
{
    "email": "john@example.com",
    "password": "123456"
}
`
----------------------------------------

Get user info /users/me
Headers: Authorization: Bearer Token

----------------------------------------

Delete User /users/me/delete
Headers: X-User-ID
         Authorization: Bearer Token

----------------------------------------

Post add found /wallet/transactions/add
Headers: Authorization: Bearer Token
Body:
`  
{
    "wallet_id": "67a08d56c65ef3f2f8c483a8",
    "user_id": "67a08d56c65ef3f2f8c483a7",
    "amount": 100,
    "type": "credit"
}
`

----------------------------------------

GET Wallet transactions /wallet/transactions
Headers: X-User-ID
         Authorization: Bearer Token

----------------------------------------

GET Wallet balance /wallet
Headers: X-User-ID
         Authorization: Bearer Token

----------------------------------------

GET Rides /rides

----------------------------------------

POST initiatie ride /rides/start
Headers: Authorization: Bearer token
`  
{
    "bike_id": "67a08dfac65ef3f2f8c483aa",
    "start_coords": [
        -33.4489,
        -70.6693
    ]
}
`
----------------------------------------
POST end ride /rides/end
Headers: Authorization: Bearer token
`  
{
    "ride_id": "67a0c5e4b417fddc74bf0ae0",
    "end_coords": [
        -33.4489,
        -70.6693
    ]
}
`
----------------------------------------