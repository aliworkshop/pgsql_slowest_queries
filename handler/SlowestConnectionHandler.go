package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type SlowQuery struct {
	Pid           string
	User          string
	QueryStart    *time.Time
	StateChange   *time.Time
	QueryTime     string
	Query         string
	State         string
	WaitEventType string
	WaitEvent     string
}

func (h *handler) SlowestConnectionHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("per_page", "10"))
	offset := (page - 1) * limit

	var data []SlowQuery
	q := h.db.Table("pg_stat_activity").Limit(limit).Offset(offset)
	q.Select("pid, user, pg_stat_activity.query_start, pg_stat_activity.state_change, pg_stat_activity.state_change - pg_stat_activity.query_start AS query_time, " +
		"query, state, wait_event_type, wait_event").Where("query_start is not null")

	statement := c.Query("statement")
	if statement != "" {
		q.Where("LOWER(query) like LOWER(?)", statement+"%")
	}

	q.Order("query_time " + c.Query("order", "desc")).Find(&data)
	if q.Error != nil {
		logrus.Warnf("error on get query : %v", q.Error)
		return c.Status(500).JSON(map[string]interface{}{
			"message": "internal server error",
		})
	}

	return c.JSON(data)
}
