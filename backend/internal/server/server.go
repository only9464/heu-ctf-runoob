package server

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"campushub-ctf/internal/config"
	"campushub-ctf/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const jwtSecret = "123456"

type Server struct {
	cfg config.Config
	db  *sql.DB
}

type authUser struct {
	ID                int     `json:"id"`
	Username          string  `json:"username"`
	DisplayName       string  `json:"displayName"`
	Role              string  `json:"role"`
	StudentNo         string  `json:"studentNo"`
	CampusCardNo      string  `json:"campusCardNo"`
	BankCardNo        string  `json:"bankCardNo"`
	CampusCardBalance float64 `json:"campusCardBalance"`
	BankCardBalance   float64 `json:"bankCardBalance"`
}

type jwtClaims struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Role        string `json:"role"`
	StudentNo   string `json:"studentNo"`
	jwt.RegisteredClaims
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string   `json:"token"`
	User  authUser `json:"user"`
}

type rechargeRequest struct {
	Amount float64 `json:"amount"`
}

type messageRequest struct {
	Content string `json:"content"`
}

type adminSummary struct {
	StudentCount  int    `json:"studentCount"`
	MessageCount  int    `json:"messageCount"`
	RechargeCount int    `json:"rechargeCount"`
	AdminUsername string `json:"adminUsername"`
}

func New(cfg config.Config, db *sql.DB) *Server {
	return &Server{cfg: cfg, db: db}
}

func (s *Server) Run() error {
	router := gin.Default()
	router.Use(corsMiddleware())

	api := router.Group("/api")
	{
		api.POST("/student/login", s.handleStudentLogin)
		api.POST("/admin/login", s.handleAdminLogin)
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}

	auth := api.Group("/")
	auth.Use(s.authRequired())
	{
		auth.GET("/auth/me", s.handleMe)
		auth.POST("/auth/logout", s.handleLogout)
	}

	student := auth.Group("/student")
	student.Use(s.studentOnly())
	{
		student.GET("/profile", s.handleStudentProfile)
		student.GET("/consumptions", s.handleStudentConsumptions)
		student.POST("/recharge", s.handleRecharge)
		student.GET("/recharges", s.handleRechargeHistory)
		student.GET("/messages", s.handleStudentMessages)
		student.POST("/messages", s.handleCreateMessage)
	}

	admin := auth.Group("/admin")
	admin.Use(s.adminOnly())
	{
		admin.GET("/summary", s.handleAdminSummary)
		admin.GET("/students", s.handleAdminStudents)
		admin.GET("/students/:id", s.handleAdminStudentDetail)
		admin.GET("/messages", s.handleAdminMessages)
		admin.DELETE("/messages/:id", s.handleAdminDeleteMessage)
		admin.POST("/reset-demo", s.handleResetDemo)
	}

	return router.Run(s.cfg.Addr)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func (s *Server) authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimSpace(strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer "))
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		claims := &jwtClaims{}
		parsedToken, err := jwt.ParseWithClaims(token, claims, func(parsedToken *jwt.Token) (any, error) {
			if _, ok := parsedToken.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %s", parsedToken.Method.Alg())
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !parsedToken.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid jwt token"})
			return
		}

		user, err := s.loadAuthUser(claims.ID, claims.Role)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token user not found"})
			return
		}

		c.Set("currentUser", user)
		c.Set("jwtClaims", claims)
		c.Next()
	}
}

func (s *Server) studentOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if currentUser(c).Role != "student" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "student role required"})
			return
		}

		c.Next()
	}
}

func (s *Server) adminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if currentUser(c).Role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin role required"})
			return
		}

		c.Next()
	}
}

func currentUser(c *gin.Context) authUser {
	value, _ := c.Get("currentUser")
	user, _ := value.(authUser)
	return user
}

func (s *Server) handleStudentLogin(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// INTENTIONALLY VULNERABLE: student login uses string concatenation for SQLi teaching.
	rawSQL := fmt.Sprintf(
		"SELECT id, username, display_name, role, student_no, campus_card_no, bank_card_no, campus_card_balance, bank_card_balance FROM users WHERE role = 'student' AND username = '%s' AND password = '%s' LIMIT 1",
		req.Username,
		req.Password,
	)

	var user authUser
	if err := s.db.QueryRow(rawSQL).Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Role,
		&user.StudentNo,
		&user.CampusCardNo,
		&user.BankCardNo,
		&user.CampusCardBalance,
		&user.BankCardBalance,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "学生账号或密码错误"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "登录失败",
			"detail": err.Error(),
		})
		return
	}

	if user.Role != "student" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "学生登录入口仅允许学生身份"})
		return
	}

	s.respondWithJWT(c, user)
}

