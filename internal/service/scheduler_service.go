package service
import (
	"log"
	"time"
	"autosmm/internal/models"
	"autosmm/internal/repository"
	"github.com/google/uuid"
)
type SchedulerService struct { postRepo *repository.PostRepository }
func NewSchedulerService(repo *repository.PostRepository) *SchedulerService { return &SchedulerService{postRepo: repo} }
func (s *SchedulerService) SchedulePost(userID uuid.UUID, req models.PostRequest) (*models.Post, error) {
	post := &models.Post{
		ID: uuid.New(), UserID: userID, Content: req.Content,
		ScheduledAt: req.ScheduledAt, Status: "scheduled",
	}
	err := s.postRepo.CreatePost(post, req.TargetAccountIDs)
	if err != nil { return nil, err }
	return post, nil
}
func (s *SchedulerService) GetUserPosts(userID uuid.UUID) ([]models.Post, error) {
	return s.postRepo.GetUserPosts(userID)
}
func (s *SchedulerService) StartWorker() {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C { s.publishDuePosts() }
	}()
}
func (s *SchedulerService) publishDuePosts() {
	posts, err := s.postRepo.GetPendingPosts()
	if err != nil { log.Printf("Error fetching pending posts: %v", err); return }
	for _, post := range posts {
		log.Printf("🚀 Publishing post %s: %s", post.ID, post.Content)
		if err := s.postRepo.MarkPostPublished(post.ID); err != nil {
			log.Printf("Error marking post %s as published: %v", post.ID, err)
		}
	}
}
