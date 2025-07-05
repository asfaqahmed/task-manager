package handlers

import (
	"strconv"
	"time"

	"task-manager/database"
	"task-manager/models"
	"task-manager/queue"

	"github.com/gofiber/fiber/v2"
)

func GetTasks(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var tasks []models.Task
	if err := database.DB.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve tasks"})
	}
	return c.JSON(tasks)
}

func CreateTask(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var task models.Task

	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	task.UserID = userID

	if err := database.DB.Create(&task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create task"})
	}

	// Schedule email reminder
	reminderTime, err := time.Parse(time.RFC3339, task.Reminder)
	if err == nil && reminderTime.After(time.Now()) {
		payload := queue.EmailTaskPayload{
			Email:   "", // TODO: fetch user email if needed
			Message: "Reminder for task: " + task.Title,
		}
		queue.EnqueueEmailTask(payload, reminderTime)
	}

	return c.JSON(task)
}

func UpdateTask(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&task).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}

	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := database.DB.Save(&task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update task"})
	}

	return c.JSON(task)
}

func DeleteTask(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Task{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete task"})
	}

	return c.JSON(fiber.Map{"message": "Task deleted"})
}
