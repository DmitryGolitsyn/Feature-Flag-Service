Окей, сделаем интенсивную программу “с руками” под твою скорость. Проект — не игрушка, а мини-продукт, который закрывает все нужные идиомы Go: конкуренция, контексты, ошибки, слойность, работа с БД, кэш, конфиги, тесты, профилирование, observability.

# Тема учебного проекта

**Feature Flags / Experimentation Service (FFS)**
Сервис фичефлагов с двумя API:

1. **Admin API** — CRUD флагов, правил таргетинга, rollout-ов.
2. **SDK API** — очень быстрый `POST /v1/evaluate` (или `GET /v1/flags/{key}`) для приложений/микросервисов.

Хранилище — **PostgreSQL** (истина), кэш — **Redis** (низкие задержки), event stream для инвалидации кэша и realtime-обновлений — по желанию **NATS** или **Kafka** (модуль “плюс”).
Особенности: мульти-тенантность, идемпотентность, правило таргетинга (по userID, процентная раскатка), строгая схема миграций, трейсинг, метрики, pprof.

Почему это идеальный учебный проект:

* **Go-идиомы**: контексты, каналы (инвалидация/фоновые воркеры), table-driven tests, error wrapping, интерфейсы на границах пакетов.
* **Производительность**: горячий путь SDK-чтений из Redis (p95↘), бенчмарки, профили.
* **Архитектура**: слои (transport → usecase/service → domain → repository), границы, низкая связность.
* **Реализм**: всё, что реально требуется в проде: миграции, healthz/readyz, конфиги, логирование, observability, CI.

---

# Техстек и принципы

* **Go 1.22+**, модули, `slog` для логов, `context.Context` везде.
* **HTTP**: `github.com/go-chi/chi/v5`
* **DB**: `pgx/v5`, миграции `golang-migrate`.
* **Cache**: `github.com/redis/go-redis/v9`
* **(Опционально)** брокер: NATS (`nats.go`) или Kafka (`segmentio/kafka-go`).
* **Config**: 12-Factor, ENV + простая обёртка (например, `caarlos0/env`).
* **Validation**: `go-playground/validator/v10` для DTO.
* **Testing**: стандартная библиотека + `testify` (assert/require), `httptest`, golden-файлы для правил.
* **Lint/QA**: `golangci-lint`, `make test`, `make bench`, `make cover`.
* **Observability**: Prometheus экспортёр (`/metrics`), OpenTelemetry (трейсы), `net/http/pprof`.
* **Докеризация**: `Dockerfile` + `docker-compose.yml` (pg/redis + сервис).
* **CI (скелет)**: `golangci-lint`, тесты, сборка образа.

---

# Архструктура репо

```
ffs/
  cmd/
    ffs-admin/       # бинарь для Admin+SDK HTTP сервера (один сервис)
  internal/
    app/
      server.go      # wiring HTTP, middlewares, routes, graceful shutdown
    transport/
      http/
        admin_handlers.go
        sdk_handlers.go
        middleware.go
        dto.go
    domain/
      flag.go        # Aggregate: Flag, Rule, Targeting
      errors.go
      eval.go        # чистая логика вычисления флага
    usecase/         # orchestration/application services
      flag_service.go
      eval_service.go
    repository/
      flag_repo.go   # интерфейсы
      pg/
        flag_repo_pg.go
      redis/
        eval_cache.go
    infra/
      db.go          # pgx pool
      redis.go
      events.go      # nats/kafka (опц.)
      migrate.go
      logger.go
      config.go
    observability/
      metrics.go     # Prometheus counters/hist, SDK latency, cache hit ratio
      tracing.go
      pprof.go
  migrations/
    001_init.sql
    002_rules.sql
  api/
    openapi.yaml
  scripts/
    devdata.sql
  Makefile
  docker-compose.yml
  Dockerfile
  Taskfile.yml (опц.)
  .golangci.yml
  README.md
```

---

# Бэклог по модулям (3 недели, быстро и по делу)

## Неделя 1 — Каркас, домен, Postgres

**Цель:** запустить базовый сервер, Admin API CRUD флагов в PostgreSQL, миграции, тесты домена.

