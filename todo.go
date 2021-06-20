package main

import "time"

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

	return
}

func (l *List) Delete(id int) {
	if id >= len(l.Items) {
		return
	}

	l.Items = append(l.Items[:id], l.Items[id:]...)

	for idx, item := range l.Items {
		item.ID = idx
	}
}
