package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"polomnik/internal/domain"
)

func (s *Store) GetSiteSettings(ctx context.Context) (domain.SiteSettings, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT site_name, full_name, tagline, description, region, departure_city,
       parent_org_name, parent_org_url, contact_phone, contact_phone_display,
       contact_email, mail_forward_to, updated_at
FROM site_settings
WHERE id = $1
`, domain.SiteSettingsID())
	return scanSiteSettings(row)
}

func (s *Store) UpsertSiteSettings(ctx context.Context, settings domain.SiteSettings) (domain.SiteSettings, error) {
	row := s.db.QueryRowContext(ctx, `
INSERT INTO site_settings (
    id, site_name, full_name, tagline, description, region, departure_city,
    parent_org_name, parent_org_url, contact_phone, contact_phone_display,
    contact_email, mail_forward_to, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10, $11,
    $12, $13, $14
)
ON CONFLICT (id) DO UPDATE SET
    site_name = EXCLUDED.site_name,
    full_name = EXCLUDED.full_name,
    tagline = EXCLUDED.tagline,
    description = EXCLUDED.description,
    region = EXCLUDED.region,
    departure_city = EXCLUDED.departure_city,
    parent_org_name = EXCLUDED.parent_org_name,
    parent_org_url = EXCLUDED.parent_org_url,
    contact_phone = EXCLUDED.contact_phone,
    contact_phone_display = EXCLUDED.contact_phone_display,
    contact_email = EXCLUDED.contact_email,
    mail_forward_to = EXCLUDED.mail_forward_to,
    updated_at = EXCLUDED.updated_at
RETURNING site_name, full_name, tagline, description, region, departure_city,
          parent_org_name, parent_org_url, contact_phone, contact_phone_display,
          contact_email, mail_forward_to, updated_at
`, domain.SiteSettingsID(),
		settings.SiteName,
		settings.FullName,
		settings.Tagline,
		settings.Description,
		settings.Region,
		settings.DepartureCity,
		settings.ParentOrgName,
		settings.ParentOrgURL,
		settings.ContactPhone,
		settings.ContactPhoneDisplay,
		settings.ContactEmail,
		domain.FormatMailForwardList(settings.MailForwardTo),
		settings.UpdatedAt,
	)
	return scanSiteSettings(row)
}

func scanSiteSettings(row scanner) (domain.SiteSettings, error) {
	var settings domain.SiteSettings
	var mailForwardRaw string
	err := row.Scan(
		&settings.SiteName,
		&settings.FullName,
		&settings.Tagline,
		&settings.Description,
		&settings.Region,
		&settings.DepartureCity,
		&settings.ParentOrgName,
		&settings.ParentOrgURL,
		&settings.ContactPhone,
		&settings.ContactPhoneDisplay,
		&settings.ContactEmail,
		&mailForwardRaw,
		&settings.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.SiteSettings{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.SiteSettings{}, fmt.Errorf("scan site settings: %w", err)
	}
	forwards, err := domain.ParseMailForwardList(mailForwardRaw)
	if err != nil {
		return domain.SiteSettings{}, err
	}
	settings.MailForwardTo = forwards
	return settings, nil
}
