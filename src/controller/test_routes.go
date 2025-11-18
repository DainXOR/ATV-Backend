package controller

import (
	"dainxor/atv/configs"
	"dainxor/atv/logger"
	"dainxor/atv/service"
	"dainxor/atv/types"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TestRoutes(router *gin.Engine) {
	testRouter := router.Group("api/test")
	{
		testRouter.GET("/get", service.Test.Get)

		testRouter.POST("/post", service.Test.Post)

		testRouter.PUT("/put", service.Test.Put)
		testRouter.PATCH("/patch", service.Test.Patch)
		testRouter.DELETE("/del", service.Test.Delete)

		testRouter.POST("/wh/send/:msg", func(c *gin.Context) {
			msg := c.Param("msg")
			logger.Info("Got:", msg)

			configs.WebHooks.SendTo("fatv.test", msg)

			c.JSON(types.Http.C200().Created(),
				types.EmptyResponse("all good"),
			)
		})

		testRouter.POST("/wh/receive", myWHReceiver)
	}
}

func myWHReceiver(c *gin.Context) {
	key := c.Query("key")
	logger.Info("Got:", key)

	var payload configs.Payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Error("Invalid payload JSON:", err)
		c.JSON(200, gin.H{"error": "invalid payload"})
		return
	}

	logger.Info("Body from wh:", payload)

	c.JSON(types.Http.C200().Created(),
		types.EmptyResponse("all good"),
	)
}

type WebhookEnvelope map[string]any

func debugWHReceiver(c *gin.Context) {
	key := c.Query("key")
	logger.Info("Got key:", key)

	// Read raw body (also allows us to log it)
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error("Failed to read raw body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot read body"})
		return
	}

	// Log headers and raw body for debugging
	logger.Debugf("Headers: %+v", c.Request.Header)
	logger.Debugf("Raw body: %s", string(rawBody))

	// Try to parse the envelope into a generic map first
	var env WebhookEnvelope
	if len(rawBody) > 0 {
		if err := json.Unmarshal(rawBody, &env); err != nil {
			// If parsing fails, log and return 400 for easier debugging
			logger.Error("Failed to unmarshal envelope JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid envelope json"})
			return
		}
	} else {
		logger.Warning("Empty request body received")
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty body"})
		return
	}

	logger.Debugf("Envelope received: %+v", env)

	// Extract payload and payload_encoding
	var innerPayloadBytes []byte
	payloadEncoding := ""
	if v, ok := env["payload_encoding"]; ok {
		if s, ok2 := v.(string); ok2 {
			payloadEncoding = s
		}
	}

	payloadVal, hasPayload := env["payload"]
	if !hasPayload || payloadVal == nil {
		// sometimes CloudAMQP might use "body" or "message" keys, log and return
		logger.Warning("No 'payload' field found in envelope; keys:", keysOfMap(env))
		c.JSON(http.StatusBadRequest, gin.H{"error": "no payload in envelope"})
		return
	}

	switch p := payloadVal.(type) {
	case string:
		// payload is a JSON string (most common)
		logger.Debug("payload is a string; length:", len(p))
		if payloadEncoding == "base64" {
			decoded, err := base64.StdEncoding.DecodeString(p)
			if err != nil {
				logger.Error("failed to base64-decode payload:", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid base64 payload"})
				return
			}
			innerPayloadBytes = decoded
		} else {
			// assume it's a JSON string: e.g. "{\"event\":\"Test\",...}"
			innerPayloadBytes = []byte(p)
		}
	case map[string]any:
		// payload already parsed into an object, re-marshal it to bytes
		b, err := json.Marshal(p)
		if err != nil {
			logger.Error("failed to re-marshal payload object:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
			return
		}
		innerPayloadBytes = b
	default:
		// other possible types (number etc) — marshal generically
		b, err := json.Marshal(p)
		if err != nil {
			logger.Error("failed to marshal payload of unknown type:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
			return
		}
		innerPayloadBytes = b
	}

	// Now we have the inner payload bytes. Unmarshal into your Payload struct.
	var msg configs.Payload
	if err := json.Unmarshal(innerPayloadBytes, &msg); err != nil {
		logger.Error("Invalid payload JSON:", err)
		logger.Debugf("Inner payload raw: %s", string(innerPayloadBytes))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid inner payload"})
		return
	}

	logger.Info("Message from RabbitMQ webhook:", msg)

	// Return 200 OK quickly
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func keysOfMap(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
