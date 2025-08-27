package routes

import (
	"net/http"
	database "subscription/db"
	"subscription/dto"
	"subscription/models"
	"time"

	"github.com/gin-gonic/gin"
)

func SubscriptionRouter(route *gin.Engine) {
	subs := route.Group("/subscriptions")
	{
		subs.GET("", getSubs)
		subs.GET("/:id", GetSub)
		subs.POST("", CreateSub)
		subs.PATCH("/:id", UpdateSub)
		subs.DELETE("/:id", DeleteSub)
		subs.GET("/total", GetTotal)
	}
}

// @Summary      Получение всех подписок
// @Description  Возвращает список всех подписок
// @Tags         subscriptions
// @Produce      json
// @Success      200  {array}   dto.SubscriptionResponse
// @Router       /subscriptions [get]
func getSubs(c *gin.Context) {
	var subscriptions []models.Subscription
	database.DB.Find(&subscriptions)

	//Конвертируем ответ в нужный формат
	var response []dto.SubscriptionResponse
	for _, sub := range subscriptions {
		response = append(response, dto.FromModel(sub))
	}

	c.JSON(http.StatusOK, response)
}

// @Summary      Создание новой подписки
// @Description  Добавляет новую подписку
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription  body      dto.SubscriptionRequest  true  "Subscription object"
// @Success      201  {object} dto.SubscriptionResponse
// @Failure      400  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /subscriptions [post]
func CreateSub(c *gin.Context) {
	var req dto.SubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := req.ToModel()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.FromModel(sub))
}

// @Summary      Получение подписки по ID
// @Description  Возвращает подписку по её ID
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      int  true  "Subscription ID"
// @Success      200  {object} dto.SubscriptionResponse
// @Failure      404  {object} map[string]string
// @Router       /subscriptions/{id} [get]
func GetSub(c *gin.Context) {
	id := c.Param("id")
	var sub models.Subscription
	if err := database.DB.First(&sub, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}
	c.JSON(http.StatusOK, dto.FromModel(sub))
}

// @Summary      Частичное обновление подписки
// @Description  Обновляет только переданные поля подписки (service_name, price, user_id, start_date, end_date)
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id   path      int                     true  "ID подписки"
// @Param        subscription body  map[string]interface{} true  "Поля для обновления"
// @Success      200  {object}  models.Subscription
// @Failure      400  {object}  map[string]string "error"
// @Failure      404  {object}  map[string]string "error"
// @Failure      500  {object}  map[string]string "error"
// @Router       /subscriptions/{id} [patch]
func UpdateSub(c *gin.Context) {
	id := c.Param("id")

	var sub models.Subscription
	if err := database.DB.First(&sub, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	// читаем только переданные поля
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// обновляем только указанные ключи
	if err := database.DB.Model(&sub).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.FromModel(sub))
}

// @Summary      Удаление подписки
// @Description  Удаляет подписку по ID
// @Tags         subscriptions
// @Param        id   path      int  true  "Subscription ID"
// @Success      200  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /subscriptions/{id} [delete]
func DeleteSub(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Subscription{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Subscription deleted"})
}

// @Summary      Подсчёт общей стоимости подписок
// @Description  Считает сумму всех подписок за выбранный период с фильтрацией по user_id и названию сервиса
// @Tags         subscriptions
// @Produce      json
// @Param        user_id       query     string  false  "ID пользователя"
// @Param        service_name  query     string  false  "Название сервиса"
// @Param        start         query     string  true   "Начало периода (MM-YYYY)"
// @Param        end           query     string  true   "Конец периода (MM-YYYY)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Router       /subscriptions/total [get]
func GetTotal(c *gin.Context) {
	var total int64

	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	startStr := c.Query("start")
	endStr := c.Query("end")

	// проверка дат
	start, err := time.Parse("01-2006", startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат дат, ожидается MM-YYYY"})
		return
	}
	end, err := time.Parse("01-2006", endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат дат, ожидается MM-YYYY"})
		return
	}

	db := database.DB.Model(&models.Subscription{})

	if userID != "" {
		db = db.Where("user_id = ?", userID)
	}
	if serviceName != "" {
		db = db.Where("service_name ILIKE ?", "%"+serviceName+"%")
	}

	db = db.Where("start_date >= ? AND (end_date IS NULL OR end_date <= ?)", start, end)

	if err := db.Select("SUM(price)").Scan(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_price": total,
	})
}
