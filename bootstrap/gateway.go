package bootstrap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"api-gateway/internal/api"
	"api-gateway/internal/services"

	"github.com/cockroachdb/pebble"
	"github.com/gin-gonic/gin"
)

// GatewayApp结构体用于网关转发相关的初始化和配置
type GatewayApp struct {
	Router            *gin.Engine
	API               *api.APIController
	PebbleDB          *pebble.DB
	downstreamService services.DownstreamServiceImpl
}

// NewGatewayApp创建并初始化用于网关转发的应用实例
func NewGatewayApp() *GatewayApp {
	return &GatewayApp{}
}

// Initialize初始化网关转发应用的各种组件
func (ga *GatewayApp) Initialize() {
	var err error
	pebblePath := filepath.Join(DB_PATH, "pebble")
	ga.PebbleDB, err = pebble.Open(pebblePath, nil)
	if err != nil {
		fmt.Printf("Error opening PebbleDB: %v", err)
		return
	}

	ga.downstreamService = services.NewDownstreamService()
	ga.Router = gin.Default()
	// ga.Router.Use(middleware.NewMiddleware().Wrap)
}

// SetupRoutes设置网关转发应用的路由
func (ga *GatewayApp) SetupRoutes() {
	ga.Router.Any("/:name/*path", func(c *gin.Context) {
		name := c.Param("name")
		fmt.Println("name =>", name)
		parts := c.Param("path")
		if len(parts) > 0 {
			serviceName := parts[1:]
			service, err := ga.downstreamService.GetByName(context.Background(), serviceName)
			if err == nil && service.URL != "" {
				backendURL, err := url.Parse(service.URL)
				if err == nil {
					// 创建一个自定义的转发器，这里可以根据需要调整转发逻辑
					proxyClient := NewProxyClient(backendURL)
					// 记录入栈流量
					// ga.Logger.LogTraffic(c.Request, true)
					ga.forwardRequest(c, proxyClient)
					// 记录出栈流量
					// ga.Logger.LogTraffic(c.Request, false)
					return
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{"message": "Service not found"})
	})

	// 没有找到对应服务，返回404错误
	ga.Router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "404"})
	})
}

// Run启动网关转发应用
func (ga *GatewayApp) Run() {
	fmt.Println("API Gateway started on :8080")
	err := ga.Router.Run(":8080")
	if err != nil {
		fmt.Printf("Error starting gateway server: %v", err)
	}
}

// closePebbleDB关闭pebble数据库
func (ga *GatewayApp) Close() {
	if ga.PebbleDB != nil {
		ga.PebbleDB.Close()
	}
}

// storeAPIInfoInPebble将API信息存储到pebble数据库中
func (ga *GatewayApp) storeAPIInfoInPebble(key, value string) error {
	return ga.PebbleDB.Set([]byte(key), []byte(value), pebble.Sync)
}

// NewProxy创建一个自定义的请求转发器，这里只是一个简单示例，可以根据实际需求扩展
func NewProxy(u *url.URL) *http.Transport {
	return &http.Transport{
		Proxy: func(request *http.Request) (*url.URL, error) {
			return u, nil
		},
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial(network, u.Host)
		},
	}
}

// NewProxyClient创建一个自定义的代理客户端
func NewProxyClient(targetURL *url.URL) *http.Client {
	transport := &http.Transport{
		Proxy: func(request *http.Request) (*url.URL, error) {
			return targetURL, nil
		},
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial(network, targetURL.Host)
		},
	}
	return &http.Client{Transport: transport}
}

type RequestInfo struct {
	Method  string
	URL     string
	Headers map[string][]string
	Body    []byte
}

type ResponseInfo struct {
	Method     string
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

// forwardRequest用于转发请求
func (ga *GatewayApp) forwardRequest(c *gin.Context, client *http.Client) {
	req := c.Request.Clone(c.Request.Context())
	requestInfo, err := extractRequestInfo(req)
	if err != nil {
		log.Printf("Error extracting request info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract request info"})
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	responseInfo, err := extractResponseInfo(resp, requestInfo.Method)
	if err != nil {
		log.Printf("Error extracting response info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract response info"})
		return
	}

	err = storeRequestInfo(ga, requestInfo)
	if err != nil {
		log.Printf("Error storing request info in Pebble: %v", err)
	}

	err = storeResponseInfo(ga, responseInfo)
	if err != nil {
		log.Printf("Error storing response info in Pebble: %v", err)
	}

	for k, vv := range resp.Header {
		for _, v := range vv {
			c.Writer.Header().Add(k, v)
		}
	}
	c.Writer.WriteHeader(resp.StatusCode)
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		fmt.Println("Error copying response body:", err)
	}
}

// 从请求中提取请求信息，包括请求体
func extractRequestInfo(req *http.Request) (RequestInfo, error) {
	var requestBody []byte
	if req.Body != nil {
		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			return RequestInfo{}, err
		}
		requestBody = reqBody
		req.Body = io.NopCloser(bytes.NewBuffer(reqBody)) // 重置请求体，以便后续转发
	}
	return RequestInfo{
		Method:  req.Method,
		URL:     req.URL.String(),
		Headers: req.Header,
		Body:    requestBody,
	}, nil
}

// 从响应中提取响应信息，包括响应体
func extractResponseInfo(resp *http.Response, reqMethod string) (ResponseInfo, error) {
	var responseBody []byte
	if resp.Body != nil {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return ResponseInfo{}, err
		}
		responseBody = respBody
		resp.Body = io.NopCloser(bytes.NewBuffer(respBody)) // 重置响应体，以便后续处理
	}
	return ResponseInfo{
		Method:     reqMethod,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       responseBody,
	}, nil
}

// 将请求信息存储到pebbleDB
func storeRequestInfo(ga *GatewayApp, requestInfo RequestInfo) error {
	timestamp := time.Now().UnixNano()
	requestKey := fmt.Sprintf("request_%d_%s_%s", timestamp, requestInfo.Method, requestInfo.URL)
	requestData, _ := json.Marshal(requestInfo)
	return ga.storeAPIInfoInPebble(requestKey, string(requestData))
}

// 将响应信息存储到pebbleDB
func storeResponseInfo(ga *GatewayApp, responseInfo ResponseInfo) error {
	timestamp := time.Now().UnixNano()
	responseKey := fmt.Sprintf("response_%d_%s_%d", timestamp, responseInfo.Method, responseInfo.StatusCode)
	responseData, _ := json.Marshal(responseInfo)
	return ga.storeAPIInfoInPebble(responseKey, string(responseData))
}
