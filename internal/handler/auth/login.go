package auth

import (
	"AIGO/internal/logic"
	"AIGO/internal/model"
	"AIGO/pkg/utils/cookie"
	"AIGO/pkg/utils/jwt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserLoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	*model.User
	AccessToken string `json:"access_token"` // 前端存local
}

func UserLoginHandler(ctx *gin.Context) {
	req := &UserLoginRequest{}
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
	loginUser, err := logic.Login(ctx, req.Username, req.Email, req.Password)
	if err != nil || loginUser == nil {
		ctx.JSON(http.StatusInternalServerError, &model.Response{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}
	// 生成token
	accessToken, refreshToken, err := jwt.GenerateTokens(loginUser.Username, time.Now())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &model.Response{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}
	data := &UserResponse{
		User:        loginUser,
		AccessToken: accessToken,
	}
	cookie.SetRefreshTokenCookie(ctx, refreshToken)
	ctx.JSON(http.StatusOK, &model.Response{
		Code: http.StatusOK,
		Msg:  "login success",
		Data: data,
	})
}
