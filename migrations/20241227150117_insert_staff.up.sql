INSERT INTO users (`email`, `name`, `role_id`, `phone_number`, `password`)
VALUES
('admin@internetbanking.com', 'Staff',  (SELECT id FROM roles WHERE name = 'staff'), '0123987456', '$2a$10$ka/P7URqQrKweoQo.9491.yNl3sBA.vm7LHlhJRtRugkM/U8jR5Dy');