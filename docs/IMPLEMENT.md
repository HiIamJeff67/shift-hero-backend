# Shift Hero Backend 實作規範

本文件定義此專案新增或修改後端功能時必須遵守的架構、依賴方向、資料流、路由組裝、安全性、錯誤處理與測試規範。

本文件以目前程式碼結構為準，不是通用 Go/Gin 範本。既有程式碼若與本文不一致，視為待逐步收斂的歷史實作，不應直接複製成新功能。

## 1. 規範用語

- **MUST / 必須**：除非有明確且記錄於 PR 的架構理由，否則不得違反。
- **SHOULD / 應**：原則上遵守；偏離時必須能說明取捨。
- **MAY / 可以**：依功能複雜度選用。
- 每次變更應保持最小範圍，不得順便重構無關模組。
- 新功能必須優先沿用現有 package、constructor、interface、DTO、Exception 與 response envelope。

## 2. 標準請求與依賴方向

功能由底層基礎往上實作，標準建置鏈為：

```text
database
  -> schemas / enums / constraints / triggers
  -> repository inputs / scopes / embedded SQL
  -> repository（可選，但複雜或可重用資料存取時必須）
  -> service
  -> controller
  -> binder
  -> module
  -> routes + middlewares / interceptors
```

不得因省略某個可選層，就把它原本負責的 persistence、商業規則或 HTTP 轉換混入其他層。

### 2.1 HTTP 請求流程

```text
Gin Engine / RouterGroup
  -> global middlewares
  -> route observability middlewares
  -> rate limit / timeout / auth / authorization / CSRF
  -> response interceptors
  -> binder
  -> controller
  -> service
  -> repository（需要資料存取抽象時）
  -> scopes / embedded SQL / repository input
  -> schema / enum
  -> database
```

依賴組裝不在請求流程內，由 module 統一負責：

```text
route -> module
module -> binder + controller + service + repositories + external services
```

### 2.2 強制依賴方向

正常情況下只允許上層依賴下層：

```text
routes
  -> modules
  -> binders / controllers
  -> services
  -> repositories
  -> inputs / scopes / sqls
  -> schemas / enums
  -> database
```

共用支援套件可由需要的層使用：

```text
exceptions, validation, contexts, options, monitor, util,
caches, tokens, cookies, emails, storages, shared/*
```

禁止反向依賴：

- repository 不得 import service、controller、binder、route。
- service 不得 import controller、binder、route，也不得依賴 `*gin.Context`。
- controller 不得直接 import repository、GORM schema 或執行 SQL。
- binder 不得呼叫 service 或 repository，只能呼叫傳入的 controller function。
- route 不得直接建立 repository query、執行商業邏輯或自行組 response。
- schema、enum、input、scope、SQL package 不得依賴 HTTP 層。
- 禁止用全域變數繞過 module constructor 注入依賴；`models.DB` 等既有基礎設施入口除外。

## 3. 目錄與責任

| 位置 | 責任 |
| --- | --- |
| `app/models/database.go` | DB 連線、migration、seed 與 DB 基礎操作 |
| `app/models/schemas/` | GORM table schema、relations、hooks、table registration |
| `app/models/schemas/enums/` | PostgreSQL/Go enum、合法值集合、轉換 |
| `app/models/schemas/constraints/` | DB constraint/index SQL 與 migration registration |
| `app/models/schemas/triggers/` | DB trigger SQL 與 migration registration |
| `app/models/inputs/` | repository 寫入模型與 partial update input |
| `app/models/scopes/` | 可組合、無商業副作用的 GORM query scope |
| `app/models/sqls/<domain>/` | 複雜或需要精確控制的 embedded SQL |
| `app/models/repositories/` | 資料存取介面與實作、DB error 轉 Exception |
| `app/dtos/` | API request/response contract |
| `app/validation/` | validator instance 與自訂 validation tags |
| `app/services/` | use case、商業規則、授權、transaction、DTO/schema 轉換 |
| `app/controllers/` | 呼叫 service，轉成固定 HTTP response |
| `app/binders/` | 讀取 header/context/body/query/path 並建立 request DTO |
| `app/modules/` | constructor wiring 與 dependency injection |
| `app/middlewares/` | request 前後的橫切行為與存取控制 |
| `app/interceptors/` | 在 controller response 寫出前修改 response |
| `app/routes/developmentroutes/` | 正式開發 API route registration |
| `app/routes/testroutes/` | E2E 專用 route registration |
| `app/exceptions/` | domain exception code、HTTP status、message |
| `app/contexts/` | Gin context field 的型別安全讀取與轉換 |
| `app/caches/` | Redis/cache 存取與 cache DTO |
| `app/graphql/`、`shared/graphql/` | GraphQL handler、resolver、schema、query |
| `app/monitor/` | log、metric、trace |
| `shared/constants/` | 跨模組固定常數 |
| `shared/types/` | 跨模組共用型別與 generic pattern |
| `shared/lib/` | 可獨立重用的基礎 library |
| `test/unit/` | 無外部服務或可隔離依賴的單元測試 |
| `test/e2e/` | 經 route 進入並驗證完整 API 行為 |

