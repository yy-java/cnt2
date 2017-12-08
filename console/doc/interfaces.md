#新配置中心接口文档

## 一、通用说明 ##
### 1. 通用返回
	{
		code:0  //状态码
		msg: "success" //状态描述，有些状态不一定有
		data:1  // 数据信息，详细见每个接口说明
	}

### 2. 状态码列表

	0：		成功
	1： 		服务异常
	2：		参数无效
	3：		未认证，未登录
	4：		权限不足
	5：		记录已存在
	6：		记录不存在
	7：		版本过时

### 3. 环境
	测试环境：http://api-test.xx.com/
	正式环境：http://api.xx.com/


## 二、APP相关接口说明 ##

## 1. 创建App

	URI: /app/create?app=app10&appType=1&name=test-name2&charger=dyf
	Params：
		app: app标记
		name: 名称
		appType: 0 服务端 9 app客户端
		charger: 负责人名字，不是uid

	Return:
		通用返回，可能的code值： 0 1 2 3 5

## 2. 删除App
	URI: /app/del/{app}
	Params：
		app: app标记 Restful风格

	Return:
		通用返回，可能的code值： 0 1 3

## 3. 更新App
	URI: /app/update?app=&appType=&name=&charger=
	Params：
		app: app标记
		name: 名称
		appType: 0 服务端 9 app客户端
		charger: 负责人名字，不是uid

	Return:
		通用返回，可能的code值： 0 1 2  3  5

## 4. 当前用户的App列表（特权用户user_auth中app字段值为all）
	URI: /app/list
	Params：
		 无

	Return:
		通用返回，可能的code值： 0 1 3
		data：
		[
		    {
		      "App": "app1",
		      "AppType": 1,
		      "Name": "test",
		      "Charger": "dyf",
		      "CreateTime": "2017-07-13T18:14:22+08:00"
		    }
		]

## 三、用户授权相关接口说明 ##

## 1. 用户授权

	URI: /userauth/create?app=app10&permission=1&uid=&username=
	Params：
		app: app标记
		permission: 1 开发 9 管理
		uid: 用户uid
		username: 用户昵称

	Return:
		通用返回，可能的code值： 0 1 2 3 4

## 2. 当前App的用户列表
	URI: /userauth/list/app1
	Params：
		 app:app1 restful风格

	Return:
		通用返回，可能的code值： 0 1 3
		data：
		[
		   {
		      "Id": 1,
		      "Uid": 50002304,
		      "Uname": "dd",
		      "App": "app1",
		      "Permission": 9 //1 开发 9 管理
		    }
		]

## 3. 删除授权
	URI: /userauth/del?id=
	Params：
		id: 授权记录id

	Return:
		通用返回，可能的code值： 0 1 2 3 4

## 4. 更新授权
	URI: /userauth/update_permission?id=&permission=1
	Params：
		id: 授权列表
		permission: 1 开发 9 管理

	Return:
		通用返回，可能的code值： 0 1 2 3 4


## 四、 配置相关接口说明 ##

## 1. 创建配置

	URI: /config/create?app=app2&profile=test&key=key&value=val&validator=validator&description=desc
	Params：
		app: app标记
		profile: 1 开发 9 管理
		key: 用户uid
		value: 用户昵称
		validator: 校验
		description: 描述

	Return:
		通用返回，可能的code值： 0 1 2 3 4

		data: 配置的ID

## 2. 更新配置

	URI: /config/update?app=app2&profile=test&key=key&value=val&validator=validator&description=desc
	Params：
		app: app标记
		profile: profile
		key: 用户uid
		value: 用户昵称
		validator: 校验
		description: 描述

	Return:
		通用返回，可能的code值： 0 1 2 3 4 5

## 3. 查询配置

	URI: /config/queryWithKey?app=app1&profile=dev&key=
	Params：
		app: app标记
		profile: profile 环境名
		key: key


	Return:
		通用返回，可能的code值： 0 1 2 3 4

		data:
			 {
			    "Id": 1,
			    "App": "app1",
			    "Profile": "dev", //环境名
			    "Key": "k1",
			    "Value": "sl",
			    "Version": 1,  // 版本号
			    "PublishedValue": "3230230423", //最近已发布配置
			    "PublishedVersion": 1, //最近已发布版本号
			    "Validator": "", //验证配置的js脚本，前端执行
			    "Modifier": "123123",
			    "ModifyTime": "2017-07-12T15:08:12+08:00",
			    "Description": "111", //修改描述
			    "Approver": "123123", // 审核人
			    "ApproveType": 1 // 0 未审核 1 已审核
			  }

## 4. 查询配置，不带Key值

	URI: /config/query?app=app1&profile=dev
	Params：
		app: app标记
		profile: profile 环境名


	Return:
		通用返回，可能的code值： 0 1 2 3 4

		data:
			 [{
			    "Id": 1,
			    "App": "app1",
			    "Profile": "dev", //环境名
			    "Key": "k1",
			    "Value": "sl",
			    "Version": 1,  // 版本号
			    "PublishedValue": "3230230423", //最近已发布配置
			    "PublishedVersion": 1, //最近已发布版本号
			    "Validator": "", //验证配置的js脚本，前端执行
			    "Modifier": "123123",
			    "ModifyTime": "2017-07-12T15:08:12+08:00",
			    "Description": "111", //修改描述
			    "Approver": "123123", // 审核人
			    "ApproveType": 1 // 0 未审核 1 已审核
			  }]

## 5. 删除配置
	URI: /config/del?app=app1&profile=dev1&key=key1
	Params：
		app: app标记
		profile: profile 环境名
		key: key

	Return:
		通用返回，可能的code值： 0 1 2 3 4

## 6. 审核

	URI: /config/approve?app=app1&profile=dev1&key=key1
	Params：
		app: app标记
		profile: profile 环境名
		key: key

	Return:
		通用返回，可能的code值： 0 1 2 3 4

## 7. 列出某个app下的所有环境配置

	URI: /config/profiles/app1
	Params：
		app: app标记 restful风格参数

	Return:
		data: //返回配置名字列表
			[
			    "dev",
			    "test"
			]

## 8. 回滚某个配置

	URI: /config/rollback?app=app1&profile=dev1&key=key1&version=
	Params：
		app: app标记
		profile: profile 环境名
		key: key
		version: 回滚到的版本

	Return:
		通用返回，可能的code值： 0 1 2 3 4


## 五、 节点接口说明 ##

## 1. 在线节点列表

	URI: /node/list?app=&profile=
	Params：
		app: app
		profile: 环境名

	Return:
		通用返回，可能的code值： 0 1 2 3 4

		data:
		{
			"xxxx":{
					NodeId:"11"
					App: "app"
					Profile： "profile"
					Pid: "pid"
					RegisterTime： 1111111111 //注册时间
				}
		}


## 六、 配置相关接口说明 ##

## 1. 发布

	URI: /publish?id=&nodes=11,22,33,44
	Params：
		id: 配置信息的id
		nodes: 需要发布的节点id，多个以英文逗号隔开

	Return:
		通用返回，可能的code值： 0 1 2 3 4
