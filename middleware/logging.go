package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     true,
	})
	logger.SetLevel(logrus.InfoLevel)
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// читаем request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody)) // восстанавливаем body

		// логируем response body
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// выполняем запрос
		c.Next()

		latency := time.Since(start)

		// собираем данные
		var reqJSON interface{}
		if err := json.Unmarshal(requestBody, &reqJSON); err != nil {
			reqJSON = string(requestBody) // fallback в строку
		}

		var respJSON interface{}
		if err := json.Unmarshal(blw.body.Bytes(), &respJSON); err != nil {
			respJSON = blw.body.String()
		}

		// собираем данные
		logData := map[string]interface{}{
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"status":   c.Writer.Status(),
			"latency":  latency.String(),
			"request":  reqJSON,
			"response": respJSON,
		}

		// логируем JSON-объект
		logger.WithFields(logrus.Fields(logData)).Info("http request")
	}
}
