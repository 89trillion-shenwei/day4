package http

import "day4/internal/router"

func Start() {
	r := router.MongoRouter()

	r.Run("localhost:8080")

}
