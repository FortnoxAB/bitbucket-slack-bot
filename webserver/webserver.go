package webserver

import (
	"net/http"

	"github.com/fortnoxab/bitbucket-slack-bot/models"
	"github.com/fortnoxab/bitbucket-slack-bot/service"
	"github.com/fortnoxab/ginprometheus"
	"github.com/gin-gonic/gin"
	"github.com/jonaz/ginlogrus"
	"github.com/sirupsen/logrus"
)

type webserver struct {
	Notifier   *service.Notifier
	Prometheus *ginprometheus.Prometheus
}

//New webserver
func New(n *service.Notifier) *webserver {
	return &webserver{
		Notifier: n,
	}
}

//Init a webserver with Gin
func (ws *webserver) Init() *gin.Engine {

	router := gin.New()

	if ws.Prometheus != nil {
		ws.Prometheus.Use(router)
	}

	router.Use(ginlogrus.New(logrus.StandardLogger(), "/health", "/metrics"), gin.Recovery())

	router.POST("/webhook/:channel", checkErr(ws.handleWebhook))
	router.POST("/webhook", checkErr(ws.handleWebhook))

	router.GET("/health", ws.healthHandler)
	return router
}

// https://confluence.atlassian.com/bitbucketserver/event-payload-938025882.html
func (ws *webserver) handleWebhook(c *gin.Context) error {
	logrus.Info(c.Request.Header)

	body := &models.WebhookBody{}
	err := c.BindJSON(body)
	if err != nil {
		return err
	}

	//logrus.Info("body", body)
	return ws.Notifier.ProcessWebhook(body)
}

func (ws *webserver) healthHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
