#### go语言练习项目

#### 1、亿级数据排序实现

```
	main.go文件：
	// 从文件中提取密码
	GetQQPasswordFromFile()
	// 二、密码排序
	QQPasswordSort()
	// 三、计算密码次数
	QQPasswordTimes()
	// 四、按照密码次数排序
	QQPasswordDictionary()
	上述四个方法实现了对文件中的文件进行提取，并进行数据排序
```

#### 2、自定义协议数据传输实现

```
Utils/sortInterface.go文件
	实现了自定义数据传输协议，可以根据不同类型数据传输，并进行数据排序处理；
	在服务器收到客户端请求，对数据做对应的处理并按照请求协议返回数据功能实现
```

#### 3、分布式数据传输

```
Utils/sortInterface.go文件
	在自定义协议数据传输实现的基础上继续开发，并能实现数据的分布式传输，
	实现master把文件切割，发送给不同的slave，在slave把数据处理并返回给master后进行数据归并
```

