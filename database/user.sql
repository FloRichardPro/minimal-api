CREATE USER 'foo_user'@'%' IDENTIFIED BY 'foo_user_dev';

GRANT SELECT, INSERT, UPDATE, DELETE ON foo_db.* TO 'foo_user'@'%';

CREATE USER 'dummy'@'%' IDENTIFIED BY 'dummy';

FLUSH PRIVILEGES;