ALTER TABLE notifications
MODIFY COLUMN type ENUM('incoming_transfer', 'outgoing_transfer', 'debt_reminder', 'debt_cancel') NOT NULL;
