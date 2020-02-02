# <center>花样新年祝福</center>>

## 一、简介

该系统用于新年的祝福。实际访问地址：

[http://blessing.lcoderfit.com/]: 

## 一、详细功能介绍

### 1.1 个性化祝福

* 通过在基地址后面添加URL参数可以实现对特定的人进行定制化祝福，并且参数可在后端添加扩展，参数选项如下：

  | URL参数 | 映射到前端显示的实际值 | 示例URL                           |
  | ------- | ---------------------- | --------------------------------- |
  | LB      | 老爸                   | blessing.lcoderfit.com/?name=LB   |
  | LM      | 老妈                   | blessing.lcoderfit.com/?name=LM   |
  | LJ      | 老姐                   | blessing.lcoderfit.com/?name=LJ   |
  | JF      | 姐夫                   | blessing.lcoderfit.com/?name=JF   |
  | GFGG    | 姑父姑姑               | blessing.lcoderfit.com/?name=GFGG |
### 1.2 播放音乐

浏览器直接输入URL会出现音乐不会自动播放的现象，但是在微信中直接访问链接则会自动播放音乐

### 1.3 技术栈及学习的新知识

#### 1.3.1 采用gin框架进行开发，其中主要涉及HTML的template和static设置，在mapName变量中添加URL参数和对应的祝福人。

```go
package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	mapName := map[string]interface{}{
		"LB": "老爸",
		"LM": "老妈",
		"LJ": "老姐",
		"JF": "姐夫",
		"GFGG": "姑父姑姑",
	}
	router := gin.Default()
    // 第一个参数为URL访问静态资源的地址，第二个参数为静态资源文件路径
    // 与StaticFS()函数作用类似，StaticFS用于一个完整的文件目录，StaticFile用于单个文件
	router.Static("/static", "./static")
	// router.StaticFS("/staticfs", http.Dir("static"))
    // 设置模板文件的路径
	router.LoadHTMLGlob("template/*")
	router.GET("/", func(c *gin.Context) {
		shortName := c.Query("name")
		realName := mapName[shortName]
		log.Println(shortName, realName)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"realName": realName,
		})
	})
	router.Run(":8050")
}
```

#### 1.3.2 在服务器对nohup.out进行重定向

```
当nohup.out文件中的内容过多时，不方便查看启动项目时的输出内容，
使用tail -n 10 nohup.out > a.out将nohup.out的末尾十行重定向到a.out文件中
然后cat a.out > nohup.out再将需要的内容从定向到nohup.out中
```

#### 1.3.3 Go如何使用redis存储高级数据结构

```
beego自带的cache只能将非string类型转化为string(json.Marshal)，然后再作为val存入redis
取出的时候先用.([]byte)断言，然后用json.UnMarshal转化为原生数据类型

存储高级数据结构：用/github.com/go-redis/redis模块
```

#### 1.3.4 Go删除map的键

```
m := map[string]int{
	"a": 1,
}
delete(m, "a")
```

清空map

```
for key := range m {
	delete(m, key)
}
```

#### 1.3.5 热启动和热升级

* 热启动使用beego的bee模块，修改源码之后会自动更新

```
go get github.com/beego/bee
```

* bee run即可启动项目

* 热升级

  ```
  to be continue:
  https://www.jianshu.com/p/b85a266a7515
  ```

#### 1.3.6 计算UV、PV、FV

* UV: 独立访客量

  ```
  第一次访问，建立session，将cookie发送并保存到本地客户端，下一次访问，判断是否存在session变量，如果存在证明已经访问过，不再累计，同一个客户端访问间隔24小时重新计算一次
  ```

* PV: 页面浏览量

  ```
  第一次访问，建立session，保存当前访问时间，并将cookie发送并保存到本地客户端，下一次访问，判断session变量是否存在，不存在则累计一次PV，存在则取出上一次访问的时间，如果已经过去了30分钟，则重新累计。
  ```

* FV: 刷新浏览量

  ```
  任何客户端每刷新一次则累计一次
  ```

* beego的session模块运行流程

  ```
  第一次向服务器发送请求，建立session，保存会话变量到设置的session引擎（redis），然后将cookie返回并保存到客户端。
  第二次访问携带上cookie，通过cookie查询session
  
  注意：不同的客户端，例如PC端浏览器（1），PC端微信（2），端1第一次调用this.GetSession("key")不存在，之后建立session，生成并保存key变量，返回cookie到端1，端2第一次访问，同样调用this.GetSession("key"),此时端2并没有session，所以端2获取的值为nil
  ```


#### 1.3.7 传递不定参

```
func main() {	
	sarr := []string{"a", "b", "c"}
	test4(sarr...)
}

// 将不定参作为参数
func test4(args ...string) {
	for _, arg := range args {
		fmt.Println(arg)
	}
}
```

#### 1.3.8 分析UV和PV等的工具

```
GA: Google Analises
```

#### 1.3.9 日志打印输出行号

```
// 设置日志输出行号
log.SetFlags(log.Lshortfile | log.LstdFlags)
```

#### 1.3.10 字符串和数字互相转换

```
// string -> int
s := "1234"
i := strconv.Atoi(s)

// int -> string
i := 1234
s := strconv.Itoa(i)
```

#### 1.3.11 json.Marshal和json.Unmarshal

```
m := map[string]int{
	"a": 1,
}
// 返回([]byte, error)
sm, err := json.Marshal(m)
if err != nil {
	log.Println(err)
	return
}
// 通过string()将[]byte转换为string
fmt.Println(string(sm))

// 参数为([]byte, &res)
// 记得传递的是一个指针
err = json.Unmarshal(sm, &m)
if err != nil {
	log.Println(err)
}
fmt.Println(m)
```

### 1.4 目录树形图

![image-20200123105659978](C:\Users\Lcoderfit\AppData\Roaming\Typora\typora-user-images\image-20200123105659978.png)

### 1.5 项目运行方式

* 用git下载该项目

  ```markdown
  git clone github.com/Lcoderfit/happy-new-year
  ```

* 该项目引入beego的bee模块，可以实现项目热启动，在项目根目录下打开cmd，输入以下命令

  ```markdown
  bee run
  ```

* 本地启动项目时，在同一局域网下访问

  ```
  1.router.Run("0.0.0.0:8050") // 0.0.0.0表示对该局域网下的其他主机进行广播
  2.本地打开cmd，用ipconfig查看该主机在该局域网下的IP
  3.该局域网下（不同主机连上同一wifi或热点）其他主机打开浏览器，输入"第二步查看的IP:8050"
  ```

### 1.6 开发过程中的错误

#### 1.6.1 invalid memory address or nil pointer dereference

* 报错原因

  ```
  声明了一个指针，但是未给该指针分配内存地址就直接进行赋值
  ```

* 报错代码

  ```
  /* 
  声明了全局变量viewNumSession,并在其他地方用到了该变量，但是在给其他地方调用前未给该全局变量分配内存空间。
  */
  var (
  	// 会话变量
  	viewNumSession *session.Manager
  )
  
  func init() {
  	BuildSession()
  }
  
  // 建立会话
  func BuildSession() {
  	var err error
  	// 下面一行错误，因为在BuildSession()函数中又新创建了一个viewNumSession变量，该变量与
  	// 一开始创建的全局变量不相同，并未给全局变量viewNumSession分配空间
  	viewNumSession, err := session.NewManager("redis", sessionConfig)
  	if err != nil {
  		log.Fatal("BuildSession->NewManager Error: ", err)
  		return
  	}
  	go viewNumSession.GC()
  }
  ```

* 解决办法

  ```
  将：
  viewNumSession, err := session.NewManager("redis", sessionConfig)
  改写成：
  viewNumSession, err = session.NewManager("redis", sessionConfig)
  ```

#### 1.6.2 interface conversion: interface {} is []uint8, not map[string]int

* 报错原因

  ```
  接口返回的是[]int8类型，但是赋值的变量为map[string]int类型
  ```

* 报错代码

  ```
  dailyUv := viewNumCache.Get("dailyuv").(map[string]int)
  ```

* 解决办法

  ```
  改成：
  dailyUv := map[string]int{}
  dailyUvByte := viewNumCache.Get("dailyuv").([]byte)
  err := json.Unmarshal(dailyUvByte, dailyUv)
  ....
  ```

#### 1.6.3 invalid character 'm' looking for beginning of value

* 报错原因

  ```
  存入redis时val对应的是非string类型的值，从redis中取出之后转码错误
  ```

* 报错代码

  ```
  
  ```

* 解决办法

  ```
  存入redis时需要将val转换成string，用json.Marshal
  ```

#### 1.6.4 json: Unmarshal(non-pointer map[string]int)

* 报错原因

  ```
  调用json.Unmarshal(data, type)时，第一个参数是[]byte类型，第二个参数需要传入一个map[string]int类型的指针，但是第二个参数未传入指针
  ```

* 报错代码

  ```
  dailyUvByte := viewNumCache.Get("dailyuv").([]byte)
  dailyUv := map[string]int{}
  // 下面一行错误，第二个参数需传入一个指针
  err = json.Unmarshal(dailyUvByte, dailyUv)
  if err != nil {
  	log.Println(err)
  	return 0, 0, err
  }
  ```

* 解决办法

  ```
  将：
  err = json.Unmarshal(dailyUvByte, dailyUv)
  改成：
  err = json.Unmarshal(dailyUvByte, &dailyUv)
  ```

#### 1.6.5 从redis中取出val，累加之后忘记存回redis

#### 1.6.6 git push错误：failed to push some refs to 

* 报错信息

  ```
  $ git push -u origin master
  
  To git@github.com:yangchao0718/cocos2d.git
  
   ! [rejected]        master -> master (non-fast-forward)
  
  error: failed to push some refs to 'git@github.com:yangchao0718/cocos2d.git
  
  hint: Updates were rejected because the tip of your current branch is behin
  
  hint: its remote counterpart. Integrate the remote changes (e.g.
  
  hint: 'git pull ...') before pushing again.
  ```

* 解决办法

  ```
  git pull --rebase origin master
  git push -u origin master -f
  git rebase --abort
  ```

  