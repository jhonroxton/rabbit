package redis_test

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"strings"
	"sync"
	"time"
)

var (
	REDIS_ERR_BUILDER = errors.New("redis builder error")
)

const (
	DEFAULT_MAX_RETRIES   int = 3   //默认最大重试次数
	DEFAULT_POOL_SIZE     int = 5   //默认连接池大小
	DEFAULT_POOL_TIME_OUT int = 20  //默认连接池超时时间(秒)
	DEFAULT_READ_TIME_OUT int = 20  //默认读取超时(秒)
	DEFAULT_DIAL_TIME_OUT int = 180 //默认连接超时时间
)

/*
	redis连接池使用
*/
type redisBuilder struct {
	isCluster   bool //是否是集群
	host        string
	passWord    string
	db          int
	maxRetries  int //最大重试次数
	poolSize    int //连接池大小
	poolTimeOut int //连接池超时
	readTimeOut int //读取超时
	dialTimeOut int //建立连接超时
}

/*
   带参数初始化
*/
func NewRedisBuilder(isCluster bool, host string, password string, db int) RedisBuilderInterface {

	redisBuilder := &redisBuilder{
		isCluster:   isCluster,
		host:        host,
		passWord:    password,
		db:          db,
		maxRetries:  DEFAULT_MAX_RETRIES,
		poolSize:    DEFAULT_POOL_SIZE,
		poolTimeOut: DEFAULT_POOL_TIME_OUT,
		readTimeOut: DEFAULT_READ_TIME_OUT,
		dialTimeOut: DEFAULT_DIAL_TIME_OUT,
	}

	return redisBuilder
}

/*
   无参初始化
*/
func NewRedisBuilderEm() RedisBuilderInterface {
	return &redisBuilder{}
}

func (r *redisBuilder) SetCluster(isCluster bool) RedisBuilderInterface {
	r.isCluster = isCluster
	return r
}

func (r *redisBuilder) GetCluster() bool {
	return r.isCluster
}

func (r *redisBuilder) SetHost(host string) RedisBuilderInterface {
	r.host = host
	return r
}

func (r *redisBuilder) GetHost() string {
	return r.host
}

func (r *redisBuilder) SetPassWord(password string) RedisBuilderInterface {
	r.passWord = password
	return r
}

func (r *redisBuilder) GetPassWord() string {
	return r.passWord
}

func (r *redisBuilder) SetDb(db int) RedisBuilderInterface {
	r.db = db
	return r
}

func (r *redisBuilder) GetBd() int {

	return r.db
}

func (r *redisBuilder) SetMaxRetries(maxRetries int) RedisBuilderInterface {

	r.maxRetries = maxRetries
	return r
}

func (r *redisBuilder) GetMaxRetries() int {
	return r.maxRetries
}

func (r *redisBuilder) SetPoolSize(poolSzie int) RedisBuilderInterface {

	r.poolSize = poolSzie
	return r
}

func (r *redisBuilder) GetPoolSize() int {
	return r.poolSize
}

func (r *redisBuilder) SetPoolTimeOut(poolTimeOut int) RedisBuilderInterface {
	r.poolTimeOut = poolTimeOut
	return r
}

func (r *redisBuilder) GetPoolTimeOut() int {
	return r.poolTimeOut
}

func (r *redisBuilder) SetReadTimeOut(readTimeOut int) RedisBuilderInterface {

	r.readTimeOut = readTimeOut
	return r
}

func (r *redisBuilder) GetReadTimeOut() int {
	return r.readTimeOut
}

func (r *redisBuilder) SetDialTimeOut(dialTimeOut int) RedisBuilderInterface {

	r.dialTimeOut = dialTimeOut
	return r
}

func (r *redisBuilder) GetDialTimeOut() int {
	return r.dialTimeOut
}

type redisPool struct {
	builder       RedisBuilderInterface
	clientCluster *redis.ClusterClient
	clientAlone   *redis.Client
	isInitialize  bool
}

var redisOnce sync.Once
var redisInstance *redisPool

func NewRedisPool() *redisPool {
	redisOnce.Do(func() {
		redisInstance = &redisPool{}
	})

	return redisInstance
}

func (rp *redisPool) SetBuilder(builder RedisBuilderInterface) *redisPool {
	rp.builder = builder
	return rp
}

