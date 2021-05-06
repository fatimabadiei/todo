package main

import (
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"strconv"
)

type Todo struct{
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	Done     bool   `json:"done"`
}

var todos = []*Todo{
	{Id: 1, Name: "Coding", Priority: 1, Done: false},
	{Id: 2, Name: "Study for Midterm", Priority: 4, Done: false},
	{Id: 3, Name: "Workout", Priority: 3, Done: false},
}

func main (){
	app := fiber.New()

	app.Use(middleware.Logger())

	//app.Get("/", func(ctx *fiber.Ctx) {
	//	ctx.Send("hello")
	//})

	SetupTodosRoutes(app)

	err := app.Listen(3000)
	if err != nil {
		panic(err)
	}
}

func SetupTodosRoutes(app *fiber.App) {
	todosRoutes := app.Group("/todos")
	todosRoutes.Get("/", GetTodos)
	todosRoutes.Get("/:id", GetTodo)
	todosRoutes.Post("/", CreatTodos)
	todosRoutes.Delete("/:id", DeleteTodo)
	todosRoutes.Patch("/:id", UpdateTodo)
}


func GetTodos(ctx *fiber.Ctx) {
	ctx.Status(fiber.StatusOK).JSON(todos)
}


func CreatTodos(ctx *fiber.Ctx) {
	type request struct {
		Name       string `json:"name"`
		Priority   int    `json:"priority"`
	}

	var body request

	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse json",
		})
		return
	}

	todo := &Todo{
		Id:        len(todos)+1,
		Name:      body.Name,
		Priority:  body.Priority,
		Done:      false,
	}

	todos = append(todos, todo)

	ctx.Status(fiber.StatusCreated).JSON(todo)
}


func GetTodo(ctx *fiber.Ctx) {
	paramId := ctx.Params("id")
	id , err := strconv.Atoi(paramId)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "cannot parse id",
		})
		return
	}

	for _, todo := range todos {
		if todo.Id == id {
			ctx.Status(fiber.StatusOK).JSON(todo)
			return
		}
	}

	ctx.Status(fiber.StatusNotFound)
}

func DeleteTodo(ctx *fiber.Ctx) {
	paramId := ctx.Params("id")
	id , err := strconv.Atoi(paramId)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "cannot parse id",
		})
		return
	}

	for i, todo := range todos {
		if todo.Id == id {
			todos = append(todos[0:i], todos[i+1:]...)
			ctx.Status(fiber.StatusNoContent)
			return
		}
	}

	ctx.Status(fiber.StatusNotFound)
}

func UpdateTodo(ctx *fiber.Ctx) {
	type request struct {
		Name       *string `json:"name"`
		Priority   *int    `json:"priority"`
		Done       *bool     `json:"done"`
	}

	paramId := ctx.Params("id")
	id , err := strconv.Atoi(paramId)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "cannot parse id",
		})
		return
	}

	var body request
	err = ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse body",
		})
		return
	}

	var todo *Todo

	for _,t := range todos {
		if t.Id == id {
			todo = t
			break
		}
	}

	if todo == nil {
		ctx.Status(fiber.StatusNotFound)
		return
	}

	if body.Name != nil {
		todo.Name = *body.Name
	}

	if body.Done != nil {
		todo.Done = *body.Done
	}

	ctx.Status(fiber.StatusOK).JSON(todo)

}

//func GetDb() {
//	conn, err := http.NewConnection(http.ConnectionConfig{
//		Endpoints: []string{"http://localhost:8529"},
//	})
//	if err != nil {
//		panic("Cannot Connect To Arangodb")
//	}
//	c, err := driver.NewClient(driver.ClientConfig{
//		Connection: conn,
//		Authentication: driver.BasicAuthentication("root", "root123***"),
//	})
//	if err != nil {
//		panic("Error happend while create NEW CLIENT")
//	}
//}