1. **Инициализация**

    * Модули, Makefile, линтеры, Docker Compose (pg+redis), базовый `cmd/ffs-admin`.
    * Конфиги: PORT, PG_DSN, REDIS_ADDR, LOG_LEVEL, OTEL_EXPORTER.

2. **Домен**

    * `domain.Flag`: ключ, описание, версия (для CAS), targeting rules:

        * env/tenant
        * allowlist/denylist
        * процентное распределение (deterministic hash(userID)%100 < rollout)
    * Чистая функция `Evaluate(flag, context) (bool, reason)` без побочек.

3. **Postgres + миграции**

    * `migrations/*`, `infra/migrate.go`.
    * `repository.pg.FlagRepo`: `Create/Update/Delete/Get/List`, optimistic locking по `version`.

4. **Admin HTTP (chi)**

    * `POST /admin/flags`, `PUT /admin/flags/{key}`, `GET /admin/flags/{key}`, `GET /admin/flags`, `DELETE /admin/flags/{key}`
    * Валидация DTO, маппинг ↔ domain, ошибки → корректные HTTP-коды.

5. **Тесты**

    * Table-driven для `domain/eval.go` (ключевой блок: проценты/списки/композиция).
    * Интеграционные для репозитория через `testcontainers` **или** docker-compose target.

**DoD:** `make up`, `make migrate`, CRUD работает, 80%+ покрытие `domain`, линтеры зелёные.

---

## Неделя 2 — SDK API, Redis-кэш, производительность

**Цель:** быстрый путь чтения, кэш-инвалидация, метрики, pprof, бенчмарки.

1. **SDK HTTP**

    * `POST /v1/evaluate` (body: tenant, userID, attributes, flagKey[]), ответ — решения по ключам.
    * Контекст с таймаутом, дедлайны, корректная отмена.

2. **Кэш**

    * Read-through для флагов: сначала Redis (hash/json), miss → PG → наполнить кэш.
    * Инвалидация: при Update/Delete — publish event (простая версия: локальный in-proc bus + TTL; прод-версия: NATS/Kafka).

3. **Метрики и профили**

    * Prometheus: `sdk_requests_total`, `sdk_duration_seconds`, `cache_hit_ratio`, `db_latency`, `errors_total`.
    * `pprof` повешен на `/debug/pprof`.
    * Бенчмарки `*_test.go` для Evaluate и SDK handler (httptest + b.N).

4. **Логи и трейсинг**

    * `slog` структурировано, кореляция request-id, sampling.
    * Otel-трейсы для SDK пути (handler → usecase → repo/cache).

**DoD:** p95 SDK < 5–10ms при cache-hit на локале; cache hit-ratio виден; pprof/метрики доступны.

---

## Неделя 3 — Конкуренция, фоновые воркеры, отказоустойчивость + бонусы

**Цель:** освоить конкурентные паттерны и эксплуатационные аспекты.

1. **Фоновые процессы**

    * Воркеры инвалидации кэша (подписка на события), worker-pool с backoff.
    * Bulk-предзагрузка флагов “горячего” тенанта в Redis при старте.

2. **Устойчивость**

    * Rate limit (token bucket для SDK пути) + circuit breaker (простая реализация на repo).
    * Idempotency-key для Admin write-операций.
    * Graceful shutdown: контексты, `server.Shutdown`, ожидание воркеров.

3. **Граница модулей**

    * Интерфейсы только на границах (`repository.FlagRepository`, `usecase.FlagService`), внутри — конкретные типы.
    * Простая manual-DI (никаких магических контейнеров; код-wiring в `app/server.go`).

4. **(Опционально) Event Streaming**

    * Поднять NATS/Kafka: публиковать `FlagUpdated/Deleted`, подписчики инвалидируют кэш на всех инстансах.
    * E2E тесты на “два инстанса сервиса” (docker-compose scale).

5. **Готовность к прод-демо**

    * Healthz/Readyz (DB, Redis, streaming).
    * OpenAPI файл в `api/openapi.yaml`.
    * Сценарии нагрузочного теста (k6 или `vegeta`) — таргет SDK API.

