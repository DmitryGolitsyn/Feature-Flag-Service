
План шагов (дальше от вашего текущего состояния)
11B (сделано/в процессе): типизация ошибок usecase, errors.Is, маппинг в HTTP-статусы.
12A. Домейн “Feature Flag” (in-memory):
ввести доменные типы Flag, Rule, Actor, валидацию;
интерфейс репозитория FlagRepository в usecase-слое;
in-memory адаптер;
usecase FlagUpsert, FlagEvaluate;
HTTP: POST /v1/flags (upsert), POST /v1/evaluate.
12B. Контракты DTO + валидация: строгие JSON-форматы, ошибки в единый формат, idempotency-key.
13A. БД (Postgres) фундамент:
модуль internal/infra/pg (подключение, миграции, migrate cmd);
репозиторий FlagRepository на SQL (вторая реализация).
13B. Транзакции и контекст: транзакционный boundary в usecase, context сквозной.
14A. Кэш Redis: read-through/ write-through для флагов, инвалидация.
14B. Метрики и алерты: HTTP p95, кэш hit/miss, DB latency, алерты на 5xx.
15A. Аутентификация/авторизация (минимум): JWT bearer, роль admin для /v1/flags.
15B. Аудит-лог: кто, что изменил (лучше Kafka “flag-events” → возможен CDC в будущем).
16A. Kafka (по требованию): продюсер “flag-updated”; консюмер “cache-invalidate”.
17A. Тесты: table-tests для usecase, http-handler tests (httptest), in-memory repo замены.
18A. CI: go vet, staticcheck, golangci-lint, go test -race.
19A. Конфиги: единый internal/config (HTTP, PG, Redis, JWT…), профили dev/stage/prod.
20A. Обзорка: логгер JSON, request-scoped поля (req_id, ip, ua), sampling.
21+. Харднинг: rate-limit, cors, pagination, SLO/SLI, pprof, readiness глубже.