package main

import (
	"fmt"
	aikuaimonitor "ikuai-mointor/internal"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jakeslee/ikuai/action"
	"github.com/urfave/cli/v2"
)

func printHelp() {
	cli.AppHelpTemplate = fmt.Sprintf(`%s
	访问http://host:7575查看
	配置文件在config/conf.ini
	配置文件格式

	[爱快路由器名称]
	  user=爱快登录用户名
	  password=爱快登录密码
	  url=访问爱快的地址

	`, cli.AppHelpTemplate)
}

func main() {
	aikuaimonitor.InitMonitor()
	printHelp()
	app := &cli.App{
		Name:  "爱快监控",
		Usage: "监控爱快的连接和总带宽",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "command",
				Aliases: []string{"c"},
				Usage:   "命令模式，即通过命令行查看数据",
			},
			&cli.IntFlag{
				Name:    "time, t",
				Aliases: []string{"t"},
				Value:   5,
				Usage:   "命令模式多久循环取一次数据集",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("help") {
				printHelp()
				return nil
			}
			if c.Bool("command") {
				t := c.Int("time")
				aikuaimonitor.ScheduleMonitorInterface(t)
				return nil
			}
			Server()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func EntryHtml() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") || strings.HasPrefix(c.Request.URL.Path, "/static") || strings.HasPrefix(c.Request.URL.Path, "/favicon.ico") {
			c.Next()
		} else {
			c.HTML(200, "index.html", nil)
		}
	}
}

func Server() {
	var cacheData = map[string]interface{}{}
	var mu sync.Mutex
	var timeLock sync.Mutex
	isStop := false
	var lastTime time.Time = time.Now()
	ticker := time.NewTicker(3 * time.Second)
	task := func(ticker *time.Ticker) {
		for {
			select {
			case <-ticker.C:
				log.Println("ticker exec")
				go func() {
					res := aikuaimonitor.Monitor.GetMonitorInterface()
					if res != nil {
						mu.Lock()
						old, isHave := cacheData["interface"].(map[string]*action.ShowMonitorInterfaceResult)
						mu.Unlock()
						if isHave {
							for key, val := range res {
								if val != nil {
									old[key] = val
								}
							}
							mu.Lock()
							cacheData["interface"] = old
							mu.Unlock()
						} else {
							mu.Lock()
							cacheData["interface"] = res
							mu.Unlock()
						}
					}
				}()
				go func() {
					res := aikuaimonitor.Monitor.GetAllMonitorLan(false)
					if res != nil {
						mu.Lock()
						old, isHave := cacheData["lanv4"].(map[string]*action.ShowMonitorResult)
						mu.Unlock()
						if isHave {
							for key, val := range res {
								if val != nil {
									old[key] = val
								}
							}
							mu.Lock()
							cacheData["lanv4"] = old
							mu.Unlock()
						} else {
							mu.Lock()
							cacheData["lanv4"] = res
							mu.Unlock()
						}
					}
				}()
				go func() {
					res := aikuaimonitor.Monitor.GetAllMonitorLan(true)
					if res != nil {
						mu.Lock()
						old, isHave := cacheData["lanv6"].(map[string]*action.ShowMonitorResult)
						mu.Unlock()
						if isHave {
							for key, val := range res {
								if val != nil {
									old[key] = val
								}
							}
							mu.Lock()
							cacheData["lanv6"] = old
							mu.Unlock()
						} else {
							mu.Lock()
							cacheData["lanv6"] = res
							mu.Unlock()
						}
					}
				}()
			}
		}
	}
	go task(ticker)

	go func() {
		for {
			// 每10s检测一次
			time.Sleep(10 * time.Second)
			timeLock.Lock()
			elapsed := time.Since(lastTime)
			log.Println("check stop task")
			if isStop {
				return
			}
			if elapsed >= 1*time.Minute {
				log.Println("sotp current task")
				isStop = true
				ticker.Stop()
			}
			timeLock.Unlock()
		}
	}()

	tickerMiddleWare := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			lastTime = time.Now()
			if isStop {
				isStop = false
				log.Println("start new task")
				ticker = time.NewTicker(3 * time.Second)
				go task(ticker)
			}
		}
	}

	// 创建Gin路由
	r := gin.Default()
	r.LoadHTMLGlob("./front/dist/*.html")
	r.Use(EntryHtml())
	// 定义API路由
	api := r.Group("/api")
	api.Use(tickerMiddleWare())
	api.GET("/interface", func(c *gin.Context) {
		// 构建JSON响应
		response := gin.H{
			"message": "Success",
			"code":    0,
			"data":    cacheData["interface"],
		}
		// 返回JSON响应
		c.JSON(http.StatusOK, response)
	})
	api.GET("/lanv4", func(c *gin.Context) {
		// 构建JSON响应
		response := gin.H{
			"message": "Success",
			"code":    0,
			"data":    cacheData["lanv4"],
		}
		// 返回JSON响应
		c.JSON(http.StatusOK, response)
	})
	api.GET("/lanv6", func(c *gin.Context) {
		// 构建JSON响应
		response := gin.H{
			"message": "Success",
			"code":    0,
			"data":    cacheData["lanv6"],
		}
		// 返回JSON响应
		c.JSON(http.StatusOK, response)
	})
	r.Static("/static", "./front/dist/static")
	r.StaticFile("/favicon.ico", "./front/dist/favicon.ico")
	r.Run(":7575")
}
