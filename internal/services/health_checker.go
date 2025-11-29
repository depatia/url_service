package services

import (
	"context"
	"health_checker/config"
	"health_checker/internal/domain/interfaces"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type HealthCheckerService struct {
	logger *logrus.Logger
	repo   interfaces.IRepository
	cfg    *config.Config
}

func NewHealthCheckerService(logger *logrus.Logger, repo interfaces.IRepository, cfg *config.Config) *HealthCheckerService {
	return &HealthCheckerService{
		logger: logger,
		repo:   repo,
		cfg:    cfg,
	}
}

// CheckSitesAvailability - проверка множества сайтов (запускается горутина для каждой проверки сайта)
func (lc *HealthCheckerService) CheckSitesAvailability(ctx context.Context, links []string) (map[string]string, int) {
	resultedLinks := make(map[string]string, len(links))

	wg := sync.WaitGroup{}
	for _, link := range links {
		wg.Add(1)
		go func(l string) { // передаем явно для избежания гонки
			defer wg.Done()
			if lc.checkSiteAvailability(ctx, link) {
				resultedLinks[link] = "available"
			} else {
				resultedLinks[link] = "not available"
			}
		}(link)
	}

	wg.Wait()

	linksID := lc.repo.AddLinks(resultedLinks)

	return resultedLinks, linksID
}

// GetSitesAvailabilityReportByLinksID - получение отчета по проверенным сайтам по id
func (lc *HealthCheckerService) GetSitesAvailabilityReportByLinksID(ctx context.Context, linksID int) map[string]string {
	healthReport := lc.repo.GetLinksByLinksID(linksID)

	return healthReport
}

// checkSiteAvailability - проверка одного сайта
func (lc *HealthCheckerService) checkSiteAvailability(ctx context.Context, siteUrl string) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, siteUrl, nil)
	if err != nil {
		lc.logger.Debugf("Failed to prepare request: %v\n", err.Error())
		return false
	}

	client := http.Client{
		Timeout: time.Second * 5, // таймаут из конфига в секундах
	}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 { // если неуспешный статускод - пишем в логи и возвращаем false
		return true
	} else {
		lc.logger.Debugf("Site response wrong http code: %d\n", resp.StatusCode)
		return false
	}
}
