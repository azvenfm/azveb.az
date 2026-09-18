
const db = require('../db');
const cron = require('node-cron');
class SchedulerService {
  async schedulePost(userId, { content, scheduledAt, targetAccountIds }) {
    const postQuery = 'INSERT INTO posts (user_id, content, scheduled_at, status) VALUES ($1, $2, $3, \'scheduled\') RETURNING id';
    const postResult = await db.query(postQuery, [userId, content, scheduledAt]);
    const postId = postResult.rows[0].id;
    for (const accId of targetAccountIds) {
      await db.query('INSERT INTO post_targets (post_id, account_id, status) VALUES ($1, $2, \'pending\')', [postId, accId]);
    }
    return { id: postId, status: 'scheduled' };
  }
  async getUserPosts(userId) {
    const result = await db.query('SELECT * FROM posts WHERE user_id = $1 ORDER BY scheduled_at ASC', [userId]);
    return result.rows;
  }
  startWorker() {
    cron.schedule('*/30 * * * * *', async () => {
      try {
        const result = await db.query("SELECT * FROM posts WHERE status = 'scheduled' AND scheduled_at <= NOW()");
        for (const post of result.rows) {
          await db.query('UPDATE posts SET status = \'published\', published_at = NOW() WHERE id = $1', [post.id]);
        }
      } catch (err) { console.error('Worker error:', err); }
    });
  }
}
module.exports = new SchedulerService();
