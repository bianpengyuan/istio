// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package perf

import (
	"encoding/json"

	istio_mixer_v1 "istio.io/api/mixer/v1"
	"istio.io/istio/mixer/pkg/attribute"
)

// Request interface is the common interface for all different types of requests.
type Request interface {
	// createRequestProtos causes the request to create one-or-more API request protos.
	getRequestProtos() []interface{}
}

// BasicReport is an implementation of Request that is used to explicitly specify a Report request.
type BasicReport struct {
	Attributes map[string]interface{}         `json:"attributes,omitempty"`
	Requests   []istio_mixer_v1.ReportRequest `json:"requests,omitempty"`
}

var _ Request = &BasicReport{}

// BasicCheck is an implementation of Request that is specified declaratively by the author.
type BasicCheck struct {
	Attributes map[string]interface{}                             `json:"attributes,omitempty"`
	Quotas     map[string]istio_mixer_v1.CheckRequest_QuotaParams `json:"quotas,omitempty"`
	Requests   []istio_mixer_v1.CheckRequest                      `json:"requests,omitempty"`
}

var _ Request = &BasicCheck{}

// BuildBasicReport ...
func BuildBasicReport(attributes map[string]interface{}) BasicReport {
	br := BasicReport{
		Attributes: attributes,
	}
	br.createRequestProtos()
	return br
}

// CreateRequest creates a request proto.
func (r *BasicReport) createRequestProtos() {
	requestBag := attribute.GetMutableBag(nil)
	for k, v := range r.Attributes {
		requestBag.Set(k, v)
	}

	var attrProto istio_mixer_v1.CompressedAttributes
	requestBag.ToProto(&attrProto, nil, 0)

	r.Requests = []istio_mixer_v1.ReportRequest{
		istio_mixer_v1.ReportRequest{
			Attributes: []istio_mixer_v1.CompressedAttributes{attrProto},
		},
	}
}

func (r BasicReport) getRequestProtos() []interface{} {
	reqs := make([]interface{}, 0, len(r.Requests))
	for _, req := range r.Requests {
		reqs = append(reqs, interface{}(&req))
	}
	return reqs
}

// MarshalJSON marshals the report as JSON.
func (r BasicReport) MarshalJSON() ([]byte, error) {
	m := make(map[string]json.RawMessage, 2)

	var err error
	m["type"], _ = json.Marshal("basicReport")

	m["attributes"], err = json.Marshal(r.Attributes)
	if err != nil {
		return nil, err
	}
	m["requests"], err = json.Marshal(r.Requests)
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

// BuildBasicCheck ...
func BuildBasicCheck(attributes map[string]interface{}, quotas map[string]istio_mixer_v1.CheckRequest_QuotaParams) BasicCheck {
	c := BasicCheck{
		Attributes: attributes,
		Quotas:     quotas,
	}

	c.createRequestProtos()
	return c
}

// CreateRequest creates a request proto.
func (c *BasicCheck) createRequestProtos() {
	requestBag := attribute.GetMutableBag(nil)
	for k, v := range c.Attributes {
		requestBag.Set(k, v)
	}

	var attrProto istio_mixer_v1.CompressedAttributes
	requestBag.ToProto(&attrProto, nil, 0)

	c.Requests = []istio_mixer_v1.CheckRequest{
		istio_mixer_v1.CheckRequest{
			Attributes: attrProto,
			Quotas:     c.Quotas,
		},
	}
}

func (c BasicCheck) getRequestProtos() []interface{} {
	reqs := make([]interface{}, 0, len(c.Requests))
	for _, req := range c.Requests {
		reqs = append(reqs, interface{}(&req))
	}
	return reqs
}

// MarshalJSON marshals the report as JSON.
func (c BasicCheck) MarshalJSON() ([]byte, error) {
	m := make(map[string]json.RawMessage, 3)

	var err error
	m["type"], _ = json.Marshal("basicCheck")

	if m["attributes"], err = json.Marshal(c.Attributes); err != nil {
		return nil, err
	}

	if c.Quotas != nil {
		if m["quotas"], err = json.Marshal(c.Quotas); err != nil {
			return nil, err
		}
	}
	m["requests"], err = json.Marshal(c.Requests)
	if err != nil {
		return nil, err
	}

	return json.Marshal(m)
}
