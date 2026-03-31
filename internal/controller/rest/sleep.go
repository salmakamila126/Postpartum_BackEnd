package rest

import (
	"Postpartum_BackEnd/internal/dto"
	"Postpartum_BackEnd/internal/errs"
	"Postpartum_BackEnd/internal/usecase"
	"Postpartum_BackEnd/pkg/timeutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SleepController struct {
	Usecase *usecase.SleepUsecase
}

func NewSleepController(u *usecase.SleepUsecase) *SleepController {
	return &SleepController{Usecase: u}
}

func (sc *SleepController) Start(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		HandleError(c, err)
		return
	}
	if err := sc.Usecase.StartSleep(userID); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.Success("sleep session started"))
}

func (sc *SleepController) End(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		HandleError(c, err)
		return
	}
	if err := sc.Usecase.EndSleep(userID); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.Success("sleep session ended"))
}

func (sc *SleepController) Manual(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	var input struct {
		Start string `json:"start" binding:"required"`
		End   string `json:"end" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		HandleError(c, errs.New(http.StatusBadRequest, "invalid request body"))
		return
	}

	start, err := timeutil.ParseRFC3339(input.Start)
	if err != nil {
		HandleError(c, errs.New(http.StatusBadRequest, "invalid start time - use RFC3339"))
		return
	}
	end, err := timeutil.ParseRFC3339(input.End)
	if err != nil {
		HandleError(c, errs.New(http.StatusBadRequest, "invalid end time - use RFC3339"))
		return
	}

	if err := sc.Usecase.AddSleepSession(userID, start, end); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.Success("sleep session added"))
}

func (sc *SleepController) Daily(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		HandleError(c, errs.New(http.StatusBadRequest, "query param 'date' is required (YYYY-MM-DD)"))
		return
	}

	date, err := timeutil.ParseDate(dateStr)
	if err != nil {
		HandleError(c, errs.New(http.StatusBadRequest, "invalid date format - use YYYY-MM-DD"))
		return
	}

	data, err := sc.Usecase.GetDailySleep(userID, date)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.Success(data))
}

func (sc *SleepController) Predict(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	preds, err := sc.Usecase.Predict(userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	res := make([]dto.SleepPredictionItem, 0, len(preds))
	for _, p := range preds {
		res = append(res, dto.SleepPredictionItem{
			Sleep: timeutil.FormatHour(p.NextSleep),
			Wake:  timeutil.FormatHour(p.NextWake),
		})
	}
	c.JSON(http.StatusOK, dto.Success(dto.SleepPredictionResponse{Predictions: res}))
}

func (sc *SleepController) Insight(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	res, err := sc.Usecase.GetTodayInsight(userID)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.Success(res))
}

func (sc *SleepController) Bulk(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	var body dto.SleepBulkRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		HandleError(c, errs.New(http.StatusBadRequest, "invalid request body"))
		return
	}
	if len(body.Sessions) == 0 {
		HandleError(c, errs.New(http.StatusBadRequest, "sessions must not be empty"))
		return
	}

	date, err := timeutil.ParseDate(body.Date)
	if err != nil {
		HandleError(c, errs.New(http.StatusBadRequest, "invalid date format - use YYYY-MM-DD"))
		return
	}

	if err := sc.Usecase.AddBulkSleepSession(userID, date, body.Sessions); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.Success("bulk sleep sessions added"))
}

func (sc *SleepController) Status(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	res, err := sc.Usecase.GetStatus(userID)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.Success(res))
}

func (sc *SleepController) History(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	res, err := sc.Usecase.GetHistory(userID)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.Success(res))
}