## 4. 新功能的實作順序

新增一個 domain/module 時，原則上依照以下順序實作：

1. 確認 API contract、權限、資料一致性與 transaction 邊界。
2. 新增 enum、schema、constraint/trigger、migration registration。
3. 新增 repository input、scope 或 embedded SQL。
4. 新增 repository interface 與 implementation。
5. 新增 request/response DTO 與 validation tag。
6. 新增 domain Exception。
7. 新增 service interface、implementation 與 transaction。
8. 新增 controller。
9. 新增 binder。
10. 新增 module wiring。
11. 新增 route 與正確 middleware/interceptor 順序。
12. 更新 OpenAPI/contract、範例、測試與相關文件。
13. 執行 format、test、build。

不是每個功能都需要新增每一層。省略某層時仍不得把該層責任塞到不相干的層。

## 5. Database、Schema 與 Migration

### 5.1 Schema

- 一張主要資料表應有一個 `<entity>_schema.go`。
- schema struct 只描述 persistence model、relation、GORM hook 與必要轉換，不放 use case。
- 欄位必須明確定義 `column`、型別、nullability、default、index/unique 或時間行為。
- ID 優先沿用 `uuid.UUID` 與資料庫 `gen_random_uuid()`。
- API JSON 命名沿用 camelCase；DB column 使用 snake_case。
- 必須實作 `TableName()`，並使用 `shared/types` 內的 table name constant。
- relation 名稱必須定義成該 schema 的 relation type constant，repository preload 不得散落 magic string。
- 關聯的 `OnUpdate`、`OnDelete` 行為必須根據 domain 明確指定。
- hook 必須是局部且可預測的 persistence 行為；跨 repository、寄信、cache、外部 API 等副作用不得放在 hook。
- 新 schema 必須加入 `app/models/schemas/migrate.go` 的 `MigratingTables`。

### 5.2 Enum

新增 enum 時必須同時完成：

1. 在 `app/models/schemas/enums/` 建立 enum type 與 constants。
2. 提供 `AllXxx` 與 `AllXxxStrings`。
3. 提供 `Name()`、`Scan()`、`Value()`、`String()`、`IsValidEnum()`；需要時提供 string converter。
4. 加入 `app/models/schemas/enums/migrate.go` 的 `MigratingEnums`。
5. 在 `app/validation/enums_validation.go` 註冊並使用對應 validation tag。
6. 補上 migration 相容性、DTO validation 與測試。

既有 enum value 不得任意 rename 或刪除。PostgreSQL enum 變更屬於資料 migration，必須評估既有資料與 rollback。

### 5.3 Constraint、Index 與 Trigger

- 能由 DB 保證的唯一性、foreign key、check、partial unique index，不能只靠 service 先查再寫。
- SQL 檔放在對應 domain 子目錄，使用 `//go:embed` 載入。
- constraint 必須加入 `MigratingConstraintSQLs`。
- trigger 必須加入 trigger migration registry。
- SQL object 命名要穩定且具 domain 意義，避免由環境或 runtime 動態產生。
- 新增 migration 時需驗證空 DB 與已有資料 DB 的行為。

