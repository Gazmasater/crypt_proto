package mexc

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"crypt_proto/domain"
	pb "crypt_proto/pb"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type Feed struct {
	debug bool
}

func NewFeed(debug bool) *Feed {
	return &Feed{debug: debug}
}

func (f *Feed) Name() string { return "MEXC" }

/* ===== proto decoder ===== */

var wrapperPool = sync.Pool{
	New: func() any { return new(pb.PushDataV3ApiWrapper) },
}

func parsePBQuote(raw []byte) (string, domain.Quote, bool) {
	w, _ := wrapperPool.Get().(*pb.PushDataV3ApiWrapper)
	defer func() {
		*w = pb.PushDataV3ApiWrapper{}
		wrapperPool.Put(w)
	}()

	if err := proto.Unmarshal(raw, w); err != nil {
		return "", domain.Quote{}, false
	}

	sym := w.GetSymbol()
	if sym == "" {
		ch := w.GetChannel()
		if i := strings.LastIndex(ch, "@"); i >= 0 && i+1 < len(ch) {
			sym = ch[i+1:]
		}
	}
	if sym == "" {
		return "", domain.Quote{}, false
	}

	if b1, ok := w.GetBody().(*pb.PushDataV3ApiWrapper_PublicBookTicker); ok && b1.PublicBookTicker != nil {
		t := b1.PublicBookTicker
		bp := t.GetBidPrice()
		ap := t.GetAskPrice()
		if bp == "" || ap == "" {
			return "", domain.Quote{}, false
		}
		bid, err1 := strconv.ParseFloat(bp, 64)
		ask, err2 := strconv.ParseFloat(ap, 64)
		if err1 != nil || err2 != nil || bid <= 0 || ask <= 0 {
			return "", domain.Quote{}, false
		}
		return sym, domain.Quote{
			Bid:    bid,
			Ask:    ask,
			BidQty: 0,
			AskQty: 0,
		}, true
	}

	if b2, ok := w.GetBody().(*pb.PushDataV3ApiWrapper_PublicAggreBookTicker); ok && b2.PublicAggreBookTicker != nil {
		t := b2.PublicAggreBookTicker

		bp := t.GetBidPrice()
		ap := t.GetAskPrice()
		bq := t.GetBidQuantity()
		aq := t.GetAskQuantity()

		if bp == "" || ap == "" {
			return "", domain.Quote{}, false
		}
		bid, err1 := strconv.ParseFloat(bp, 64)
		ask, err2 := strconv.ParseFloat(ap, 64)
		if err1 != nil || err2 != nil || bid <= 0 || ask <= 0 {
			return "", domain.Quote{}, false
		}

		var bidQty, askQty float64
		if bq != "" {
			if v, err := strconv.ParseFloat(bq, 64); err == nil {
				bidQty = v
			}
		}
		if aq != "" {
			if v, err := strconv.ParseFloat(aq, 64); err == nil {
				askQty = v
			}
		}

		return sym, domain.Quote{
			Bid:    bid,
			Ask:    ask,
			BidQty: bidQty,
			AskQty: askQty,
		}, true
	}

	return "", domain.Quote{}, false
}

func (f *Feed) dlog(format string, args ...any) {
	if f.debug {
		log.Printf(format, args...)
	}
}

/* ===== WS ===== */

