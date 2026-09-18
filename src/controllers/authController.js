
const authService = require('../services/authService');
exports.register = async (req, res) => {
  try {
    const { username, email, password, full_name } = req.body;
    const user = await authService.register(username, email, password, full_name);
    res.status(201).json({ message: 'User created successfully', userId: user.id });
  } catch (err) { res.status(500).json({ error: err.message }); }
};
exports.login = async (req, res) => {
  try {
    const { username, password } = req.body;
    const result = await authService.login(username, password);
    res.json(result);
  } catch (err) { res.status(401).json({ error: err.message }); }
};
