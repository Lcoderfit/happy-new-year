package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"happy-new-year/utils"
	"log"
	"net/http"
	"time"
)

/*
import (
	"net/http"
	"log"
	"io/ioutil"
	"fmt"
)


func main() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("index.html")
		if err != nil {
			log.Fatal("test1->ParseFile Error: ", err)
			return
		}
		t.Execute(res, nil)
	})
	log.Println("Start Server")
	http.ListenAndServe(":80", nil)
	log.Println("Server end")
}
*/

func main() {
	// 设置日志输出行号
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	mapName := map[string]interface{}{
		"LB":   "老爸",
		"LM":   "老妈",
		"LJ":   "老姐",
		"JF":   "姐夫",
		"GFGG": "姑父姑姑",
		"FG":   "峰哥",
		"PJ":   "萍姐",
		"DG":   "德哥",
		"SSSS": "叔叔婶婶",
		"JG":   "俊哥",
		"XY":	"小姨",
		"XS":   "小施",
	}
	router := gin.Default()
	router.Static("/static", "./static")
	//router.StaticFS("/staticfs", http.Dir("static"))
	router.LoadHTMLGlob("template/*")
	router.GET("/", func(c *gin.Context) {
		// 获取URl参数
		shortName := c.Query("name")
		realName := mapName[shortName]

		// 打印带时间的信息
		utils.PrintDateAndMessage(shortName, realName)
		// 计算刷新访问量
		fv, err := utils.CountFlushView()
		if err != nil {
			log.Fatal("CountFlushView Error", err)
		}
		utils.PrintDateAndMessage("FV: ", fv)

		// 计算页面访问量
		viewTimeUnix, pv, todayPv, err := utils.CountPageView(c.Writer, c.Request)
		if err != nil {
			log.Fatal("CountPageView Error: ", err)
		}
		viewTime := utils.SwitchTimeStampToData(viewTimeUnix)
		nowTime := utils.SwitchTimeStampToData(time.Now().Unix())
		// 重新累计PV所需的时间
		reCountTime := 30 - (time.Now().Unix() - viewTimeUnix)/60
		str := fmt.Sprintf("PVSum: %d, TodayPV: %d, RecountTime: %dM, ViewTime: %s, NowTime: %s",
			pv, todayPv, reCountTime, viewTime, nowTime)
		utils.PrintDateAndMessage(str)

		// 计算独立访问量
		// *gin.Context.Writer: http.ResponseWriter
		// *gin.Context.Request: *http.Request
		uvSum, todayUv, err := utils.CountUniqueView(c.Writer, c.Request)
		if err != nil {
			log.Fatal("CountUniqueView Error")
		}
		str = fmt.Sprintf("UVSum: %d, TodayUV: %d", uvSum, todayUv)
		utils.PrintDateAndMessage(str)

		c.HTML(http.StatusOK, "index.html", gin.H{
			"realName": realName,
		})
	})
	router.Run("0.0.0.0:8050")
}
