package bootstrap

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"

	"api-gateway/internal/api"
	"api-gateway/internal/downstream"
	"api-gateway/internal/global"
	"api-gateway/internal/services"
	"api-gateway/pkg/db"

	"github.com/cockroachdb/pebble"
	"github.com/gin-gonic/gin"
)

// GatewayApp结构体用于网关转发相关的初始化和配置
type GatewayApp struct {
	Router     *gin.Engine
	API        *api.APIController
	Downstream *downstream.DownstreamServiceHandler
	PebbleDB   *pebble.DB
	// Logger     *utils.TrafficLogger
}

// NewGatewayApp创建并初始化用于网关转发的应用实例
func NewGatewayApp() *GatewayApp {
	return &GatewayApp{}
}

// Initialize初始化网关转发应用的各种组件
func (ga *GatewayApp) Initialize() {
	var err error
	global.DB, err = db.NewDB()
	if err != nil {
		fmt.Printf("Error initializing database: %v", err)
		return
	}

	ga.PebbleDB, err = pebble.Open("", nil)
	if err != nil {
		fmt.Printf("Error opening PebbleDB: %v", err)
		return
	}

	downstreamService := services.NewDownstreamService()
	ga.Downstream = downstream.NewDownstreamServiceHandler(downstreamService)

	// trafficService := utils.NewTrafficService(global.DB)
	// ga.Logger = utils.NewTrafficLogger(trafficService)

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
			service, err := ga.Downstream.GetService(context.Background(), serviceName)
			if err == nil && service.URL != "" {
				backendURL, err := url.Parse(service.URL)
				if err == nil {
					apiInfo := fmt.Sprintf("Service: %s, URL: %s", serviceName, service.URL)
					err = ga.storeAPIInfoInPebble(serviceName, apiInfo)
					if err != nil {
						fmt.Printf("Error storing API info in Pebble: %v", err)
					}
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
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

// forwardRequest用于转发请求
func (ga *GatewayApp) forwardRequest(c *gin.Context, client *http.Client) {
	req := c.Request.Clone(c.Request.Context())
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

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
