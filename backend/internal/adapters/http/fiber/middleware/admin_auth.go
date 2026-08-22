package middleware

import (
	"github.com/gofiber/fiber/v2"

	"polomnik/internal/adapters/http/fiber/dto"
	"polomnik/internal/application"
	"polomnik/internal/domain"
)

const (
	LocalsManagementPrincipal = "managementPrincipal"
)

func AdminAuth(adminRoles *application.AdminRoleService, adminToken string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if adminRoles == nil && adminToken == "" {
			return c.Status(fiber.StatusServiceUnavailable).JSON(dto.ErrorEnvelope{
				Error: dto.ErrorBody{
					Code:    "SERVICE_UNAVAILABLE",
					Message: "Management API не настроен",
				},
			})
		}

		var principal application.ManagementPrincipal
		var err error
		if adminRoles != nil {
			principal, err = adminRoles.AuthenticateHeader(c.Get("X-Admin-Token"), c.Get("X-Admin-Session"))
		} else if c.Get("X-Admin-Token") == adminToken && adminToken != "" {
			principal = application.ManagementPrincipal{
				FullAdmin:   true,
				Permissions: domain.AllPermissions(),
			}
		} else {
			err = domain.ErrInvalidCredentials
		}
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorEnvelope{
				Error: dto.ErrorBody{
					Code:    "UNAUTHORIZED",
					Message: "Неверный токен администратора",
				},
			})
		}

		c.Locals(LocalsManagementPrincipal, principal)
		return c.Next()
	}
}

func RequirePermission(adminRoles *application.AdminRoleService, perm domain.Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		principal, ok := c.Locals(LocalsManagementPrincipal).(application.ManagementPrincipal)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorEnvelope{
				Error: dto.ErrorBody{Code: "UNAUTHORIZED", Message: "Нужна авторизация администратора"},
			})
		}
		if adminRoles != nil && !adminRoles.HasPermission(principal, perm) {
			return c.Status(fiber.StatusForbidden).JSON(dto.ErrorEnvelope{
				Error: dto.ErrorBody{Code: "FORBIDDEN", Message: "Недостаточно прав"},
			})
		}
		if adminRoles == nil && !principal.FullAdmin {
			return c.Status(fiber.StatusForbidden).JSON(dto.ErrorEnvelope{
				Error: dto.ErrorBody{Code: "FORBIDDEN", Message: "Недостаточно прав"},
			})
		}
		return c.Next()
	}
}

func ManagementPrincipalFromCtx(c *fiber.Ctx) (application.ManagementPrincipal, bool) {
	principal, ok := c.Locals(LocalsManagementPrincipal).(application.ManagementPrincipal)
	return principal, ok
}