func (rp *redisPool) GetBuilder() RedisBuilderInterface {
	return rp.builder
}

func (rp *redisPool) IsCluster() bool {
	return rp.builder.GetCluster()
}

func (rp *redisPool) GetClusterClient() *redis.ClusterClient {
	return rp.clientCluster
}

func (rp *redisPool) GetAloneClinet() *redis.Client {
	return rp.clientAlone
}

func (rp *redisPool) GetClient() (*redis.Client, *redis.ClusterClient) {
	if rp.builder.GetCluster() {
		return nil, rp.clientCluster
	}
	return rp.clientAlone, nil
}

func getIsCluster() (bool, error) {
	pool := NewRedisPool()
	err := pool.Init()
	if err != nil {
		return false, err
	}
	return pool.IsCluster(), nil
}

func getClientCluster() (*redis.ClusterClient, error) {
	pool := NewRedisPool()
	err := pool.Init()
	if err != nil {
		return nil, err
	}
	return pool.GetClusterClient(), nil
}

func getClientAlone() (*redis.Client, error) {
	pool := NewRedisPool()
	err := pool.Init()
	if err != nil {
		return nil, err
	}
	return pool.GetAloneClinet(), nil
}

func (rp *redisPool) Init() error {
	if !rp.isInitialize {
		if rp.builder == nil {
			return REDIS_ERR_BUILDER
		}
		//集群版本初始化
		if rp.builder.GetCluster() {
			clinet, err := initRedisCluster(rp.builder.GetHost(), rp.builder.GetPassWord(), rp.builder.GetMaxRetries(), rp.builder.GetPoolSize(), rp.builder.GetDialTimeOut(), rp.builder.GetPoolTimeOut(), rp.builder.GetReadTimeOut())
			if err != nil {
				return err
			}
			rp.isInitialize = true
			rp.clientCluster = clinet
		} else { //单体版本
			client, err := initRedisAlone(rp.builder.GetHost(), rp.builder.GetPassWord(), rp.builder.GetBd(), rp.builder.GetMaxRetries(), rp.builder.GetPoolSize(), rp.builder.GetDialTimeOut(), rp.builder.GetPoolTimeOut(), rp.builder.GetReadTimeOut())
			if err != nil {
				return err
			}
			rp.isInitialize = true
			rp.clientAlone = client
		}
	}
	return nil
}

/*
	初始化集群版 redis

*/

func initRedisCluster(host string, password string, maxRetries int, poolSize int, dialTimeOut int, poolTimeOut int, readTimeOut int) (*redis.ClusterClient, error) {
	hostDns := strings.Split(host, ",")

	redisClinet := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:           hostDns,
		Password:        password,
		MaxRetries:      maxRetries,
		PoolSize:        poolSize,
		DialTimeout:     time.Second * time.Duration(dialTimeOut),
		PoolTimeout:     time.Second * time.Duration(poolTimeOut),
		ReadTimeout:     time.Second * time.Duration(readTimeOut),
		MinRetryBackoff: 3 * time.Second, // 最小重试间隔
		MaxRetryBackoff: 5 * time.Second, //最大重试间隔
	})
	_, err := redisClinet.Ping().Result()
	if err != nil {
		return nil, errors.New("redis集群连接失败：" + err.Error())
	}
	fmt.Println("redis集群连接成功~~~~~~~~~~~~~~~~~~~")
	return redisClinet, nil
}

/*
	初始化单机版
*/
func initRedisAlone(host string, password string, db int, maxRetries int, poolSize int, dialTimeOut int, poolTimeOut int, readTimeout int) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:            host,
		Password:        password,
		DB:              db,
		MaxRetries:      maxRetries,
		PoolSize:        poolSize,
		MinIdleConns:    10,
		DialTimeout:     time.Second * time.Duration(dialTimeOut),
		PoolTimeout:     time.Second * time.Duration(poolTimeOut),
		ReadTimeout:     time.Second * time.Duration(readTimeout),
		MinRetryBackoff: 3 * time.Second, // 最小重试间隔
		MaxRetryBackoff: 5 * time.Second, //最大重试间隔
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		return nil, errors.New("redis单机版连接失败：" + err.Error())
	}
	fmt.Println("============================redis单机版连接成\n")
	return redisClient, nil
}
