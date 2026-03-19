package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(10 * time.Minute)

	for i := 0; i < 20; i++ {
		if pingErr := db.Ping(); pingErr == nil {
			return db, nil
		} else {
			err = pingErr
			time.Sleep(2 * time.Second)
		}
	}

	return nil, fmt.Errorf("ping database: %w", err)
}

func EnsureSchema(db *sql.DB) error {
	statements := []string{
		`DROP TABLE IF EXISTS messages`,
		`DROP TABLE IF EXISTS recharge_records`,
		`DROP TABLE IF EXISTS meal_records`,
		`DROP TABLE IF EXISTS submissions`,
		`DROP TABLE IF EXISTS registrations`,
		`DROP TABLE IF EXISTS score_records`,
		`DROP TABLE IF EXISTS activities`,
		`DROP TABLE IF EXISTS users`,
		`CREATE TABLE users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(64) NOT NULL UNIQUE,
			display_name VARCHAR(128) NOT NULL,
			role VARCHAR(16) NOT NULL,
			password VARCHAR(128) NOT NULL,
			student_no VARCHAR(32) NOT NULL DEFAULT '',
			campus_card_no VARCHAR(32) NOT NULL DEFAULT '',
			bank_card_no VARCHAR(32) NOT NULL DEFAULT '',
			campus_card_balance DECIMAL(10,3) NOT NULL DEFAULT 0,
			bank_card_balance DECIMAL(10,3) NOT NULL DEFAULT 0,
			phone VARCHAR(32) NOT NULL DEFAULT '',
			dormitory VARCHAR(64) NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE meal_records (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			canteen_window VARCHAR(64) NOT NULL,
			dish_name VARCHAR(128) NOT NULL,
			amount DECIMAL(10,2) NOT NULL,
			consumed_at DATETIME NOT NULL
		)`,
		`CREATE TABLE recharge_records (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			raw_amount DECIMAL(10,3) NOT NULL,
			credited_amount DECIMAL(10,3) NOT NULL,
			note VARCHAR(255) NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL
		)`,
		`CREATE TABLE messages (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			student_name VARCHAR(128) NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME NOT NULL
		)`,
	}

	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}

	return nil
}

func SeedDemoData(db *sql.DB, _ string) error {
	statements := []string{
		`INSERT INTO users (id, username, display_name, role, password, student_no, campus_card_no, bank_card_no, campus_card_balance, bank_card_balance, phone, dormitory) VALUES
			(1, 'student001', '董小含', 'student', 'Meal2026#001', '2024001', 'CARD-001', '6222024001001', 42.500, 680.220, '13800000001', '1栋101'),
			(2, 'student002', '李小娜', 'student', 'Meal2026#002', '2024002', 'CARD-002', '6222024001002', 58.200, 512.800, '13800000002', '1栋102'),
			(3, 'student003', '王小茜', 'student', 'Meal2026#003', '2024003', 'CARD-003', '6222024001003', 16.800, 900.560, '13800000003', '2栋201'),
			(4, 'student004', '杨小昊', 'student', 'Meal2026#004', '2024004', 'CARD-004', '6222024001004', 88.000, 355.120, '13800000004', '2栋202'),
			(5, 'student005', '孙小胤', 'student', 'Meal2026#005', '2024005', 'CARD-005', '6222024001005', 23.100, 430.000, '13800000005', '3栋301'),
			(6, 'student006', '王小毅', 'student', 'Meal2026#006', '2024006', 'CARD-006', '6222024001006', 12.000, 288.880, '13800000006', '3栋302'),
			(7, 'admin', '食堂系统管理员', 'admin', 'admin123', '', 'ADMIN-CARD-001', 'ADMIN-BANK-001', 0, 0, '13900000000', '办公楼A-301')`,
		`INSERT INTO meal_records (user_id, canteen_window, dish_name, amount, consumed_at) VALUES
			(1, '带美二楼红烧麻辣面', '肉沫汤面', 12.50, DATE_SUB(NOW(), INTERVAL 2 DAY)),
			(1, '天美一楼', '红油刀削面', 11.00, DATE_SUB(NOW(), INTERVAL 1 DAY)),
			(2, '带美二楼', '重庆鸡公煲', 18.50, DATE_SUB(NOW(), INTERVAL 3 DAY)),
			(2, '小美一楼', '章丘炒坤', 8.80, DATE_SUB(NOW(), INTERVAL 1 DAY)),
			(3, '小美二楼', '成都小面', 15.00, DATE_SUB(NOW(), INTERVAL 4 DAY)),
			(4, '带美一楼', '馄饨', 19.90, DATE_SUB(NOW(), INTERVAL 2 DAY)),
			(5, '小美四楼', '热干面', 9.50, DATE_SUB(NOW(), INTERVAL 1 DAY)),
			(6, '杏苑', '自助餐', 7.00, DATE_SUB(NOW(), INTERVAL 5 DAY))`,
		`INSERT INTO recharge_records (user_id, raw_amount, credited_amount, note, created_at) VALUES
			(1, 20.000, 20.000, '银行卡向校园卡转入 20 元', DATE_SUB(NOW(), INTERVAL 7 DAY)),
			(2, 30.000, 30.000, '银行卡向校园卡转入 30 元', DATE_SUB(NOW(), INTERVAL 6 DAY)),
			(4, 50.000, 50.000, '银行卡向校园卡转入 50 元', DATE_SUB(NOW(), INTERVAL 5 DAY))`,
		`INSERT INTO messages (user_id, student_name, content, created_at) VALUES
			(1, '董小含', '肉沫汤面神中神（手动狗头）。', DATE_SUB(NOW(), INTERVAL 2 DAY)),
			(2, '孙小胤', '鸡公煲，鸡公煲，经过哦滴胃~~~。', DATE_SUB(NOW(), INTERVAL 1 DAY))`,
	}

	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}

	return nil
}

func ResetDemoData(db *sql.DB, uploadDir string) error {
	if err := EnsureSchema(db); err != nil {
		return err
	}

	return SeedDemoData(db, uploadDir)
}
