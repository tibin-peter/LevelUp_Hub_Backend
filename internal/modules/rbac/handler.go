package rbac

import "github.com/gofiber/fiber/v2"

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// POST /api/admin/permissions
func (h *Handler) CreatePermission(c *fiber.Ctx) error {
	var req PermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	if err := h.svc.AddNewPermission(req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Permission created"})
}

// PATCH /api/admin/roles/toggle
func (h *Handler) ToggleRolePermission(c *fiber.Ctx) error {
	var req AssignPermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	if err := h.svc.TogglePermission(req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Role updated successfully"})
}

func (h *Handler) CreateRole(c *fiber.Ctx) error {

	var req CreateRoleRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid body",
		})
	}

	if err := h.svc.CreateRole(req); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "role created",
	})
}

func (h *Handler) GetPermissionsByRole(c *fiber.Ctx) error {

	role := c.Params("role")

	data, err := h.svc.GetPermissionsByRole(role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}

func (h *Handler) ListRoles(c *fiber.Ctx) error {
	roles, err := h.svc.GetAllRoles()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": roles})
}

func (h *Handler) ListPermissions(c *fiber.Ctx) error {
	permissions, err := h.svc.GetAllPermissions()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": permissions})
}
