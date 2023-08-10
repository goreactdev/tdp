package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/ton-developer-program/internal/database"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

const timeoutDuration = 2 * time.Minute

func checkHealth(heartbeat time.Time) bool {
    // Check if the time since the last heartbeat is less than our timeout.
    if time.Since(heartbeat) < timeoutDuration {
        // The goroutine has sent a heartbeat recently enough, so it's healthy.
        return true
    } else {
        // The goroutine hasn't sent a heartbeat recently enough, so it's unhealthy.
        return false
    }
}


func (app *application) getMetaData(url string) (*database.Metadata, database.JSONB, error)  {
    var metadata struct {
        Name		*string `json:"name"`
        Description	*string `json:"description"`
        Image		*string `json:"image"`
    }


    client := &http.Client{}

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        app.logger.Warning(fmt.Sprintf("error creating request: %v", err))
        return (*database.Metadata)(&metadata), nil, err
    }

    var body []byte

    // add timeout
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    req = req.WithContext(ctx)
    resp, err := client.Do(req)

    // if context deadline exceeded set metadata to nil
    if errors.Is(err, context.DeadlineExceeded) {
        app.logger.Warning(fmt.Sprintf("context deadline exceeded: %v", err))
        return (*database.Metadata)(&metadata), nil, nil
    }

    // if error is not nil and not context deadline exceeded return error
    if err != nil && !errors.Is(err, context.DeadlineExceeded) {
        app.logger.Warning(fmt.Sprintf("error getting metadata: %v", err))
        return (*database.Metadata)(&metadata), nil, err
    }

    defer resp.Body.Close()		


    // if response failed set body to empty
    if resp.StatusCode != http.StatusOK {
        body = []byte("")
    } else {
        body, err = io.ReadAll(resp.Body)
        if err != nil {
            app.logger.Warning(fmt.Sprintf("error reading body: %v", err))
            return (*database.Metadata)(&metadata), nil, err
        }
    }

    err = json.Unmarshal(body, &metadata)
    if err != nil {
        app.logger.Warning(fmt.Sprintf("error unmarshaling metadata: %v", err))
        return (*database.Metadata)(&metadata), nil, err
    }
    
    var bodyInterface database.JSONB

    err = json.Unmarshal(body, &bodyInterface)
    if err != nil {
            app.logger.Warning(fmt.Sprintf("error unmarshalling body: %v", err))
            bodyInterface = nil
            return (*database.Metadata)(&metadata), bodyInterface, err
    }

    return &database.Metadata{
        Name: metadata.Name,
        Description: metadata.Description,
        Image: metadata.Image,
    }, bodyInterface, nil
}



func (app *application) getWallet() *wallet.Wallet {
	words := strings.Split(app.config.App.SeedPhrase, "_")
	w, err := wallet.FromSeed(app.tonLiteClient, words, wallet.V4R2)
	if err != nil {
		panic(err)
	}
	return w
}


func (app *application) SkipError(err error, t *asynq.Task) error {

    var skipErrors = []string{
        "duplicate key value violates unique constraint",
    }

    for _, skipError := range skipErrors {
        if strings.Contains(err.Error(), skipError) {
            t.ResultWriter().Write([]byte("skipped: " + err.Error()))
            return nil
        }
    }
    
    return err
}