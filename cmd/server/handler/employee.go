package handler

import (
	"net/http"
	"strconv"

	"github.com/extmatperez/meli_bootcamp_go_w6-2/cmd/server/handler/requests"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/internal/employee"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/logging"
	"github.com/extmatperez/meli_bootcamp_go_w6-2/pkg/web"
	"github.com/gin-gonic/gin"
)

type Employee struct {
	employeeService employee.Service
}

func NewEmployee(e employee.Service) *Employee {
	return &Employee{
		employeeService: e,
	}
}

// Get godoc
// @Summary     Get employee by ID
// @Tags        Employees
// @Description Retrieves existing employee by ID from database
// @Produce     json
// @Param       id  path     int               true "Employee id"
// @Success     200 {object} web.response      "Employee"
// @Failure     400 {object} web.errorResponse "Invalid id type"
// @Failure     404 {object} web.errorResponse "Employee not found"
// @Failure     500 {object} web.errorResponse "Connection to database error"
// @Router      /api/v1/employees/{id} [get]
func (e *Employee) Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			logging.Log("invalid id")
			web.Error(ctx, http.StatusBadRequest, "invalid id")
			return
		}

		obtainedEmployee, err := e.employeeService.Get(ctx, id)

		if err != nil {
			logging.Log(err)
			switch err.Error() {
			case employee.ErrEmployeeNotFound.Error():
				web.Error(ctx, http.StatusNotFound, employee.ErrEmployeeNotFound.Error())
			default:
				web.Error(ctx, http.StatusInternalServerError, err.Error())
			}
			return
		}

		web.Success(ctx, http.StatusOK, obtainedEmployee)
	}
}

// GetAll godoc
// @Summary     List employees
// @Tags        Employees
// @Description Lists all existing employees from database
// @Produce     json
// @Success     200 {object} web.response      "List of employees"
// @Failure     500 {object} web.errorResponse "Connection to database error"
// @Router      /api/v1/employees [get]
func (employee *Employee) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		employees, err := employee.employeeService.GetAll(ctx)

		if err != nil {
			logging.Log(err)
			web.Error(ctx, http.StatusInternalServerError, err.Error())
			return
		}

		if employees == nil {
			web.Success(ctx, http.StatusOK, []domain.Employee{})
			return
		}

		web.Success(ctx, http.StatusOK, employees)
	}
}

// Create godoc
// @Summary     Create employee
// @Tags        Employees
// @Description Creates a new employee in database
// @Accept      json
// @Produce     json
// @Param       employee body     requests.EmployeeDTOPost true "Employee to be stored"
// @Success     201      {object} web.response             "Employee created"
// @Failure     409      {object} web.errorResponse        "Employee ID already exists error"
// @Failure     422      {object} web.errorResponse        "Missing field or type casting error"
// @Failure     500      {object} web.errorResponse        "Connection to database error"
// @Router      /api/v1/employees [post]
func (e *Employee) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req requests.EmployeeDTOPost

		if err := ctx.ShouldBindJSON(&req); err != nil {
			logging.Log(err)
			web.Error(ctx, http.StatusUnprocessableEntity, err.Error()) // TODO: fmt.Sprintf("el campo %s es requerido", strings.Split(err.Error(), "'")[3])
			return
		}

		newEmployee := domain.Employee{CardNumberID: *req.CardNumberID, FirstName: *req.FirstName, LastName: *req.LastName, WarehouseID: *req.WarehouseID}

		newEmployee, err := e.employeeService.Save(ctx, newEmployee)

		if err != nil {
			logging.Log(err)
			switch err.Error() {
			case employee.ErrEmployeeAlreadyExists.Error():
				web.Error(ctx, http.StatusConflict, employee.ErrEmployeeAlreadyExists.Error())
			case employee.ErrEmployeeNotSaved.Error():
				web.Error(ctx, http.StatusInternalServerError, employee.ErrEmployeeNotSaved.Error())
			default:
				web.Error(ctx, http.StatusInternalServerError, err.Error())
			}
			return
		}

		web.Success(ctx, http.StatusCreated, newEmployee)
	}
}

// Update godoc
// @Summary     Update employee
// @Tags        Employees
// @Description Updates information of an existing employee in database
// @Accept      json
// @Produce     json
// @Param       id       path     int                       true "Employee id"
// @Param       employee body     requests.EmployeeDTOPatch true "Employee to update"
// @Success     200      {object} web.response              "Employee updated"
// @Failure     400      {object} web.errorResponse         "Invalid id type"
// @Failure     404      {object} web.errorResponse         "Employee not found"
// @Failure     422      {object} web.errorResponse         "Missing field or type casting error"
// @Failure     500      {object} web.errorResponse         "Connection to database error"
// @Router      /api/v1/employees/{id} [patch]
func (e *Employee) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req requests.EmployeeDTOPatch
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			logging.Log("invalid id")
			web.Error(ctx, http.StatusBadRequest, "invalid id")
			return
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			logging.Log(err)
			web.Error(ctx, http.StatusUnprocessableEntity, err.Error())
			return
		}

		employeeToUpdate := domain.Employee{ID: id, FirstName: req.FirstName, LastName: req.LastName, WarehouseID: req.WarehouseID}

		employeeToUpdate, err = e.employeeService.Update(ctx, employeeToUpdate)

		if err != nil {
			logging.Log(err)
			switch err.Error() {
			case employee.ErrEmployeeNotFound.Error():
				web.Error(ctx, http.StatusNotFound, employee.ErrEmployeeNotFound.Error())
			case employee.ErrEmployeeNotUpdated.Error():
				web.Error(ctx, http.StatusInternalServerError, employee.ErrEmployeeNotUpdated.Error())
			default:
				web.Error(ctx, http.StatusInternalServerError, err.Error())
			}
			return
		}

		web.Success(ctx, http.StatusOK, employeeToUpdate)
	}
}

// Delete godoc
// @Summary     Delete employee
// @Tags        Employees
// @Description Deletes an existing employee from database
// @Produce     json
// @Param       id path int true "Employee id"
// @Success     204
// @Failure     400 {object} web.errorResponse "Invalid id type"
// @Failure     404 {object} web.errorResponse "Employee not found"
// @Failure     500 {object} web.errorResponse "Connection to dabatase error"
// @Router      /api/v1/employees/{id} [delete]
func (e *Employee) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			logging.Log("invalid ip")
			web.Error(ctx, http.StatusBadRequest, "invalid id")
			return
		}

		err = e.employeeService.Delete(ctx, id)

		if err != nil {
			logging.Log(err)
			switch err.Error() {
			case employee.ErrEmployeeNotFound.Error():
				web.Error(ctx, http.StatusNotFound, employee.ErrEmployeeNotFound.Error())
			default:
				web.Error(ctx, http.StatusInternalServerError, err.Error())
			}
			return
		}

		web.Success(ctx, http.StatusNoContent, "")
	}
}
