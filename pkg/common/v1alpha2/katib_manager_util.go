/*

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

package v1alpha2

import (
	"os"
	"context"

	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	api_pb "github.com/kubeflow/katib/pkg/api/v1alpha2"
	"google.golang.org/grpc"
)

const (
	KatibManagerServiceIPEnvName        = "KATIB_MANAGER_PORT_6789_TCP_ADDR"
	KatibManagerServicePortEnvName      = "KATIB_MANAGER_PORT_6789_TCP_PORT"
	KatibManagerServiceNamespaceEnvName = "KATIB_MANAGER_NAMESPACE"
	KatibManagerService                 = "katib-manager"
	KatibManagerPort                    = "6789"
	ManagerAddr                         = KatibManagerService + ":" + KatibManagerPort
)

type katibManagerClientAndConn struct {
	Conn *grpc.ClientConn
	KatibManagerClient api_pb.ManagerClient
}

func GetManagerAddr() string {
	ns := os.Getenv(experimentsv1alpha2.DefaultKatibNamespaceEnvName)
	if len(ns) == 0 {
		addr := os.Getenv(KatibManagerServiceIPEnvName)
		port := os.Getenv(KatibManagerServicePortEnvName)
		if len(addr) > 0 && len(port) > 0 {
			return addr + ":" + port
		} else {
			return ManagerAddr
		}
	} else {
		return KatibManagerService + "." + ns + ":" + KatibManagerPort
	}
}

func getKatibManagerClientAndConn() (*katibManagerClientAndConn, error) {
	addr := GetManagerAddr()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	kcc := &katibManagerClientAndConn {
		Conn: conn,
		KatibManagerClient: api_pb.NewManagerClient(conn),
	}
	return kcc, nil
}

func RegisterExperiment(request *api_pb.RegisterExperimentRequest) (*api_pb.RegisterExperimentReply, error) {
	ctx := context.Background()
	kcc, err := getKatibManagerClientAndConn()
	if err != nil {
		return nil, err
	}
	defer kcc.Conn.Close()
	kc := kcc.KatibManagerClient
	return kc.RegisterExperiment(ctx, request)
}
