package handlers

import (
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jmcerez0/fiber-demo/models"
	"github.com/jmcerez0/fiber-demo/utils"
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

	hashedPassword, err := utils.HashPassword(u.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user := models.User{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  hashedPassword,
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

	if err := utils.ComparePassword(user.Password, u.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Incorrect email or password.",
		})
	}

	token, err := utils.GetToken(user)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Incorrect email or password.",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * 3),
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}

func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User
	utils.DB.Find(&users)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"users": users,
	})
}
