package handlers

import (
	"fmt"
	"time"

	"github.com/DimTur/multi_user_rw_sys/models"
)

// AddUsers simulates adding users and messages to message channels
func AddUsers(channels []chan models.Message, tokens []string, numMsg int) {
	for i := 1; i <= numMsg; i++ {
		for user, ch := range channels {
			msg := models.Message{
				Token:  tokens[user],
				FileID: fmt.Sprintf("file_%d", user),
				Data:   fmt.Sprintf("Message %d from user %d", i, user),
			}
			ch <- msg
		}
		time.Sleep(1 * time.Second)
	}

	for _, ch := range channels {
		close(ch)
	}
}
