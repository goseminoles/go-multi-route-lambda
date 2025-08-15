package main

import (
    "context"
    "log"
    "os"

    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
    echo "github.com/labstack/echo/v4"
)

var echoAdapter *echoadapter.EchoLambdaV2

func main() {
    // Get stage from environment variable
    env := os.Getenv("ENV")

    // Initialize Echo
    e := echo.New()

    // Group routes under stage prefix
    g := e.Group("")
    if env != "" {
        g = e.Group("/" + env)
    }

    g.GET("/", func(c echo.Context) error {
        return c.JSON(200, map[string]string{"status": "ok"})
    })

    g.GET("/health", func(c echo.Context) error {
        return c.JSON(200, map[string]string{"status": "healthy"})
    })

    g.POST("/quote", func(c echo.Context) error {
        return c.JSON(201, map[string]interface{}{
            "quote_number": "1234567890",
            "carrier":      "USCD",
            "rate":         299.99,
        })
    })

    // Adapter for Lambda (HTTP API v2)
    echoAdapter = echoadapter.NewV2(e)

    // Local server vs Lambda
    if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
        log.Println("Running locally at http://localhost:8080")
        e.Logger.Fatal(e.Start(":8080"))
    } else {
        lambda.Start(Handler)
    }
}

// Lambda handler for HTTP API v2
func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
    return echoAdapter.ProxyWithContext(ctx, req)
}
