package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	"github.com/sirupsen/logrus"
)

func configurationForProxy(secure bool) *datadog.Configuration {
	proxyHost := "localhost"
	protocol := "http"
	logrus.Debugf("Datadog client config: secure=%t, host=%s", secure, proxyHost)

	conf := datadog.NewConfiguration()
	conf.Servers = []datadog.ServerConfiguration{
		{
			URL:         protocol + "://{site}",
			Description: "No description provided",
			Variables: map[string]datadog.ServerVariable{
				"site": {
					Description:  "The regional site for Datadog customers.",
					DefaultValue: proxyHost,
					EnumValues: []string{
						"api.datadoghq.com",
						"api.us3.datadoghq.com",
						"api.us5.datadoghq.com",
						"api.datadoghq.eu",
						"api.ddog-gov.com",
						proxyHost,
					},
				},
			},
		},
		{
			URL:         "{protocol}://{name}",
			Description: "No description provided",
			Variables: map[string]datadog.ServerVariable{
				"name": {
					Description:  "Full site DNS name.",
					DefaultValue: proxyHost,
				},
				"protocol": {
					Description:  "The protocol for accessing the API.",
					DefaultValue: protocol,
				},
			},
		},
		{
			URL:         protocol + "://{site}",
			Description: "No description provided",
			Variables: map[string]datadog.ServerVariable{
				"site": {
					Description:  "Any Datadog deployment.",
					DefaultValue: proxyHost,
				},
			},
		},
	}
	return conf
}

func datadogSetup(ctx context.Context, proxyHost string, secure bool) (context.Context, *datadog.APIClient) {
	ctx = datadog.NewDefaultContext(ctx)
	logrus.WithContext(ctx).Debug("Setting up datadog client")
	configuration := configurationForProxy(secure)
	logrus.WithContext(ctx).Debug("Set datadog client config successfully")
	datadogClient := datadog.NewAPIClient(configuration)
	logrus.WithContext(ctx).Debug("Set up datadog client successfully")
	return ctx, datadogClient
}

func PushMetrics(body datadogV2.MetricPayload) error {
	ctx, datadogClient := datadogSetup(context.Background(), "", false)

	payload, resp, err := datadogV2.NewMetricsApi(datadogClient).SubmitMetrics(ctx, body, *datadogV2.NewSubmitMetricsOptionalParameters())
	if err != nil {
		logrus.Errorf("Error when calling `MetricsApi.SubmitMetrics`: %v\n", err)
		logrus.Debugf("Full HTTP response: %v\n", resp)
		payloadContent, _ := json.MarshalIndent(payload, "", "  ")
		logrus.Debugf("Payload from `MetricsApi.SubmitMetrics`:\n%s\n", payloadContent)
		clientConfig := datadogClient.GetConfig()
		logrus.Debugf("Datadog Client Configuration:\n%v\n", clientConfig)
	}
	return err
}

func makeEmptyMetricBody() datadogV2.MetricPayload {
	return datadogV2.MetricPayload{
		Series: []datadogV2.MetricSeries{
			{
				Metric: "system-alive",
				Type:   datadogV2.METRICINTAKETYPE_GAUGE.Ptr(),
				Points: []datadogV2.MetricPoint{
					{
						Timestamp: datadog.PtrInt64(time.Now().Unix()),
						Value:     datadog.PtrFloat64(float64(1)),
					},
				},
			},
		},
	}
}

func main() {
	err := PushMetrics(makeEmptyMetricBody())
	logrus.Error(err)
}
