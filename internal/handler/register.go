package handler

import (
	"AIGO/internal/logic"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserRegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UserRegisterHandler(ctx *gin.Context) {
	req := &UserRegisterRequest{}
	// 前端参数绑定
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, &Response{
			Code: http.StatusBadRequest,
			Msg:  "invailid params",
			Data: nil,
		})
		return
	}
	// 业务逻辑
	if err := logic.Register(req.Username, req.Email, req.Password); err != nil {
		ctx.JSON(http.StatusInternalServerError, &Response{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, &Response{
		Code: http.StatusOK,
		Msg:  "register success",
		Data: nil,
	})
}
