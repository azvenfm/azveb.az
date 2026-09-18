
const schedulerService = require('../services/schedulerService');
exports.createPost = async (req, res) => {
  try {
    const userId = req.user.userId;
    const post = await schedulerService.schedulePost(userId, req.body);
    res.status(201).json(post);
  } catch (err) { res.status(500).json({ error: err.message }); }
};
exports.listPosts = async (req, res) => {
  try {
    const userId = req.user.userId;
    const posts = await schedulerService.getUserPosts(userId);
    res.json(posts);
  } catch (err) { res.status(500).json({ error: err.message }); }
};
