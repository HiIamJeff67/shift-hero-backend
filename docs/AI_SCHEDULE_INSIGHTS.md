# AI Schedule Insights

這項功能把排班資料轉成可展示的 AI 主管 briefing，但不讓 LLM 自行計算工時。

## Workflow

```text
PostgreSQL schedule data
  -> deterministic_schedule_analyzer
     - coverage / unfilled headcount
     - weekly hours / workload spread
     - short rest / consecutive work days
     - night and weekend shifts
     - availability conflicts
     - active swap pressure
  -> LangChainGo Workforce Auditor Agent
  -> LangChainGo Sequential Chain
  -> Executive Narrator Agent
  -> JSON or SSE token stream
```

預設模型為 `openai/gpt-oss-20b:free`，可透過 `OPEN_ROUTER_MODEL` 覆寫。指定模型在尚未輸出 token 前失敗時，workflow 會自動重試 `openrouter/free`。OpenRouter free model 仍有流量與可用性限制，正式環境應改用穩定的付費模型。

## Environment

```env
OPEN_ROUTER_API_KEY=...
OPEN_ROUTER_MODEL=openai/gpt-oss-20b:free
```

Docker Compose 已把這兩個變數傳入 API container。

## API

只有 company Manager 可以讀取全體員工洞察，日期範圍最多 31 天。

一般 JSON：

```http
POST /api/development/v1/companies/{companyId}/ai/scheduleInsights
Authorization: Bearer <accessToken>
Content-Type: application/json
```

真正 streaming：

```http
POST /api/development/v1/companies/{companyId}/ai/scheduleInsights/stream
Authorization: Bearer <accessToken>
Content-Type: application/json
Accept: text/event-stream
```

共用 body：

```json
{
  "startAt": "2026-06-08T00:00:00+08:00",
  "endAt": "2026-06-15T00:00:00+08:00",
  "locale": "zh-TW",
  "focus": "優先檢查缺班、公平性與過勞風險"
}
```

SSE events：

```text
stage  deterministic analyzer 已開始
token  最終 briefing 的文字片段
done   完整結果、metrics、model 與 workflow metadata
error  已開始串流後發生的安全錯誤資訊
```

因為瀏覽器原生 `EventSource` 不支援 POST body 與 Authorization header，前端應使用 `fetch()` 讀取 response stream，或使用符合 SSE 規格的 parser library。

## Monthly AI Usage

AI 額度以「成功完成一次使用者要求」計算，不直接以 token 計算。一次 schedule insight 即使內部執行 auditor、narrator 或 fallback model，仍只扣一筆；模型失敗、timeout 或 streaming 中斷會補回預扣額度。

| Plan | Monthly generations |
| --- | ---: |
| Free | 5 |
| Pro | 30 |
| Premium | 100 |
| Ultimate | 300 |
| Enterprise | 1000 |

成功回應的 `aiUsage` 會提供 `used`、`limit`、`remaining`、`resetAt`。JSON API 超額時回傳 `429 AIUsageLimitExceeded`；SSE 因 headers 已送出，會透過 `error` event 傳回相同 exception。額度於每月 1 日 `00:00:00 UTC` 重置。

## Demo Ideas

下一批適合 LLM、且比單純 summary 更有展示效果的功能：

1. **What-if 排班模擬器**：主管輸入「Alice 週三請假」，系統複製目前班表、用 deterministic solver 找替代方案，再由 AI 解釋成本、公平性與風險差異。
2. **自然語言排班 Copilot**：輸入「下週五晚上至少三位 Staff，避免任何人連上六天」，LLM 只負責把語句轉成結構化 constraints，真正排班仍由演算法處理。
3. **Swap Matchmaker**：替換班請求找出最合適的三位候選人，AI 產生不洩漏私人資訊的邀請訊息與理由。
4. **Team Pulse Timeline**：每週保存 deterministic 指標，讓 AI 解釋「這週和上週相比發生什麼變化」，前端可呈現動畫化風險時間線。
5. **Manager Voice Briefing**：把目前的 briefing 接到 TTS，產生 30 秒語音晨會摘要；視覺上可同步標亮被提到的班次。

其中最推薦先做 **What-if 排班模擬器**。它同時展示 LLM 理解需求、系統實際推演、前後方案比較與可解釋 AI，實用性和新奇度都比聊天框高。

## Privacy

送往 OpenRouter 的內容包含排班統計與員工 display name，不包含內部 user UUID、email、token、密碼或 user profile。若部署到正式環境，應再加入：

- 公司層級 AI opt-in
- display name pseudonymization
- retention / provider privacy policy 告知
- prompt 與輸出稽核，但禁止記錄 secret 或完整個資
