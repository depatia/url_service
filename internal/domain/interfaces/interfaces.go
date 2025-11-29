package interfaces

import "context"

type IRepository interface {
	AddLinks(newLinks map[string]string) int
	GetLinksByLinksID(linksID int) map[string]string
}

type IHealthCheckerService interface {
	CheckSitesAvailability(ctx context.Context, links []string) (map[string]string, int)
	GetSitesAvailabilityReportByLinksID(ctx context.Context, linksID int) map[string]string
}

type IPdfService interface {
	GeneratePdf(ctx context.Context, text string) ([]byte, error)
}