**DoD:** два инстанса сервиса + брокер → кэш синхронно инвалидируется, SDK стабильный под нагрузкой.

---

# Контрольные точки (быстро проверять, что ты освоил)

* **Go-идиомы**:
  table-driven tests, контекст в каждой публичной функции, error wrapping `fmt.Errorf("...: %w", err)`, zero-values-friendly структуры, интерфейсы узкие и только на границах.
* **Конкуренция**:
  worker-pool, backpressure (буферизированные каналы), context cancellation, `select` с таймаутами, `errgroup` (из `x/sync/errgroup`) для параллельных запросов к repo+cache.
* **Производительность**:
  pprof CPU/heap, flamegraph, бенчмарки, уменьшение аллокаций, батч-загрузки.
* **Эксплуатация**:
  миграции вперёд/назад, observability, rate limit, graceful shutdown.

---

# Примерные фрагменты (чтобы старт был мгновенным)

**Интерфейс репозитория**

```go
// internal/repository/flag_repo.go
package repository

import "context"

type Flag struct {
    Key       string
    Version   int64
    Rules     []byte // JSON rules; в domain парсим в структуру
    Tenant    string
    UpdatedAt int64
}

type FlagRepository interface {
    Get(ctx context.Context, tenant, key string) (Flag, error)
    List(ctx context.Context, tenant string) ([]Flag, error)
    Create(ctx context.Context, f Flag) error
    Update(ctx context.Context, f Flag) error // optimistic by Version
    Delete(ctx context.Context, tenant, key string, version int64) error
}
```

**Table-driven тест домена**

```go
func TestEvaluate(t *testing.T) {
  cases := []struct{
    name string
    flag Flag
    ctx  EvalContext
    want bool
  }{
    {"allow user", flagAllow("u1"), EvalContext{UserID:"u1"}, true},
    {"50 percent", flagPercent(50), EvalContext{UserID:"u9"}, true}, // детерминированный хэш
  }
  for _, c := range cases {
    t.Run(c.name, func(t *testing.T) {
      got, _ := Evaluate(c.flag, c.ctx)
      if got != c.want { t.Fatalf("want %v got %v", c.want, got) }
    })
  }
}
```

**Маршруты (chi)**

```go
r := chi.NewRouter()
r.Use(middleware.RequestID, middleware.RealIP, middleware.Recoverer, loggerMiddleware)
r.Route("/admin", func(r chi.Router) {
    r.Post("/flags", h.CreateFlag)
    r.Get("/flags/{key}", h.GetFlag)
    r.Put("/flags/{key}", h.UpdateFlag)
    r.Delete("/flags/{key}", h.DeleteFlag)
    r.Get("/flags", h.ListFlags)
})
r.Post("/v1/evaluate", h.Evaluate)
r.Get("/healthz", h.Healthz)
```

---

# Мини-план на каждый день (если хочешь turbo-режим)

**День 1–2:** каркас, конфиг, миграции, домен, CRUD в PG
**День 3:** Admin API + юнит/интеграционные тесты
**День 4:** SDK API, Redis read-through
**День 5:** метрики, pprof, бенчи, доводим p95
**День 6–7:** воркеры инвалидации, rate limit, graceful shutdown
**День 8+:** брокер событий, E2E, нагрузочные, CI-скелет

---

# Критерии “овладел” (самопроверка)

* Можешь за 10–15 минут добавить новое правило таргетинга с тестами.
* По flamegraph видишь “узкое место” и устраняешь 1–2 аллокации.
* Умеешь локально воспроизвести “дрожание” p95 и снизить его (кэш/пулы/батчи).
* Чётко разводишь domain/usecase/repo/transport и не протаскиваешь лишние зависимости.
* Любую write-операцию делаешь идемпотентной и понимаешь, где нужен CAS.

---

Хочешь — я сразу сгенерирую скелет репозитория (папки, Makefile, `docker-compose.yml`, базовые файлы, заглушки кода и тестов), чтобы ты просто сделал `make up && make migrate && make test` и поехал.
