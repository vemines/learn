use myDB;


CREATE DATABASE myDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

SET AUTOCOMMIT = OFF;
COMMIT;
ROLLBACK;

RENAME TABLE employees TO workers;

CREATE TABLE employees (
    employee_id INT PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50), 
    hourly_pay DECIMAL(5, 2),
		job VARCHAR(50),
    hire_date DATE,
    supervisor_id INT,
    CONSTRAINT chk_hourly_pay CHECK (hourly_pay >= 10.00)
);
DROP TABLE employees;

ALTER TABLE employees ADD phone_number VARCHAR(15);
ALTER TABLE employees RENAME COLUMN phone_number TO email;
ALTER TABLE employees MODIFY COLUMN email VARCHAR(15);
ALTER TABLE employees  MODIFY email VARCHAR(100) AFTER last_name;
ALTER TABLE employees  MODIFY email VARCHAR(100) FIRST;
ALTER TABLE employees DROP COLUMN email;
ALTER TABLE employees ADD CONSTRAINT chk_hourly_pay CHECK (hourly_pay >= 10.00);
ALTER TABLE employees DROP CHECK chk_hourly_pay;

INSERT INTO employees VALUES
	(1, "Eugene", "Krabs", 25.50, "manager", "2023-01-02", NULL),
	(2, "Squidward", "Tentacles", 15.00, "cashier", "2023-01-03", 5), 
	(3, "Spongebob", "Squarepants", 12.50, "cook", "2023-01-04", 5), 
	(4, "Patrick", "Star", 12.50, "cook","2023-01-05", 5), 
	(5, "Sandy", "Cheeks", 17.25, "asst. manager","2023-01-06", 1);
    
INSERT INTO employees (employee_id, first_name, last_name, hourly_pay, job, hire_date,supervisor_id) VALUES 
	(6, "Sheldon", "Plankton", 20, "janitor","2023-01-07", 5);

SELECT * FROM employees;

SELECT first_name, last_name FROM employees;
SELECT * FROM employees WHERE employee_id = 1;
SELECT * FROM employees WHERE first_name = "Spongebob";
SELECT * FROM employees WHERE hourly_pay >= 15;
SELECT hire_date, first_name FROM employees WHERE hire_date <= "2023-01-03";
SELECT * FROM employees WHERE employee_id != 1;
SELECT * FROM employees WHERE supervisor_id IS NULL;
SELECT * FROM employees WHERE supervisor_id IS NOT NULL;

SELECT * FROM employees WHERE hire_date > '2023-01-02' AND job = "cook";
SELECT * FROM employees WHERE job = 'Cook' OR job = 'Cashier';
SELECT * FROM employees WHERE NOT job = 'Manager';
SELECT * FROM employees WHERE NOT job = 'Manager' AND NOT job = "asst. manager";
SELECT * FROM employees WHERE hire_date BETWEEN '2023-01-04' AND '2023-01-07';
SELECT * FROM employees WHERE job IN ("cook", "cashier", "janitor");

SELECT * FROM employees WHERE first_name LIKE "s%";
SELECT * FROM employees WHERE last_name LIKE "%r";
SELECT * FROM employees WHERE hire_date LIKE "2023%";
SELECT * FROM employees WHERE job LIKE "_ook";
SELECT * FROM employees WHERE hire_date LIKE "____-01-__";
SELECT * FROM employees WHERE job LIKE "_a%";

SELECT a.first_name, a.last_name, CONCAT(b.first_name," ", b.last_name) AS "reports_to"
FROM employees AS a
INNER JOIN employees AS b 
ON a.supervisor_id = b.employee_id;

SELECT first_name, last_name, hourly_pay, (SELECT AVG(hourly_pay) FROM employees) AS avg_pay FROM employees;

UPDATE employees SET hourly_pay = 10.25 WHERE employee_id = 6;

DELETE FROM employees WHERE employee_id = 6;

CREATE VIEW employee_attendance AS 
SELECT first_name, last_name FROM employees;

SELECT * FROM employee_attendance;

CREATE TABLE testTime(
 my_date DATE,
     my_time TIME,
     my_datetime DATETIME
);

INSERT INTO testTime VALUES(CURRENT_DATE(), CURRENT_TIME(), NOW());
INSERT INTO testTime VALUES(CURRENT_DATE() + 1, CURRENT_TIME(), NOW());
SELECT * FROM testTime;

CREATE TABLE products (
    product_id INT,
    product_name varchar(25) UNIQUE,
    price DECIMAL(4, 2) NOT NULL DEFAULT 0
);

ALTER TABLE products ADD CONSTRAINT UNIQUE (product_name);
ALTER TABLE products MODIFY price DECIMAL(4, 2) NOT NULL;
ALTER TABLE products ALTER price SET DEFAULT 0;

INSERT INTO products VALUES 
	(100, 'hamburger', 3.99),
	(101, 'fries', 1.89),
	(102, 'soda', 1.00),
	(103, "ice cream", 1.49);

INSERT INTO products (product_id, product_name) VALUES 
	(104, "straw"), 
	(105, "napkin"), 
	(106, "fork"), 
	(107, "spoon");    

SELECT * FROM products;

CREATE TABLE customers (
	customer_id INT PRIMARY KEY AUTO_INCREMENT,
	first_name VARCHAR(50),
	last_name VARCHAR(50)
);

