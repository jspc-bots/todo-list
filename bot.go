package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jspc/bottom"
	"github.com/lrstanley/girc"
)

var (
	tfmt = "2006. 01. 02. 15:04"
)

type Bot struct {
	bottom bottom.Bottom
	Lists  *Lists
	tz     *time.Location
}

type handlerFunc func(groups [][]byte) error

func New(user, password, server string, verify bool, f, timezone string) (b Bot, err error) {
	b.tz, err = time.LoadLocation(timezone)
	if err != nil {
		return
	}

	b.Lists, err = LoadLists(f)
	if err != nil {
		return
	}

	b.bottom, err = bottom.New(user, password, server, verify)
	if err != nil {
		return
	}

	b.bottom.Client.Handlers.Add(girc.CONNECTED, func(c *girc.Client, e girc.Event) {
		c.Cmd.Join(Chan)
	})

	router := bottom.NewRouter()
	router.AddRoute(`todo\:\s+\"(.*)\"`, b.add)
	router.AddRoute(`edit\s+todo\s+(\d+)\:\s+\"(.*)\"`, b.edit)
	router.AddRoute(`mark\s+todo\s+(\d+)`, b.mark)
	router.AddRoute(`delete\s+todo\s+(\d+)`, b.delete)
	router.AddRoute(`show\s+todo\s+list`, b.show)

	b.bottom.Middlewares.Push(router)

	return
}

func (b Bot) add(u, channel string, groups []string) (err error) {
	b.Lists.locker.Lock()
	defer b.unlockAndSave()

	list, ok := b.Lists.Items[channel]
	if !ok {
		list = NewList()
		b.Lists.Items[channel] = list
	}

	list.Create(groups[1])

	b.bottom.Client.Cmd.Messagef(u, "added %q to %s todo list", groups[1], channel)

	return
}

func (b Bot) edit(_, channel string, groups []string) (err error) {
	b.Lists.locker.Lock()
	defer b.unlockAndSave()

	l, i, err := b.getListAndId(channel, groups[1])
	if err != nil {
		return
	}

	l.Update(i, groups[2])

	return
}

func (b Bot) mark(_, channel string, groups []string) (err error) {
	b.Lists.locker.Lock()
	defer b.unlockAndSave()

	l, i, err := b.getListAndId(channel, groups[1])
	if err != nil {
		return
	}

	l.Finish(i)

	return
}

func (b Bot) delete(_, channel string, groups []string) (err error) {
	b.Lists.locker.Lock()
	defer b.unlockAndSave()

	l, i, err := b.getListAndId(channel, groups[1])
	if err != nil {
		return
	}

	l.Delete(i)

	return
}

func (b Bot) show(_, channel string, groups []string) (err error) {
	l, _, err := b.getListAndId(channel, "-1")
	if err != nil {
		return
	}

	b.bottom.Client.Cmd.Messagef(channel, "%s todo list", channel)
	for _, item := range l.Items {
		if item.Done {
			b.bottom.Client.Cmd.Messagef(channel, "%d    ☑️  %s    (%s)", item.ID, item.Title, item.MarkedDone.In(b.tz).Format(tfmt))
		} else {
			b.bottom.Client.Cmd.Messagef(channel, "%d    🚫 %s    (%s)", item.ID, item.Title, item.CreatedAt.In(b.tz).Format(tfmt))
		}
	}

	return
}

func (b Bot) unlockAndSave() (err error) {
	// Unlock always, even if saving fails
	// - it's better to have a todo list that doesn't save, than have everything
	//   fail totally
	defer b.Lists.locker.Unlock()

	return b.Lists.Save()
}

func (b Bot) getListAndId(channel, id string) (l *List, idInt int, err error) {
	l, ok := b.Lists.Items[channel]
	if !ok {
		err = fmt.Errorf("there is no todo list registered for %q, add a todo item to create it", channel)

		return
	}

	idInt, err = strconv.Atoi(id)

	return
}
