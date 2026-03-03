package auth

import (
	"AIGO/internal/logic"
	"AIGO/internal/model"
	"AIGO/pkg/utils/email"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserRegisterRequest struct {
	Username         string `json:"username"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	VerificationCode string `json:"verification_code"`
}

func UserRegisterHandler(ctx *gin.Context) {
	req := &UserRegisterRequest{}
	// 前端参数绑定
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, &model.Response{
			Code: http.StatusBadRequest,
			Msg:  "invailid params",
			Data: nil,
		})
		return
	}
	// 业务逻辑
	if err := logic.Register(ctx, req.Username, req.Email, req.Password, req.VerificationCode); err != nil {
		ctx.JSON(http.StatusInternalServerError, &model.Response{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, &model.Response{
		Code: http.StatusOK,
		Msg:  "register success",
		Data: nil,
	})
	// 发送欢迎邮件
	email.SendWelcomeEmail(req.Email, req.Username)
}
