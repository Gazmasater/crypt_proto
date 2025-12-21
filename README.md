apikey = "4333ed4b-cd83-49f5-97d1-c399e2349748"
secretkey = "E3848531135EDB4CCFDA0F1BC14CD274"
IP = ""
Название API-ключа = "Arb"
Доступы = "Чтение"



sudo systemctl mask sleep.target suspend.target hibernate.target hybrid-sleep.target



wbs-api.mexc.com/ws 


[https://edis-global.vercel.app/ru/vps-hosting/singapore-singapore
](https://sg.edisglobal.com/)



git pull --rebase origin privat
git push origin privat


BOOK_INTERVAL=100ms
SYMBOLS_FILE=triangles_markets.csv
DEBUG=false


import (
    // ...
    "net/http"
    _ "net/http/pprof"
)


   // pprof HTTP-сервер
    go func() {
        log.Println("pprof on http://localhost:6060/debug/pprof/")
        if err := http.ListenAndServe("localhost:6060", nil); err != nil {
            log.Printf("pprof server error: %v", err)
        }
    }()


	go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30


(pprof) top        # показать топ функций по CPU
(pprof) top10
(pprof) list parsePBWrapperMid   # подробный разбор одной функции
(pprof) quit


go tool pprof http://localhost:6060/debug/pprof/heap


(pprof) top
(pprof) top -cum
(pprof) list parsePBWrapperMid
(pprof) quit


models/market_data.go
package models

// MarketData хранит данные с биржи для одного инструмента
type MarketData struct {
	Exchange  string  `json:"exchange"`  // название биржи
	Symbol    string  `json:"symbol"`    // торговая пара, например BTC-USDT
	Bid       float64 `json:"bid"`       // лучшая цена покупки
	Ask       float64 `json:"ask"`       // лучшая цена продажи
	Timestamp int64   `json:"timestamp"` // метка времени в миллисекундах
}

models/signal.go
package models

// Signal представляет сигнал арбитража
type Signal struct {
	ExchangeStart string  `json:"exchange_start"` // биржа первой сделки
	ExchangeMid   string  `json:"exchange_mid"`   // биржа второй сделки
	ExchangeEnd   string  `json:"exchange_end"`   // биржа третьей сделки
	SymbolStart   string  `json:"symbol_start"`   // первая монета
	SymbolMid     string  `json:"symbol_mid"`     // вторая монета
	SymbolEnd     string  `json:"symbol_end"`     // третья монета
	ProfitPercent float64 `json:"profit_percent"` // ожидаемая прибыль в процентах
	Amount        float64 `json:"amount"`         // объём сделки
	Timestamp     int64   `json:"timestamp"`      // время создания сигнала
}



package collector

import "arb_project/models"

// Collector — интерфейс для любого коллектора биржи
type Collector interface {
	// Start запускает сбор данных и отправку в канал dataCh
	Start(dataCh chan<- models.MarketData) error
	
	// Stop останавливает сбор данных
	Stop() error
}



