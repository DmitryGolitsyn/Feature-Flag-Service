package app

// Application — корневой объект бизнес-логики (usecase-слой).
// Здесь будут все сценарии: Ping, UserCreate, OrderProcess, NotificationSend, etc.
type Application struct {
	Ping  *PingUsecase
	Echo  *EchoUsecase
	Flags struct {
		Upsert   *FlagUpsertUsecase
		Evaluate *FlagEvaluateUsecase
	}
}

func NewWithRepo(repo FlagRepository) *Application {
	return &Application{
		Ping: NewPingUsecase(),
		Echo: NewEchoUsecase(),
		Flags: struct {
			Upsert   *FlagUpsertUsecase
			Evaluate *FlagEvaluateUsecase
		}{
			Upsert:   NewFlagUpsertUsecase(repo),
			Evaluate: NewFlagEvaluateUsecase(repo),
		},
	}
}

// New создаёт приложение и все его usecase.
// Пока они пустые, но скоро появится логика и зависимости (БД, Kafka, кеши).
func New() *Application {
	return &Application{
		Ping: NewPingUsecase(),
		Echo: NewEchoUsecase(),
	}
}
