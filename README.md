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




package collector

import (
	"context"
	"log"
	"strings"
	"time"

	"crypt_proto/configs"
	"crypt_proto/pkg/models"
	pb "crypt_proto/pb"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type MEXCCollector struct {
	ctx      context.Context
	cancel   context.CancelFunc
	conn     *websocket.Conn
	symbols  []string
	lastData map[string]*models.MarketData
}

func NewMEXCCollector(symbols []string) *MEXCCollector {
	ctx, cancel := context.WithCancel(context.Background())
	return &MEXCCollector{
		ctx:      ctx,
		cancel:   cancel,
		symbols:  symbols,
		lastData: make(map[string]*models.MarketData, len(symbols)),
	}
}

// Имя биржи
func (c *MEXCCollector) Name() string {
	return "MEXC"
}

// Старт
func (c *MEXCCollector) Start(out chan<- models.MarketData) error {
	conn, _, err := websocket.DefaultDialer.Dial(configs.MEXC_WS, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	log.Println("[MEXC] connected")

	if err := c.subscribeAll(); err != nil {
		return err
	}

	go c.pingLoop()
	go c.readLoop(out)
	return nil
}

// Стоп
func (c *MEXCCollector) Stop() error {
	c.cancel()
	if c.conn != nil {
		_ = c.conn.Close()
	}
	return nil
}

// ----------------- Внутренние методы -----------------

// Подписка на все пары чанками по N
func (c *MEXCCollector) subscribeAll() error {
	chunkSize := 25
	chunks := chunkSymbols(c.symbols, chunkSize)

	for _, chunk := range chunks {
		params := make([]string, 0, len(chunk))
		for _, s := range chunk {
			params = append(params, "spot@public.aggre.bookTicker.v3.api.pb@100ms@"+s)
		}
		sub := map[string]interface{}{
			"method": "SUBSCRIPTION",
			"params": params,
		}
		if err := c.conn.WriteJSON(sub); err != nil {
			return err
		}
	}
	return nil
}

// Пинг
func (c *MEXCCollector) pingLoop() {
	t := time.NewTicker(configs.MEXC_PING_INTERVAL)
	defer t.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-t.C:
			_ = c.conn.WriteMessage(websocket.PingMessage, []byte("hb"))
		}
	}
}

// Основной цикл чтения
func (c *MEXCCollector) readLoop(out chan<- models.MarketData) {
	_ = c.conn.SetReadDeadline(time.Now().Add(configs.MEXC_READ_TIMEOUT))

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		mt, raw, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("[MEXC] read error: %v", err)
			return
		}
		_ = c.conn.SetReadDeadline(time.Now().Add(configs.MEXC_READ_TIMEOUT))

		if mt != websocket.BinaryMessage {
			continue
		}

		var wrap pb.PushDataV3ApiWrapper
		if err := proto.Unmarshal(raw, &wrap); err != nil {
			continue
		}

		if md := c.handleWrapper(&wrap); md != nil {
			out <- *md
		}
	}
}

// Преобразуем protobuf → MarketData
func (c *MEXCCollector) handleWrapper(wrap *pb.PushDataV3ApiWrapper) *models.MarketData {
	body := wrap.GetBody()
	pa, ok := body.(*pb.PushDataV3ApiWrapper_PublicAggreBookTicker)
	if !ok {
		return nil
	}

	bt := pa.PublicAggreBookTicker

	symbol := wrap.GetSymbol()
	if symbol == "" {
		ch := wrap.GetChannel()
		if ch != "" {
			parts := strings.Split(ch, "@")
			symbol = parts[len(parts)-1]
		}
	}
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return nil
	}

	// Проверка изменений
	last, exists := c.lastData[symbol]
	if exists && last.Bid == bt.GetBidPrice() && last.Ask == bt.GetAskPrice() &&
		last.BidSize == bt.GetBidQuantity() && last.AskSize == bt.GetAskQuantity() {
		return nil
	}

	md := &models.MarketData{
		Exchange:  "MEXC",
		Symbol:    symbol,
		BidStr:    bt.GetBidPrice(),
		AskStr:    bt.GetAskPrice(),
		BidSizeStr: bt.GetBidQuantity(),
		AskSizeStr: bt.GetAskQuantity(),
		Timestamp: time.Now().UnixMilli(),
	}

	c.lastData[symbol] = md
	return md
}

// ----------------- Вспомогательные функции -----------------

// Разбить слайс символов на чанки
func chunkSymbols(src []string, size int) [][]string {
	if len(src) == 0 || size <= 0 {
		return nil
	}

	var out [][]string
	for len(src) > size {
		out = append(out, src[:size])
		src = src[size:]
	}
	out = append(out, src)
	return out
}




