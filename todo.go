package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
	"sync"
	"time"
)

type Lists struct {
	Items  map[string]*List
	locker sync.Mutex

	// f is a cache for the loadpath; it will not
	// persist, but is useful for saving to disk
	// (so *Lists.Save() knows where to put files)
	f string
}

func LoadLists(f string) (l *Lists, err error) {
	l = new(Lists)
	l.f = f

	file, err := os.ReadFile(f)
	if err != nil {
		if os.IsNotExist(err) {
			l.Items = make(map[string]*List)
			err = nil
		}
		return
	}

	b := bytes.NewBuffer(file)
	dec := gob.NewDecoder(b)

	err = dec.Decode(l)

	return
}

func (l *Lists) Save() (err error) {
	defer func() {
		if err != nil {
			log.Print(err.Error())
		}
	}()

	b := bytes.Buffer{}
	dec := gob.NewEncoder(&b)

	err = dec.Encode(l)
	if err != nil {
		return
	}

	return os.WriteFile(l.f, b.Bytes(), 0600)
}

type List struct {
	Items []*Item
}

type Item struct {
	ID         int
	Done       bool
	Title      string
	CreatedAt  time.Time
	MarkedDone time.Time
}

func NewList() *List {
	return &List{
		Items: make([]*Item, 0),
	}
}

func (l *List) Create(title string) (i *Item) {
	i = &Item{
		ID:        len(l.Items),
		Done:      false,
		Title:     title,
		CreatedAt: time.Now(),
	}

	l.Items = append(l.Items, i)

	return
}

func (l *List) Read(id int) (i *Item) {
	if id >= len(l.Items) {
		return
	}

	return l.Items[id]
}

func (l *List) Update(id int, title string) {
	if id >= len(l.Items) {
		return
	}

	l.Items[id].Title = title
}

func (l *List) Finish(id int) {
	if id >= len(l.Items) {
		return
	}

	l.Items[id].Done = true
	l.Items[id].MarkedDone = time.Now()
}

func (l *List) Delete(id int) {
	if id >= len(l.Items) {
		return
	}

	if len(l.Items) < 2 {
		l.Items = make([]*Item, 0)

		return
	}

	// This method will get slower as more todo items exist
	// we should revist this then (or add a delete-range option
	// to the bot)
	l.Items = append(l.Items[:id], l.Items[id+1:]...)

	for idx, item := range l.Items {
		item.ID = idx
	}
}
