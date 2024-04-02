package controllers

import (
	"GDSC/initializers"
	"GDSC/models"
	"time"

	"github.com/gofiber/fiber/v2"
	// "github.com/google/uuid"
)

func GetMe(c *fiber.Ctx) error {

	user, ok := c.Locals("user").(models.UserResponse)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve user from context",
		})
	}

	
	var documents []models.Task

	
	if err := initializers.DB.Where("user_id = ? OR user_id IS ''", user.ID).Find(&documents).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve documents for the user",
		})
	}

	// Prepare response
	response := fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"documents": documents,
			"user":      user,
		},
	}

	
	return c.Status(fiber.StatusOK).JSON(response)
}




func CreateTask(c *fiber.Ctx) error {
	
	user, ok := c.Locals("user").(models.UserResponse)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve user from context",
		})
	}

	
	if user.Role != "Admin" && user.Role != "Manager" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Only admin or manager can create tasks",
		})
	}

	
	var payload models.TaskInput
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	
	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": errors})
	}

	// Create new task
	newTask := models.Task{
		Title:       payload.Title,
		Description: payload.Description,
		Status:      payload.Status,
	
	}

	// Set UserID to nil if user is admin or manager
	if user.Role == "Admin" || user.Role == "Manager" {
		newTask.UserID = ""
	} else {
		// Otherwise, set UserID to user's ID
		newTask.UserID = user.ID.String()
	}

	// Save task to database
	result := initializers.DB.Create(&newTask)
	if result.Error != nil {
		return result.Error
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   newTask,
	})
}
func UpdateTask(c *fiber.Ctx) error {
	var payload *models.TaskInput
	Id := c.Params("taskid")

	user, ok := c.Locals("user").(models.UserResponse)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve user from context",
		})
	}

		
	if user.Role != "Admin" && user.Role != "Manager" {
		
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Only admin or manager users can delete tasks",
		})
	}


	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Validate payload
	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": errors})
	}

	// Check if the task ID is provided
	if Id != "" {
		// If task ID is provided, it's an update operation
		var existingTask models.Task
		result := initializers.DB.Where("id = ?", Id).First(&existingTask)
		if result.Error != nil {
			// If the task is not found, return an error
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Task not found",
			})
		}

		// Update existing task fields
		existingTask.Title = payload.Title
		existingTask.Description = payload.Description
		existingTask.Status = payload.Status // Update status if needed
		existingTask.Deadline = time.Now().AddDate(0, 0, 7)
		// Save the updated task
		result = initializers.DB.Save(&existingTask)
		if result.Error != nil {
			return result.Error
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Task updated successfully",
			"data":    existingTask,
		})
	}

	// If task ID is not provided, it's a create operation
	demoTask := models.Task{
		Title:       payload.Title,
		Description: payload.Description,
		UserID:      user.ID.String(),
		Status:      payload.Status,              // Set initial status as "To Do"
		Deadline:    time.Now().AddDate(0, 0, 7), // Use provided deadline
		CreatedAt:   time.Now(),                  // Set creation timestamp
		UpdatedAt:   time.Now(),                  // Set update timestamp
	}

	// Create the task
	result := initializers.DB.Create(&demoTask)
	if result.Error != nil {
		return result.Error
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Task created successfully",
		"data":    demoTask,
	})
}

func DeleteTask(c *fiber.Ctx) error {

	user, ok := c.Locals("user").(models.UserResponse)
	if !ok {
		// If user information cannot be retrieved from the context, return an error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve user from context",
		})
	}

	
	if user.Role != "Admin" && user.Role != "Manager" {
		
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Only admin or manager users can delete tasks",
		})
	}


	taskId := c.Params("id")

	if taskId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Task ID is required",
		})
	}

	var existingTask models.Task
	result := initializers.DB.Where("id = ?", taskId).First(&existingTask)
	if result.Error != nil {

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Task not found",
		})
	}

	if err := initializers.DB.Delete(&existingTask).Error; err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete task",
		})
	}

	//
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Task deleted successfully",
		"data":    existingTask,
	})
}


func UpdateUserRole(c *fiber.Ctx) error {
	
	user, ok := c.Locals("user").(models.UserResponse)

	if !ok {
	
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve user from context",
		})
	}

	
	if user.Role != "Admin" {
		
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Only admin users can update user roles",
		})
	}


	var payload struct {
		Email string `json:"email" validate:"required,email"`
		Role  string `json:"role" validate:"required"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}


	// Find the user by email
	var userToUpdate models.User
	result := initializers.DB.Where("email = ?", payload.Email).First(&userToUpdate)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "User not found"})
	}

	
	userToUpdate.Role = payload.Role
	result = initializers.DB.Save(&userToUpdate)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to update user role"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "User role updated successfully"})
}
