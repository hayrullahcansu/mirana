package utils

import (
	"math/rand"
	"time"

	"bitbucket.org/digitdreamteam/mirana/core/mdl"
	"bitbucket.org/digitdreamteam/mirana/utils/que"
)

func GetSingleDeckPack() *que.Queue {
	queue := que.Init()
	// var a = make([]interface{}, len(mdl.CardValues)*len(mdl.CardTypes)) // or slice := make([]int, elems)

	var a []*mdl.Card
	// var indexer = 0
	for _, cardValue := range mdl.CardValues {
		for _, cardType := range mdl.CardTypes {
			c := mdl.NewCardData(cardType, cardValue)
			a = append(a, c)
			// a[indexer] =
		}
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	for _, v := range a {
		queue.Enqueue(v)
	}
	return queue

}

func GetAmericanPack() *que.Queue {
	queue := que.Init()
	// var a = make([]interface{}, len(mdl.CardValues)*len(mdl.CardTypes)) // or slice := make([]int, elems)

	var a []*mdl.Card
	// var indexer = 0
	for _, cardValue := range mdl.CardValues {
		for _, cardType := range mdl.CardTypes {
			c := mdl.NewCardData(cardType, cardValue)
			a = append(a, c)
			// a[indexer] =
		}
	}
	for _, cardValue := range mdl.CardValues {
		for _, cardType := range mdl.CardTypes {
			c := mdl.NewCardData(cardType, cardValue)
			a = append(a, c)
			// a[indexer] =
		}
	}
	for _, cardValue := range mdl.CardValues {
		for _, cardType := range mdl.CardTypes {
			c := mdl.NewCardData(cardType, cardValue)
			a = append(a, c)
			// a[indexer] =
		}
	}
	for _, cardValue := range mdl.CardValues {
		for _, cardType := range mdl.CardTypes {
			c := mdl.NewCardData(cardType, cardValue)
			a = append(a, c)
			// a[indexer] =
		}
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	for _, v := range a {
		queue.Enqueue(v)
	}
	return queue

}
