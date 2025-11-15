package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterMetrics(engine *gin.Engine) error {
	reg := prometheus.NewRegistry()
	collectorsList := []prometheus.Collector{
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	}
	for _, v := range collectorsList {
		if err := reg.Register(v); err != nil {
			return err
		}
	}

	engine.GET("/metrics", gin.WrapH(promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		ErrorHandling:     promhttp.ContinueOnError,
		EnableOpenMetrics: true,
		Registry:          reg,
	})))
	return nil
}
