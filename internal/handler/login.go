package handler

import (
	"AIGO/internal/logic"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserLoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UserLoginHandler(ctx *gin.Context) {
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
	loginUser, err := logic.Login(req.Username, req.Email, req.Password)
	if err != nil || loginUser == nil {
		ctx.JSON(http.StatusInternalServerError, &Response{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, &Response{
		Code: http.StatusOK,
		Msg:  "login success",
		Data: loginUser,
	})
}
