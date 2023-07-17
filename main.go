package main

import (
	"fmt"
	aikuaimonitor "ikuai-mointor/internal"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
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
	// 创建Gin路由
	r := gin.Default()
	r.LoadHTMLGlob("./front/dist/*.html")
	r.Use(EntryHtml())
	// 定义API路由
	api := r.Group("/api")
	api.GET("/interface", func(c *gin.Context) {
		// 构建JSON响应
		response := gin.H{
			"message": "Success",
			"code":    0,
			"data":    aikuaimonitor.Monitor.GetMonitorInterface(),
		}
		// 返回JSON响应
		c.JSON(http.StatusOK, response)
	})
	r.Static("/static", "./front/dist/static")
	r.StaticFile("/favicon.ico", "./front/dist/favicon.ico")
	r.Run(":7574")
}