目前啟動 migration 順序是 enum、特定 relation migration、table、trigger、constraint；新 migration 必須尊重被依賴物件先建立的順序。

### 5.4 Seed

- seed 只放必要預設資料或明確標示的 example data。
- seed SQL 放在 `app/models/seeds/<domain>_seeds/`，使用 embedded SQL。
- 必須加入 `SeedingDefaultDataSQLs`，且盡可能可重複執行或能安全偵測重複資料。
- 不得把 production secret、個資或環境限定 ID 寫入 seed。

## 6. Input、Scope、SQL 與 Repository

### 6.1 Repository Input

- `app/models/inputs/` 是 service 到 repository 的 persistence input，不是 HTTP DTO。
- create/update input 只包含該資料操作允許寫入的欄位。
- update 欄位用 pointer 表示「有無提供」，partial update 沿用 `PartialUpdateInput[T]`。
- 不得直接將外部 request body 當成 GORM update map；必須先經 DTO validation 與允許欄位映射。

### 6.2 Scope

- 可重用的 filter、preload、排序、soft-delete 條件應放 `app/models/scopes/`。
- scope 應為純 query composition，形式優先採用 `func(*gorm.DB) *gorm.DB`。
- scope 不得執行 `Create`、`Update`、`Delete`、`Commit` 或外部副作用。
- 權限條件若能安全地併入 query，可做成命名清楚的 scope；最終授權決策仍由 service 負責。

### 6.3 Embedded SQL

- GORM 無法清楚表達、效能敏感、需要 CTE/locking/atomic update 時可以使用 SQL。
- SQL 放在 `app/models/sqls/<domain>/*.sql`，並由同 package 的 `sql.go` 使用 `//go:embed`。
- 所有使用者輸入必須透過 bind parameters 傳入，不得字串拼接。
- table/column 名若必須動態選擇，只能來自程式內 whitelist。
- 跨多個 use case 重用的 SQL 應由 repository 封裝；不可散落在 controller/binder/route。

### 6.4 Repository

repository 是可選層，但符合以下任一條件時必須建立：

- 同一 entity 有多個資料操作。
- query 會被多個 service/use case 使用。
- 需要 mock、transaction DB 注入、locking、preload、scope 或一致的 DB error mapping。
- 資料查詢已足以形成清楚的 domain vocabulary。

只有非常局部、一次性且不會重用的簡單資料操作，service 才可以直接使用 GORM。新模組應優先採 repository；`scheduling_service.go` 等既有直接 GORM 寫法不應視為預設範本。

repository 必須：

- 同時提供 `<Entity>RepositoryInterface` 與 concrete implementation。
- 使用 `New<Entity>Repository()` constructor。
- 每個 method 接受 `opts ...options.RepositoryOptions`，並以 `options.ParseRepositoryOptions` 取得 DB。
- 預設使用 `models.DB`，service transaction 則以 `options.WithDB(tx)` 或 `WithTransactionDB(tx)` 傳入。
- 回傳 schema、repository projection type、primitive result 與 `*exceptions.Exception`。
- 將 GORM/driver error 轉成 domain Exception，不把裸 `error` 傳到 controller。
- 精確區分 not found、duplicate、no changes 與真正的 DB failure。
- 對需要一致性保護的讀取明確使用 transaction 與適當 locking。
- 不得開始或 commit 跨 repository transaction。
- 不得寄信、寫 HTTP response、操作 cookie/token 或執行 use case。

repository method 應使用 domain 行為命名，例如 `GetOneById`、`GetManyByUserId`、`UpdateReviewState`，避免模糊的 `HandleData`、`Process`。

## 7. DTO 與 Validation

### 7.1 Request DTO

request DTO 必須沿用：

```go
type Request[H, C, B, P any] struct {
    Header        H
    ContextFields C
    Body          B
    Param         P
}
```

- `Header`：只放 endpoint 真正需要的 header。
- `ContextFields`：只放由可信 middleware 建立的欄位，例如 authenticated user ID。
- `Body`：JSON body 或 query binding 的資料。
- `Param`：path parameter。
- 不使用的區塊填 `any`。
- request 與 response DTO 放在同一 domain 的 `<domain>_dto.go`。
- API request/response 不得直接暴露 GORM schema，尤其不得暴露 password、refresh token、內部 ID 或 relation。

