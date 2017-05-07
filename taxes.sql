DROP DATABASE IF EXISTS tax;

CREATE DATABASE tax;

USE tax;

DROP TABLE IF EXISTS Tax;

CREATE TABLE Tax (id int NOT NULL AUTO_INCREMENT, Municipality varchar(40), PeriodStart DATE, PeriodEnd DATE, Rate DECIMAL(5, 2), PRIMARY KEY (id));

CREATE USER 'tax'@'%' IDENTIFIED BY 'tax';
GRANT ALL ON *.* TO 'tax'@'%';

CREATE USER 'tax'@'localhost' IDENTIFIED BY 'tax';
GRANT ALL ON *.* TO 'tax'@'localhost';

-- INSERT INTO Tax(Municipality, PeriodStart, PeriodEnd, Rate) values ('vilnius', '2017-01-01', '2017-12-31', 0.2);

-- INSERT INTO Tax(Municipality, PeriodStart, PeriodEnd, Rate) values ('vilnius', '2017-05-01', '2017-05-31', 0.4);

-- INSERT INTO Tax(Municipality, PeriodStart, PeriodEnd, Rate) values ('vilnius', '2017-01-01', '2017-01-01', 0.1);

-- INSERT INTO Tax(Municipality, PeriodStart, PeriodEnd, Rate) values ('vilnius', '2017-12-25', '2017-12-25', 0.1);

