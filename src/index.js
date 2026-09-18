
const express = require('express');
const cors = require('cors');
require('dotenv').config();
const authController = require('./controllers/authController');
const schedulerController = require('./controllers/schedulerController');
const authMiddleware = require('./middlewares/authMiddleware');
const schedulerService = require('./services/schedulerService');

const app = express();
app.use(cors());
app.use(express.json());

// Health Check
app.get('/health', (req, res) => res.json({ status: 'healthy', timestamp: new Date() }));

// Auth Routes
app.post('/api/v1/auth/register', authController.register);
app.post('/api/v1/auth/login', authController.login);

// Protected Routes
app.use('/api/v1', authMiddleware);
app.post('/api/v1/scheduler/posts', schedulerController.createPost);
app.get('/api/v1/scheduler/posts', schedulerController.listPosts);

const PORT = process.env.PORT || 8080;
app.listen(PORT, () => {
  console.log(`🚀 AutoSMM Node Backend running on port ${PORT}...`);
  schedulerService.startWorker();
});
