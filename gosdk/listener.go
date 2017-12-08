package gosdk

type ConfigListener interface {
	HandlePutEvent(config *Config) error    //处理配置更新事件，返回error表示处理失败
	HandleDeleteEvent(config *Config) error //处理删除事件，返回error表示处理失败
}
