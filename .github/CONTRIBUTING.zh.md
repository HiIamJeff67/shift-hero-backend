[English](./CONTRIBUTING.md) | [繁體中文](./CONTRIBUTING.zh.md)

# go-start-monolithic-kit 貢獻指南

感謝你願意參與 **go-start-monolithic-kit**（由 **Notezy** 開發的後端 template 架構）。

## 本文件適用範圍

這份 `.github` 規範是針對本 repository（`go-start-monolithic-kit`）的協作規則，
不直接約束由此 template 建立出的下游專案。

## 權利與授權

- 本 repository 的架構與實作由 Notezy 開發並持有完整控制權。
- Notezy 以 [`Apache-2.0`](../LICENSE) 公開授權此 repository。
- 你提交的貢獻一旦合併，將可在本 repository 內以 Apache-2.0 發布。

## 開發要求

1. `cp .env.example .env`
2. `go mod tidy`
3. `go build ./...`
4. 針對你的變更執行相對應測試。

## Pull Request 流程

1. 由 `main` 開出功能/修復分支。
2. commit 需聚焦且邏輯清楚。
3. 開 PR 時請使用 `.github/PULL_REQUEST_TEMPLATE/` 內對應語言模板。
4. PR 內容至少需包含：
   - 問題描述
   - 變更範圍與設計取捨
   - 驗證證據（build/test/log）

## 品質要求

- 除非是刻意的 repo 專案調整，否則命名與行為應維持 template 中性。
- 避免把不相關重構塞進同一個 PR。
- 行為或流程變更時必須更新文件。
- 發送 review 前須確保 `go build ./...` 可通過。

## 安全性回報

請勿在公開 issue/PR 直接揭露漏洞細節。
請依 [`SECURITY.md`](./SECURITY.md) 流程私下回報。

## 聯絡方式

- 維護者 / repository 聯絡信箱：`thenotezy@gmail.com`
