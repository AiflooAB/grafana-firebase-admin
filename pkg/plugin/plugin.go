package plugin

import (
	"context"
	"encoding/json"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

var (
	_ backend.QueryDataHandler      = (*FirebaseAdminDatasource)(nil)
	_ backend.CheckHealthHandler    = (*FirebaseAdminDatasource)(nil)
	_ instancemgmt.InstanceDisposer = (*FirebaseAdminDatasource)(nil)
)

type settings struct {
}

// NewFirebaseAdminDatasource creates a new datasource instance.
func NewFirebaseAdminDatasource(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	credentials := settings.DecryptedSecureJSONData["credentialsJSON"]
	s := struct {
		ProjectID string `json:"projectID"`
	}{}
	json.Unmarshal(settings.JSONData, &s)

	projectID := s.ProjectID

	opt := option.WithCredentialsJSON([]byte(credentials))

	config := firebase.Config{
		ProjectID: projectID,
	}

	app, err := firebase.NewApp(context.Background(), &config, opt)
	if err != nil {
		log.DefaultLogger.Error("error initializing app", "error", err)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		log.DefaultLogger.Error("error initializing auth client", "error", err)
	}

	return &FirebaseAdminDatasource{
		client: client,
	}, nil
}

type FirebaseAdminDatasource struct {
	client *auth.Client
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewFirebaseAdminDatasource factory function.
func (d *FirebaseAdminDatasource) Dispose() {
	// Clean up datasource instance resources.
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *FirebaseAdminDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	log.DefaultLogger.Info("QueryData called", "request", req)

	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

func (d *FirebaseAdminDatasource) query(ctx context.Context, _ backend.PluginContext, _ backend.DataQuery) backend.DataResponse {
	response := backend.DataResponse{}

	var uids, displayNames, emails, phoneNumbers []string

	iter := d.client.Users(ctx, "")
	for {
		u, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.DefaultLogger.Error("error listing users", "error", err)
		}
		uids = append(uids, u.UID)
		displayNames = append(displayNames, u.DisplayName)
		emails = append(emails, u.Email)
		phoneNumbers = append(phoneNumbers, u.PhoneNumber)
	}

	response.Frames = data.Frames{
		data.NewFrame(
			"response",
			data.NewField("uid", data.Labels{}, uids),
			data.NewField("displayName", data.Labels{}, displayNames),
			data.NewField("email", data.Labels{}, emails),
			data.NewField("phoneNumber", data.Labels{}, phoneNumbers),
		),
	}

	return response
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *FirebaseAdminDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	log.DefaultLogger.Info("CheckHealth called", "request", req)

	var status = backend.HealthStatusOk
	var message = "Data source is working"

	_, err := d.client.GetUser(ctx, "")
	if !auth.IsUserNotFound(err) {
		status = backend.HealthStatusError
		message = err.Error()
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}
