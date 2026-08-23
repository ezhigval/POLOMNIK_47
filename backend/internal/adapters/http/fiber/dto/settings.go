package dto

import (
	"palomnik/internal/application"
	"palomnik/internal/domain"
)

type NotificationRecipientDTO struct {
	Channel string `json:"channel"`
	Address string `json:"address"`
}

type NotificationRecipientStatusDTO struct {
	Channel string `json:"channel"`
	Address string `json:"address"`
	Event   string `json:"event"`
	Ready   bool   `json:"ready"`
	Status  string `json:"status"`
	Label   string `json:"label"`
}

type NotificationEventBlockDTO struct {
	Kind       string                           `json:"kind"`
	Title      string                           `json:"title"`
	Recipients []NotificationRecipientStatusDTO `json:"recipients"`
}

type NotificationChannelDTO struct {
	ID         string `json:"id"`
	Configured bool   `json:"configured"`
	Label      string `json:"label"`
}

type NotificationSettingsResponse struct {
	Channels []NotificationChannelDTO    `json:"channels"`
	Events   []NotificationEventBlockDTO `json:"events"`
}

type NotificationSettingsUpsertRequest struct {
	Events map[string][]NotificationRecipientDTO `json:"events"`
}

func ToNotificationSettings(view application.NotificationSettingsView) NotificationSettingsResponse {
	channels := make([]NotificationChannelDTO, 0, len(view.Channels))
	for _, ch := range view.Channels {
		channels = append(channels, NotificationChannelDTO{
			ID:         string(ch.ID),
			Configured: ch.Configured,
			Label:      ch.Label,
		})
	}
	events := make([]NotificationEventBlockDTO, 0, len(view.Events))
	for _, event := range view.Events {
		recipients := make([]NotificationRecipientStatusDTO, 0, len(event.Recipients))
		for _, item := range event.Recipients {
			recipients = append(recipients, NotificationRecipientStatusDTO{
				Channel: string(item.Channel.ID),
				Address: item.Address,
				Event:   string(item.Event),
				Ready:   item.Ready,
				Status:  item.Status,
				Label:   item.Channel.Label,
			})
		}
		events = append(events, NotificationEventBlockDTO{
			Kind:       string(event.Kind),
			Title:      event.Title,
			Recipients: recipients,
		})
	}
	return NotificationSettingsResponse{Channels: channels, Events: events}
}

func ParseNotificationEvents(req NotificationSettingsUpsertRequest) (map[domain.NotificationEventKind][]domain.NotificationRecipient, error) {
	out := make(map[domain.NotificationEventKind][]domain.NotificationRecipient)
	for kindRaw, list := range req.Events {
		kind := domain.NotificationEventKind(kindRaw)
		recipients := make([]domain.NotificationRecipient, 0, len(list))
		for _, item := range list {
			recipients = append(recipients, domain.NotificationRecipient{
				Channel: domain.NotificationChannel(item.Channel),
				Address: item.Address,
			})
		}
		out[kind] = recipients
	}
	return out, nil
}

type SiteSettingsResponse struct {
	SiteName            string `json:"site_name"`
	FullName            string `json:"full_name"`
	Tagline             string `json:"tagline"`
	Description         string `json:"description"`
	Region              string `json:"region"`
	DepartureCity       string `json:"departure_city"`
	ParentOrgName       string `json:"parent_org_name"`
	ParentOrgURL        string `json:"parent_org_url"`
	ContactPhone        string `json:"contact_phone"`
	ContactPhoneDisplay string `json:"contact_phone_display"`
	ContactEmail        string `json:"contact_email"`
	MailForwardTo       string `json:"mail_forward_to"`
}

type SiteSettingsUpsertRequest struct {
	SiteName            string `json:"site_name"`
	FullName            string `json:"full_name"`
	Tagline             string `json:"tagline"`
	Description         string `json:"description"`
	Region              string `json:"region"`
	DepartureCity       string `json:"departure_city"`
	ParentOrgName       string `json:"parent_org_name"`
	ParentOrgURL        string `json:"parent_org_url"`
	ContactPhone        string `json:"contact_phone"`
	ContactPhoneDisplay string `json:"contact_phone_display"`
	ContactEmail        string `json:"contact_email"`
	MailForwardTo       string `json:"mail_forward_to"`
}

func ToSiteSettings(settings domain.SiteSettings) SiteSettingsResponse {
	return SiteSettingsResponse{
		SiteName:            settings.SiteName,
		FullName:            settings.FullName,
		Tagline:             settings.Tagline,
		Description:         settings.Description,
		Region:              settings.Region,
		DepartureCity:       settings.DepartureCity,
		ParentOrgName:       settings.ParentOrgName,
		ParentOrgURL:        settings.ParentOrgURL,
		ContactPhone:        settings.ContactPhone,
		ContactPhoneDisplay: settings.ContactPhoneDisplay,
		ContactEmail:        settings.ContactEmail,
		MailForwardTo:       domain.FormatMailForwardList(settings.MailForwardTo),
	}
}

