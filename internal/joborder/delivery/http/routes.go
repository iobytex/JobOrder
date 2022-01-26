package http

func (joh *jobOrderHandler) JobOrderMapRoute(){
	joh.group.POST("/create/user",joh.CreateUser())
	joh.group.GET("/check/user",joh.middleware.LoginHandler)

	joh.group.GET("/refresh_token",joh.middleware.RefreshHandler)

	categories := joh.group.Group("/categories")
	{
		categories.Use(joh.middleware.MiddlewareFunc())
		{
			categories.POST("/add", joh.AddCategories())
			categories.GET("/all", joh.GetCategories())
			categories.DELETE("/delete", joh.DeleteCategories())
			categories.PATCH("/update", joh.UpdateCategory())
		}
	}

	product := joh.group.Group("/products",joh.AddProducts())
	{
		product.Use(joh.middleware.MiddlewareFunc())
		{
			product.POST("/add", joh.AddProducts())
			product.GET("/all", joh.GetProducts())
			product.DELETE("/delete", joh.DeleteProducts())
			product.PATCH("/update", joh.UpdateProduct())
		}
	}

	order := joh.group.Group("/order")
	{
		order.Use(joh.middleware.MiddlewareFunc())
		{
			order.POST("/add", joh.SetOrder())
			order.POST("/add/items", joh.AddOrderItems())
			order.GET("/all/items", joh.GetOrderItems())
		}
	}

}
