package file

import (
	"AIGO/internal/logic"
	"AIGO/internal/middleware"
	"AIGO/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	FileReq struct {
		SessionID string `json:"session_id" form:"session_id"`
	}

	FileResp struct {
		FileIDs []string `json:"file_ids"`
	}
)

const FileFormKey = "file"

func FileHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// fileHeader, err := ctx.FormFile(FileFormKey)
		// 多文件上传
		form, err := ctx.MultipartForm()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, &model.Response{
				Code: http.StatusBadRequest,
				Msg:  "invalid file: " + err.Error(),
				Data: nil,
			})
			return
		}
		fileHeaders := form.File[FileFormKey]
		if len(fileHeaders) == 0 {
			ctx.JSON(http.StatusBadRequest, &model.Response{
				Code: http.StatusBadRequest,
				Msg:  "no file",
				Data: nil,
			})
			return
		}
		var req FileReq
		if err = ctx.ShouldBind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, &model.Response{
				Code: http.StatusBadRequest,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		// 由于开启了异步索引文件 所以这里的文件路径是立即返回的 但是索引文件的过程是异步的
		fileIds, err := logic.SaveUpLoadFiles(ctx, ctx.GetString(middleware.AuthUserKey), req.SessionID, fileHeaders)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, &model.Response{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}
		ctx.JSON(http.StatusOK, &model.Response{
			Code: http.StatusOK,
			Msg:  "success",
			Data: &FileResp{
				FileIDs: fileIds,
			},
		})
	}
}
