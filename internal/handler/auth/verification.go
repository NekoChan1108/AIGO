package auth

import (
	"AIGO/config"
	"AIGO/internal/logic"
	"AIGO/internal/model"
	"AIGO/pkg/db"
	"AIGO/pkg/utils/email"
	regx "AIGO/pkg/utils/regex"
	"AIGO/pkg/utils/verification"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type VerificationCodeRequest struct {
	Email string `json:"email"`
}

func VerificationCodeHandler(ctx *gin.Context) {
	var req VerificationCodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, &model.Response{
			Code: http.StatusBadRequest,
			Data: nil,
			Msg:  "invailid params",
		})
		return
	}
	// 生成验证码
	code := verification.GenerateVerificationCode()
	// 存入redis 多10秒来保证验证码的有效期
	d := time.Duration(config.Cfg.EmailCfg.Expiration)*time.Minute + time.Second*10
	cmd := db.RedisDB.Set(ctx, logic.EmailVerificationCacheKey+req.Email, code, d)
	if err := cmd.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, &model.Response{
			Code: http.StatusInternalServerError,
			Data: nil,
			Msg:  "send verification code failed",
		})
		return
	}
	// 邮箱正则校验 不通过不存 redis
	if ok, err := regx.EamilRegex(req.Email); !ok || err != nil {
		ctx.JSON(http.StatusBadRequest, &model.Response{
			Code: http.StatusBadRequest,
			Data: nil,
			Msg:  "invailid email",
		})
		return
	}
	// 发送验证码
	if err := email.SendVerificationEmail(req.Email, code); err != nil {
		ctx.JSON(http.StatusInternalServerError, &model.Response{
			Code: http.StatusInternalServerError,
			Data: nil,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, &model.Response{
		Code: http.StatusOK,
		Data: nil,
		Msg:  "send verification code success",
	})
}
