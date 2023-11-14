package handlers

import (
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmcerez0/fiber-demo/models"
	"github.com/jmcerez0/fiber-demo/utils"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func SignUp(c *fiber.Ctx) error {
	type User struct {
		FirstName string `json:"first_name" xml:"first_name" form:"first_name" validate:"required"`
		LastName  string `json:"last_name" xml:"last_name" form:"last_name" validate:"required"`
		Email     string `json:"email" xml:"email" form:"email" validate:"required,email"`
		Password  string `json:"password" xml:"password" form:"password" validate:"required,min=8"`
	}

	u := new(User)

	if err := c.BodyParser(u); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(u); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user := models.User{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  string(hash),
	}

	result := utils.DB.Create(&user)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "Error 1062 (23000)") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": result.Error.Error(),
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": result.Error.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User created successfully.",
	})
}

func SignIn(c *fiber.Ctx) error {
	type User struct {
		Email    string `json:"email" xml:"email" form:"email" validate:"required,email"`
		Password string `json:"password" xml:"password" form:"password" validate:"required"`
	}

	u := new(User)

	if err := c.BodyParser(u); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(u); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var user models.User
	utils.DB.First(&user, "email = ?", u.Email)

	if user.ID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Incorrect email or password.",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Incorrect email or password.",
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"name": user.FirstName + " " + user.LastName,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 24 * 3).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Incorrect email or password.",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 3),
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": tokenString,
	})
}

func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User
	utils.DB.Find(&users)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"users": users,
	})
}