### 7.2 Validation

- binder 負責「能否解析」，service 負責 `validation.Validator.Struct(reqDto)`。
- service method 必須先 validation，再執行資料查詢、cache、寄信或其他副作用。
- validation tag 必須與欄位 optional/nullable 語意一致。
- optional pointer 優先使用 `omitnil`；非 pointer optional value 使用 `omitempty`。
- enum 使用已註冊的 `is...` validation tag。
- 跨欄位規則與商業規則放 service，不應塞入 binder。
- 新增自訂 validation tag 時，要在 `app/validation/` 對應檔案註冊並補測試。

### 7.3 Partial Update

- partial update 必須明確區分「未提供」、「更新值」與「設為 null」。
- 沿用 `PartialUpdateDto[T]` / `PartialUpdateInput[T]` 與既有 preprocess helper。
- 不得以 zero value 推測使用者是否提供欄位。
- 沒有任何有效變更時回傳 domain `NoChanges` Exception。

## 8. Service

service 是商業邏輯與 transaction 的主要邊界。

每個 service：

- 定義 `<Domain>ServiceInterface`。
- concrete struct 只持有 DB、repository 與必要外部 service dependency。
- 使用 `New<Domain>Service(...)` constructor 注入依賴。
- method 第一個參數使用 `context.Context`，不得接收 `*gin.Context`。
- DB 操作從 `s.db.WithContext(ctx)` 開始，確保 timeout/cancellation 能往下傳。
- 回傳 response DTO 或 domain result 與 `*exceptions.Exception`。
- validation 必須在 method 開頭完成。
- 授權必須在 mutation/query 前完成，不能只依賴前端或 route 隱藏按鈕。
- 負責 DTO、input、schema、response DTO 之間的顯式 mapping。
- 不得直接寫 HTTP response、header、cookie。

建議 method 內順序：

```text
validate request DTO
-> derive normalized values/defaults
-> attach request context to DB
-> authenticate/authorize domain resource
-> begin transaction（需要時）
-> repository/query operations
-> external/cache side effects（依一致性策略安排）
-> commit
-> map response DTO
```

### 8.1 Transaction

- 多個寫入必須原子完成時，由 service 開 transaction。
- transaction 內所有 repository 呼叫必須傳同一個 `tx`。
- 任一錯誤都必須 rollback；commit error 必須轉成 domain Exception。
- 必須處理 panic rollback，且不得吞掉 panic 後假裝成功。
- transaction 內避免執行慢速外部 API、寄信或不受 DB rollback 控制的副作用。
- 若外部副作用必須與 DB 一致，應採 outbox/job/idempotency 設計，而不是假設 rollback 能撤回外部操作。
- commit 後才可用一般 DB 重新讀取最終結果。

### 8.2 授權

- middleware 的 role/plan 是平台級粗粒度權限。
- company membership、manager、resource owner 等資源級權限必須在 service 再驗證。
- 共用授權邏輯可放同 domain 的 guard/helper，並優先透過 repository 存取資料。
- 不得接受 client 傳入的 user ID 作為目前登入者；登入者 ID 必須來自 AuthMiddleware 寫入的 context。

## 9. Controller

controller 只負責 service 與 HTTP response 的轉接：

```go
resDto, exception := c.service.Operation(ctx.Request.Context(), reqDto)
if exception != nil {
    exception.Log().SafelyAbortAndResponseWithJSON(ctx)
    return
}

ctx.JSON(http.StatusOK, gin.H{
    "success":   true,
    "data":      resDto,
    "exception": nil,
})
```

規則：

- 定義 `<Domain>ControllerInterface`。
- constructor 只注入 service interface。
- 不做 DTO validation、授權、transaction、GORM query 或商業計算。
- Exception 一律走 `Log().SafelyAbortAndResponseWithJSON(ctx)`。
- 成功 response 必須維持固定 envelope：`success`、`data`、`exception`。
- 目前 API 主要使用 `http.StatusOK`；若要引入 `201`、`204` 等狀態，必須同步更新 contract、client 與測試，不得單點改變。