INSERT INTO customers (first_name, last_name) VALUES  
	("Fred", "Fish"),
	("Larry", "Lobster"),
	("Bubble", "Bass"),
	("Poppy", "Puff");
	
SELECT * FROM customers;

SELECT * FROM customers LIMIT 2;
SELECT * FROM customers ORDER BY last_name DESC LIMIT 3;
SELECT * FROM customers ORDER BY last_name ASC LIMIT 3;
SELECT * FROM customers LIMIT 1, 3; // Skip to 1 and limit 3

-- Indexes are used to find values within a specific column more quickly, but UPDATE slower
CREATE INDEX last_name_idx ON customers (last_name);
CREATE INDEX last_name_first_name_idx ON customers (last_name, first_name);
DROP INDEX last_name_idx ON customers;

CREATE TABLE transactions(
	transaction_id INT PRIMARY KEY AUTO_INCREMENT,
	amount DECIMAL(5, 2),
    customer_id INT,
	order_date DATETIME DEFAULT NOW(),
    FOREIGN KEY (customer_id) REFERENCES customers(customer_id)
    ON DELETE SET NULL -- ON DELETE CASCADE
);
ALTER TABLE transactions ADD CONSTRAINT PRIMARY KEY (transaction_id);
ALTER TABLE transactions AUTO_INCREMENT = 1000;
ALTER TABLE transactions ADD CONSTRAINT fk_customer_id FOREIGN KEY (customer_id) REFERENCES customers(customer_id) ON DELETE SET NULL;
ALTER TABLE transactions ADD CONSTRAINT fk_customer_id FOREIGN KEY (customer_id) REFERENCES customers(customer_id) ON DELETE CASCADE;

INSERT INTO transactions VALUES  
	(1000, 4.99, 3, "2023-01-01"),
	(1001, 2.89, 2, "2023-01-01"),
	(1002, 3.38, 3, "2023-01-02"),
	(1003, 4.99, 1, "2023-01-02"),
	(1004, 1.00, NULL, "2023-01-03"),
	(1005, 2.49, 4, "2023-01-03"),
	(1006, 5.48, NULL, "2023-01-03");

SELECT * FROM transactions;

SELECT * FROM transactions ORDER BY amount DESC, customer_id DESC;

SELECT SUM(amount), order_date FROM transactions GROUP BY order_date WITH ROLLUP;
SELECT COUNT(transaction_id) AS "# of orders", order_date FROM transactions GROUP BY order_date WITH ROLLUP;
SELECT COUNT(transaction_id) AS "# of orders", customer_id FROM transactions GROUP BY customer_id WITH ROLLUP;

SELECT COUNT(amount), customer_id FROM transactions GROUP BY customer_id HAVING COUNT(amount) > 1 AND customer_id IS NOT NULL;

-- INNER JOIN selects records that have a matching key in both tables.
SELECT * 
FROM transactions INNER JOIN customers
ON transactions.customer_id = customers.customer_id;

-- LEFT JOIN returns all records from the left table 
-- and the matching records from the right table
SELECT *
FROM transactions LEFT JOIN customers
ON transactions.customer_id = customers.customer_id;

-- RIGHT JOIN returns all records from the right table 
-- and the matching records from the left table
SELECT *
FROM transactions RIGHT JOIN customers
ON transactions.customer_id = customers.customer_id;

SELECT COUNT(amount) as count FROM transactions;
SELECT MAX(amount) AS maximum FROM transactions;
SELECT MIN(amount) AS minimum FROM transactions;
SELECT AVG(amount) AS average FROM transactions;
SELECT SUM(amount) AS sum FROM transactions;
SELECT CONCAT(first_name, " ", last_name) AS full_name FROM employees;

-- NO DUPLICATES
SELECT first_name, last_name FROM employees
UNION
SELECT first_name, last_name FROM customers;

-- DUPLICATES OK
SELECT first_name, last_name FROM employees
UNION ALL
SELECT first_name, last_name FROM customers;

DELIMITER $$
CREATE PROCEDURE find_customer(IN f_name VARCHAR(50), IN l_name VARCHAR(50))  
BEGIN  
 SELECT *
 FROM customers
 WHERE first_name = f_name AND last_name = l_name;
END $$
DELIMITER ;

CALL find_customer("Larry", "Lobster");

DROP PROCEDURE find_customer;

CREATE TRIGGER before_hourly_pay_update BEFORE UPDATE ON employees 
FOR EACH ROW
SET NEW.salary = (NEW.hourly_pay * 2080);

CREATE TRIGGER before_hourly_pay_insert BEFORE INSERT ON employees 
FOR EACH ROW
SET NEW.salary = (NEW.hourly_pay * 2080);

CREATE TRIGGER after_salary_delete AFTER DELETE ON employees 
FOR EACH ROW
UPDATE expenses SET expense_total = expense_total - OLD.salary WHERE expense_name = "salaries";

CREATE TRIGGER after_salary_insert AFTER INSERT ON employees 
FOR EACH ROW
UPDATE expenses SET expense_total = expense_total + NEW.salary WHERE expense_name = "salaries";

CREATE TRIGGER after_salary_update AFTER UPDATE ON employees 
FOR EACH ROW
UPDATE expenses SET expense_total = expense_total + (NEW.salary - OLD.salary) WHERE expense_name = "salaries";
