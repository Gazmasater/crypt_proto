Если оставить только нужное:

p99 execution latency
Micro-volatility (100 мс)
Fill ratio
Capture rate
Inventory drift




Название API
9623527002

696935c42a6dcd00013273f2
b348b686-55ff-4290-897b-02d55f815f65




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




go run -race main.go


GOMAXPROCS=8 go run -race main.go



(pprof) list queue                     
Total: 11.62MB
ROUTINE ======================== crypt_proto/internal/queue.(*MemoryStore).Run in /home/gaz358/myprog/crypt_proto/internal/queue/in_memory_queue.go
         0     6.09MB (flat, cum) 52.46% of Total
         .          .     73:func (s *MemoryStore) Run() {
         .          .     74:   for md := range s.batch {
         .     6.09MB     75:           s.apply(md)
         .          .     76:   }
         .          .     77:}
         .          .     78:
         .          .     79:func (s *MemoryStore) Push(md *models.MarketData) {
         .          .     80:   select {
ROUTINE ======================== crypt_proto/internal/queue.(*MemoryStore).apply in /home/gaz358/myprog/crypt_proto/internal/queue/in_memory_queue.go
         0     6.09MB (flat, cum) 52.46% of Total
         .          .     96:func (s *MemoryStore) apply(md *models.MarketData) {
         .          .     97:   key := md.Exchange + "|" + md.Symbol
         .          .     98:
         .          .     99:   buf, ok := s.buffers[key]
         .          .    100:   if !ok {
         .     6.09MB    101:           buf = NewRingBuffer(s.BufSize)
         .          .    102:           s.buffers[key] = buf
         .          .    103:   }
         .          .    104:
         .          .    105:   buf.Push(Quote{
         .          .    106:           Bid:       md.Bid,
ROUTINE ======================== crypt_proto/internal/queue.NewMemoryStore in /home/gaz358/myprog/crypt_proto/internal/queue/in_memory_queue.go
  553.04kB   553.04kB (flat, cum)  4.65% of Total
         .          .     64:func NewMemoryStore(bufSize int) *MemoryStore {
         .          .     65:   return &MemoryStore{
         .          .     66:           buffers: make(map[string]*RingBuffer),
  553.04kB   553.04kB     67:           batch:   make(chan *models.MarketData, 10_000),
         .          .     68:           BufSize: bufSize,
         .          .     69:   }
         .          .     70:}
         .          .     71:
         .          .     72:// Run — писатель, один поток
ROUTINE ======================== crypt_proto/internal/queue.NewRingBuffer in /home/gaz358/myprog/crypt_proto/internal/queue/in_memory_queue.go
    6.09MB     6.09MB (flat, cum) 52.46% of Total
         .          .     25:func NewRingBuffer(size int) *RingBuffer {
         .          .     26:   r := &RingBuffer{
    3.09MB     3.09MB     27:           data: make([]*atomic.Pointer[Quote], size),
         .          .     28:           size: size,
         .          .     29:           pos:  0,
         .          .     30:   }
         .          .     31:   for i := 0; i < size; i++ {
       3MB        3MB     32:           var ptr atomic.Pointer[Quote]
         .          .     33:           r.data[i] = &ptr
         .          .     34:   }
         .          .     35:   return r
         .          .     36:}
         .          .     37:
(pprof) 



type RingBuffer struct {
    data []Quote
    size uint64
    pos  uint64 // atomic
}

func (r *RingBuffer) Push(q Quote) {
    i := atomic.AddUint64(&r.pos, 1) - 1
    r.data[i%r.size] = q
}