## 10. Binder

binder 是 Gin 與 typed request DTO 的邊界。

- 定義 `<Domain>BinderInterface` 與 `New<Domain>Binder()`。
- 每個 bind method 接收 `types.ControllerFunc[*dtos.XxxReqDto]` 並回傳 `gin.HandlerFunc`。
- 依來源使用正確 API：
  - JSON：`ctx.ShouldBindJSON(&reqDto.Body)`
  - query：`ctx.ShouldBindQuery(&reqDto.Body)`
  - path：`ctx.Param(...)` 後做明確型別轉換
  - context：使用 `app/contexts` converter
- authenticated user/context field 必須先取出並放進 `ContextFields`。
- path UUID 解析失敗要回 domain bad request/invalid DTO Exception。
- bind/parse 失敗後必須 abort response 並立刻 `return`。
- binder 不做 repository query、權限判斷、transaction 或 response DTO mapping。
- 同一 domain 重複的 path parsing 可以抽成未 export helper。

## 11. Module 與 Dependency Injection

每個公開 API domain 應有 `app/modules/<domain>_module.go`：

```text
construct repositories/external services
-> construct service
-> construct binder
-> construct controller
-> return module{Binder, Controller}
```

規則：

- module 是 object graph 的唯一組裝位置。
- route 只建立 module 並使用其 `Binder`、`Controller`。
- controller 依賴 service interface；service 依賴 repository/external service interface。
- constructor 不得執行資料查詢、migration 或 request-specific side effect。
- module 對外只暴露 route registration 需要的元件，通常是 `Binder` 與 `Controller`。

## 12. Routes、Middleware 與 Interceptor

### 12.1 Global RouterGroup

`DevelopmentRouterGroup` 目前固定套用：

```text
SanitizeXForwardedForMiddleware
-> CORSMiddleware
-> DomainWhiteListMiddleware
```

- health route 在 API group 外，保持輕量且不可依賴 DB/Redis 才能回應。
- 新 route 必須註冊於對應 `<domain>_route.go`。
- 新 domain 的 configure function 必須加入 `ConfigureDevelopmentRoutes()`。
- E2E 專用 route 放 `testroutes`，不得意外註冊到 production/development router。

### 12.2 Endpoint 標準順序

公開且需要登入的 REST endpoint 原則上依序為：

```text
1. ApplyTracerMiddleware
2. ApplyMeterMiddleware
3. UnauthorizedRateLimitMiddleware（在辨識使用者前限制 IP/fingerprint）
4. TimeoutMiddleware
5. AuthMiddleware
6. AuthorizedRateLimitMiddleware（若該 endpoint 啟用使用者級限制）
7. UserRoleMiddleware / UserPlanMiddleware（若需要）
8. CSRFMiddleware（需要 CSRF 保護的 state-changing request）
9. ShareableResponseWriterInterceptor(...)
10. Binder -> Controller
```

不需登入的 endpoint：

```text
tracer -> meter -> unauthorized rate limit -> timeout -> binder -> controller
```

注意：

- `AuthorizedRateLimitMiddleware` 必須在 `AuthMiddleware` 後，因為它需要 user ID。
- `CSRFMiddleware` 必須在 `AuthMiddleware` 後。
- role/plan middleware 必須在 `AuthMiddleware` 後。
- `TimeoutMiddleware` 應包住後續 handler 與 response interceptor。
- `ShareableResponseWriterInterceptor` 必須在 binder/controller 前。
- `RefreshTokenInterceptor` 需要 AuthMiddleware 建立 token context。
- `EmbeddedInterceptor` 需要 authenticated public ID。
- 是否同時套用 unauthorized 與 authorized limiter 由 endpoint 風險決定，但順序不得顛倒。
- `RepositionMiddleware(fronts, backs, handler)` 只負責依序串接，不會自動修正錯誤順序。

### 12.3 CSRF

- 以 cookie 驗證且會改變狀態的 `POST`、`PUT`、`PATCH`、`DELETE` 應套用 CSRF。
- 純 Bearer API 若明確不依賴 cookie，可依 threat model 例外，但必須在 contract/PR 說明。
- 不得只因目前某條舊 route 未套 CSRF，就複製該缺口到新 route。