func (f *Feed) runPublicBookTickerWS(
	ctx context.Context,
	wg *sync.WaitGroup,
	connID int,
	symbols []string,
	interval string,
	out chan<- domain.Event,
) {
	defer wg.Done()

	const (
		baseRetry = 2 * time.Second
		maxRetry  = 30 * time.Second
	)

	urlWS := "wss://wbs-api.mexc.com/ws"

	topics := make([]string, 0, len(symbols))
	for _, s := range symbols {
		topics = append(topics, "spot@public.aggre.bookTicker.v3.api.pb@"+interval+"@"+s)
	}

	retry := baseRetry

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		conn, _, err := websocket.DefaultDialer.Dial(urlWS, nil)
		if err != nil {
			log.Printf("[MEXC WS #%d] dial err: %v (retry in %v)", connID, err, retry)
			time.Sleep(retry)
			if retry < maxRetry {
				retry *= 2
				if retry > maxRetry {
					retry = maxRetry
				}
			}
			continue
		}
		log.Printf("[MEXC WS #%d] connected to %s (symbols: %d)", connID, urlWS, len(symbols))
		retry = baseRetry

		_ = conn.SetReadDeadline(time.Now().Add(90 * time.Second))

		var lastPing time.Time
		conn.SetPongHandler(func(appData string) error {
			rtt := time.Since(lastPing)
			f.dlog("[MEXC WS #%d] Pong через %v", connID, rtt)
			return conn.SetReadDeadline(time.Now().Add(90 * time.Second))
		})

		stopPing := make(chan struct{})
		go func() {
			t := time.NewTicker(45 * time.Second)
			defer t.Stop()
			for {
				select {
				case <-t.C:
					lastPing = time.Now()
					if err := conn.WriteControl(websocket.PingMessage, []byte("hb"), time.Now().Add(5*time.Second)); err != nil {
						f.dlog("[MEXC WS #%d] ping error: %v", connID, err)
						return
					}
				case <-stopPing:
					return
				}
			}
		}()

		sub := map[string]any{
			"method": "SUBSCRIPTION",
			"params": topics,
			"id":     time.Now().Unix(),
		}
		if err := conn.WriteJSON(sub); err != nil {
			log.Printf("[MEXC WS #%d] subscribe send err: %v", connID, err)
			close(stopPing)
			_ = conn.Close()
			time.Sleep(retry)
			continue
		}
		log.Printf("[MEXC WS #%d] SUB -> %d topics", connID, len(topics))

		for {
			mt, raw, err := conn.ReadMessage()
			if err != nil {
				log.Printf("[MEXC WS #%d] read err: %v (reconnect)", connID, err)
				break
			}

			switch mt {
			case websocket.TextMessage:
				if f.debug {
					var tmp any
					if err := json.Unmarshal(raw, &tmp); err == nil {
						j, _ := json.Marshal(tmp)
						f.dlog("[MEXC #%d TEXT] %s", connID, string(j))
					} else {
						f.dlog("[MEXC #%d TEXT RAW] %s", connID, string(raw))
					}
				}
			case websocket.BinaryMessage:
				sym, q, ok := parsePBQuote(raw)
				if !ok {
					continue
				}
				ev := domain.Event{
					Symbol: sym,
					Bid:    q.Bid,
					Ask:    q.Ask,
					BidQty: q.BidQty,
					AskQty: q.AskQty,
				}
				select {
				case out <- ev:
				case <-ctx.Done():
					close(stopPing)
					_ = conn.Close()
					return
				}
			default:
			}
		}

		close(stopPing)
		_ = conn.Close()
		time.Sleep(retry)
		if retry < maxRetry {
			retry *= 2
			if retry > maxRetry {
				retry = maxRetry
			}
		}
	}
}

// Start реализует интерфейс MarketDataFeed.
func (f *Feed) Start(
	ctx context.Context,
	wg *sync.WaitGroup,
	symbols []string,
	interval string,
	out chan<- domain.Event,
) {
	const maxPerConn = 25
	chunks := make([][]string, 0)
	for i := 0; i < len(symbols); i += maxPerConn {
		j := i + maxPerConn
		if j > len(symbols) {
			j = len(symbols)
		}
		chunks = append(chunks, symbols[i:j])
	}
	log.Printf("[MEXC] будем использовать %d WS-подключений", len(chunks))

	for idx, chunk := range chunks {
		wg.Add(1)
		go f.runPublicBookTickerWS(ctx, wg, idx, chunk, interval, out)
	}
}
