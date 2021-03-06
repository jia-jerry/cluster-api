/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package phases

import (
	"fmt"

	"github.com/golang/glog"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/clusterdeployer/bootstrap"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/clusterdeployer/clusterclient"
)

func CreateBootstrapCluster(provisioner bootstrap.ClusterProvisioner, cleanupBootstrapCluster bool, clientFactory clusterclient.Factory) (clusterclient.Client, func(), error) {
	glog.Info("Creating bootstrap cluster")

	cleanupFn := func() {}
	if err := provisioner.Create(); err != nil {
		return nil, cleanupFn, fmt.Errorf("could not create bootstrap control plane: %v", err)
	}

	if cleanupBootstrapCluster {
		cleanupFn = func() {
			glog.Info("Cleaning up bootstrap cluster.")
			provisioner.Delete()
		}
	}

	bootstrapKubeconfig, err := provisioner.GetKubeconfig()
	if err != nil {
		return nil, cleanupFn, fmt.Errorf("unable to get bootstrap cluster kubeconfig: %v", err)
	}
	bootstrapClient, err := clientFactory.NewClientFromKubeconfig(bootstrapKubeconfig)
	if err != nil {
		return nil, cleanupFn, fmt.Errorf("unable to create bootstrap client: %v", err)
	}

	return bootstrapClient, cleanupFn, nil
}
