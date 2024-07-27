package main

import "gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/infra/database"

type config struct {
	databasePool []database.Config
}

func newConfig() config {
	return config{
		databasePool: []database.Config{
			// индекс это номер шарда,
			// индексы важны и их нельзя менять
			0: {
				DSN: "postgresql://ws-8-user-1:ws-8-pass-1@localhost:8431/ws-8-db-1",
			},
			1: {
				DSN: "postgresql://ws-8-user-2:ws-8-pass-2@localhost:8432/ws-8-db-2",
			},
		},
	}
}
