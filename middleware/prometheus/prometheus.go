package prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"sync"
	"time"
)

const (
	metricsPath = "/metrics"
)

var (
	// httpHistogram prometheus 模型
	httpHistogram = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "http_server",
			Subsystem: "",
			Name:      "requests_summary",
			Help:      "Summary of response latency (seconds) of http handlers.",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		}, []string{"method", "code", "uri"})
	Registry = prometheus.NewRegistry()
)

func init() {

	Registry.MustRegister(httpHistogram)
	//	Registry.MustRegister(httpCounter)
}

// handlerPath 定义采样路由struct
type handlerPath struct {
	sync.Map
}

// get 获取path
func (hp *handlerPath) get(handler string) string {
	v, ok := hp.Load(handler)
	if !ok {
		return ""
	}
	return v.(string)
}

// set 保存path到sync.Map
func (hp *handlerPath) set(ri gin.RouteInfo) {
	hp.Store(ri.Handler, ri.Path)
}

// GinPrometheus gin调用Prometheus的struct
type GinPrometheus struct {
	engine  *gin.Engine
	ignored map[string]bool
	pathMap *handlerPath
	updated bool
}
type Option func(*GinPrometheus)

// Ignore 添加忽略的路径
func Ignore(path ...string) Option {
	return func(gp *GinPrometheus) {
		for _, p := range path {
			gp.ignored[p] = true
		}
	}
}

// New new gin prometheus
func New(e *gin.Engine, options ...Option) *GinPrometheus {
	if e == nil {
		return nil
	}
	gp := &GinPrometheus{
		engine: e,
		ignored: map[string]bool{
			metricsPath: true,
			//			faviconPath: true,
		},
		pathMap: &handlerPath{},
	}

	for _, o := range options {
		o(gp)
	}
	return gp
}

// updatePath 更新path
func (gp *GinPrometheus) updatePath() {
	gp.updated = true
	for _, ri := range gp.engine.Routes() {
		gp.pathMap.set(ri)
	}
}

// Middleware set gin middleware
func (gp *GinPrometheus) PromeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !gp.updated {
			gp.updatePath()
		}
		// 过滤请求
		if gp.ignored[c.Request.URL.String()] {
			c.Next()
			return
		}
		start := time.Now()
		c.Next()
		httpHistogram.WithLabelValues(
			c.Request.Method,
			strconv.Itoa(c.Writer.Status()),
			gp.pathMap.get(c.HandlerName()),
		).Observe(time.Since(start).Seconds())
	}
}
