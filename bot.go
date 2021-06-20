package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jspc/bottom"
	"github.com/lrstanley/girc"
	"github.com/olekukonko/tablewriter"
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
	router.AddRoute(`delete\s+todo\s+(\d+)"`, b.delete)
	router.AddRoute(`show\s+todo\s+list"`, b.show)

	b.bottom.Middlewares.Push(router)

	return
}

func (b Bot) add(_, channel string, groups []string) (err error) {
	b.Lists.Locker.Lock()
	defer b.unlockAndSave()

	list, ok := b.Lists.Items[channel]
	if !ok {
		list = NewList()
		b.Lists.Items[channel] = list
	}

	list.Create(groups[1])

	return
}

func (b Bot) edit(_, channel string, groups []string) (err error) {
	b.Lists.Locker.Lock()
	defer b.unlockAndSave()

	l, i, err := b.getListAndId(channel, groups[1])
	if err != nil {
		return
	}

	l.Update(i, groups[2])

	return
}

func (b Bot) mark(_, channel string, groups []string) (err error) {
	b.Lists.Locker.Lock()
	defer b.Lists.Locker.Unlock()

	l, i, err := b.getListAndId(channel, groups[1])
	if err != nil {
		return
	}

	l.Finish(i)

	return
}

func (b Bot) delete(_, channel string, groups []string) (err error) {
	b.Lists.Locker.Lock()
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

	sb := strings.Builder{}

	table := tablewriter.NewWriter(&sb)
	table.SetHeader([]string{"", "Item", "Date"})

	for _, item := range l.Items {
		if item.Done {
			table.Append([]string{"☑️", item.Title, item.MarkedDone.In(b.tz).String()})
		} else {
			table.Append([]string{"", item.Title, item.CreatedAt.In(b.tz).String()})
		}
	}

	table.Render()

	b.bottom.Client.Cmd.Messagef(channel, "%s todo list", channel)

	for _, line := range strings.Split(sb.String(), "\n") {
		b.bottom.Client.Cmd.Message(channel, line)
	}

	return
}

func (b Bot) unlockAndSave() (err error) {
	// Unlock always, even if saving fails
	// - it's better to have a todo list that doesn't save, than have everything
	//   fail totally
	defer b.Lists.Locker.Unlock()

	return b.Lists.Save()
}

func (b Bot) getListAndId(channel, id string) (l *List, idInt int, err error) {
	l, ok := b.Lists.Items[channel]
	if !ok {
		err = fmt.Errorf("there is no todo list registered for %q, try adding one", channel)

		return
	}

	idInt, err = strconv.Atoi(id)

	return
}
