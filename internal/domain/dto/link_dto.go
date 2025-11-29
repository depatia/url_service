package dto

import (
	"fmt"
	"health_checker/pkg/url_normalizer"
)

type (
	CheckSitesAvailabilityRequest struct {
		Links []string `json:"links"`
	}

	CheckSitesAvailabilityResponse struct {
		Links    map[string]string `json:"links"`
		LinksNum int               `json:"links_num"`
	}

	GenerateReportRequest struct {
		LinksList []int `json:"links_list"`
	}
)

func (r *CheckSitesAvailabilityRequest) Validate() error {
	if len(r.Links) == 0 {
		return fmt.Errorf("field 'links' is empty")
	}

	if len(r.Links) > 25 {
		return fmt.Errorf("too many links")
	}

	for i, link := range r.Links {
		newLink, err := url_normalizer.NormalizeURL(link)
		if err != nil {
			return fmt.Errorf("failed to normalize url: %v", err.Error())
		}
		r.Links[i] = newLink
	}

	return nil
}

func (r GenerateReportRequest) Validate() error {
	if len(r.LinksList) == 0 {
		return fmt.Errorf("field 'links_list' is empty")
	}

	return nil
}
