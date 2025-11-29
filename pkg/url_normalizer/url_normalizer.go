package url_normalizer

import (
	"net/url"
	"strings"
)

// NormalizeURL - нормализует переданный урл в формат "https://google.com"
func NormalizeURL(siteUrl string) (string, error) {
	siteUrl = strings.TrimSpace(siteUrl)

	// Добавляем временный протокол для парсинга, если его нет
	if !strings.Contains(siteUrl, "://") {
		siteUrl = "https://" + siteUrl
	}

	u, err := url.Parse(siteUrl)
	if err != nil {
		return "", err
	}

	// Всегда устанавливаем https
	if u.Scheme != "http" {
		u.Scheme = "https"
	}

	return u.String(), nil
}
