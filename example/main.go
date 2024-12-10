package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/burke/gocronometer"
    "github.com/burke/cronometerflux"
)

func getYesterday() string {
    now := time.Now()
    yesterday := now.AddDate(0, 0, -1)
    return yesterday.Format("2006-01-02")
}

func main() {
    yesterday := getYesterday()
    username := flag.String("username", os.Getenv("CRONOMETER_USER"), "Cronometer username (or set CRONOMETER_USER)")
    password := flag.String("password", os.Getenv("CRONOMETER_PASS"), "Cronometer password (or set CRONOMETER_PASS)")
    startDate := flag.String("start", yesterday, "Start date (YYYY-MM-DD)")
    endDate := flag.String("end", yesterday, "End date (YYYY-MM-DD)")
    flag.Parse()

    if *username == "" || *password == "" {
        log.Fatal("Username and password are required. Set via flags or CRONOMETER_USER/CRONOMETER_PASS environment variables")
    }

    start, err := time.Parse("2006-01-02", *startDate)
    if err != nil {
        log.Fatalf("Invalid start date format, must be YYYY-MM-DD: %v", err)
    }

    end, err := time.Parse("2006-01-02", *endDate)
    if err != nil {
        log.Fatalf("Invalid end date format, must be YYYY-MM-DD: %v", err)
    }

    if end.Before(start) {
        log.Fatal("End date must not be before start date")
    }

    client := gocronometer.NewClient(nil)
    ctx := context.Background()
    
    if err := client.Login(ctx, *username, *password); err != nil {
        log.Fatalf("Failed to login: %v", err)
    }
    defer client.Logout(ctx)

    servings, err := client.ExportServingsParsedWithLocation(ctx, start, end, time.Local)
    if err != nil {
        log.Fatalf("Failed to fetch servings: %v", err)
    }

    for _, line := range cronometerflux.FormatServings(servings) {
        fmt.Println(line)
    }
}