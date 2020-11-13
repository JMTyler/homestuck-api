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
	if err := fcm.Subscribe(token); err != nil {
		// TODO: Gotta start using log.Fatal() and its ilk.
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *Service) Unsubscribe(token string) error {
	if err := fcm.Unsubscribe(token); err != nil {
		// TODO: Gotta start using log.Fatal() and its ilk.
		fmt.Println(err)
		return err
	}

	return nil
}

func (_ *Service) GetStories() ([]*db.Story, error) {
	return new(db.Story).FindAll(), nil
}