### 12.4 Timeout 與 Gin Context

- timeout 必須依 endpoint 實際工作量設定，不可無理由放大。
- service/repository/external call 必須使用 request `context.Context`。
- 不得在 handler 結束後持有或異步存取 `*gin.Context`。
- 啟動 goroutine 時只能傳必要 immutable data 與標準 context，並明確處理 cancellation、panic 與資源釋放。
- controller/service 不得繞過 shareable response writer 直接寫原始 socket。

### 12.5 Interceptor

- 需要讀寫 response body 的 interceptor 必須透過 `ShareableResponseWriterInterceptor` 註冊。
- interceptor factory signature 維持 `func(responseWriterKey string) gin.HandlerFunc`。
- interceptor 應先 `ctx.Next()` 讓後續 handler 產生 response，再安全修改 buffer。
- timeout、已寫出 response、HTTP error 狀態時不得再次改寫成功 body。
- 修改 response body 後必須維持合法 JSON 與固定 envelope。

## 13. Exception 與錯誤處理

- 每個 domain 使用 `app/exceptions/<domain>_exception.go`。
- 新 domain 必須取得未衝突的 subdomain code 與 prefix。
- Exception code 必須穩定；client 已依賴的 code/reason 不得任意重用或改義。
- 預期錯誤使用具體 reason，例如 `NotFound`、`Forbidden`、`Duplicate...`、`NoChanges`。
- 未預期 DB/外部錯誤要保留 trace/log 所需 origin，但 response 不得洩漏 secret、token、SQL、DSN、個資或內部 stack。
- 不得在各層直接 `ctx.JSON` 自製錯誤格式。
- 不得只 log 後繼續執行可能產生錯誤成功 response。
- 下層 Exception 往上傳時，可用 `WithOrigin`、`WithDetails` 增加診斷資訊，但不得破壞既有 code/reason。
- `panic` 只用於真正不可恢復狀態；一般 validation、not found、conflict 一律回 Exception。

固定失敗 envelope：

```json
{
  "success": false,
  "data": null,
  "exception": {
    "code": 0,
    "reason": "Reason",
    "prefix": "Domain",
    "message": "Safe message",
    "status": 400,
    "details": null,
    "error": null
  }
}
```

## 14. Context、Auth、Token 與 Cookie

- context key 必須定義於 `shared/types/context_field_name.go`，不得散落 string literal。
- 寫入與讀取 context 的型別必須一致。
- binder 讀 context 必須使用 `app/contexts` helper，不自行 type assertion。
- AuthMiddleware 是 user identity、role、plan、token refresh state 的唯一可信 HTTP 建立點。
- service 不信任 request body/header 內聲稱的 role、plan、user ID。
- access/refresh/CSRF token 必須沿用 `app/tokens`、`app/cookies`、`app/caches`。
- 不得把 token、password、auth code、cookie value 寫入一般 log 或 Exception details。
- token refresh response 必須透過 `RefreshTokenInterceptor`，不得每個 controller 各自實作。

## 15. Cache、Email、Storage 與外部服務

- cache key、TTL、serialization 應集中在 `app/caches`，service 只表達 use case。
- cache miss 與 cache infrastructure failure 必須區分；不得把 Redis 暫時失敗一律當成資料不存在。
- DB mutation 後需同步考慮 cache invalidation/update。
- 外部服務必須由 interface + constructor 注入 service，方便測試與替換。
- 外部 call 必須使用 timeout/cancellation，並轉成 domain Exception。
- email/storage 不得由 controller、binder、repository 直接呼叫。
- 非關鍵通知優先在 DB commit 後執行；需要可靠投遞時使用 durable job/outbox。

## 16. Observability

- 每條正式 endpoint 必須有 tracer 與 request meter。
- span name 使用穩定的 lowerCamelCase operation name。
- metric name 優先加入 `app/monitor/metrics` 的集中 registry，避免 route 內散落字串。
- service、repository 與 middleware 記錄錯誤時應保留 trace context。
- log 必須包含可診斷欄位，但不得包含 secret、完整 token、password、auth code 或不必要個資。
- 不使用臨時 `fmt.Println` 作為正式 observability；沿用 `app/monitor/logs`。
- 高基數資料（user ID、email、URL query）不得直接作為 metric name/label。

