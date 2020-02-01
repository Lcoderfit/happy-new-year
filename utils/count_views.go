package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

)

// 刷新浏览量
func CountFlushView() (string, error) {
	if viewNumCache.IsExist("fv") {
		err := viewNumCache.Incr("fv")
		if err != nil {
			return "", err
		}
	} else {
		err := viewNumCache.Put("fv", "1", 30*24*60*60*time.Second)
		if err != nil {
			return "", err
		}
	}
	fv := string(viewNumCache.Get("fv").([]byte))
	return fv, nil
}

// 独立访客量
func CountUniqueView(w http.ResponseWriter, r *http.Request) (int, int, error) {
	// 启动session
	sess, _ := viewNumSession.SessionStart(w, r)
	defer sess.SessionRelease(w)

	// 按天存储UV
	now := time.Now()
	uvField := fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day())
	dailyUv := map[string]int{}

	if sess.Get("uv") == nil {
		fmt.Println("session == nil")
		if viewNumCache.IsExist("uv") {
			fmt.Println("add UV")
			// 增加总UV
			err := viewNumCache.Incr("uv")
			if err != nil {
				log.Println(err)
				return 0, 0, err
			}

			// 增加该天的UV
			err = AddDailyVal("dailyuv", dailyUv, uvField)
			if err != nil {
				log.Println(err)
				return 0, 0, err
			}
		} else {
			fmt.Println("set UV")
			// 设置总UV
			err := viewNumCache.Put("uv", 1, 30*24*60*60*time.Second)
			if err != nil {
				log.Println(err)
				return 0, 0, err
			}
			// 存储每一天的UV
			dailyUv[uvField] = 1
			// 将dailyUv存入redis
			err = StoreKvToRedis("dailyuv", dailyUv, 30*24*60*60*time.Second)
			if err != nil {
				log.Println(err)
				return 0, 0, err
			}
		}
	} else {
		err := GetViewsFromRedis("dailyuv", dailyUv)
		if err != nil {
			log.Println(err)
			return 0, 0, err
		}
	}
	if err := sess.Set("uv", time.Now().Unix()); err != nil {
		log.Println(err)
		return 0, 0, err
	}
	if sess.Get("uv") != nil {
		fmt.Println("session != nil ")
	}

	// 获取总uv
	uvStr := string(viewNumCache.Get("uv").([]byte))
	uvSum, err := strconv.Atoi(uvStr)
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}
	// 获取当天的uv
	todayUv := dailyUv[uvField]

	return uvSum, todayUv, nil
}

// 页面浏览量
func CountPageView(w http.ResponseWriter, r *http.Request) (int64, int, int, error) {
	// 启动session
	sess, _ := viewNumSession.SessionStart(w, r)
	defer sess.SessionRelease(w)
	// 按天存储PV
	now := time.Now()
	pvField := fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day())
	// 用于保存每天的PV数
	dailyPv := map[string]int{}

	// 会话为空，则表示该用户之前没访问过该网页，PV增加1，
	// 会话不为空，且该用户距离上次访问网页时间间隔半个小时，PV增加1
	if sess.Get("pv") != nil {
		//// 获取存储每一天PV的字典
		//dailyPvByte := viewNumCache.Get("dailypv").([]byte)
		//err := json.Unmarshal(dailyPvByte, &dailyPv)
		//if err != nil {
		//	log.Println(err)
		//	return -1, 0, 0, err
		//}
		// 判断是否重新累计PV
		now := time.Now().Unix()
		lastTime := sess.Get("pv").(int64)
		if (now - lastTime) > 30*60 {
			// 增加总PV
			viewNumCache.Incr("pv")
			// 增加日PV
			err := AddDailyVal("dailypv", dailyPv, pvField)
			if err != nil {
				log.Println(err)
				return -1, 0, 0, err
			}
			// 更新浏览的时间
			sess.Set("pv", now)
		} else {
			// 获取dailyPv
			err := GetViewsFromRedis("dailypv", dailyPv)
			if err != nil {
				log.Println(err)
				return -1, 0, 0, err
			}
		}
	} else {
		sess.Set("pv", time.Now().Unix())

		if viewNumCache.IsExist("pv") {
			// 增加总pv量
			err := viewNumCache.Incr("pv")
			if err != nil {
				log.Println(err)
				return -1, 0, 0, err
			}
			// 增加日PV
			err = AddDailyVal("dailypv", dailyPv, pvField)
			if err != nil {
				log.Println(err)
				return -1, 0, 0, err
			}
		} else {
			// 保存总PV
			err := viewNumCache.Put("pv", "1", 30*24*60*60*time.Second)
			if err != nil {
				return -1, 0, 0, err
			}

			// 保存每一天的PV
			dailyPv[pvField] = 1
			err = StoreKvToRedis("dailypv", dailyPv, 30*24*60*60*time.Second)
			if err != nil {
				log.Println(err)
				return -1, 0, 0, err
			}
		}
	}
	// 获取总PV
	pvStr := string(viewNumCache.Get("pv").([]byte))
	pv, _ := strconv.Atoi(pvStr)
	// 获取每一天的PV
	todayPv := dailyPv[pvField]
	viewTimeUnix := sess.Get("pv").(int64)
	return viewTimeUnix, pv, todayPv, nil
}

// 将每一天的PV或UV对应的map从redis中取出
// ？？？？？？？？？？？？？？？？？？？？？val参数出入指针才会改变原始值
func GetViewsFromRedis(key string, val map[string]int) error {
	dailyViewsByte := viewNumCache.Get(key).([]byte)
	if err := json.Unmarshal(dailyViewsByte, &val); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// 将非string的类型作为key存储到redis中
func StoreKvToRedis(key string, val map[string]int, expire time.Duration) error {
	valStr, err := json.Marshal(val)
	if err != nil {
		log.Println(err)
		return err
	}
	err = viewNumCache.Put(key, valStr, expire)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// 将作为redis的val的map类型取出，并对map的值进行累加后存入redis
func AddDailyVal(key string, val map[string]int, today string) error {
	// 取出键对应的每天的浏览量map
	err := GetViewsFromRedis(key, val)
	if err != nil {
		log.Println(err)
		return err
	}
	val[today]++

	// 将map重新存入redis
	return StoreKvToRedis(key, val, 30*24*60*60*time.Second)
}

// 将time结构体转换成timeStr
func ChangeDateToStr(tim time.Time) string {
	timeByte, err := json.Marshal(tim)
	if err != nil {
		log.Fatal(err)
	}
	timeStr := string(timeByte)
	return timeStr[1:11] + "-" + timeStr[12:20]
}

// 将时间戳转换为日期
func SwitchTimeStampToData(unixTime int64) string {
	timeStr := time.Unix(unixTime, 0).Format("2006-01-02 15:04:05")
	return timeStr
}

// 打印带时间的信息
func PrintDateAndMessage(args ...interface{}) {
	timeStr := ChangeDateToStr(time.Now())
	fmt.Println(timeStr + " : ", args)
}