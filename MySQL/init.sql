USE mydb;

DROP TABLE IF EXISTS employee;

CREATE TABLE employee (
    uid INT(10) NOT NULL AUTO_INCREMENT,
    name VARCHAR(64) NULL DEFAULT NULL,
    PRIMARY KEY (uid)
);

COMMIT;