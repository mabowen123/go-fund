package handlers

import (
	"fmt"
	"fund/mysql"
	"fund/reptile"
	"fund/xRedis"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func Reptile(c *gin.Context) {
	reptile.Run(true)
	c.JSON(http.StatusOK, Success)
}

type reptileData struct {
	Id                uint
	FundId            string
	FundName          string
	Share             float64
	Rate              string
	ActualRate        string
	EstimatedEarnings interface{}
	ActualEarnings    interface{}
	UpdateAt          string
}

func ReptileData(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userId := claims["id"]
	var fund []mysql.UserFund
	mysql.Db.Where("user_id = ? ", userId).Find(&fund)
	data := make(map[string]interface{})
	for key, item := range fund {
		fundData := xRedis.HGetAll(xRedis.FundName(item.FundId))
		res := reptileData{
			Id:             item.ID,
			FundId:         item.FundId,
			FundName:       fundData["name"],
			Share:          item.Share,
			Rate:           fundData["gszzl"] + "%",
			ActualRate:     "暂未更新",
			ActualEarnings: "暂未更新",
			UpdateAt:       fundData["gztime"],
		}

		gsz, er1 := strconv.ParseFloat(fundData["gsz"], 32)
		if er1 == nil {
			dwjz, er2 := strconv.ParseFloat(fundData["dwjz"], 32)
			if er2 == nil {
				res.EstimatedEarnings = Decimal((gsz - dwjz) * item.Share)
			} else {
				res.EstimatedEarnings = "暂未更新"
			}

			rActualDate, _ := time.Parse("2006-01-02", fundData["actual_date"])
			if (time.Now().Day() == rActualDate.Day()) &&
				(time.Now().Month() == rActualDate.Month() &&
					(time.Now().Year() == rActualDate.Year())) {
				res.ActualRate = fundData["actual_gszl"]
				actualGsz, er3 := strconv.ParseFloat(fundData["actual_gsz"], 32)
				if er3 == nil {
					res.ActualEarnings = Decimal((actualGsz - dwjz) * item.Share)
				} else {
					res.ActualEarnings = "暂未更新"
				}
			}
		}

		data[fmt.Sprintf("%v", key)] = res
	}
	c.JSON(http.StatusOK, Success.WithData(data))
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

type fundParams struct {
	FundId string `form:"fund_id" json:"fund_id" binding:"required"`
	Share  string `form:"share" json:"share" binding:"required"`
}

func FundAdd(c *gin.Context) {
	var params fundParams
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, Fail.WithMsg(err.Error()))
		return
	}
	claims := jwt.ExtractClaims(c)
	userId := claims["id"]
	share, _ := strconv.ParseFloat(params.Share, 32)
	createData := &mysql.UserFund{
		UserId: int(userId.(float64)),
		FundId: params.FundId,
		Share:  share,
		Amount: 0,
	}

	res := mysql.Db.Create(createData)
	if res.RecordNotFound() {
		c.JSON(http.StatusOK, Fail.WithMsg("创建失败"))
		return
	}
	reptile.Run(false)
	c.JSON(http.StatusOK, Success)
}

type fundDelParams struct {
	Id uint `form:"id" json:"id" binding:"required"`
}

func FundDel(c *gin.Context) {
	var params fundDelParams
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, Fail.WithMsg(err.Error()))
		return
	}
	claims := jwt.ExtractClaims(c)
	userId := claims["id"]
	userFund := &mysql.UserFund{
		ID:     params.Id,
		UserId: int(userId.(float64)),
	}
	mysql.Db.Delete(userFund)
	c.JSON(http.StatusOK, Success)
}
