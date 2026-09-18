package repository
import (
	"database/sql"
	"github.com/google/uuid"
	"autosmm/internal/models"
	_ "github.com/lib/pq"
)
type PostRepository struct { db *sql.DB }
func NewPostRepository(db *sql.DB) *PostRepository { return &PostRepository{db: db} }
func (r *PostRepository) CreatePost(post *models.Post, targetAccountIDs []uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil { return err }
	query := `INSERT INTO posts (id, user_id, content, scheduled_at, status) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(query, post.ID, post.UserID, post.Content, post.ScheduledAt, post.Status)
	if err != nil { tx.Rollback(); return err }
	for _, accID := range targetAccountIDs {
		targetQuery := `INSERT INTO post_targets (post_id, account_id, status) VALUES ($1, $2, 'pending')`
		_, err = tx.Exec(targetQuery, post.ID, accID)
		if err != nil { tx.Rollback(); return err }
	}
	return tx.Commit()
}
func (r *PostRepository) GetUserPosts(userID uuid.UUID) ([]models.Post, error) {
	rows, err := r.db.Query(`SELECT id, user_id, content, scheduled_at, status FROM posts WHERE user_id = $1 ORDER BY scheduled_at ASC`, userID)
	if err != nil { return nil, err }
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Content, &p.ScheduledAt, &p.Status); err != nil { return nil, err }
		posts = append(posts, p)
	}
	return posts, nil
}
func (r *PostRepository) GetPendingPosts() ([]models.Post, error) {
	rows, err := r.db.Query(`SELECT id, user_id, content, scheduled_at, status FROM posts WHERE status = 'scheduled' AND scheduled_at <= NOW()`)
	if err != nil { return nil, err }
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Content, &p.ScheduledAt, &p.Status); err != nil { return nil, err }
		posts = append(posts, p)
	}
	return posts, nil
}
func (r *PostRepository) MarkPostPublished(postID uuid.UUID) error {
	_, err := r.db.Exec(`UPDATE posts SET status = 'published', published_at = NOW() WHERE id = $1`, postID)
	return err
}
