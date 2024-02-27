package main

import (
	"database/sql"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/cmd/server/routes"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// NO MODIFICAR
	// had to add "?parseTime=true" for time support on product_record
	db, err := sql.Open("mysql", "meli_sprint_user:Meli_Sprint#123@/melisprint?parseTime=true")
	if err != nil {
		panic(err)
	}

	logging.InitLog(db)

	eng := gin.Default()

	router := routes.NewRouter(eng, db)
	router.MapRoutes()

	if err := eng.Run(); err != nil {
		panic(err)
	}
}
