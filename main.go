package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/harold1194/todo-golang/models"
	"github.com/harold1194/todo-golang/storage"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type User struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type Attendance struct {
	FullName  *string `json:"fullname"`
	Address   *string `json:"address"`
	Degree    *string `json:"degree"`
	Year      *string `json:"year"`
	Block     *string `json:"block"`
	Subject   *string `json:"subject"`
	Date      *string `json:"date"`
	StartTime *string `json:"startTime"`
	EndTime   *string `json:"endTime"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateAttendance(context *fiber.Ctx) error {
	attendance := Attendance{}

	err := context.BodyParser(&attendance)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(attendance).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create attendance"})
		return nil
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "attendance has added"})
	return nil
}

func (r *Repository) DeleteAttendance(context *fiber.Ctx) error {
	attendaceModel := models.Attendance{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "id cannot be empty"})
		return nil
	}

	err := r.DB.Delete(attendaceModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not delete attendance"})
		return nil
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Attendance Deleted."})
	return nil
}

func (r *Repository) GetAttendance(context *fiber.Ctx) error {
	attendanceModel := &[]models.Attendance{}

	err := r.DB.Find(attendanceModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not get attendance data",
		})
		return nil
	}
	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "student successfully fetch", "data": attendanceModel},
	)
	return nil
}

func (r *Repository) GetAttendanceByID(context *fiber.Ctx) error {
	id := context.Params("id")
	attendaceModel := &models.Attendance{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot found",
		})
		return nil
	}

	fmt.Println("the ID is", id)

	err := r.DB.Where("id = ?", id).First(attendaceModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not get attendance",
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "attendance id fetched successfully",
		"data":    attendaceModel,
	})
	return nil
}

func (r *Repository) Register(context *fiber.Ctx) error {
	user := User{}

	err := context.BodyParser(&user)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&user).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not register the user"})
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Successfully Registered",
	})
	return nil
}

func (r *Repository) GetUsers(context *fiber.Ctx) error {
	userModel := &[]models.User{}

	err := r.DB.Find(userModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get users data"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user fetch successfully",
		"data":    userModel,
	})
	return nil
}

func (r *Repository) GetUserByID(context *fiber.Ctx) error {
	id := context.Params("id")
	userModel := &models.User{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be found",
		})
		return nil
	}

	fmt.Println("the ID is", id)

	err := r.DB.Where("id = ?", id).First(userModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not get user",
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user id successfully fetched",
		"data":    userModel,
	})
	return nil
}

func (r *Repository) DeleteUser(context *fiber.Ctx) error {
	userModel := models.User{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "id cannot be empty"})
		return nil
	}

	err := r.DB.Delete(userModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not delete user"})
		return nil
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user delete successfully",
	})
	return nil
}

// GenerateJWT generates a JSON Web Token (JWT) with the provided user ID.
func GenerateJWT(userID uint) (string, error) {
	// Define the claims for the JWT
	claims := jwt.MapClaims{
		"user_id": userID,
		// Add other desired claims...
	}

	// Create a new JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate the token string
	jwtSecret := []byte("your_secret_key") // Replace with your secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (r *Repository) Login(context *fiber.Ctx) error {
	user := User{}

	err := context.BodyParser(&user)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	// Perform authentication logic
	// Example:
	userModel := &models.User{}
	err = r.DB.Where("email = ? AND password = ?", user.Email, user.Password).First(userModel).Error
	if err != nil {
		context.Status(http.StatusUnauthorized).JSON(
			&fiber.Map{"message": "invalid email or password"})
		return err
	}

	// Assuming authentication is successful, generate a token
	// Example:
	token, err := GenerateJWT(userModel.ID)
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "failed to generate token"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "login successful",
		"token":   token,
	})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/register", r.Register)
	api.Post("/login", r.Login)
	api.Delete("delete_user/:id", r.DeleteUser)
	api.Get("/get_users/:id", r.GetUserByID)
	api.Get("/users", r.GetUsers)

	// attendance
	api.Post("/create_attendance", r.CreateAttendance)
	api.Delete("delete_attendance/:id", r.DeleteAttendance)
	api.Get("/get_attendance/:id", r.GetAttendanceByID)
	api.Get("/attendance", r.GetAttendance)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not load the database")
	}
	err = models.MigrateUsers(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	// attendance
	if err != nil {
		log.Fatal("could not load the database")
	}
	err = models.MigrateAttendance(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