func SiteSettingsFromRequest(req SiteSettingsUpsertRequest) (domain.SiteSettings, error) {
	forwards, err := domain.ParseMailForwardList(req.MailForwardTo)
	if err != nil {
		return domain.SiteSettings{}, err
	}
	return domain.SiteSettings{
		SiteName:            req.SiteName,
		FullName:            req.FullName,
		Tagline:             req.Tagline,
		Description:         req.Description,
		Region:              req.Region,
		DepartureCity:       req.DepartureCity,
		ParentOrgName:       req.ParentOrgName,
		ParentOrgURL:        req.ParentOrgURL,
		ContactPhone:        req.ContactPhone,
		ContactPhoneDisplay: req.ContactPhoneDisplay,
		ContactEmail:        req.ContactEmail,
		MailForwardTo:       forwards,
	}, nil
}

type ManagementLoginRequest struct {
	Role     string `json:"role"`
	Password string `json:"password"`
}

type ManagementLoginResponse struct {
	Token       string   `json:"token"`
	FullAdmin   bool     `json:"full_admin"`
	RoleID      string   `json:"role_id,omitempty"`
	RoleName    string   `json:"role_name,omitempty"`
	Permissions []string `json:"permissions"`
}

type ManagementSessionResponse struct {
	FullAdmin   bool     `json:"full_admin"`
	RoleID      string   `json:"role_id,omitempty"`
	RoleName    string   `json:"role_name,omitempty"`
	Permissions []string `json:"permissions"`
}

type AdminRoleResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

type AdminRoleTemplateResponse struct {
	ID          string   `json:"id"`
	Label       string   `json:"label"`
	RoleName    string   `json:"role_name"`
	Permissions []string `json:"permissions"`
}

type AdminRoleCreateRequest struct {
	Name        string   `json:"name"`
	Password    string   `json:"password"`
	Permissions []string `json:"permissions"`
}

type AdminRoleUpdateRequest struct {
	Password    string   `json:"password"`
	Permissions []string `json:"permissions"`
}

type AdminRoleAssignRequest struct {
	UserID string `json:"user_id"`
}

type AdminRoleAssignmentResponse struct {
	RoleID    string `json:"role_id"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

func ToAdminRole(role domain.AdminRole) AdminRoleResponse {
	perms := make([]string, 0, len(role.Permissions))
	for _, p := range role.Permissions {
		perms = append(perms, string(p))
	}
	return AdminRoleResponse{
		ID:          role.ID.String(),
		Name:        role.Name,
		Permissions: perms,
		CreatedAt:   role.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   role.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
}

func ToAdminRoleTemplate(template domain.RoleTemplate) AdminRoleTemplateResponse {
	perms := make([]string, 0, len(template.Permissions))
	for _, p := range template.Permissions {
		perms = append(perms, string(p))
	}
	return AdminRoleTemplateResponse{
		ID:          template.ID,
		Label:       template.Label,
		RoleName:    template.RoleName,
		Permissions: perms,
	}
}

func PermissionsFromStrings(raw []string) []domain.Permission {
	out := make([]domain.Permission, 0, len(raw))
	for _, item := range raw {
		out = append(out, domain.Permission(item))
	}
	return out
}

func ToManagementLogin(result application.ManagementLoginResult) ManagementLoginResponse {
	perms := make([]string, 0, len(result.Permissions))
	for _, p := range result.Permissions {
		perms = append(perms, string(p))
	}
	resp := ManagementLoginResponse{
		Token:       result.Token,
		FullAdmin:   result.FullAdmin,
		RoleName:    result.RoleName,
		Permissions: perms,
	}
	if result.RoleID.String() != "00000000-0000-0000-0000-000000000000" {
		resp.RoleID = result.RoleID.String()
	}
	return resp
}

func ToManagementSession(p application.ManagementPrincipal) ManagementSessionResponse {
	perms := make([]string, 0, len(p.Permissions))
	for _, item := range p.Permissions {
		perms = append(perms, string(item))
	}
	resp := ManagementSessionResponse{
		FullAdmin:   p.FullAdmin,
		RoleName:    p.RoleName,
		Permissions: perms,
	}
	if p.RoleID.String() != "00000000-0000-0000-0000-000000000000" {
		resp.RoleID = p.RoleID.String()
	}
	return resp
}
