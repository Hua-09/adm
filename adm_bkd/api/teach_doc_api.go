package api

import (
	"net/http"

	"adm_bkd/config"
	"adm_bkd/utils/api_result"
	"adm_bkd/utils/err_mgr"
	"adm_bkd/utils/storage"

	"github.com/gin-gonic/gin"
)

// RegisterTeachDocAPI registers teaching-document routes.
func RegisterTeachDocAPI(r *gin.Engine, cfg *config.Config) {
	group := r.Group("/teach-doc")
	group.POST("/upload", uploadTeachDoc(cfg))
	group.GET("/list", listTeachDocs(cfg))
	group.DELETE("/:id", deleteTeachDoc(cfg))
}

func uploadTeachDoc(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest,
				api_result.Fail(err_mgr.ErrInvalidParam, err_mgr.ErrStr(err_mgr.ErrInvalidParam)))
			return
		}

		dest, err := storage.SaveUpload(cfg.Storage.RootDir, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				api_result.Fail(err_mgr.ErrStorageWrite, err_mgr.ErrStr(err_mgr.ErrStorageWrite)))
			return
		}

		c.JSON(http.StatusOK, api_result.OK(gin.H{"path": dest}))
	}
}

func listTeachDocs(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		docs, err := storage.ListDocs(cfg.Storage.RootDir)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				api_result.Fail(err_mgr.ErrStorageRead, err_mgr.ErrStr(err_mgr.ErrStorageRead)))
			return
		}
		c.JSON(http.StatusOK, api_result.OK(docs))
	}
}

func deleteTeachDoc(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := storage.DeleteDoc(cfg.Storage.RootDir, id); err != nil {
			c.JSON(http.StatusInternalServerError,
				api_result.Fail(err_mgr.ErrStorageDelete, err_mgr.ErrStr(err_mgr.ErrStorageDelete)))
			return
		}
		c.JSON(http.StatusOK, api_result.OK(nil))
	}
}
