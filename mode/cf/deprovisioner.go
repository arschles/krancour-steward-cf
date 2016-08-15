package cf

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/web"
	"github.com/juju/loggo"
)

type backendDeprovisionResp struct {
	Operation string `json:"operation"`
}

type deprovisioner struct {
	logger loggo.Logger
	cl     *RESTClient
}

func (d deprovisioner) Deprovision(instanceID, serviceID, planID string) (*mode.DeprovisionResponse, error) {
	query := url.Values(map[string][]string{})
	query.Add(serviceIDQueryKey, serviceID)
	query.Add(planIDQueryKey, planID)
	req, err := d.cl.Delete(d.logger, query, "v2", "service_instances", instanceID)
	if err != nil {
		return nil, err
	}
	res, err := d.cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted {
		return nil, web.ErrUnexpectedResponseCode{
			URL:      req.URL.String(),
			Expected: http.StatusOK,
			Actual:   res.StatusCode,
		}
	}
	deproResp := new(backendDeprovisionResp)
	if err := json.NewDecoder(res.Body).Decode(deproResp); err != nil {
		return nil, err
	}
	return &mode.DeprovisionResponse{Status: res.StatusCode, Operation: deproResp.Operation}, nil
}

// NewDeprovisioner creates a new CloudFoundry-broker-backed deprovisioner implementation
func NewDeprovisioner(logger loggo.Logger, cl *RESTClient) mode.Deprovisioner {
	return deprovisioner{logger: logger, cl: cl}
}