package handlers
import (
	"net/http"
	"autosmm/internal/models"
	"autosmm/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)
type SchedulerHandler struct { schedulerService *service.SchedulerService }
func NewSchedulerHandler(s *service.SchedulerService) *SchedulerHandler { return &SchedulerHandler{schedulerService: s} }
func (h *SchedulerHandler) CreatePost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uID := userID.(uuid.UUID)
	var req models.PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post, err := h.schedulerService.SchedulePost(uID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to schedule post"})
		return
	}
	c.JSON(http.StatusCreated, post)
}
func (h *SchedulerHandler) ListPosts(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uID := userID.(uuid.UUID)
	posts, err := h.schedulerService.GetUserPosts(uID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}
	c.JSON(http.StatusOK, posts)
}
