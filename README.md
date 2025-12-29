# URL Shortener - Go Backend

## Mục lục

- [Mô tả bài toán](#mô-tả-bài-toán)
- [Features](#features-đã-implement)
- [Demo](#demo)
- [Cách chạy project](#cách-chạy-project)
- [Cây thư mục Project](#cây-thư-mục-project)
- [API Documentation](#api-documentation)
- [Kiến trúc & Thiết kế](#kiến-trúc--thiết-kế)
- [Quyết định kỹ thuật](#quyết-định-kỹ-thuật)
- [Trade-offs](#trade-offs)
- [Thử thách & Giải pháp](#thử-thách--giải-pháp)
- [Hạn chế & Hướng cải thiện](#hạn-chế--hướng-cải-thiện-trong-tương-lai)

---

# Mô tả bài toán

## Hiểu bài toán

Bài toán yêu cầu xây dựng một backend service rút gọn URL tương tự như Bitly hoặc TinyURL.
Hệ thống cho phép người dùng chuyển đổi một URL dài thành một URL ngắn, dễ chia sẻ.
Khi truy cập URL ngắn, hệ thống sẽ redirect người dùng về URL gốc và đồng thời ghi nhận lượt click.
Ngoài các chức năng cơ bản, hệ thống cần được thiết kế để xử lý concurrency, đảm bảo hiệu năng,
và có khả năng mở rộng khi traffic tăng cao.

## Use cases thực tế:

- Marketing campaigns: Tracking hiệu quả từng kênh
- Social media: Links gọn gàng hơn trên Twitter, Instagram
- Print materials: QR codes với short URLs
- Analytics: Phân tích hành vi người dùng

---

## Features đã implement

### Core Features

- Tạo short URL từ long URL với collision handling
- Redirect từ short URL về original URL
- Liệt kê danh sách URLs với pagination
- Xem thông tin chi tiết của URL (clicks, created time)

### Advanced Features

- **Detailed Analytics**: Track IP, User Agent, Device Type, Country, Referer
- **Click Statistics**: Phân tích chi tiết từng lượt click
- **Async Click Tracking**: Không làm chậm redirect
- **Pagination**: Hỗ trợ phân trang cho tất cả list endpoints
- **Error Handling**: Xử lý lỗi đầy đủ với status codes chuẩn
- **Rate Limiting** : Tránh Spam Api quá nhiều (60 request/phút), simple, in-memory

---

## Demo

### Link demo (deployed)

```
🔗 https://golang-url-shortener-1u52.onrender.com
```

### API Endpoints

```
POST   /api/url/shorten                     - Tạo short URL
GET    /:short_code                     - Redirect về URL gốc
GET    /api/url                         - Danh sách URLs (pagination)
GET    /api/url/:url_id/stats           - Analytics chi tiết
GET    /api/url/:url_id/stats/count     - Click count
GET    /api/metrics                     - Key Metrics cho phân tích
GET    /health                          - Health check
```

### Ví dụ sử dụng

```bash
# 1. Tạo short URL
curl -X POST http://localhost:8080/api/url/shorten \
  -H "Content-Type: application/json" \
  -d '{"long_url": "https://github.com/golang/go"}'

# Response:
{
  "short_url": "http://localhost:8080/abc123"
}

# 2. Sử dụng short URL
curl -L http://localhost:8080/abc123
# → Redirect 302 về https://github.com/golang/go

# 3. Xem danh sách URLs
curl "http://localhost:8080/api/urls?limit=10&page=0"

# 4. Xem analytics
curl "http://localhost:8080/api/stats/1?limit=20&page=0"
```

---

## Cách chạy project

### Chạy với Docker

```bash
# Clone repository
git clone https://github.com/Vyzz1/golang-url-shortener url-shortener
cd url-shortener


docker-compose up -d

docker-compose logs -f app

# Test

curl http://localhost:8080/health
```

---

## Cây thư mục Project

```
golang-url-shortener/
├── .air.toml                    # Config cho Air
├── .env                         # Environment
├── docker-compose.yml           # Docker Compose config
├── Dockerfile                   # Docker image config
├── go.mod                       # Go module dependencies
├── go.sum                       # Go module checksums
├── main.go                      # Entry point của application
├── Makefile                     # Build automation scripts
├── README.md                    # Documentation
├── sqlc.yaml                    # sqlc configuration
├── wait-for-it.sh              # Script đợi database ready
│
├── api/                         # HTTP handlers & routing
│   ├── metrics.go              # Handler cho /api/metrics
│   ├── server.go               # Gin server setup & routes
│   └── url.go                  # Handlers cho URL operations
│
├── db/                          # Database layer
│   ├── migrations/             # Database migration files
│   │   ├── 000001_init_source.up.sql
│   │   ├── 000001_init_source.down.sql
│   │   ├── 000002_unique_orignal_url.up.sql
│   │   ├── 000002_unique_orignal_url.down.sql
│   │   ├── 000003_add_click_columns.up.sql
│   │   └── 000003_add_click_columns.down.sql
│   │
│   ├── query/                  # SQL queries (sqlc input)
│   │   ├── clicks.sql         # Click tracking queries
│   │   ├── stats.sql          # Statistics queries
│   │   └── urls.sql           # URL CRUD queries
│   │
│   └── sqlc/                   # Generated Go code (sqlc output)
│       ├── clicks.sql.go      # Generated click queries
│       ├── db.go              # Database connection
│       ├── models.go          # Generated models
│       ├── querier.go         # Query interface
│       ├── stats.sql.go       # Generated stats queries
│       ├── store.go           # Store implementation
│       └── urls.sql.go        # Generated URL queries
│
├── middlewares/                 # Gin middlewares
│   ├── cors.go                 # CORS middleware
│   └── rate-limit.go           # Rate limiting middleware
│
├── tmp/                         # Temporary files (Air hot reload)
│   ├── build-errors.log
│   └── main.exe
│
└── utils/                       # Utility functions
    ├── click.go                # Click tracking utilities
    ├── config.go               # Config loading
    └── shortcode.go            # Short code generation
```

**Giải thích cấu trúc:**

- **`main.go`**: Entry point - khởi tạo config, database, server
- **`api/`**: HTTP layer - Gin handlers, routing, request/response
- **`db/migrations/`**: SQL migration files cho versioning database schema
- **`db/query/`**: Raw SQL queries - input cho sqlc
- **`db/sqlc/`**: Type-safe Go code được generate từ SQL queries
- **`middlewares/`**: Reusable middlewares (CORS, rate limit)
- **`utils/`**: Helper functions (config, short code generation, click tracking)
- **`tmp/`**: Build artifacts cho hot reload

---

## API Documentation

### Base URL

```
http://localhost:8080
```

### 1. Tạo Short URL

**Endpoint:** `POST /api/url/shorten`

**Request:**

```json
{
  "long_url": "https://example.com/very/long/path"
}
```

**Response:** `200 OK`

```json
{
  "short_url": "http://localhost:8080/abc123"
}
```

**Error Responses:**

- `400 Bad Request`: URL không hợp lệ
- `409 Conflict`: Không thể tạo unique code (retry)
- `500 Internal Server Error`: Lỗi server

---

### 2. Redirect về URL gốc

**Endpoint:** `GET /:short_code`

**Example:** `GET /abc123`

**Response:** `302 Found`

```
Location: https://example.com/very/long/path
```

**Behavior:**

- Redirect về URL gốc
- Tự động track click (async, không làm chậm redirect)
- Lưu thông tin: IP, User Agent, Device Type, Country, Referer

**Error Responses:**

- `404 Not Found`: Short code không tồn tại
- `500 Internal Server Error`: Lỗi server

---

### 3. Danh sách URLs

**Endpoint:** `GET /api/url`

**Query Parameters:**

- `limit` (optional): Số items per page, default = 10, max = 100
- `page` (optional): Page number, bắt đầu từ 0, default = 0

**Example:** `GET /api/url?limit=20&page=0`

**Response:** `200 OK`

```json
{
  "content": [
    {
      "id": 1,
      "short_code": "abc123",
      "original_url": "https://example.com",
      "click_count": 42,
      "tiny_url": "http://localhost:8080/abc123",
      "created_at": "2024-12-28T10:00:00Z"
    }
  ],
  "current_page": 0,
  "total_count": 100,
  "is_first": true,
  "is_last": false,
  "is_next": true,
  "is_previous": false
}
```

---

### 4. Analytics chi tiết

**Endpoint:** `GET /api/url/:url_id/stats`

**Query Parameters:**

- `limit` (optional): Số items per page, default = 10
- `page` (optional): Page number, default = 0

**Example:** `GET /api/url/5/stats?limit=100&page=10`

**Response:** `200 OK`

```json
{
  "content": [
    {
      "id": 1,
      "clicked_at": "2024-12-28T10:15:30Z",
      "ip_address": "192.168.1.100",
      "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)...",
      "device_type": "desktop",
      "country": "VN",
      "referer": "https://google.com"
    }
  ],
  "current_page": 0,
  "total_count": 42,
  "is_first": true,
  "is_last": false
}
```

---

### 5. Click Count

**Endpoint:** `GET /api/url/:url_id/stats/count`

**Example:** `GET /api/url/4/stats/count`

**Response:** `200 OK`

```json
{
  "click_count": 42
}
```

---

### 6. Metrics

**Endpoint:** `GET /api/metrics`

**Response:** `200 OK`

```json
{
  "total_urls": 42,
  "total_clicks": 1000,
  "urls_created_today": 10,
  "clicks_today": 5,
  "top_urls:": [
    {
      "short_code": "abc",
      "original_url": "https://www.google.com/webhp?hl=vi",
      "clicks": 4,
      "tiny_url": "http://localhost:8000/abc"
    }
  ]
}
```

---

### 7. Health Check

**Endpoint:** `GET /health`

**Response:** `200 OK`

```json
{
  "status": "healthy",
  "time": 1703764800
}
```

---

## Kiến trúc & Thiết kế

### Tech Stack

| Component           | Technology     | Lý do chọn                                  |
| ------------------- | -------------- | ------------------------------------------- |
| **Language**        | Go 1.21        | Performance cao, concurrency tốt, type-safe |
| **Framework**       | Gin            | Fast, minimal, production-ready             |
| **Database**        | PostgreSQL 16  | ACID, reliability, powerful indexing        |
| **Query Builder**   | sqlc           | Type-safe queries, compile-time checking    |
| **Database Driver** | pgx/v5         | Fastest PostgreSQL driver cho Go            |
| **Migrations**      | golang-migrate | Standard migration tool                     |

### Database Schema

![Database Schema](https://res.cloudinary.com/dl8h3byxa/image/upload/v1766939672/url_shortener_ac07bv.png)

**Giải thích thiết kế:**

1. **`urls.short_code`**: UNIQUE constraint → đảm bảo không duplicate
2. **`urls.click_count`**: Denormalized counter → fast read, không cần JOIN
3. **Separate `clicks` table**: Chi tiết analytics mà không làm chậm table chính
4. **Indexes**: Optimize cho queries hay dùng (lookup, list, analytics)

## Quyết định kỹ thuật

### 1. Tại sao chọn PostgreSQL thay vì NoSQL?

**PostgreSQL:**

- ACID transactions → đảm bảo data consistency
- UNIQUE constraints → prevent duplicates at database level
- Powerful indexing (B-tree) → fast lookups
- Relations → dễ mở rộng (users, teams, permissions)
- Mature ecosystem với Go

**MongoDB/NoSQL:**

- Overkill cho schema đơn giản này
- ACID transactions phức tạp hơn
- Không cần flexibility của document-based

→ **Quyết định:** PostgreSQL cho primary storage, có thể thêm Redis cache sau

---

### 2. Thuật toán generate short code: Random Base62

**Implementation:**

```go
const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateShortCode(length int) (string, error) {
    result := make([]byte, length)
    for i := 0; i < length; i++ {
        num, err := rand.Int(rand.Reader, big.NewInt(62))
        if err != nil {
            return "", err
        }
        result[i] = base62Chars[num.Int64()]
    }
    return string(result), nil
}
```

**Các options đã xem xét:**

| Algorithm               | Collision? | Security          | Shortest? | Complexity |
| ----------------------- | ---------- | ----------------- | --------- | ---------- |
| Random Base62           | 0.0003%    | High              | No        | Simple     |
| Auto-increment + Base62 | Never      | Low (predictable) | Yes       | Medium     |
| Hash (MD5)              | 0.2%       | Medium            | No        | Medium     |
| Snowflake ID            | Never      | Medium            | No        | Complex    |

**Tại sao chọn Random Base62?**

- **Simple**: Dễ implement, ít bug
- **Secure**: Dùng `crypto/rand` → unpredictable
- **Collision cực thấp**: 62^7 = 3.5 trillion combinations
  - Với 1M URLs: collision probability = 0.00003%
  - Retry 3-5 lần là đủ
- **Distributed-friendly**: Không cần coordination giữa servers
- **Cons**: Không guarantee shortest code (trade-off chấp nhận được)

**Math:**

```
Length 7: 62^7 = 3,521,614,606,208 combinations
Với 1 million URLs:
- Probability of collision ≈ 1M / 3.5T ≈ 0.00003%
- Cần ~60 million URLs mới có 50% chance collision
```

---

### 3. Xử lý Concurrency & Duplicates

#### Vấn đề 1: Duplicate short codes

**Scenario:**

```
T1: Request A generates "abc123"
T2: Request B generates "abc123" (random trùng!)
T3: Both try INSERT → CONFLICT
```

**Giải pháp:**

```go
const maxRetries = 5

for attempt := 0; attempt < maxRetries; attempt++ {
    shortCode := GenerateShortCode(7)

    _, err := store.CreateURL(ctx, db.CreateURLParams{
        ShortCode:   shortCode,
        OriginalUrl: longUrl,
    })

    if err == nil {
        break // Success!
    }

    // Check if duplicate key error (PostgreSQL code 23505)
    if isDuplicateKeyError(err) {
        continue // Retry với code mới
    }

    return err // Real error
}
```

**Tại sao approach này tốt?**

- Database UNIQUE constraint đảm bảo atomicity
- Collision rate cực thấp (0.0003%) → ít khi retry
- Đơn giản, dễ debug

#### Vấn đề 2: Click counter race condition

**Scenario:**

```
T1: Read click_count = 100
T2: Read click_count = 100
T3: Write click_count = 101
T4: Write click_count = 101 ❌ (should be 102)
```

**Giải pháp: Atomic SQL increment**

```sql
UPDATE urls
SET click_count = click_count + 1
WHERE short_code = $1;
```

- Database đảm bảo atomicity
- Không cần application-level locking
- Async execution (không block redirect)

#### Vấn đề 3: Duplicate original URLs

**Quyết định: CHO PHÉP duplicate original URLs**

**Lý do:**

1. Users có thể muốn nhiều short links cho cùng URL (campaigns, A/B testing)
2. Đơn giản hóa implementation
3. Performance: Không cần query TEXT column trước khi insert

**Trade-off:**

- Pros: Simple, flexible, fast
- Cons: Database lớn hơn
- Future: Có thể add optional deduplication với query parameter

---

### 4. Validation & Security

#### URL Validation

```go
func isValidURL(raw string) bool {
    // Parse URL
    u, err := url.Parse(raw)
    if err != nil {
        return false
    }

    // Must have http/https scheme
    if u.Scheme != "http" && u.Scheme != "https" {
        return false
    }

    // Must have host
    if u.Host == "" {
        return false
    }

    // Block localhost/private IPs (security)
    if isPrivateHost(u.Host) {
        return false
    }

    return true
}
```

**Các edge cases được handle:**

- Empty URL
- Invalid scheme (ftp://, javascript:, data:)
- Localhost/private IPs (127.0.0.1, 192.168.x.x)
- URLs > 2048 characters
- Malformed URLs

#### SQL Injection Protection

- Dùng **sqlc** → parameterized queries
- Never concatenate user input vào SQL
- pgx/v5 automatically escapes parameters

#### Open Redirect Protection

**Vấn đề:** Attacker có thể tạo short link đến phishing site

**Mitigations hiện tại:**

- Block localhost/private IPs
- Validate URL scheme (chỉ http/https)
- URL validation strict

**Future enhancements:**

- Google Safe Browsing API integration
- User-reported spam system

---

### 5. Tại sao dùng sqlc + pgx/v5 thay vì GORM?

| Feature                 | sqlc + pgx/v5  | GORM                 |
| ----------------------- | -------------- | -------------------- |
| **Type safety**         | Compile-time   | Runtime              |
| **Performance**         | ~30% faster    | Slower (reflection)  |
| **SQL control**         | Raw SQL        | ORM magic            |
| **Learning curve**      | Medium         | Easy                 |
| **PostgreSQL features** | Full support   | Limited              |
| **Debugging**           | Easy (see SQL) | Hard (generated SQL) |

**Tại sao chọn sqlc + pgx/v5?**

- **Type-safe**: Errors at compile-time, không phải runtime
- **Performance**: pgx là fastest PostgreSQL driver
- **Full control**: Viết raw SQL, optimize queries dễ dàng

---

### 6. Analytics: Clicks table riêng

**Quyết định: Tách `clicks` table riêng thay vì chỉ có counter**

**Schema:**

```sql
CREATE TABLE clicks (
    id BIGSERIAL PRIMARY KEY,
    url_id BIGINT REFERENCES urls(id),
    clicked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address TEXT,
    user_agent TEXT,
    device_type TEXT,
    country TEXT,
    referer TEXT
);
```

**Tại sao?**

- Rich analytics: clicks over time, geography, devices
- Không làm chậm table `urls`
- Có thể partition by time sau này
- Show advanced database design skills

**Trade-offs:**

- Pros: Detailed insights, scalable
- Cons: More storage, complex queries
- Denormalized counter (`urls.click_count`) cho fast reads

**Async tracking:**

```go
go func() {
    bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Insert click details
    store.InsertClick(bgCtx, clickData)

    // Increment counter
    store.IncrementClickCount(bgCtx, shortCode)
}()
```

→ Không block redirect, user experience tốt

---

## Trade-offs

### 1. Random generation vs Auto-increment

**Chọn:** Random Base62

| Aspect        | Random Base62      | Auto-increment    |
| ------------- | ------------------ | ----------------- |
| Collision     | Possible (0.0003%) | Never             |
| Security      | Unpredictable      | Predictable (bad) |
| Shortest code | No                 | Yes               |
| Distributed   | Easy               | Need coordination |

**Lý do:** Security và simplicity quan trọng hơn shortest code

---

### 2. Allow duplicate URLs vs Deduplicate

**Chọn:** Allow duplicates

**Reasoning:**

- Đơn giản, không cần check trước khi insert
- Flexible: Users có thể tạo nhiều links cho cùng URL
- Performance: Không query TEXT column
- Trade-off: Database lớn hơn (storage rẻ, chấp nhận được)

---

### 3. Clicks table vs Chỉ counter

**Chọn:** Cả hai (denormalized counter + detailed clicks)

**Reasoning:**

- `click_count` trong `urls`: Fast reads, đơn giản
- `clicks` table: Detailed analytics
- Best of both worlds!

**Cost:**

- More storage
- Insert overhead (nhưng async → không ảnh hưởng UX)

---

### 4. Sync vs Async click tracking

**Chọn:** Async

```go
go func() {
    // Track click in background
}()
ctx.Redirect(302, originalURL) // Don't wait
```

**Reasoning:**

- Redirect phải nhanh (<50ms) → UX tốt
- Click tracking có thể chậm (100-200ms)
- Trade-off: Có thể mất vài clicks nếu server crash (acceptable)

---

## Thử thách & Giải pháp

### Challenge 1: Concurrency - Duplicate short codes

**Problem:**
2 requests cùng lúc generate "abc123" → cả hai INSERT → conflict!

**Solution:**

```go
// Database UNIQUE constraint
CREATE UNIQUE INDEX idx_short_code ON urls(short_code);

// Application retry logic
for retries := 0; retries < 5; retries++ {
    shortCode := GenerateShortCode(7)
    err := db.Insert(shortCode, url)
    if err == nil { break }
    if isDuplicateKeyError(err) { continue }
    return err
}
```

**Learned:**

- Database constraints > Application-level checks
- Let database handle atomicity
- Retry logic phải có max attempts

---

### Challenge 2: Click tracking làm chậm redirect

**Problem:**

```go
// Bad approach
store.InsertClick(ctx, clickData)           // Takes 100ms
store.IncrementClickCount(ctx, shortCode)   // Takes 50ms
ctx.Redirect(302, url)                       // User waits 150ms!
```

**Solution:**

```go
// Good approach
go func() {
    bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    store.InsertClick(bgCtx, clickData)
    store.IncrementClickCount(bgCtx, shortCode)
}()
ctx.Redirect(302, url) // Instant!
```

**Learned:**

- User experience > Perfect data
- Async for non-critical operations
- Context phải tạo mới trong goroutine (không dùng request context)

---

### Challenge 3: pgx/v5 với pgtype.Text, pgtype.Int8

**Problem:**

```go
// sqlc generates nullable types
type Click struct {
    IpAddress pgtype.Text // Not string!
    UrlID     pgtype.Int8 // Not int64!
}
```

**Solution:**

```go
// Construct nullable types properly
clickData := db.InsertClickParams{
    UrlID:     pgtype.Int8{Int64: urlRecord.ID, Valid: true},
    IpAddress: pgtype.Text{String: ip, Valid: true},
}

// Access values
if click.IpAddress.Valid {
    ip := click.IpAddress.String
}
```

**Learned:**

- PostgreSQL NULLs cần special handling
- pgx/v5 type system is type-safe nhưng verbose
- Trade-off: Verbosity for safety

# Hạn chế & Hướng cải thiện trong tương lai

### Hạn chế hiện tại

- Chưa sử dụng Redis cache cho các short URL được truy cập nhiều, nên mọi request redirect vẫn phải truy vấn database
- Đang sử dụng offset-based pagination, có thể gây chậm khi số lượng bản ghi rất lớn
- Rate limiting hiện tại chỉ ở mức in-memory, chưa hỗ trợ môi trường distributed
- Click tracking được xử lý bất đồng bộ (async), nên có khả năng mất một số lượt click nếu server bị crash

### Hướng cải tiến trong tương lai

- Tích hợp Redis để cache các short URL phổ biến, giúp giảm tải cho database
- Chuyển sang cursor-based pagination để tối ưu hiệu năng khi dữ liệu lớn
- Bổ sung tính năng URL expiration và custom alias theo nhu cầu người dùng (hiện chưa implement)
- Sử dụng background worker hoặc message queue để xử lý click tracking một cách ổn định hơn

### Hướng tới production-ready

- Triển khai service theo mô hình stateless và scale ngang phía sau load balancer
- Bổ sung centralized logging và monitoring (ví dụ: Prometheus, Grafana)
- Áp dụng rate limiting phân tán bằng Redis
- Thiết lập CI/CD pipeline với automated tests để đảm bảo chất lượng code
