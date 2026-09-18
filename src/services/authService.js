
const bcrypt = require('bcrypt');
const jwt = require('jsonwebtoken');
const db = require('../db');
class AuthService {
  async register(username, email, password, fullName) {
    const hashedPassword = await bcrypt.hash(password, 10);
    const query = 'INSERT INTO users (username, email, password_hash, full_name, role, plan) VALUES ($1, $2, $3, $4, \'user\', \'freemium\') RETURNING id';
    const result = await db.query(query, [username, email, hashedPassword, fullName]);
    return result.rows[0];
  }
  async login(username, password) {
    const query = 'SELECT * FROM users WHERE username = $1';
    const result = await db.query(query, [username]);
    const user = result.rows[0];
    if (!user || !(await bcrypt.compare(password, user.password_hash))) throw new Error('Invalid credentials');
    const token = jwt.sign({ userId: user.id, role: user.role }, process.env.JWT_SECRET || 'secret', { expiresIn: '72h' });
    return { token, user: { id: user.id, username: user.username, role: user.role } };
  }
}
module.exports = new AuthService();
