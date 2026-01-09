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




gaz358@gaz358-BOD-WXX9:~/myprog/crypt_proto/cmd/arb$ go run .
2026/01/09 22:08:45 pprof on http://localhost:6060/debug/pprof/
2026/01/09 22:08:47 [KuCoin WS 0] connected
2026/01/09 22:08:48 [KuCoin WS 0] subscribed MANA-USDT
2026/01/09 22:08:48 [KuCoin WS 1] connected
2026/01/09 22:08:48 [KuCoin WS 0] subscribed WIN-TRX
2026/01/09 22:08:48 [KuCoin WS 0] subscribed DASH-BTC
2026/01/09 22:08:48 [KuCoin WS 0] subscribed SOL-KCS
2026/01/09 22:08:49 [KuCoin WS 0] subscribed SKL-USDT
2026/01/09 22:08:49 [KuCoin WS 0] subscribed XDC-BTC
2026/01/09 22:08:49 [KuCoin WS 0] subscribed SHIB-DOGE
2026/01/09 22:08:49 [KuCoin WS 0] subscribed ONT-ETH
2026/01/09 22:08:49 [KuCoin WS 0] subscribed SNX-ETH
2026/01/09 22:08:49 [KuCoin WS 0] subscribed BCHSV-ETH
2026/01/09 22:08:49 [KuCoin WS 0] subscribed ZIL-ETH
2026/01/09 22:08:49 [KuCoin WS 1] subscribed REQ-USDT
2026/01/09 22:08:49 [KuCoin WS 2] connected
2026/01/09 22:08:49 [KuCoin] started with 3 WS
2026/01/09 22:08:49 [Main] KuCoinCollector started
2026/01/09 22:08:49 [KuCoin WS 0] subscribed AVA-BTC
2026/01/09 22:08:49 [KuCoin WS 1] subscribed PEPE-KCS
2026/01/09 22:08:50 [KuCoin WS 0] subscribed LINK-USDT
2026/01/09 22:08:50 [KuCoin WS 1] subscribed KNC-BTC
2026/01/09 22:08:50 [KuCoin WS 0] subscribed BCHSV-BTC
2026/01/09 22:08:50 [KuCoin WS 1] subscribed TRVL-USDT
2026/01/09 22:08:50 [KuCoin WS 0] subscribed PAXG-USDT
2026/01/09 22:08:50 [KuCoin WS 1] subscribed VET-ETH
2026/01/09 22:08:50 [KuCoin WS 0] subscribed XRP-ETH
2026/01/09 22:08:50 [KuCoin WS 1] subscribed RSR-USDT
2026/01/09 22:08:50 [KuCoin WS 0] subscribed OM-BTC
2026/01/09 22:08:50 [KuCoin WS 1] subscribed LTC-BTC
2026/01/09 22:08:50 [KuCoin WS 0] subscribed AVA-USDT
2026/01/09 22:08:50 [KuCoin WS 1] subscribed BAX-BTC
2026/01/09 22:08:50 [KuCoin WS 0] subscribed ETH-EUR
2026/01/09 22:08:50 [KuCoin WS 1] subscribed HBAR-USDT
2026/01/09 22:08:50 [KuCoin WS 0] subscribed COTI-USDT
2026/01/09 22:08:50 [KuCoin WS 1] subscribed LTC-KCS
2026/01/09 22:08:51 [KuCoin WS 0] subscribed USDT-DAI
2026/01/09 22:08:51 [KuCoin WS 2] subscribed RLC-BTC
2026/01/09 22:08:51 [KuCoin WS 1] subscribed HYPE-KCS
2026/01/09 22:08:51 [KuCoin WS 0] subscribed CFX-BTC
2026/01/09 22:08:51 [KuCoin WS 2] subscribed XMR-ETH
2026/01/09 22:08:51 [KuCoin WS 1] subscribed ADA-KCS
2026/01/09 22:08:51 [KuCoin WS 0] subscribed ETH-BRL
2026/01/09 22:08:51 [KuCoin WS 2] subscribed WBTC-BTC
2026/01/09 22:08:51 [KuCoin WS 1] subscribed DAG-BTC
2026/01/09 22:08:51 [KuCoin WS 0] subscribed DASH-ETH
2026/01/09 22:08:51 [KuCoin WS 2] subscribed BTC-BRL
2026/01/09 22:08:51 [KuCoin WS 1] subscribed PERP-USDT
2026/01/09 22:08:51 [KuCoin WS 0] subscribed XTZ-USDT
2026/01/09 22:08:51 [KuCoin WS 2] subscribed IOTX-ETH
2026/01/09 22:08:51 [KuCoin WS 1] subscribed XMR-USDT
2026/01/09 22:08:51 [KuCoin WS 0] subscribed SCRT-USDT
2026/01/09 22:08:51 [KuCoin WS 2] subscribed IOTA-BTC
2026/01/09 22:08:51 [KuCoin WS 1] subscribed XDC-USDT
2026/01/09 22:08:51 [KuCoin WS 0] subscribed OGN-USDT
2026/01/09 22:08:51 [KuCoin WS 2] subscribed CFX-USDT
2026/01/09 22:08:51 [KuCoin WS 1] subscribed DAG-ETH
2026/01/09 22:08:51 [KuCoin WS 0] subscribed INJ-USDT
2026/01/09 22:08:51 [KuCoin WS 2] subscribed SOL-USDT
2026/01/09 22:08:51 [KuCoin WS 1] subscribed XMR-BTC
2026/01/09 22:08:51 [KuCoin WS 0] subscribed A-USDT
2026/01/09 22:08:51 [KuCoin WS 2] subscribed NEAR-BTC
2026/01/09 22:08:51 [KuCoin WS 1] subscribed KLV-TRX
2026/01/09 22:08:52 [KuCoin WS 0] subscribed BDX-USDT
2026/01/09 22:08:52 [KuCoin WS 2] subscribed PAXG-BTC
2026/01/09 22:08:52 [KuCoin WS 1] subscribed DOT-BTC
2026/01/09 22:08:52 [KuCoin WS 0] subscribed CKB-USDT
2026/01/09 22:08:52 [KuCoin WS 2] subscribed KAS-USDT
2026/01/09 22:08:52 [KuCoin WS 1] subscribed AR-USDT
2026/01/09 22:08:52 [KuCoin WS 0] subscribed TRAC-USDT
2026/01/09 22:08:52 [KuCoin WS 2] subscribed ICP-USDT
2026/01/09 22:08:52 [KuCoin WS 1] subscribed VSYS-BTC
2026/01/09 22:08:52 [KuCoin WS 0] subscribed RSR-BTC
2026/01/09 22:08:52 [KuCoin WS 2] subscribed HYPE-USDT
2026/01/09 22:08:52 [KuCoin WS 1] subscribed CRO-USDT
2026/01/09 22:08:52 [KuCoin WS 0] subscribed TRX-BTC
2026/01/09 22:08:52 [KuCoin WS 2] subscribed LTC-USDT
2026/01/09 22:08:52 [KuCoin WS 1] subscribed PERP-BTC
2026/01/09 22:08:52 [KuCoin WS 0] subscribed BDX-BTC
2026/01/09 22:08:52 [KuCoin WS 2] subscribed TRX-USDT
2026/01/09 22:08:52 [KuCoin WS 1] subscribed VET-USDT
2026/01/09 22:08:52 [KuCoin WS 0] subscribed ETC-ETH
2026/01/09 22:08:52 [KuCoin WS 2] subscribed POND-BTC
2026/01/09 22:08:52 [KuCoin WS 1] subscribed ALGO-USDT
2026/01/09 22:08:52 [KuCoin WS 0] subscribed ANKR-BTC
2026/01/09 22:08:52 [KuCoin WS 2] subscribed AAVE-USDT
2026/01/09 22:08:52 [KuCoin WS 1] subscribed NFT-TRX
2026/01/09 22:08:53 [KuCoin WS 0] subscribed DOGE-KCS
2026/01/09 22:08:53 [KuCoin WS 2] subscribed DOGE-BTC
2026/01/09 22:08:53 [KuCoin WS 1] subscribed XCN-BTC
2026/01/09 22:08:53 [KuCoin WS 0] subscribed CKB-BTC
2026/01/09 22:08:53 [KuCoin WS 2] subscribed ANKR-USDT
2026/01/09 22:08:53 [KuCoin WS 1] subscribed BNB-BTC
2026/01/09 22:08:53 [KuCoin WS 0] subscribed TRAC-ETH
2026/01/09 22:08:53 [KuCoin WS 2] subscribed ZEC-USDT
2026/01/09 22:08:53 [KuCoin WS 1] subscribed ZIL-USDT
2026/01/09 22:08:53 [KuCoin WS 0] subscribed KNC-ETH
2026/01/09 22:08:53 [KuCoin WS 2] subscribed NEO-USDT
2026/01/09 22:08:53 [KuCoin WS 1] subscribed XRP-KCS
2026/01/09 22:08:53 [KuCoin WS 0] subscribed BAX-USDT
2026/01/09 22:08:53 [KuCoin WS 2] subscribed IOST-ETH
2026/01/09 22:08:53 [KuCoin WS 1] subscribed LYX-USDT
2026/01/09 22:08:53 [KuCoin WS 0] subscribed EWT-BTC
2026/01/09 22:08:53 [KuCoin WS 2] subscribed ALGO-ETH
2026/01/09 22:08:53 [KuCoin WS 1] subscribed XYO-ETH
2026/01/09 22:08:53 [KuCoin WS 0] subscribed BTC-EUR
2026/01/09 22:08:53 [KuCoin WS 2] subscribed SXP-USDT
2026/01/09 22:08:53 [KuCoin WS 1] subscribed BCHSV-USDT
2026/01/09 22:08:53 [KuCoin WS 0] subscribed NKN-USDT
2026/01/09 22:08:53 [KuCoin WS 2] subscribed ERG-BTC
2026/01/09 22:08:53 [KuCoin WS 1] subscribed CSPR-ETH
2026/01/09 22:08:54 [KuCoin WS 0] subscribed OM-USDT
2026/01/09 22:08:54 [KuCoin WS 2] subscribed AVAX-BTC
2026/01/09 22:08:54 [KuCoin WS 1] subscribed XLM-ETH
2026/01/09 22:08:54 [KuCoin WS 0] subscribed TWT-USDT
2026/01/09 22:08:54 [KuCoin WS 2] subscribed WIN-BTC
2026/01/09 22:08:54 [KuCoin WS 1] subscribed XRP-USDT
2026/01/09 22:08:54 [KuCoin WS 0] subscribed ADA-USDT
2026/01/09 22:08:54 [KuCoin WS 2] subscribed ENJ-ETH
2026/01/09 22:08:54 [KuCoin WS 1] subscribed ERG-USDT
2026/01/09 22:08:54 [KuCoin WS 0] subscribed ONT-USDT
2026/01/09 22:08:54 [KuCoin WS 2] subscribed MOVR-ETH
2026/01/09 22:08:54 [KuCoin WS 1] subscribed BTC-DAI
2026/01/09 22:08:54 [KuCoin WS 0] subscribed XDC-ETH
2026/01/09 22:08:54 [KuCoin WS 2] subscribed NFT-USDT
2026/01/09 22:08:54 [KuCoin WS 1] subscribed XYO-USDT
2026/01/09 22:08:54 [KuCoin WS 0] subscribed WIN-USDT
2026/01/09 22:08:54 [KuCoin WS 2] subscribed ALGO-BTC
2026/01/09 22:08:54 [KuCoin WS 1] subscribed XLM-USDT
2026/01/09 22:08:54 [KuCoin WS 0] subscribed ETH-USDT
2026/01/09 22:08:54 [KuCoin WS 2] subscribed GAS-USDT
2026/01/09 22:08:54 [KuCoin WS 1] subscribed SCRT-BTC
2026/01/09 22:08:54 [KuCoin WS 0] subscribed KCS-ETH
2026/01/09 22:08:54 [KuCoin WS 2] subscribed NKN-BTC
2026/01/09 22:08:54 [KuCoin WS 1] subscribed SUPER-USDT
2026/01/09 22:08:54 [KuCoin WS 0] subscribed VRA-USDT
2026/01/09 22:08:54 [KuCoin WS 2] subscribed TRVL-BTC
2026/01/09 22:08:54 [KuCoin WS 1] subscribed SXP-BTC
2026/01/09 22:08:55 [KuCoin WS 0] subscribed EGLD-USDT
2026/01/09 22:08:55 [KuCoin WS 2] subscribed RUNE-BTC
2026/01/09 22:08:55 [KuCoin WS 1] subscribed WAN-BTC
2026/01/09 22:08:55 [KuCoin WS 0] subscribed AAVE-BTC
2026/01/09 22:08:55 [KuCoin WS 2] subscribed SUI-KCS
2026/01/09 22:08:55 [KuCoin WS 1] subscribed DOT-USDT
2026/01/09 22:08:55 [KuCoin WS 0] subscribed WBTC-USDT
2026/01/09 22:08:55 [KuCoin WS 2] subscribed CSPR-USDT
2026/01/09 22:08:55 [KuCoin WS 1] subscribed CHZ-USDT
2026/01/09 22:08:55 [KuCoin WS 0] subscribed BTC-USDT
2026/01/09 22:08:55 [KuCoin WS 2] subscribed DAG-USDT
2026/01/09 22:08:55 [KuCoin WS 1] subscribed SNX-BTC
2026/01/09 22:08:55 [KuCoin WS 0] subscribed SNX-USDT
2026/01/09 22:08:55 [KuCoin WS 2] subscribed A-BTC
2026/01/09 22:08:55 [KuCoin WS 1] subscribed AVAX-USDT
2026/01/09 22:08:55 [KuCoin WS 0] subscribed ETH-DAI
2026/01/09 22:08:55 [KuCoin WS 2] subscribed STX-BTC
2026/01/09 22:08:55 [KuCoin WS 1] subscribed ASTR-BTC
2026/01/09 22:08:55 [KuCoin WS 0] subscribed ASTR-USDT
2026/01/09 22:08:55 [KuCoin WS 2] subscribed ATOM-USDT
2026/01/09 22:08:55 [KuCoin WS 1] subscribed ONT-BTC
2026/01/09 22:08:55 [KuCoin WS 0] subscribed RLC-USDT
2026/01/09 22:08:55 [KuCoin WS 2] subscribed KCS-USDT
2026/01/09 22:08:55 [KuCoin WS 1] subscribed STORJ-USDT
2026/01/09 22:08:56 [KuCoin WS 0] subscribed ETC-BTC
2026/01/09 22:08:56 [KuCoin WS 2] subscribed USDT-BRL
2026/01/09 22:08:56 [KuCoin WS 1] subscribed XLM-BTC
2026/01/09 22:08:56 [KuCoin WS 0] subscribed KRL-USDT
2026/01/09 22:08:56 [KuCoin WS 2] subscribed KLV-USDT
2026/01/09 22:08:56 [KuCoin WS 1] subscribed XYO-BTC
2026/01/09 22:08:56 [KuCoin WS 0] subscribed ELA-BTC
2026/01/09 22:08:56 [KuCoin WS 2] subscribed EGLD-BTC
2026/01/09 22:08:56 [KuCoin WS 1] subscribed BCH-USDT
2026/01/09 22:08:56 [KuCoin WS 0] subscribed DASH-USDT
2026/01/09 22:08:56 [KuCoin WS 2] subscribed DOGE-USDT
2026/01/09 22:08:56 [KuCoin WS 1] subscribed DGB-BTC
2026/01/09 22:08:56 [KuCoin WS 0] subscribed VRA-BTC
2026/01/09 22:08:56 [KuCoin WS 2] subscribed ETH-BTC
2026/01/09 22:08:56 [KuCoin WS 1] subscribed IOST-USDT
2026/01/09 22:08:56 [KuCoin WS 0] subscribed LINK-BTC
2026/01/09 22:08:56 [KuCoin WS 2] subscribed TRAC-BTC
2026/01/09 22:08:56 [KuCoin WS 1] subscribed KAS-BTC
2026/01/09 22:08:56 [KuCoin WS 0] subscribed REQ-BTC
2026/01/09 22:08:56 [KuCoin WS 2] subscribed ATOM-ETH
2026/01/09 22:08:56 [KuCoin WS 1] subscribed PEPE-USDT
2026/01/09 22:08:56 [KuCoin WS 0] subscribed XTZ-BTC
2026/01/09 22:08:56 [KuCoin WS 2] subscribed BAX-ETH
2026/01/09 22:08:56 [KuCoin WS 1] subscribed POND-USDT
2026/01/09 22:08:57 [KuCoin WS 0] subscribed ZEC-BTC
2026/01/09 22:08:57 [KuCoin WS 2] subscribed TEL-ETH
2026/01/09 22:08:57 [KuCoin WS 1] subscribed ENJ-USDT
2026/01/09 22:08:57 [KuCoin WS 0] subscribed LTC-ETH
2026/01/09 22:08:57 [KuCoin WS 2] subscribed SKL-BTC
2026/01/09 22:08:57 [KuCoin WS 1] subscribed DGB-ETH
2026/01/09 22:08:57 [KuCoin WS 0] subscribed RUNE-USDT
2026/01/09 22:08:57 [KuCoin WS 2] subscribed ONE-BTC
2026/01/09 22:08:57 [KuCoin WS 1] subscribed ICP-BTC
2026/01/09 22:08:57 [KuCoin WS 0] subscribed AR-BTC
2026/01/09 22:08:57 [KuCoin WS 2] subscribed OGN-BTC
2026/01/09 22:08:57 [KuCoin WS 1] subscribed WAVES-USDT
2026/01/09 22:08:57 [KuCoin WS 0] subscribed DYP-USDT
2026/01/09 22:08:57 [KuCoin WS 2] subscribed TEL-BTC
2026/01/09 22:08:57 [KuCoin WS 1] subscribed SHIB-USDT
2026/01/09 22:08:57 [KuCoin WS 0] subscribed TEL-USDT
2026/01/09 22:08:57 [KuCoin WS 2] subscribed HBAR-BTC
2026/01/09 22:08:57 [KuCoin WS 1] subscribed TWT-BTC
2026/01/09 22:08:57 [KuCoin WS 0] subscribed VET-BTC
2026/01/09 22:08:57 [KuCoin WS 2] subscribed ICX-ETH
2026/01/09 22:08:57 [KuCoin WS 1] subscribed FET-USDT
2026/01/09 22:08:57 [KuCoin WS 0] subscribed COTI-BTC
2026/01/09 22:08:57 [KuCoin WS 2] subscribed WAVES-BTC
2026/01/09 22:08:57 [KuCoin WS 1] subscribed CHZ-BTC
2026/01/09 22:08:57 [KuCoin WS 0] subscribed XCN-USDT
2026/01/09 22:08:57 [KuCoin WS 2] subscribed KRL-BTC
2026/01/09 22:08:57 [KuCoin WS 1] subscribed ELA-USDT
2026/01/09 22:08:58 [KuCoin WS 0] subscribed MANA-ETH
2026/01/09 22:08:58 [KuCoin WS 2] subscribed IOTX-BTC
2026/01/09 22:08:58 [KuCoin WS 1] subscribed FET-BTC
2026/01/09 22:08:58 [KuCoin WS 0] subscribed KNC-USDT
2026/01/09 22:08:58 [KuCoin WS 2] subscribed STORJ-ETH
2026/01/09 22:08:58 [KuCoin WS 1] subscribed BNB-USDT
2026/01/09 22:08:58 [KuCoin WS 0] subscribed IOTX-USDT
2026/01/09 22:08:58 [KuCoin WS 2] subscribed ADA-BTC
2026/01/09 22:08:58 [KuCoin WS 1] subscribed FET-ETH
2026/01/09 22:08:58 [KuCoin WS 0] subscribed MOVR-USDT
2026/01/09 22:08:58 [KuCoin WS 2] subscribed TRX-ETH
2026/01/09 22:08:58 [KuCoin WS 1] subscribed VSYS-USDT
2026/01/09 22:08:58 [KuCoin WS 0] subscribed USDT-EUR
2026/01/09 22:08:58 [KuCoin WS 2] subscribed A-ETH
2026/01/09 22:08:58 [KuCoin WS 1] subscribed KLV-BTC
2026/01/09 22:08:58 [KuCoin WS 0] subscribed DOT-KCS
2026/01/09 22:08:58 [KuCoin WS 2] subscribed DGB-USDT
2026/01/09 22:08:58 [KuCoin WS 1] subscribed SUI-USDT
2026/01/09 22:08:58 [KuCoin WS 0] subscribed BNB-KCS
2026/01/09 22:08:58 [KuCoin WS 2] subscribed STX-USDT
2026/01/09 22:08:58 [KuCoin WS 1] subscribed KCS-BTC
2026/01/09 22:08:58 [KuCoin WS 0] subscribed CRO-BTC
2026/01/09 22:08:58 [KuCoin WS 1] subscribed BCH-BTC
2026/01/09 22:08:59 [KuCoin WS 0] subscribed DYP-ETH
2026/01/09 22:08:59 [KuCoin WS 1] subscribed IOTA-USDT
2026/01/09 22:08:59 [KuCoin WS 0] subscribed ONE-USDT
2026/01/09 22:08:59 [KuCoin WS 1] subscribed NEAR-USDT
2026/01/09 22:08:59 [KuCoin WS 0] subscribed INJ-BTC
2026/01/09 22:08:59 [KuCoin WS 1] subscribed AVA-ETH
2026/01/09 22:08:59 [KuCoin WS 1] subscribed NEO-BTC
2026/01/09 22:08:59 [KuCoin WS 1] subscribed SUPER-BTC
2026/01/09 22:08:59 [KuCoin WS 1] subscribed ATOM-BTC
2026/01/09 22:08:59 [KuCoin WS 1] subscribed LYX-ETH
2026/01/09 22:08:59 [KuCoin WS 1] subscribed ICX-USDT
2026/01/09 22:09:00 [KuCoin WS 1] subscribed ETC-USDT
2026/01/09 22:09:00 [KuCoin WS 1] subscribed WAN-USDT
2026/01/09 22:09:00 [KuCoin WS 1] subscribed GAS-BTC
2026/01/09 22:09:00 [KuCoin WS 1] subscribed EWT-USDT
2026/01/09 22:09:00 [KuCoin WS 1] subscribed XRP-BTC





