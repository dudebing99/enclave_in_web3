package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	"io/ioutil"
	"strings"
	"time"
)

type responseWriter struct {
	gin.ResponseWriter
	b *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	//向一个bytes.buffer中写一份数据来为获取body使用
	w.b.Write(b)
	//完成gin.Context.Writer.Write()原有功能
	return w.ResponseWriter.Write(b)
}

// 日志记录到文件
func LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 请求 Body
		requestBody, err := c.GetRawData()
		if err != nil {
			panic(err.Error())
		}

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody)) // 关键点
		writer := responseWriter{
			c.Writer,
			bytes.NewBuffer([]byte{}),
		}
		c.Writer = writer

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求 IP
		clientIP := c.ClientIP()

		//responseHeader := c.Writer.Header()
		//responseBodySize := c.Writer.Size()
		responseBody := writer.b.String()

		debug := viper.GetBool("gateway.debug")
		if !debug {
			if strings.Contains(reqUri, "/key/") {
				responseBody = "******"
			}
		}

		// 日志格式
		glog.Infof("| %3d | %13v | %15s | %s | %s | %s -> %s\n",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
			string(requestBody),
			responseBody,
		)
	}
}

// 日志记录到 MongoDB
func LoggerToMongo() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志记录到 ES
func LoggerToES() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志记录到 MQ
func LoggerToMQ() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
