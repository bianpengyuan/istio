// Copyright 2019 Istio Authors
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

package tokenmanager

import (
	"errors"
	"fmt"

	"istio.io/istio/pkg/bootstrap/platform"
	"istio.io/istio/security/pkg/stsservice"
	"istio.io/istio/security/pkg/stsservice/tokenmanager/google"
)

const (
	// GoogleTokenExchange is the name of the google token exchange service.
	GoogleTokenExchange = "GoogleTokenExchange"
)

// Plugin provides common interfaces for specific token exchange services.
type Plugin interface {
	ExchangeToken(parameters stsservice.StsRequestParameters) ([]byte, error)
	DumpPluginStatus() ([]byte, error)
}

type TokenManager struct {
	plugin Plugin
}

type Config struct {
	TrustDomain string
}

// GcpProjectInfo stores GCP project information, including project number,
// project ID, cluster location, cluster name
type GcpProjectInfo struct {
	Number          string
	id              string
	cluster         string
	clusterLocation string
}

func getGcpProjectNumber() GcpProjectInfo {
	fmt.Println("try to get project number")
	info := GcpProjectInfo{}
	if platform.IsGCP() {
		md := platform.NewGCP().Metadata()
		fmt.Println("is on gcp")
		fmt.Println(md)
		if projectNum, found := md[platform.GCPProjectNumber]; found {
			fmt.Println("project num " + projectNum)
			info.Number = projectNum
		}
		if projectID, found := md[platform.GCPProject]; found {
			fmt.Println("projectID " + projectID)
			info.id = projectID
		}
		if clusterName, found := md[platform.GCPCluster]; found {
			fmt.Println("cluster name " + clusterName)
			info.cluster = clusterName
		}
		if clusterLocation, found := md[platform.GCPLocation]; found {
			fmt.Println("cluster location " + clusterLocation)
			info.clusterLocation = clusterLocation
		}
	} else {
		fmt.Println("is not on gcp")
	}
	return info
}

// CreateTokenManager creates a token manager with specified type and returns
// that token manager
func CreateTokenManager(tokenManagerType string, config Config) stsservice.TokenManager {
	fmt.Println("create token manager " + tokenManagerType)
	tm := &TokenManager{
		plugin: nil,
	}
	switch tokenManagerType {
	case GoogleTokenExchange:
		fmt.Println("is google token exchange")
		if projectInfo := getGcpProjectNumber(); len(projectInfo.Number) > 0 {
			gKEClusterURL := fmt.Sprintf("https://container.googleapis.com/v1/projects/%s/locations/%s/clusters/%s",
				projectInfo.id, projectInfo.clusterLocation, projectInfo.cluster)
			if p, err := google.CreateTokenManagerPlugin(config.TrustDomain, projectInfo.Number, gKEClusterURL); err == nil {
				tm.plugin = p
			}
		}
	}
	return tm
}

func (tm *TokenManager) GenerateToken(parameters stsservice.StsRequestParameters) ([]byte, error) {
	if tm.plugin != nil {
		return tm.plugin.ExchangeToken(parameters)
	}
	return nil, errors.New("no plugin is found")
}

func (tm *TokenManager) DumpTokenStatus() ([]byte, error) {
	if tm.plugin != nil {
		return tm.plugin.DumpPluginStatus()
	}
	return nil, errors.New("no plugin is found")
}

// SetPlugin sets token exchange plugin for testing purposes only.
func (tm *TokenManager) SetPlugin(p Plugin) {
	tm.plugin = p
}
