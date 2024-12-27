ALTER TABLE authentications
ADD COLUMN device_token VARCHAR(255) AFTER refresh_token;
