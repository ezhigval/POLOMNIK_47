package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"palomnik/internal/adapters/http/fiber/dto"
	appmiddleware "palomnik/internal/adapters/http/fiber/middleware"
	"palomnik/internal/domain"
)

func (h *Handler) ManagementLogin(c *fiber.Ctx) error {
	if h.adminRoles == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "SERVICE_UNAVAILABLE", Message: "Management API не настроен"})
	}
	var req dto.ManagementLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"})
	}
	result, err := h.adminRoles.Login(c.Context(), req.Role, req.Password)
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.ManagementLoginResponse]{Data: dto.ToManagementLogin(result)})
}

func (h *Handler) ManagementSession(c *fiber.Ctx) error {
	principal, ok := appmiddleware.ManagementPrincipalFromCtx(c)
	if !ok {
		return writeAppError(c, &AppError{Status: 401, Code: "UNAUTHORIZED", Message: "Нужна авторизация администратора"})
	}
	return c.JSON(dto.DataEnvelope[dto.ManagementSessionResponse]{Data: dto.ToManagementSession(principal)})
}

func (h *Handler) ManagementGetNotificationSettings(c *fiber.Ctx) error {
	svc := h.notificationSettings()
	if svc == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "SERVICE_UNAVAILABLE", Message: "Настройки уведомлений недоступны"})
	}
	view, err := svc.Settings(c.Context())
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.NotificationSettingsResponse]{Data: dto.ToNotificationSettings(view)})
}

func (h *Handler) ManagementUpdateNotificationSettings(c *fiber.Ctx) error {
	svc := h.notificationSettings()
	if svc == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "SERVICE_UNAVAILABLE", Message: "Настройки уведомлений недоступны"})
	}
	var req dto.NotificationSettingsUpsertRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"})
	}
	events, err := dto.ParseNotificationEvents(req)
	if err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные получатели"})
	}
	view, err := svc.UpdateRouting(c.Context(), events)
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.NotificationSettingsResponse]{Data: dto.ToNotificationSettings(view)})
}

func (h *Handler) ManagementGetSiteSettings(c *fiber.Ctx) error {
	if h.siteSettings == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "SERVICE_UNAVAILABLE", Message: "Настройки сайта недоступны"})
	}
	settings, err := h.siteSettings.Settings(c.Context())
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.SiteSettingsResponse]{Data: dto.ToSiteSettings(settings)})
}

func (h *Handler) ManagementUpdateSiteSettings(c *fiber.Ctx) error {
	if h.siteSettings == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "SERVICE_UNAVAILABLE", Message: "Настройки сайта недоступны"})
	}
	var req dto.SiteSettingsUpsertRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"})
	}
	input, err := dto.SiteSettingsFromRequest(req)
	if err != nil {
		return respondError(c, err, MapError)
	}
	settings, err := h.siteSettings.Update(c.Context(), input)
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.SiteSettingsResponse]{Data: dto.ToSiteSettings(settings)})
}

func (h *Handler) GetPublicSiteSettings(c *fiber.Ctx) error {
	if h.siteSettings == nil {
		return c.JSON(dto.DataEnvelope[dto.SiteSettingsResponse]{Data: dto.SiteSettingsResponse{}})
	}
	settings, err := h.siteSettings.Public(c.Context())
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.SiteSettingsResponse]{Data: dto.ToSiteSettings(settings)})
}

func (h *Handler) ManagementListRoles(c *fiber.Ctx) error {
	if h.adminRoles == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "SERVICE_UNAVAILABLE", Message: "Роли недоступны"})
	}
	roles, err := h.adminRoles.ListRoles(c.Context())
	if err != nil {
		return respondError(c, err, MapError)
	}
	out := make([]dto.AdminRoleResponse, 0, len(roles))
	for _, role := range roles {
		out = append(out, dto.ToAdminRole(role))
	}
	return c.JSON(dto.DataEnvelope[[]dto.AdminRoleResponse]{Data: out})
}

func (h *Handler) ManagementCreateRole(c *fiber.Ctx) error {
	if h.adminRoles == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "SERVICE_UNAVAILABLE", Message: "Роли недоступны"})
	}
	var req dto.AdminRoleCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"})
	}
	role, err := h.adminRoles.CreateRole(c.Context(), req.Name, req.Password, dto.PermissionsFromStrings(req.Permissions))
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.Status(fiber.StatusCreated).JSON(dto.DataEnvelope[dto.AdminRoleResponse]{Data: dto.ToAdminRole(role)})
}

func (h *Handler) ManagementUpdateRole(c *fiber.Ctx) error {
	if h.adminRoles == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "SERVICE_UNAVAILABLE", Message: "Роли недоступны"})
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный идентификатор роли"})
	}
	var req dto.AdminRoleUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"})
	}
	var perms []domain.Permission
	if req.Permissions != nil {
		perms = dto.PermissionsFromStrings(req.Permissions)
	}
	role, err := h.adminRoles.UpdateRole(c.Context(), id, perms, req.Password)
	if err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[dto.AdminRoleResponse]{Data: dto.ToAdminRole(role)})
}

func (h *Handler) ManagementDeleteRole(c *fiber.Ctx) error {
	if h.adminRoles == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "SERVICE_UNAVAILABLE", Message: "Роли недоступны"})
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный идентификатор роли"})
	}
	if err := h.adminRoles.DeleteRole(c.Context(), id); err != nil {
		return respondError(c, err, MapError)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) ManagementAssignRoleUser(c *fiber.Ctx) error {
	if h.adminRoles == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "SERVICE_UNAVAILABLE", Message: "Роли недоступны"})
	}
	roleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный идентификатор роли"})
	}
	var req dto.AdminRoleAssignRequest
	if err := c.BodyParser(&req); err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректные данные запроса"})
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Укажите UUID пользователя"})
	}
	if err := h.adminRoles.AssignUser(c.Context(), roleID, userID); err != nil {
		return respondError(c, err, MapError)
	}
	return c.JSON(dto.DataEnvelope[map[string]string]{Data: map[string]string{"status": "ok"}})
}

func (h *Handler) ManagementUnassignRoleUser(c *fiber.Ctx) error {
	if h.adminRoles == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "SERVICE_UNAVAILABLE", Message: "Роли недоступны"})
	}
	roleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный идентификатор роли"})
	}
	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный UUID пользователя"})
	}
	if err := h.adminRoles.UnassignUser(c.Context(), roleID, userID); err != nil {
		return respondError(c, err, MapError)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) ManagementListRoleAssignments(c *fiber.Ctx) error {
	if h.adminRoles == nil {
		return writeAppError(c, &AppError{Status: 503, Code: "SERVICE_UNAVAILABLE", Message: "Роли недоступны"})
	}
	roleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return writeAppError(c, &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный идентификатор роли"})
	}
	items, err := h.adminRoles.ListAssignments(c.Context(), roleID)
	if err != nil {
		return respondError(c, err, MapError)
	}
	out := make([]dto.AdminRoleAssignmentResponse, 0, len(items))
	for _, item := range items {
		out = append(out, dto.AdminRoleAssignmentResponse{
			RoleID:    item.RoleID.String(),
			UserID:    item.UserID.String(),
			CreatedAt: item.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return c.JSON(dto.DataEnvelope[[]dto.AdminRoleAssignmentResponse]{Data: out})
}
