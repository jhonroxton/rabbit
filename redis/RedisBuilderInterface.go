package redis_test

type RedisBuilderInterface interface {
	// SetCluster 设置是否是集群模式，默认为false 单体模式 true为集群模式
	SetCluster(isCluster bool) RedisBuilderInterface
	GetCluster() bool

	// SetHost 设置Ip 主要为Ip:端口
	SetHost(host string) RedisBuilderInterface
	GetHost() string

	// SetPassWord 设置密码
	SetPassWord(password string) RedisBuilderInterface
	GetPassWord() string

	// SetDb 使用设置的数据库
	SetDb(db int) RedisBuilderInterface
	GetBd() int

	// SetMaxRetries 设置最大重试次数
	SetMaxRetries(maxRetries int) RedisBuilderInterface
	GetMaxRetries() int

	// SetPoolSize 设置连接池大小
	SetPoolSize(poolSzie int) RedisBuilderInterface
	GetPoolSize() int

	// SetPoolTimeOut 设置连接池超时时间
	SetPoolTimeOut(poolTimeOut int) RedisBuilderInterface
	GetPoolTimeOut() int

	// SetReadTimeOut 设置读取超时时间
	SetReadTimeOut(readTimeOut int) RedisBuilderInterface
	GetReadTimeOut() int

	// SetDialTimeOut 建立连接超时时间
	SetDialTimeOut(dialTimeOut int) RedisBuilderInterface
	GetDialTimeOut() int
}