func (s *Server) handleAdminLogin(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user authUser
	row := s.db.QueryRow(
		`SELECT id, username, display_name, role, student_no, campus_card_no, bank_card_no, campus_card_balance, bank_card_balance FROM users WHERE role = 'admin' AND username = ? AND password = ? LIMIT 1`,
		req.Username,
		req.Password,
	)

	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Role,
		&user.StudentNo,
		&user.CampusCardNo,
		&user.BankCardNo,
		&user.CampusCardBalance,
		&user.BankCardBalance,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "管理员账号或密码错误"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s.respondWithJWT(c, user)
}

func (s *Server) respondWithJWT(c *gin.Context, user authUser) {
	token, err := generateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "签发 JWT 失败"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		Token: token,
		User:  user,
	})
}

func generateJWT(user authUser) (string, error) {
	claims := jwtClaims{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Role:        user.Role,
		StudentNo:   user.StudentNo,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (s *Server) handleMe(c *gin.Context) {
	user := currentUser(c)
	if user.Role == "student" {
		student, err := s.currentStudentFromJWT(c)
		if err == nil {
			user = student
		}
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (s *Server) handleLogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "JWT 为无状态认证，前端已清除本地 token"})
}

func (s *Server) handleStudentProfile(c *gin.Context) {
	// INTENTIONALLY VULNERABLE: the endpoint fully trusts studentNo in the JWT.
	// With weak secret 123456, attackers can forge studentNo to read another student.
	item, err := s.studentProfileFromJWT(c)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, errStudentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "学生不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profile": item,
	})
}

