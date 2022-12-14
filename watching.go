package controller

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/chanwit/sa-64-example/entity"
	"github.com/gin-gonic/gin"
)

// POST /watch_videos
func CreateWatchVideo(c *gin.Context) {

	var watchvideo entity.WatchVideo
	var resolution entity.Resolution
	var playlist entity.Playlist
	var video entity.Video

	// ผลลัพธ์ที่ได้จากขั้นตอนที่ 8 จะถูก bind เข้าตัวแปร watchVideo
	if err := c.ShouldBindJSON(&watchvideo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 9: ค้นหา video ด้วย id
	if tx := entity.DB().Where("id = ?", watchvideo.VideoID).First(&video); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "video not found"})
		return
	}

	// 10: ค้นหา resolution ด้วย id
	if tx := entity.DB().Where("id = ?", watchvideo.ResolutionID).First(&resolution); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resolution not found"})
		return
	}

	// 11: ค้นหา playlist ด้วย id
	if tx := entity.DB().Where("id = ?", watchvideo.PlaylistID).First(&playlist); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "playlist not found"})
		return
	}
	// 12: สร้าง WatchVideo
	wv := entity.WatchVideo{
		Resolution:  resolution,             // โยงความสัมพันธ์กับ Entity Resolution
		Video:       video,                  // โยงความสัมพันธ์กับ Entity Video
		Playlist:    playlist,               // โยงความสัมพันธ์กับ Entity Playlist
		WatchedTime: watchvideo.WatchedTime, // ตั้งค่าฟิลด์ watchedTime
	}

	// ขั้นตอนการ validate ที่นำมาจาก unit test
	if _, err := govalidator.ValidateStruct(wv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 13: บันทึก
	if err := entity.DB().Create(&wv).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": wv})
}

// GET /watchvideo/:id
func GetWatchVideo(c *gin.Context) {
	var watchvideo entity.WatchVideo
	id := c.Param("id")
	if err := entity.DB().Preload("Resolution").Preload("Playlist").Preload("Video").Raw("SELECT * FROM watch_videos WHERE id = ?", id).Find(&watchvideo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": watchvideo})
}

// GET /watch_videos
func ListWatchVideos(c *gin.Context) {
	var watchvideos []entity.WatchVideo
	if err := entity.DB().Preload("Resolution").Preload("Playlist").Preload("Video").Raw("SELECT * FROM watch_videos").Find(&watchvideos).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": watchvideos})
}

// DELETE /watch_videos/:id
func DeleteWatchVideo(c *gin.Context) {
	id := c.Param("id")
	if tx := entity.DB().Exec("DELETE FROM watch_videos WHERE id = ?", id); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "watchvideo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": id})
}

// PATCH /watch_videos
func UpdateWatchVideo(c *gin.Context) {
	var watchvideo entity.WatchVideo
	if err := c.ShouldBindJSON(&watchvideo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if tx := entity.DB().Where("id = ?", watchvideo.ID).First(&watchvideo); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "watchvideo not found"})
		return
	}

	if err := entity.DB().Save(&watchvideo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": watchvideo})
}
