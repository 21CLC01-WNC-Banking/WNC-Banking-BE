CREATE TABLE notifications (
   id INT AUTO_INCREMENT PRIMARY KEY,
   type ENUM('incoming_transfer', 'outgoing_transfer', 'debt_reminder') NOT NULL,
   title VARCHAR(255) NOT NULL,
   content TEXT NOT NULL,
   is_seen BOOLEAN DEFAULT FALSE,
   user_id INT NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   deleted_at TIMESTAMP NULL DEFAULT NULL,
   FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);