func (s *Server) handleStudentConsumptions(c *gin.Context) {
	student, err := s.currentStudentFromJWT(c)
	if err != nil {
		if errors.Is(err, errStudentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "学生不存在"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	rows, err := s.db.Query(
		`SELECT id, canteen_window, dish_name, amount, consumed_at
		 FROM meal_records
		 WHERE user_id = ?
		 ORDER BY consumed_at DESC, id DESC`,
		student.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	items := make([]gin.H, 0)
	for rows.Next() {
		var id int
		var canteenWindow string
		var dishName string
		var amount float64
		var consumedAt time.Time
		if err := rows.Scan(&id, &canteenWindow, &dishName, &amount, &consumedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		items = append(items, gin.H{
			"id":            id,
			"canteenWindow": canteenWindow,
			"dishName":      dishName,
			"amount":        amount,
			"consumedAt":    consumedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (s *Server) handleRecharge(c *gin.Context) {
	var req rechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student, err := s.currentStudentFromJWT(c)
	if err != nil {
		if errors.Is(err, errStudentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "学生不存在"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// INTENTIONALLY VULNERABLE:
	// should require positive integer amount and deduct the same amount from bank card,
	// but this logic rounds sub-1 values up to 1 and turns negatives into credits.
	creditedAmount := math.Abs(req.Amount)
	if creditedAmount < 1 {
		creditedAmount = 1
	}

	if _, err := s.db.Exec(
		`UPDATE users
		 SET campus_card_balance = campus_card_balance + ?,
		     bank_card_balance = bank_card_balance - ?
		 WHERE id = ?`,
		creditedAmount,
		req.Amount,
		student.ID,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if _, err := s.db.Exec(
		`INSERT INTO recharge_records (user_id, raw_amount, credited_amount, note, created_at) VALUES (?, ?, ?, ?, NOW())`,
		student.ID,
		req.Amount,
		creditedAmount,
		"教学环境：银行卡向校园卡转账时存在非法金额入账漏洞",
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var campusCardBalance float64
	var bankBalance float64
	_ = s.db.QueryRow(`SELECT campus_card_balance, bank_card_balance FROM users WHERE id = ?`, student.ID).Scan(&campusCardBalance, &bankBalance)

	c.JSON(http.StatusOK, gin.H{
		"message":           "转入成功",
		"rawAmount":         req.Amount,
		"creditedAmount":    creditedAmount,
		"campusCardBalance": campusCardBalance,
		"bankBalance":       bankBalance,
		"note":              "0.001 会让校园卡至少增加 1 元，-1 会导致校园卡加钱且银行卡也增加 1 元。",
	})
}

func (s *Server) handleRechargeHistory(c *gin.Context) {
	student, err := s.currentStudentFromJWT(c)
	if err != nil {
		if errors.Is(err, errStudentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "学生不存在"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	rows, err := s.db.Query(
		`SELECT id, raw_amount, credited_amount, note, created_at
		 FROM recharge_records
		 WHERE user_id = ?
		 ORDER BY id DESC`,
		student.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	items := make([]gin.H, 0)
	for rows.Next() {
		var id int
		var rawAmount float64
		var creditedAmount float64
		var note string
		var createdAt time.Time
		if err := rows.Scan(&id, &rawAmount, &creditedAmount, &note, &createdAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		items = append(items, gin.H{
			"id":             id,
			"rawAmount":      rawAmount,
			"creditedAmount": creditedAmount,
			"note":           note,
			"createdAt":      createdAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (s *Server) handleStudentMessages(c *gin.Context) {
	if _, err := s.currentStudentFromJWT(c); err != nil {
		if errors.Is(err, errStudentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "学生不存在"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	rows, err := s.db.Query(
		`SELECT id, student_name, content, created_at
		 FROM messages
		 ORDER BY id DESC`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	items := make([]gin.H, 0)
	for rows.Next() {
		var id int
		var studentName string
		var content string
		var createdAt time.Time
		if err := rows.Scan(&id, &studentName, &content, &createdAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		items = append(items, gin.H{
			"id":          id,
			"studentName": studentName,
			"content":     content,
			"createdAt":   createdAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (s *Server) handleCreateMessage(c *gin.Context) {
	var req messageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student, err := s.currentStudentFromJWT(c)
	if err != nil {
		if errors.Is(err, errStudentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "学生不存在"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if _, err := s.db.Exec(
		`INSERT INTO messages (user_id, student_name, content, created_at) VALUES (?, ?, ?, NOW())`,
		student.ID,
		student.DisplayName,
		req.Content,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "留言已提交"})
}

func (s *Server) handleAdminSummary(c *gin.Context) {
	var summary adminSummary
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM users WHERE role = 'student'`).Scan(&summary.StudentCount)
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM messages`).Scan(&summary.MessageCount)
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM recharge_records`).Scan(&summary.RechargeCount)
	summary.AdminUsername = "admin"

	c.JSON(http.StatusOK, gin.H{"summary": summary})
}

func (s *Server) handleAdminStudents(c *gin.Context) {
	rows, err := s.db.Query(
		`SELECT id, username, display_name, student_no, campus_card_balance, bank_card_balance, phone, dormitory
		 FROM users
		 WHERE role = 'student'
		 ORDER BY id`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	items := make([]gin.H, 0)
	for rows.Next() {
		var id int
		var username string
		var displayName string
		var studentNo string
		var campusCardBalance float64
		var bankBalance float64
		var phone string
		var dormitory string
		if err := rows.Scan(&id, &username, &displayName, &studentNo, &campusCardBalance, &bankBalance, &phone, &dormitory); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		items = append(items, gin.H{
			"id":                id,
			"username":          username,
			"displayName":       displayName,
			"studentNo":         studentNo,
			"campusCardBalance": campusCardBalance,
			"bankBalance":       bankBalance,
			"phone":             phone,
			"dormitory":         dormitory,
		})
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (s *Server) handleAdminStudentDetail(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var detail struct {
		ID                int     `json:"id"`
		Username          string  `json:"username"`
		DisplayName       string  `json:"displayName"`
		StudentNo         string  `json:"studentNo"`
		CampusCardNo      string  `json:"campusCardNo"`
		BankCardNo        string  `json:"bankCardNo"`
		CampusCardBalance float64 `json:"campusCardBalance"`
		BankCardBalance   float64 `json:"bankBalance"`
		Phone             string  `json:"phone"`
		Dormitory         string  `json:"dormitory"`
	}
	err = s.db.QueryRow(
		`SELECT id, username, display_name, student_no, campus_card_no, bank_card_no, campus_card_balance, bank_card_balance, phone, dormitory
		 FROM users
		 WHERE id = ? AND role = 'student'`,
		targetID,
	).Scan(
		&detail.ID,
		&detail.Username,
		&detail.DisplayName,
		&detail.StudentNo,
		&detail.CampusCardNo,
		&detail.BankCardNo,
		&detail.CampusCardBalance,
		&detail.BankCardBalance,
		&detail.Phone,
		&detail.Dormitory,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "学生不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"student": detail})
}

func (s *Server) handleAdminMessages(c *gin.Context) {
	rows, err := s.db.Query(
		`SELECT id, student_name, content, created_at
		 FROM messages
		 ORDER BY id DESC`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	items := make([]gin.H, 0)
	for rows.Next() {
		var id int
		var studentName string
		var content string
		var createdAt time.Time
		if err := rows.Scan(&id, &studentName, &content, &createdAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		items = append(items, gin.H{
			"id":          id,
			"studentName": studentName,
			"content":     content,
			"createdAt":   createdAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (s *Server) handleAdminDeleteMessage(c *gin.Context) {
	messageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message id"})
		return
	}

	result, err := s.db.Exec(`DELETE FROM messages WHERE id = ?`, messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "留言不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "留言已删除"})
}

func (s *Server) handleResetDemo(c *gin.Context) {
	if err := database.ResetDemoData(s.db, s.cfg.UploadDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "演示环境已重置，默认账号、双卡余额与留言数据已恢复"})
}

func (s *Server) loadAuthUser(id int, role string) (authUser, error) {
	var user authUser
	err := s.db.QueryRow(
		`SELECT id, username, display_name, role, student_no, campus_card_no, bank_card_no, campus_card_balance, bank_card_balance
		 FROM users
		 WHERE id = ? AND role = ? LIMIT 1`,
		id,
		role,
	).Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Role,
		&user.StudentNo,
		&user.CampusCardNo,
		&user.BankCardNo,
		&user.CampusCardBalance,
		&user.BankCardBalance,
	)
	if err != nil {
		return authUser{}, err
	}

	return user, nil
}

var errStudentNotFound = errors.New("student from jwt not found")

func currentClaims(c *gin.Context) jwtClaims {
	value, _ := c.Get("jwtClaims")
	claims, _ := value.(*jwtClaims)
	if claims == nil {
		return jwtClaims{}
	}

	return *claims
}

func (s *Server) currentStudentFromJWT(c *gin.Context) (authUser, error) {
	claims := currentClaims(c)
	if claims.Role != "student" || strings.TrimSpace(claims.StudentNo) == "" {
		return authUser{}, errStudentNotFound
	}

	var user authUser
	err := s.db.QueryRow(
		`SELECT id, username, display_name, role, student_no, campus_card_no, bank_card_no, campus_card_balance, bank_card_balance
		 FROM users
		 WHERE role = 'student' AND student_no = ? LIMIT 1`,
		claims.StudentNo,
	).Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Role,
		&user.StudentNo,
		&user.CampusCardNo,
		&user.BankCardNo,
		&user.CampusCardBalance,
		&user.BankCardBalance,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return authUser{}, errStudentNotFound
		}
		return authUser{}, err
	}

	return user, nil
}

func (s *Server) studentProfileFromJWT(c *gin.Context) (gin.H, error) {
	claims := currentClaims(c)
	if claims.Role != "student" || strings.TrimSpace(claims.StudentNo) == "" {
		return nil, errStudentNotFound
	}

	var item struct {
		ID                int     `json:"id"`
		Username          string  `json:"username"`
		DisplayName       string  `json:"displayName"`
		StudentNo         string  `json:"studentNo"`
		CampusCardNo      string  `json:"campusCardNo"`
		BankCardNo        string  `json:"bankCardNo"`
		CampusCardBalance float64 `json:"campusCardBalance"`
		BankCardBalance   float64 `json:"bankBalance"`
		Phone             string  `json:"phone"`
		Dormitory         string  `json:"dormitory"`
	}

	err := s.db.QueryRow(
		`SELECT id, username, display_name, student_no, campus_card_no, bank_card_no, campus_card_balance, bank_card_balance, phone, dormitory
		 FROM users
		 WHERE role = 'student' AND student_no = ? LIMIT 1`,
		claims.StudentNo,
	).Scan(
		&item.ID,
		&item.Username,
		&item.DisplayName,
		&item.StudentNo,
		&item.CampusCardNo,
		&item.BankCardNo,
		&item.CampusCardBalance,
		&item.BankCardBalance,
		&item.Phone,
		&item.Dormitory,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errStudentNotFound
		}
		return nil, err
	}

	return gin.H{
		"id":                item.ID,
		"username":          item.Username,
		"displayName":       item.DisplayName,
		"studentNo":         item.StudentNo,
		"campusCardNo":      item.CampusCardNo,
		"bankCardNo":        item.BankCardNo,
		"campusCardBalance": item.CampusCardBalance,
		"bankBalance":       item.BankCardBalance,
		"phone":             item.Phone,
		"dormitory":         item.Dormitory,
	}, nil
}

func requestedUserID(c *gin.Context, fallback int) int {
	rawID := c.Query("id")
	if rawID == "" {
		return fallback
	}

	id, err := strconv.Atoi(rawID)
	if err != nil || id <= 0 {
		return fallback
	}

	return id
}
