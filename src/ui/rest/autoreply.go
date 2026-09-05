package rest

import (
	"strconv"

	domainChatStorage "github.com/aldinokemal/go-whatsapp-web-multidevice/domains/chatstorage"
	"github.com/aldinokemal/go-whatsapp-web-multidevice/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type AutoReply struct {
	Repo domainChatStorage.IChatStorageRepository
}

func InitRestAutoReply(app fiber.Router, repo domainChatStorage.IChatStorageRepository) AutoReply {
	rest := AutoReply{Repo: repo}

	// AutoReply endpoints
	app.Get("/api/v1/autoreply", rest.ListAutoReplies)
	app.Post("/api/v1/autoreply", rest.StoreAutoReply)
	app.Delete("/api/v1/autoreply/:id", rest.DeleteAutoReply)

	return rest
}

func (controller *AutoReply) ListAutoReplies(c *fiber.Ctx) error {
	replies, err := controller.Repo.GetAutoReplies()
	utils.PanicIfNeeded(err)

	return c.JSON(utils.ResponseData{
		Status:  200,
		Code:    "SUCCESS",
		Message: "Success get auto replies",
		Results: replies,
	})
}

func (controller *AutoReply) StoreAutoReply(c *fiber.Ctx) error {
	var request domainChatStorage.AutoReply
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(utils.ResponseData{
			Status:  400,
			Code:    "BAD_REQUEST",
			Message: "Invalid request body",
		})
	}

	err := controller.Repo.StoreAutoReply(&request)
	utils.PanicIfNeeded(err)

	return c.JSON(utils.ResponseData{
		Status:  200,
		Code:    "SUCCESS",
		Message: "Success store auto reply",
		Results: request,
	})
}

func (controller *AutoReply) DeleteAutoReply(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(utils.ResponseData{
			Status:  400,
			Code:    "BAD_REQUEST",
			Message: "Invalid ID",
		})
	}

	err = controller.Repo.DeleteAutoReply(id)
	utils.PanicIfNeeded(err)

	return c.JSON(utils.ResponseData{
		Status:  200,
		Code:    "SUCCESS",
		Message: "Success delete auto reply",
	})
}
