package redis_test

import (
	"github.com/go-redis/redis"
)

type BaseRedisModel struct{}

/**
判断是否是集群
*/
func (rm *BaseRedisModel) IsCluster() (bool, error) {
	return getIsCluster()
}

/**
获取单机客户端
*/
func (rm *BaseRedisModel) GetAloneClient() (*redis.Client, error) {
	return getClientAlone()
}

/**
获取集群客户端
*/
func (rm *BaseRedisModel) GetClusterClient() (*redis.ClusterClient, error) {
	return getClientCluster()
}

func (rm *BaseRedisModel) GetClient() (redis.Cmdable, error) {
	isCluster, err := getIsCluster()
	if err != nil {
		return nil, err
	}
	if isCluster {
		client, err := getClientCluster()
		if err != nil {
			return nil, err
		}
		return client, nil
	} else {
		client, err := getClientAlone()
		if err != nil {
			return nil, err
		}
		return client, nil
	}

}
