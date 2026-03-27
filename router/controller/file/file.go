package file

import (
	"GopherAI/common/code"
	"GopherAI/router/controller"
	"GopherAI/router/service/file"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadFileResponse struct {
	FilePath string `json:"file_path,omitempty"`
	response.Response
}

func UploadRagFile(c *gin.Context) {
	res := new(UploadFileResponse)
	uploadedFile, err := c.FormFile("file")
	if err != nil {
		log.Println("FormFile fail ", err)
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	// 得到用户名
	username := c.GetString("userName")
	if username == "" {
		log.Println("Username not found in context")
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidToken))
		return
	}

	// 得到sessionId
	sessionId := c.PostForm("sessionId")
	if sessionId == "" {
		log.Println("SessionId not found in context")
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidToken))
		return
	} else {
		fmt.Println(sessionId)
	}

	// 将上传的文件保存到服务器
	filePath, err := file.UploadsFiles(c, username, sessionId, uploadedFile)

	if err != nil {
		log.Println("UploadFile fail ", err)
		c.JSON(http.StatusOK, res.CodeOf(code.CodeServerBusy))
		return
	}

	res.Success()
	res.FilePath = filePath
	c.JSON(http.StatusOK, res)
}
