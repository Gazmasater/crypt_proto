package mexc

import (
	"strconv"
	"strings"
	"sync"

	"crypt_proto/domain"
	pb "crypt_proto/pb"

	"google.golang.org/protobuf/proto"
)

/* =========================  PROTO DECODER  ========================= */

var wrapperPool = sync.Pool{
	New: func() any { return new(pb.PushDataV3ApiWrapper) },
}

func ParsePBQuote(raw []byte) (string, domain.Quote, bool) {
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

	// PublicBookTicker
	if b1, ok := w.GetBody().(*pb.PushDataV3ApiWrapper_PublicBookTicker); ok && b1.PublicBookTicker != nil {
		t := b1.PublicBookTicker
		return parseQuoteFromStrings(sym, t.GetBidPrice(), t.GetAskPrice(), "", "")
	}

	// PublicAggreBookTicker
	if b2, ok := w.GetBody().(*pb.PushDataV3ApiWrapper_PublicAggreBookTicker); ok && b2.PublicAggreBookTicker != nil {
		t := b2.PublicAggreBookTicker
		return parseQuoteFromStrings(
			sym,
			t.GetBidPrice(),
			t.GetAskPrice(),
			t.GetBidQuantity(),
			t.GetAskQuantity(),
		)
	}

	return "", domain.Quote{}, false
}

func parseQuoteFromStrings(sym, bp, ap, bq, aq string) (string, domain.Quote, bool) {
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
