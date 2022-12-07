/*
 * Copyright 1999-2020 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

func main() {
	sc := []constant.ServerConfig{
		{
			IpAddr: "{serverAddr}",
			Port:   8848,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         "{namespace}", //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogRollingConfig:    &lumberjack.Logger{MaxSize: 10},
		LogLevel:            "debug",
		OpenKMS:			true,
		RegionId:			"{regionId}",
		AccessKey:			"{ak}",
		SecretKey:			"{sk}",
	}

	// create config client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}

	//publish config
	//config key=dataId+group+namespaceId
	_, err = client.PublishConfig(vo.ConfigParam{
		DataId:  "cipher-test-data",
		Group:   "test-group",
		Content: "hello world!",
	})
	_, err = client.PublishConfig(vo.ConfigParam{
		DataId:  "cipher-test-data-2",
		Group:   "test-group",
		Content: "hello world!",
	})

	if err != nil {
		fmt.Printf("PublishConfig err:%+v \n", err)
	}

	time.Sleep(1 * time.Second)

	//get config
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: "cipher-test-data",
		Group:  "test-group",
	})
	fmt.Println("GetConfig,config :" + content)


	//Listen config change,key=dataId+group+namespaceId.
	err = client.ListenConfig(vo.ConfigParam{
		DataId: "test-data",
		Group:  "test-group",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
		},
	})

	err = client.ListenConfig(vo.ConfigParam{
		DataId: "test-data-2",
		Group:  "test-group",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
		},
	})

	time.Sleep(9999 * time.Second)

}
