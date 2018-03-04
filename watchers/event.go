package watchers

import "github.com/kkentzo/tagger/utils"

type Event struct {
	Names *utils.Set
}

func NewEvent() Event {
	return Event{
		Names: utils.NewSet([]string{}),
	}
}
