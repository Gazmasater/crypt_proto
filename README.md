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




2026/01/07 23:33:27 [KuCoin WS 0] connected
2026/01/07 23:33:27 [KuCoin WS DEBUG] {"id":"1767818006746572747","type":"welcome"}
2026/01/07 23:33:29 [KuCoin WS DEBUG] {"id":"1767818008993833749","type":"ack"}
2026/01/07 23:33:29 [KuCoin WS DEBUG] {"topic":"/market/ticker:TRX-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.000003273","bestAskSize":"152.35","bestBid":"0.000003271","bestBidSize":"1009.52","price":"0.000003273","sequence":"869598256","size":"6.47","time":1767817813625}}
2026/01/07 23:33:29 [KuCoin WS 1] connected
2026/01/07 23:33:29 [KuCoin WS DEBUG] {"id":"1767818008282024385","type":"welcome"}
2026/01/07 23:33:29 [KuCoin WS DEBUG] {"id":"1767818009114515888","type":"ack"}
2026/01/07 23:33:29 [KuCoin WS DEBUG] {"topic":"/market/ticker:ZEC-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.0052463","bestAskSize":"0.1205","bestBid":"0.0052393","bestBidSize":"0.1205","price":"0.0052333","sequence":"992543244","size":"0.4924","time":1767817883813}}
2026/01/07 23:33:29 [KuCoin WS DEBUG] {"id":"1767818009233955960","type":"ack"}
2026/01/07 23:33:29 [KuCoin WS DEBUG] {"topic":"/market/ticker:ZEC-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.0052463","bestAskSize":"0.1205","bestBid":"0.0052393","bestBidSize":"0.1205","price":"0.0052333","sequence":"992543248","size":"0.4924","time":1767817883813}}
2026/01/07 23:33:29 [KuCoin WS DEBUG] {"topic":"/market/ticker:HBAR-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.000001349","bestAskSize":"1210.0482","bestBid":"0.000001347","bestBidSize":"2040.158","price":"0.000001351","sequence":"835905446","size":"4.4582","time":1767815825464}}
2026/01/07 23:33:29 [KuCoin WS DEBUG] {"id":"1767818009353628587","type":"ack"}
2026/01/07 23:33:29 [KuCoin WS DEBUG] {"topic":"/market/ticker:ZEC-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.0052463","bestAskSize":"0.1205","bestBid":"0.0052387","bestBidSize":"0.1205","price":"0.0052333","sequence":"992543250","size":"0.4924","time":1767817883813}}
2026/01/07 23:33:29 [KuCoin WS DEBUG] {"id":"1767818009473514254","type":"ack"}
2026/01/07 23:33:29 [KuCoin WS DEBUG] {"topic":"/market/ticker:ZEC-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.0052463","bestAskSize":"0.1205","bestBid":"0.0052393","bestBidSize":"0.1205","price":"0.0052333","sequence":"992543252","size":"0.4924","time":1767817883813}}
2026/01/07 23:33:29 [KuCoin WS DEBUG] {"id":"1767818009593942145","type":"ack"}
2026/01/07 23:33:29 [KuCoin WS DEBUG] {"topic":"/market/ticker:ZEC-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.0052463","bestAskSize":"0.1205","bestBid":"0.0052393","bestBidSize":"0.1205","price":"0.0052333","sequence":"992543256","size":"0.4924","time":1767817883813}}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"id":"1767818009714143821","type":"ack"}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"topic":"/market/ticker:ZEC-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.0052463","bestAskSize":"0.1205","bestBid":"0.0052393","bestBidSize":"0.1205","price":"0.0052333","sequence":"992543257","size":"0.4924","time":1767817883813}}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"topic":"/market/ticker:LYX-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.4697","bestAskSize":"37.7171","bestBid":"0.4656","bestBidSize":"4.6088","price":"0.4676","sequence":"1393444469","size":"13.8885","time":1767817991972}}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"id":"1767818009834310922","type":"ack"}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"topic":"/market/ticker:TEL-ETH","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.0000012499","bestAskSize":"178753.9","bestBid":"0.0000012341","bestBidSize":"2160","price":"0.0000012419","sequence":"1189087643","size":"1836.33","time":1767816023876}}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"topic":"/market/ticker:ZEC-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.0052463","bestAskSize":"0.1205","bestBid":"0.0052393","bestBidSize":"0.1205","price":"0.0052333","sequence":"992543260","size":"0.4924","time":1767817883813}}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"id":"1767818009953847754","type":"ack"}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"topic":"/market/ticker:TEL-ETH","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.0000012499","bestAskSize":"20000","bestBid":"0.0000012342","bestBidSize":"20000","price":"0.0000012419","sequence":"1189087651","size":"1836.33","time":1767816023876}}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"topic":"/market/ticker:HBAR-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.000001349","bestAskSize":"1210.0482","bestBid":"0.000001347","bestBidSize":"2040.158","price":"0.000001351","sequence":"835905447","size":"4.4582","time":1767815825464}}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"id":"1767818010075941190","type":"ack"}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"topic":"/market/ticker:LYX-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.4697","bestAskSize":"37.7171","bestBid":"0.4655","bestBidSize":"128.017","price":"0.4676","sequence":"1393444473","size":"13.8885","time":1767817991972}}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"topic":"/market/ticker:ZEC-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.0052463","bestAskSize":"0.1205","bestBid":"0.0052393","bestBidSize":"0.1205","price":"0.0052333","sequence":"992543261","size":"0.4924","time":1767817883813}}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"topic":"/market/ticker:HBAR-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.000001349","bestAskSize":"1210.0482","bestBid":"0.000001347","bestBidSize":"2040.158","price":"0.000001351","sequence":"835905448","size":"4.4582","time":1767815825464}}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"topic":"/market/ticker:SOL-KCS","type":"message","subject":"trade.ticker","data":{"bestAsk":"11.903","bestAskSize":"0.607","bestBid":"11.885","bestBidSize":"0.607","price":"11.897","sequence":"377972183","size":"0.099","time":1767816472441}}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"topic":"/market/ticker:TEL-ETH","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.00000125","bestAskSize":"401385.17","bestBid":"0.0000012342","bestBidSize":"20000","price":"0.0000012419","sequence":"1189087654","size":"1836.33","time":1767816023876}}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"id":"1767818010193790546","type":"ack"}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"topic":"/market/ticker:SOL-KCS","type":"message","subject":"trade.ticker","data":{"bestAsk":"11.903","bestAskSize":"0.607","bestBid":"11.885","bestBidSize":"0.607","price":"11.897","sequence":"377972184","size":"0.099","time":1767816472441}}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"id":"1767818010313830845","type":"ack"}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"topic":"/market/ticker:HBAR-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.000001349","bestAskSize":"1210.0482","bestBid":"0.000001347","bestBidSize":"3282.2985","price":"0.000001351","sequence":"835905449","size":"4.4582","time":1767815825464}}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"id":"1767818010434114059","type":"ack"}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"topic":"/market/ticker:ZEC-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.0052463","bestAskSize":"0.1205","bestBid":"0.0052393","bestBidSize":"0.1205","price":"0.0052333","sequence":"992543262","size":"0.4924","time":1767817883813}}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"id":"1767818010473225408","type":"ack"}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"topic":"/market/ticker:HBAR-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.000001349","bestAskSize":"2024.9018","bestBid":"0.000001347","bestBidSize":"2040.158","price":"0.000001351","sequence":"835905464","size":"4.4582","time":1767815825464}}
2026/01/07 23:33:30 [KuCoin WS 2] connected
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"id":"1767818009753505587","type":"welcome"}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"id":"1767818010554364673","type":"ack"}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"id":"1767818010597674956","type":"ack"}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"topic":"/market/ticker:LYX-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.4697","bestAskSize":"37.7171","bestBid":"0.4655","bestBidSize":"128.017","price":"0.4676","sequence":"1393444474","size":"13.8885","time":1767817991972}}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"topic":"/market/ticker:ZEC-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.0052463","bestAskSize":"0.1205","bestBid":"0.0052393","bestBidSize":"0.1205","price":"0.0052333","sequence":"992543263","size":"0.4924","time":1767817883813}}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"id":"1767818010674388180","type":"ack"}
2026/01/07 23:33:30 [KuCoin WS DEBUG] {"id":"1767818010713583751","type":"ack"}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:LYX-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.4697","bestAskSize":"37.7171","bestBid":"0.4655","bestBidSize":"128.017","price":"0.4676","sequence":"1393444475","size":"13.8885","time":1767817991972}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"id":"1767818010794177549","type":"ack"}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"id":"1767818010833735275","type":"ack"}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:TEL-ETH","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.00000125","bestAskSize":"401385.17","bestBid":"0.0000012341","bestBidSize":"2160","price":"0.0000012419","sequence":"1189087655","size":"1836.33","time":1767816023876}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"id":"1767818010914378375","type":"ack"}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"id":"1767818010953301757","type":"ack"}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:NEAR-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.00001893","bestAskSize":"825.0669","bestBid":"0.0000189","bestBidSize":"88.1937","price":"0.00001891","sequence":"1236703432","size":"6.3189","time":1767817813623}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"id":"1767818011034234018","type":"ack"}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:ZEC-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.0052463","bestAskSize":"0.1205","bestBid":"0.0052393","bestBidSize":"0.1205","price":"0.0052333","sequence":"992543274","size":"0.4924","time":1767817883813}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:TEL-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.00000004294","bestAskSize":"4838","bestBid":"0.00000004244","bestBidSize":"26519","price":"0.0000000427","sequence":"2567521058","size":"875.99","time":1767817712459}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"id":"1767818011073419410","type":"ack"}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:TEL-ETH","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.00000125","bestAskSize":"401385.17","bestBid":"0.0000012343","bestBidSize":"810.11","price":"0.0000012419","sequence":"1189087660","size":"1836.33","time":1767816023876}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:SNX-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.00000544","bestAskSize":"148.5","bestBid":"0.00000539","bestBidSize":"148.5","price":"0.00000542","sequence":"411308693","size":"47.73","time":1767800340004}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:ZEC-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.0052463","bestAskSize":"0.1205","bestBid":"0.0052399","bestBidSize":"0.1205","price":"0.0052333","sequence":"992543280","size":"0.4924","time":1767817883813}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"id":"1767818011154334030","type":"ack"}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:HBAR-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.000001349","bestAskSize":"1413.7616","bestBid":"0.000001347","bestBidSize":"2040.158","price":"0.000001351","sequence":"835905467","size":"4.4582","time":1767815825464}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"id":"1767818011193510025","type":"ack"}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:KAS-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.0489","bestAskSize":"2204.9664","bestBid":"0.04889","bestBidSize":"2144.8104","price":"0.04889","sequence":"2728312172","size":"85.9848","time":1767818006295}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:TEL-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.00000004294","bestAskSize":"4838","bestBid":"0.00000004244","bestBidSize":"26519","price":"0.0000000427","sequence":"2567521060","size":"875.99","time":1767817712459}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:NEAR-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.00001893","bestAskSize":"825.0669","bestBid":"0.0000189","bestBidSize":"88.1937","price":"0.00001891","sequence":"1236703433","size":"6.3189","time":1767817813623}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:ELA-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"1.0863","bestAskSize":"61.4","bestBid":"1.0815","bestBidSize":"66.1434","price":"1.0853","sequence":"2646470460","size":"19.9971","time":1767817870149}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:KAS-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.000000539","bestAskSize":"9672.263","bestBid":"0.000000536","bestBidSize":"2244.3487","price":"0.000000533","sequence":"276237234","size":"19.0074","time":1767815823727}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"id":"1767818011273446559","type":"ack"}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"id":"1767818011313274640","type":"ack"}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:WIN-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.00002896","bestAskSize":"2416000","bestBid":"0.00002887","bestBidSize":"14639609.9999","price":"0.0000288","sequence":"878711992","size":"20000","time":1767817586369}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:TEL-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.00000004294","bestAskSize":"4838","bestBid":"0.00000004245","bestBidSize":"28679","price":"0.0000000427","sequence":"2567521067","size":"875.99","time":1767817712459}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:KAS-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.0489","bestAskSize":"2204.9664","bestBid":"0.04889","bestBidSize":"2144.8104","price":"0.04889","sequence":"2728312174","size":"85.9848","time":1767818006295}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:NEAR-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.00001893","bestAskSize":"825.0669","bestBid":"0.0000189","bestBidSize":"88.1937","price":"0.00001891","sequence":"1236703435","size":"6.3189","time":1767817813623}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:AAVE-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"166.226","bestAskSize":"0.1621","bestBid":"166.22","bestBidSize":"0.2867","price":"166.242","sequence":"8401242937","size":"0.603","time":1767818008647}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:ELA-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"1.0863","bestAskSize":"61.4","bestBid":"1.0816","bestBidSize":"43.98","price":"1.0853","sequence":"2646470461","size":"19.9971","time":1767817870149}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:TRX-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.000003273","bestAskSize":"152.35","bestBid":"0.000003271","bestBidSize":"1532.64","price":"0.000003273","sequence":"869598258","size":"6.47","time":1767817813625}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:HBAR-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.000001349","bestAskSize":"1413.7616","bestBid":"0.000001347","bestBidSize":"2040.158","price":"0.000001351","sequence":"835905477","size":"4.4582","time":1767815825464}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"id":"1767818011393899295","type":"ack"}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:NEAR-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.00001893","bestAskSize":"642.8169","bestBid":"0.0000189","bestBidSize":"73.6536","price":"0.00001891","sequence":"1236703438","size":"6.3189","time":1767817813623}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:WIN-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.00002896","bestAskSize":"690000","bestBid":"0.0000289","bestBidSize":"1730000","price":"0.0000288","sequence":"878711999","size":"20000","time":1767817586369}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"id":"1767818011433805552","type":"ack"}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:ELA-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"1.0863","bestAskSize":"61.4","bestBid":"1.0812","bestBidSize":"24.7","price":"1.0853","sequence":"2646470467","size":"19.9971","time":1767817870149}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:AAVE-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"166.226","bestAskSize":"0.8146","bestBid":"166.22","bestBidSize":"0.2867","price":"166.242","sequence":"8401242940","size":"0.603","time":1767818008647}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:A-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.178","bestAskSize":"10.61","bestBid":"0.1777","bestBidSize":"2899.16","price":"0.1789","sequence":"113775489","size":"6.79","time":1767815435580}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:HBAR-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.000001349","bestAskSize":"1413.7616","bestBid":"0.000001347","bestBidSize":"2040.158","price":"0.000001351","sequence":"835905480","size":"4.4582","time":1767815825464}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"id":"1767818011513725097","type":"ack"}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:TEL-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.00000004294","bestAskSize":"4838","bestBid":"0.00000004245","bestBidSize":"28679","price":"0.0000000427","sequence":"2567521069","size":"875.99","time":1767817712459}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:WIN-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.00002896","bestAskSize":"690000","bestBid":"0.0000289","bestBidSize":"1730000","price":"0.0000288","sequence":"878712002","size":"20000","time":1767817586369}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:NEAR-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.00001893","bestAskSize":"825.0669","bestBid":"0.0000189","bestBidSize":"15.4932","price":"0.00001891","sequence":"1236703443","size":"6.3189","time":1767817813623}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:KAS-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.0489","bestAskSize":"2204.9664","bestBid":"0.04889","bestBidSize":"2144.8104","price":"0.04889","sequence":"2728312176","size":"85.9848","time":1767818006295}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"id":"1767818011553969669","type":"ack"}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:ELA-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"1.0863","bestAskSize":"61.4","bestBid":"1.0812","bestBidSize":"24.7","price":"1.0853","sequence":"2646470478","size":"19.9971","time":1767817870149}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:AAVE-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"166.225","bestAskSize":"0.6525","bestBid":"166.22","bestBidSize":"0.2867","price":"166.242","sequence":"8401242945","size":"0.603","time":1767818008647}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:HBAR-BTC","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.000001349","bestAskSize":"1210.0482","bestBid":"0.000001347","bestBidSize":"2040.158","price":"0.000001351","sequence":"835905484","size":"4.4582","time":1767815825464}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:LYX-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.4697","bestAskSize":"37.7171","bestBid":"0.4655","bestBidSize":"128.017","price":"0.4676","sequence":"1393444477","size":"13.8885","time":1767817991972}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:WIN-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"0.00002896","bestAskSize":"690000","bestBid":"0.0000289","bestBidSize":"1730000","price":"0.0000288","sequence":"878712005","size":"20000","time":1767817586369}}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"id":"1767818011634233687","type":"ack"}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"id":"1767818011673415630","type":"ack"}
2026/01/07 23:33:31 [KuCoin WS DEBUG] {"topic":"/market/ticker:ELA-USDT","type":"message","subject":"trade.ticker","data":{"bestAsk":"1.0863","bestAskSize":"61.4","bestBid":"1.0812","bestBidSize":"24.7","price":"1.0853","sequence":"2646470492","size":"19.9971","time":1767817870149}}











