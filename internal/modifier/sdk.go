/**
# Copyright (c) 2024, Shanghai Iluvatar CoreX Semiconductor Co., Ltd.
# All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may
# not use this file except in compliance with the License. You may obtain
# a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
**/

package modifier

import (
	"context"
	"fmt"
	"time"

	"gitee.com/deep-spark/ix-container-runtime/internal/config/image"
	"gitee.com/deep-spark/ix-container-runtime/internal/oci"
	pb "gitee.com/deep-spark/ix-container-runtime/internal/sdk"
	"github.com/opencontainers/runtime-spec/specs-go"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type sdkModifier struct {
	client pb.SdkServiceClient
	conn   *grpc.ClientConn
	ctx    context.Context
	Cancel context.CancelFunc

	Change image.VisibleSdk
}

const (
	visibleSdkEnvvar   = "COREX_IMAGE"
	pathEnv            = "PATH"
	ldPathEnv          = "LD_LIBRARY_PATH"
	defaultDestination = "/usr/local/corex"

	pathAdded   = "/usr/local/corex/bin"
	ldPathAdded = "/usr/local/corex/lib64"
)

func (c *sdkModifier) QueryCache(imageName string) (string, string, string, error) {
	req := &pb.QueryCacheRequest{ImageName: imageName, Method: ""}
	resp, err := c.client.QueryCache(c.ctx, req)
	if err != nil {
		return "", "", "", err
	}
	if resp.Destination == "" {
		log.Printf("No such Destination")
	} else {
		log.Printf("Response: %s", resp.Destination)
	}
	return resp.Destination, resp.Method, resp.Type, nil
}

func (c *sdkModifier) PrepareCache(imageName string) (bool, string, error) {
	req := &pb.PrepareCacheRequest{ImageName: imageName, Method: "", Filter: []string{"type=sdk", "vendor=iluvatarcorex"}}
	resp, err := c.client.PrepareCache(c.ctx, req)
	if err != nil {
		return false, "", err
	}

	if resp.Status == "ok" {
		return false, resp.Method, nil
	} else {
		return true, resp.Method, nil
	}
}

func (s sdkModifier) Modify(spec *specs.Spec) error {
	log.Printf("entry modfiy\n")
	var destination_corex_dir string

	pathVal := s.Change.Path()
	pathVal = fmt.Sprintf("%v=%v:%v", pathEnv, pathAdded, pathVal)
	ldpathVal := s.Change.LdPath()
	if ldpathVal == "" {
		ldpathVal = fmt.Sprintf("%v=%v", ldPathEnv, ldPathAdded)
	} else {
		ldpathVal = fmt.Sprintf("%v=%v:%v", ldPathEnv, ldPathAdded, ldpathVal)
	}
	//var needPull bool
	image := s.Change.Name()
	if image == "" {
		return nil
	}
	destination, _, imageType, err := s.QueryCache(image)
	if err != nil {
		log.Printf("query cache failed %v", image)
		return err
	}

	if destination == "" {
		_, _, err = s.PrepareCache(image)
		if err != nil {
			log.Printf("PrepareCache failed %v", image)
			return err
		} else {
			log.Printf("Cache not exists, call prepare to pull image %v", image)
			return fmt.Errorf("Cache not exists, call prepare to pull image %v", image)
		}
	}
	if imageType != "sdk" {
		return fmt.Errorf("Image type is not sdk, real type:%v\n", imageType)
	}

	destination_corex_dir = destination

	spec.Mounts = append(spec.Mounts,
		specs.Mount{Destination: defaultDestination,
			Source: destination_corex_dir,
			Type:   "linux",
			Options: []string{
				"ro",
				"nosuid",
				"nodev",
				"bind"}})

	log.Printf("---> pathval :%v  ldpathval:%v\n", pathVal, ldpathVal)
	spec.Process.Env = append(spec.Process.Env, pathVal)
	spec.Process.Env = append(spec.Process.Env, ldpathVal)

	return nil
}

func NewSdkModifier(ig image.CUDA) oci.SpecModifier {
	var err error
	ret := sdkModifier{}

	connectPath := "unix://" + ig.Cfg.SdkSocketPath
	log.Printf("connect address:%v", connectPath)
	ret.conn, err = grpc.NewClient(connectPath, grpc.WithInsecure())
	if err != nil {
		log.Printf("open sdk local failed\n")
		return nil
	}
	ret.client = pb.NewSdkServiceClient(ret.conn)
	ret.ctx, ret.Cancel = context.WithTimeout(context.Background(), time.Second)
	ret.Change = ig.SdkFromEnvvars(visibleSdkEnvvar, pathEnv, ldPathEnv)

	return ret
}