## 17. API 與命名

- Go package 名維持小寫且語意單一。
- 檔名使用 snake_case：`company_join_request_repository.go`。
- exported type/function 使用 PascalCase；local variable/method 使用 camelCase。
- interface、constructor、實作命名沿用：
  - `XxxRepositoryInterface` / `NewXxxRepository`
  - `XxxServiceInterface` / `NewXxxService`
  - `XxxControllerInterface` / `NewXxxController`
  - `XxxBinderInterface` / `NewXxxBinder`
  - `XxxModule` / `NewXxxModule`
- DTO 使用 `ReqDto`、`ResDto`；repository persistence payload 使用 `Input`。
- JSON 欄位使用 camelCase；DB column 使用 snake_case。
- route path 需與既有 client contract 一致。若改善舊 route 命名，應新增版本或相容 alias，不可直接破壞既有 route。
- 不得在新名稱延續已知拼字錯誤。

## 18. Security

- 所有外部輸入都視為不可信，必須經 bind、type conversion、validation、authorization。
- SQL 一律 parameterized。
- request body 大小依 endpoint 風險使用 `MaxContextSizeMiddleware`。
- upload 必須驗證 MIME、大小、檔名與儲存路徑，不能只信副檔名。
- CORS 與 domain whitelist 必須由環境設定，不得為了方便直接允許任意 production origin。
- proxy header 只在 trusted proxy 設定正確時可信。
- mutation 必須防止 IDOR：即使 ID 合法，也要驗證登入者對該 resource 的權限。
- 密碼只儲存強 hash，不得回傳或記錄。
- duplicate/error mapping 不應依賴會暴露敏感 SQL 的 response。
- 機密只能由 env/secret manager 取得，不得 commit 到 source、testdata 或 docs。

## 19. GraphQL

- GraphQL schema 放 `shared/graphql/schemas`，fragment/query 放對應目錄。
- generated code 不手動修改；依 `infra/graphql/gqlgen.yaml` 與 Makefile 重新產生。
- resolver 應保持薄層，商業邏輯仍放 service。
- resolver 不得直接複製 REST controller 邏輯或繞過 service 授權。
- GraphQL error 必須轉為既有 Exception/GQL error 格式。
- schema 變更需同步 query/fragment、generated code 與測試。

## 20. 測試規範

### 20.1 Unit Test

- pure helper、validation、guard、mapping 與 service business rule 應有 unit test。
- util 測試沿用 `test/unit/criteria.md` 與 JSON testdata pattern。
- table-driven test 必須涵蓋成功、邊界、非法輸入與 error branch。
- service 測試應以 interface fake/mock 隔離 repository 與外部服務。
- transaction、locking、constraint 等 DB 行為不能只靠 mock，需補 integration/E2E。

### 20.2 E2E

- route contract、middleware、binder、controller、service、DB 整合行為放 `test/e2e/<domain>`。
- E2E 必須驗證 HTTP status、固定 envelope、data、exception、header/cookie 與必要 DB side effect。
- authentication endpoint 要驗證 access/refresh/CSRF token 的 cookie/header 行為。
- 新增 route 時至少測：
  - 正常成功
  - malformed JSON/query/path
  - DTO validation failure
  - unauthenticated
  - unauthorized/forbidden
  - not found
  - conflict/duplicate
  - transaction rollback 或關鍵 failure path
- 平行測試不得共享會互相覆寫的固定資料；testdata 必須可隔離。
- 測試結束要清理建立的資料，不得依賴測試執行順序。

### 20.3 驗證命令

提交前至少執行：

```bash
gofmt -w <changed-go-files>
go test ./...
go build ./...
```

若完整 E2E 需要 Docker/PostgreSQL/Redis，必須另外執行對應 Make target，並在 PR 記錄未執行項目與原因。

## 21. 禁止事項

以下做法不得出現在新功能：

