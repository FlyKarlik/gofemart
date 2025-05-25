package handler

import (
	"net/http"
	"time"

	"github.com/FlyKarlik/gofemart/internal/delivery/http/response"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type respObj struct {
	Uptime   string `json:"uptime"`
	DateTime string `json:"datetime"`
}

func (h *Handler) Ping(c *gin.Context) {
	trace := otel.Tracer("handler/ping")
	_, span := trace.Start(c.Request.Context(), "PingHandler")
	defer span.End()

	span.SetAttributes(
		attribute.String("handler", "Ping"),
		attribute.String("method", c.Request.Method),
		attribute.String("path", c.FullPath()),
	)

	resp := &respObj{
		Uptime:   time.Since(h.startup).String(),
		DateTime: time.Now().Format(time.RFC1123),
	}

	response.New(c, http.StatusOK, true, resp, nil)
}
