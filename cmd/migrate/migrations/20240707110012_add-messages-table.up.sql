CREATE TABLE IF NOT EXISTS messages (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `conversation_id` INT UNSIGNED,
    `sender_id` INT UNSIGNED,
    `content` TEXT NOT NULL,
    `sentAt` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `readAt` TIMESTAMP NULL,

    FOREIGN KEY (conversation_id) REFERENCES conversations(id),
    FOREIGN KEY (sender_id) REFERENCES users(id)
);
