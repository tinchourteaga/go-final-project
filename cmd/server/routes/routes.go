package routes

import (
	"database/sql"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/record/product_record"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/record/report_record"
	"os"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/buyer"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/carry"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/employee"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/inbound_order"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/locality"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/product"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/productBatch"
	purchaseorders "github.com/extmatperez/meli_bootcamp_go_w6-2/internal/purchase_orders"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/section"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/seller"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/warehouse"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router interface {
	MapRoutes()
}

type router struct {
	eng *gin.Engine
	rg  *gin.RouterGroup
	db  *sql.DB
}

func NewRouter(eng *gin.Engine, db *sql.DB) Router {
	return &router{eng: eng, db: db}
}

func (r *router) MapRoutes() {
	r.setGroup()
	r.buildSwaggerRoutes()
	r.buildSellerRoutes()
	r.buildProductRoutes()
	r.buildProductRecordsRoutes()
	r.buildSectionRoutes()
	r.buildWarehouseRoutes()
	r.buildEmployeeRoutes()
	r.buildBuyerRoutes()
	r.buildPurchaseOrderRoutes()
	r.buildProductBatchRoutes()
	r.buildInboundOrderRoutes()
	r.buildCarryRoutes()
	r.buildLocalityRoutes()
}

func (r *router) setGroup() {
	r.rg = r.eng.Group("/api/v1")
}

func (r *router) buildSwaggerRoutes() {
	_ = godotenv.Load()
	docs.SwaggerInfo.Host = os.Getenv("HOST")
	r.eng.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (r *router) buildSellerRoutes() {
	repo := seller.NewRepository(r.db)
	service := seller.NewService(repo)
	handler := handler.NewSeller(service)
	sell := r.rg.Group("/sellers")

	sell.POST("", handler.Create())
	sell.GET("", handler.GetAll())
	sell.GET("/:id", handler.Get())
	sell.PATCH("/:id", handler.Update())
	sell.DELETE("/:id", handler.Delete())
}

func (r *router) buildProductRoutes() {
	productRepository := product.NewRepository(r.db)
	productService := product.NewService(productRepository)
	productHandler := handler.NewProduct(productService)
	productGroup := r.rg.Group("/products")
	productGroup.DELETE("/:id", productHandler.Delete())
	productGroup.PATCH("/:id", productHandler.PartialUpdate())
	productGroup.POST("/", productHandler.Create())
	productGroup.GET("/:id", productHandler.Get())
	productGroup.GET("/", productHandler.GetAll())

	reportRecordRepository := report_record.NewRepository(r.db)
	reportRecordService := report_record.NewService(reportRecordRepository)
	reportRecordHandler := handler.NewReportRecord(reportRecordService)
	productGroup.GET("/reportRecords", reportRecordHandler.GetReportRecords())
}

func (r *router) buildProductRecordsRoutes() {
	productRecordRepository := product_record.NewRepository(r.db)
	productRecordService := product_record.NewService(productRecordRepository)
	productRecordHandler := handler.NewProductRecord(productRecordService)
	productRecordGroup := r.rg.Group("/productRecords")
	productRecordGroup.POST("/", productRecordHandler.Create())
}

func (r *router) buildSectionRoutes() {
	repo := section.NewRepository(r.db)
	service := section.NewService(repo)
	handler := handler.NewSection(service)
	sec := r.rg.Group("/sections")

	sec.GET("/", handler.GetAll())
	sec.GET("/:id", handler.Get())
	sec.POST("/", handler.Create())
	sec.PATCH("/:id", handler.Update())
	sec.DELETE("/:id", handler.Delete())
	sec.GET("/reportProducts", handler.GetSectionProducts())

}

func (r *router) buildWarehouseRoutes() {
	repo := warehouse.NewRepository(r.db)
	service := warehouse.NewService(repo)
	controller := handler.NewWarehouse(service)
	warehouseRouter := r.rg.Group("/warehouses")
	warehouseRouter.GET("/", controller.GetAll)
	warehouseRouter.GET("/:id", controller.Get)
	warehouseRouter.POST("/", controller.Create)
	warehouseRouter.PATCH("/:id", controller.Update)
	warehouseRouter.DELETE("/:id", controller.Delete)
}

func (router *router) buildEmployeeRoutes() {
	repoEmployee := employee.NewRepository(router.db)
	serviceEmployee := employee.NewService(repoEmployee)
	handlerEmployee := handler.NewEmployee(serviceEmployee)
	employeesRoutesGroup := router.rg.Group("/employees")

	employeesRoutesGroup.GET("/", handlerEmployee.GetAll())
	employeesRoutesGroup.GET("/:id", handlerEmployee.Get())
	employeesRoutesGroup.POST("/", handlerEmployee.Create())
	employeesRoutesGroup.PATCH("/:id", handlerEmployee.Update())
	employeesRoutesGroup.DELETE("/:id", handlerEmployee.Delete())

	repoInboundOrder := inbound_order.NewRepository(router.db)
	serviceInboundOrder := inbound_order.NewService(repoInboundOrder)
	handlerInboundOrder := handler.NewInboundOrder(serviceInboundOrder)

	employeesRoutesGroup.GET("/reportInboundOrders", handlerInboundOrder.GetAllEmployeesInboundOrders())
	employeesRoutesGroup.GET("/reportInboundOrders/:id", handlerInboundOrder.GetEmployeeInboundOrders())
}

func (r *router) buildBuyerRoutes() {
	repo := buyer.NewRepository(r.db)
	service := buyer.NewService(repo)
	handler := handler.NewBuyer(service)
	sec := r.rg.Group("/buyers")
	sec.GET("/", handler.GetAll())
	sec.GET("/:id", handler.Get())
	sec.POST("/", handler.Create())
	sec.PATCH("/:id", handler.Update())
	sec.DELETE("/:id", handler.Delete())
}

func (r *router) buildPurchaseOrderRoutes() {
	repo := purchaseorders.NewRepository(r.db)
	service := purchaseorders.NewService(repo)
	handler := handler.NewPurchaseOrders(service)

	sec := r.rg.Group("/purchase_orders")
	sec.POST("/", handler.CreateOrder())

	rep := r.rg.Group("/reportPurchaseOrder")
	rep.GET("", handler.GetAllOrdersByBuyers())
}

func (r *router) buildProductBatchRoutes() {
	repo := productbatch.NewRepository(r.db)
	service := productbatch.NewService(repo)
	handler := handler.NewProductBatch(service)
	group := r.rg.Group("/productBatches")
	group.POST("/", handler.Create())
}

func (router *router) buildInboundOrderRoutes() {
	repo := inbound_order.NewRepository(router.db)
	service := inbound_order.NewService(repo)
	handler := handler.NewInboundOrder(service)
	inboundOrdersRoutesGroup := router.rg.Group("/inboundOrders")

	inboundOrdersRoutesGroup.POST("/", handler.Create())
}

func (r *router) buildCarryRoutes() {
	repo := carry.NewRepository(r.db)
	service := carry.NewService(repo)
	controller := handler.NewCarry(service)
	carryRouter := r.rg.Group("/carries")
	carryRouter.POST("/", controller.Save)
}

func (r *router) buildLocalityRoutes() {
	repo := locality.NewRepository(r.db)
	service := locality.NewService(repo)
	handler := handler.NewLocality(service)
	loc := r.rg.Group("/localities")

	loc.POST("", handler.Create())
	loc.GET("/:id", handler.Get())
	loc.GET("/reportSellers", handler.GetReportSellers())
	loc.GET("/reportCarries", handler.GetReportCarries())
}
