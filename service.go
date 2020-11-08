package homestuck_watcher

import (
	"fmt"
	"github.com/JMTyler/homestuck-watcher/internal/db"
	"github.com/JMTyler/homestuck-watcher/internal/fcm"
)

type Service struct {
	Body map[string]interface{}
}

func (s *Service) Subscribe(token string) error {
	if err := fcm.Subscribe("v1", token); err != nil {
		// TODO: Gotta start using log.Fatal() and its ilk.
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *Service) Unsubscribe(token string) error {
	if err := fcm.Unsubscribe("v1", token); err != nil {
		// TODO: Gotta start using log.Fatal() and its ilk.
		fmt.Println(err)
		return err
	}

	return nil
}

func (_ *Service) GetStories() ([]map[string]interface{}, error) {
	stories := new(db.Story).FindAll("v1")
	scrubbed := make([]map[string]interface{}, len(stories))
	for i, model := range stories {
		scrubbed[i] = model.Scrub("v1")
	}

	return scrubbed, nil
}
