CREATE TABLE `users`(
    `id` INT(11) AUTO_INCREMENT PRIMARY KEY,
    `name` varchar(256) NOT NULL,
    `phone` varchar(256) NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP()
);

INSERT INTO `users` (`name`, `phone`) VALUES('Daniel', '+79099898988');

INSERT INTO `users` (`name`, `phone`) VALUES('Jamie', '+38908767562');