- route/controller 直接操作 GORM 或 SQL。
- repository 依賴 service 或 Gin。
- service 接收 `*gin.Context`。
- binder 內做商業授權或 DB query。
- controller 內 validation、transaction 或 DTO/schema mapping。
- 直接回傳 schema 作為 API response。
- 使用 client body 的 user ID 取代 authenticated context user ID。
- transaction 中有 repository 偷偷 commit。
- SQL 字串拼接使用者輸入。
- 忽略 `RowsAffected`、commit error、bind error 或 context cancellation。
- 發生 Exception 後沒有立即 return。
- 隨意更動 Exception code/reason 或 response envelope。
- 將 secret、token、password、DSN、SQL error、stack trace 回傳給 client。
- 為了通過測試而移除 auth、CSRF、validation、constraint 或 timeout。
- 複製舊模組的技術債，並以「現有程式也這樣寫」作為理由。

## 22. 新增 Module Checklist

### Contract

- [ ] HTTP method/path、request、response、status、Exception 已定義。
- [ ] authentication、role/plan、resource-level permission、CSRF 已定義。
- [ ] OpenAPI/contract/example 已更新。

### Data

- [ ] enum/schema/relations/table name 已完成。
- [ ] table 已加入 `MigratingTables`。
- [ ] enum 已加入 `MigratingEnums` 與 validator。
- [ ] constraint/index/trigger 已註冊。
- [ ] repository input/scope/SQL 已按需要建立。
- [ ] DB constraint 與 service validation 沒有互相矛盾。

### Application

- [ ] repository interface/implementation 已完成，或有合理理由省略。
- [ ] request/response DTO 與 validation 已完成。
- [ ] domain Exception code/reason 已完成且未衝突。
- [ ] service validation、authorization、transaction、mapping 已完成。
- [ ] controller 僅處理 service 與 response。
- [ ] binder 正確綁定 context/body/query/path。
- [ ] module 完成依賴注入。

### Route

- [ ] route 已註冊到對應 configure function。
- [ ] tracer、meter、rate limit、timeout 順序正確。
- [ ] AuthMiddleware 在所有需要 identity 的 middleware 前。
- [ ] state-changing cookie request 已套 CSRF。
- [ ] interceptor 由 ShareableResponseWriterInterceptor 管理。
- [ ] binder/controller 是 route 最後的 handler。

### Quality

- [ ] unit/integration/E2E 覆蓋主要成功與失敗路徑。
- [ ] logs/metrics/traces 不含敏感資料。
- [ ] cache invalidation、external side effect 與 transaction 邊界已確認。
- [ ] `gofmt`、`go test ./...`、`go build ./...` 已通過。
- [ ] 沒有修改無關檔案或提交環境 secret。

## 23. Code Review Checklist

Reviewer 應優先檢查：

1. 依賴方向是否被破壞。
2. binder、controller、service、repository 責任是否混層。
3. middleware/interceptor 順序是否會造成 auth、CSRF、timeout 或 response rewrite 漏洞。
4. resource-level authorization 是否完整，是否存在 IDOR。
5. transaction 是否涵蓋所有必要寫入，rollback/commit error 是否處理。
6. DB constraint 是否足以抵抗並發，而非只做先查再寫。
7. DTO 是否洩漏 schema 或敏感欄位。
8. Exception 是否穩定、安全且可診斷。
9. context cancellation 是否一路傳到 DB/外部服務。
10. 測試是否驗證 failure path，而不只 happy path。

## 24. 現況例外與收斂原則

- `app/models/scopes/` 目前尚未建立實際 scope；新需求有可重用 query composition 時應開始使用，不必等待全專案重構。
- `app/services/scheduling_service.go` 等既有 service 有直接 GORM 操作。新增大型 scheduling 功能時應逐步抽 repository，不應繼續擴大單檔與直接 query 範圍。
- 部分舊 route 尚未一致使用 authorized rate limit 或 CSRF。新 route 依本文 threat model 實作，舊 route 另開相容性修正。
- 部分既有命名與 REST path 屬歷史 contract。除非有版本化或相容 alias，不直接破壞 client。
- 收斂時以小步驟、可測試、保持 API 相容為原則，不進行與需求無關的全面重寫。
