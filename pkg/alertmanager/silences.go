package alertmanager

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	urlpkg "net/url"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/prometheus/alertmanager/api/v2/client"
	"github.com/prometheus/alertmanager/api/v2/client/silence"
	"github.com/prometheus/alertmanager/api/v2/models"
	"github.com/prometheus/alertmanager/pkg/labels"
)

func GetAlertmanagerSilances(url string) ([]*models.GettableSilence, error) {
	log.Printf("Fetching silances from %s", url)

	parsedURL, err := urlpkg.Parse(url)
	if err != nil {
		return nil, err
	}

	alertmanager := client.NewHTTPClientWithConfig(
		strfmt.NewFormats(),
		&client.TransportConfig{
			Host:     parsedURL.Host,
			BasePath: "/api/v2",
			Schemes:  []string{parsedURL.Scheme},
		})

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	httpClient := http.Client{
		Transport: transport,
	}

	params := silence.GetSilencesParams{
		Context:    context.TODO(),
		HTTPClient: &httpClient,
	}

	// Call the Alertmanager API to retrieve the list of silences
	resp, err := alertmanager.Silence.GetSilences(&params)
	if err != nil {
		return nil, err
	}

	// Extract the silences from the API response
	silences := resp.Payload

	return silences, nil

}

func GetAlertmanagerSilencesByMatchers(url, key, value string) ([]*models.GettableSilence, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	parsedURL, err := urlpkg.Parse(url)
	if err != nil {
		return nil, err
	}
	transport.TLSHandshakeTimeout = 10 * time.Second

	alertmanager := client.NewHTTPClientWithConfig(
		strfmt.NewFormats(),
		&client.TransportConfig{
			Host:     parsedURL.Host,
			BasePath: "/api/v2",
			Schemes:  []string{parsedURL.Scheme},
		})

	matcher, err := labels.NewMatcher(labels.MatchRegexp, key, value)
	if err != nil {
		log.Printf("cannot create matcher with error %s", err)
	}

	param := silence.NewGetSilencesParams().
		WithFilter([]string{matcher.String()})

	// Call the Alertmanager API to retrieve the list of silences
	resp, err := alertmanager.Silence.GetSilences(param)
	if err != nil {
		return nil, err
	}

	// Extract the silences from the API response
	silences := resp.Payload

	return silences, nil
}

func PostAlertmanagerSilences(url string, matchers []*models.Matcher, startsAt *strfmt.DateTime, endsAt *strfmt.DateTime, createdBy string, comment string) error {

	parsedURL, err := urlpkg.Parse(url)
	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	httpClient := http.Client{
		Transport: transport,
	}

	params := silence.NewPostSilencesParams().WithContext(context.Background())

	var amSilence = models.PostableSilence{
		Silence: models.Silence{
			CreatedBy: &createdBy,
			Comment:   &comment,
			Matchers:  matchers,
			StartsAt:  startsAt,
			EndsAt:    endsAt,
		},
	}

	params.SetSilence(&amSilence)
	params.SetHTTPClient(&httpClient)

	alertmanager := client.NewHTTPClientWithConfig(
		strfmt.NewFormats(),
		&client.TransportConfig{
			Host:     parsedURL.Host,
			BasePath: "/api/v2",
			Schemes:  []string{parsedURL.Scheme},
		})
	// Call the Alertmanager API to retrieve the list of silences
	resp, err := alertmanager.Silence.PostSilences(params)
	if err != nil {
		return err
	}

	if resp.IsSuccess() {
		log.Printf("Silance pushed sucessfully to: %s", url)
	}

	return nil
}

//AreSilencesEqual is used to compare alertManager Silences
func AreSilencesEqual(silence1, silence2 models.GettableSilence) bool {
	if len(silence1.Matchers) != len(silence2.Matchers) {
		return false
	}

	for _, m1 := range silence1.Matchers {
		found := false
		for _, m2 := range silence2.Matchers {
			if strings.Compare(*m1.Name, *m2.Name) == 0 &&
				strings.Compare(*m1.Value, *m2.Value) == 0 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
