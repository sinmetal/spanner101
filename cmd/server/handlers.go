package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/sinmetal/spanner101/data"
	stores1 "github.com/sinmetal/spanner101/pattern1/stores"
	stores2 "github.com/sinmetal/spanner101/pattern2/stores"
	stores3 "github.com/sinmetal/spanner101/pattern3/stores"
)

type Handlers struct {
	OrdersStore1 *stores1.OrdersStore
	OrdersStore2 *stores2.OrdersStore
	OrdersStore3 *stores3.OrdersStore
}

func (h *Handlers) Insert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := data.RandomUserID()
	orderID := fmt.Sprintf("o%s", time.Now().Format("2006-0102-150405"))
	detailCount := rand.Intn(10) + 1

	var details1 []*stores1.OrderDetail
	var details2 []*stores2.OrderDetail
	var details3 []*stores3.OrderDetail
	for i := 0; i < detailCount; i++ {
		item := data.RandomItem()
		quantity := rand.Int63n(1000) + 1
		details1 = append(details1, &stores1.OrderDetail{
			OrderID:       orderID,
			OrderDetailID: int64(i + 1),
			ItemID:        item.ItemID,
			Price:         item.Price,
			Quantity:      quantity,
			CommitedAt:    spanner.CommitTimestamp,
		})
		details2 = append(details2, &stores2.OrderDetail{
			OrderID:       orderID,
			OrderDetailID: int64(i + 1),
			ItemID:        item.ItemID,
			Price:         item.Price,
			Quantity:      quantity,
			CommitedAt:    spanner.CommitTimestamp,
		})
		details3 = append(details3, &stores3.OrderDetail{
			OrderID:       orderID,
			OrderDetailID: int64(i + 1),
			ItemID:        item.ItemID,
			Price:         item.Price,
			Quantity:      quantity,
			CommitedAt:    spanner.CommitTimestamp,
		})
	}
	resultCh := make(chan string)
	go func() {
		_, err := h.OrdersStore1.Insert(ctx, userID, orderID, details1)
		if err != nil {
			msg := fmt.Sprintf("failed OrdersStore1.Insert() err=%s", err)
			log.Println(msg)
			resultCh <- msg
		}
		resultCh <- fmt.Sprintf("done OrdersStore1.Insert() OrderID=%s", orderID)
	}()
	go func() {
		_, err := h.OrdersStore2.Insert(ctx, userID, orderID, details2)
		if err != nil {
			msg := fmt.Sprintf("failed OrdersStore2.Insert() err=%s", err)
			log.Println(msg)
			resultCh <- msg
		}
		resultCh <- fmt.Sprintf("done OrdersStore2.Insert() OrderID=%s", orderID)
	}()
	go func() {
		_, err := h.OrdersStore3.Insert(ctx, userID, orderID, details3)
		if err != nil {
			msg := fmt.Sprintf("failed OrdersStore3.Insert() err=%s", err)
			log.Println(msg)
			resultCh <- msg
		}
		resultCh <- fmt.Sprintf("done OrdersStore3.Insert() OrderID=%s", orderID)
	}()
	var results []string
	for i := 0; i < 3; i++ {
		ret := <-resultCh
		results = append(results, ret)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Println(err)
	}
}
