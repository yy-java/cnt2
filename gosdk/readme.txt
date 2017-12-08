1. 安装：go get -insecure github.com/yy-java/cnt2/cnt2goclient

2. 开始
    1）、 cnt2.Start()方法启动配置中心，并返回一个cnt2Service实例

    type ClientConfig struct {
    	Endpoints     []string  //必填
    	DialTimeout   time.Duration //可选，默认是30秒
    	LocalFilePath string //可选，指定本地缓存的存储路径，便于在配置中心不可用时能使用上一次的配置，默认为：/data/file/cnt2/
    	App           string //必填
    	Profile       string //必填
    	EnableCommon  bool   //可选，默认为false，是否开启公共配置，需要自己创建为"common"的profile。非common覆盖common的同名key配置。
    }

    cnt2Service, err := cnt2.Start(&ClientConfig{Endpoints: []string{"61.147.187.152:2379", "61.147.187.142:2379", "61.147.187.150:2379"}, App: "demo", Profile: "development"})


    2）、 实现cnt2.Listener接口
    type TestListenter struct{}
    //通知key的增加或者修改事件
    func (t *TestListenter) HandlePutEvent(config *Config) error {
    	fmt.Printf("put key: %s ; newValue: %s; version:%s \n", config.Key, config.Value, config.Version)
    	return nil
    }
    //通知key的删除事件
    func (t *TestListenter) HandleDeleteEvent(config *Config) error {
    	fmt.Printf("delete app:%s, profile:%s, key: %s \n", config.App, config.Profile, config.Key)
    	return nil
    }

    3）、 注册监听器，支持一个监听器监听多个Key
     cnt2Service.RegisterListener(&TestListenter{}, "my_config_1", "test_config")

    4)、  移除监听器
     cnt2Service.UnRegisterListener(&TestListenter{}, "my_config_1", "test_config")

3. 测试：
	go test -run Test_GetConfig
	
	go test -run Test_NewClient