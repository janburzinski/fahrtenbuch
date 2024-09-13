package handlers

import "github.com/gofiber/fiber/v2"

type OrganisationHandler struct{}

func NewOrganisationHandler() *OrganisationHandler {
	return &OrganisationHandler{}
}

func (oh *OrganisationHandler) CreateOrganisation(c *fiber.Ctx) {

}
