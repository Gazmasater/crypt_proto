git rm --cached cmd/arb/metrics/arb_metrics.csv

echo "cmd/arb/metrics/*.csv" >> .gitignore

git add .gitignore
git commit --amend --no-edit


git push origin new_arh --force



git filter-branch --force --index-filter \
'git rm --cached --ignore-unmatch cmd/arb/metrics/arb_metrics.csv' \
--prune-empty --tag-name-filter cat -- new_arh


rm -rf .git/refs/original/
git reflog expire --expire=now --all
git gc --prune=now --aggressive


git push origin new_arh --force



func (ws *kucoinWS) bootstrapAllBooks(ctx context.Context) error {
	const workers = 8

	start := time.Now()
	client := &http.Client{Timeout: httpTimeout}

	jobs := make(chan string, len(ws.symbols))
	doneCh := make(chan string, len(ws.symbols))
	errCh := make(chan error, 1)

	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()

		for symbol := range jobs {
			if ctx.Err() != nil {
				return
			}

			if err := ws.bootstrapBook(ctx, client, symbol); err != nil {
				select {
				case errCh <- fmt.Errorf("bootstrap %s: %w", symbol, err):
				default:
				}
				return
			}

			select {
			case doneCh <- symbol:
			case <-ctx.Done():
				return
			}
		}
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker()
	}

	for _, symbol := range ws.symbols {
		jobs <- symbol
	}
	close(jobs)

	total := len(ws.symbols)
	completed := 0

	waitCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitCh)
	}()

	for completed < total {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case err := <-errCh:
			return err

		case symbol := <-doneCh:
			completed++
			if completed%10 == 0 || completed == total {
				log.Printf("[KuCoin WS %d] bootstrap progress %d/%d last=%s\n",
					ws.id, completed, total, symbol)
			}

		case <-waitCh:
			if completed == total {
				log.Printf("[KuCoin WS %d] bootstrap finished %d/%d in %v\n",
					ws.id, completed, total, time.Since(start))
				return nil
			}
			return fmt.Errorf("bootstrap stopped early: %d/%d", completed, total)
		}
	}

	log.Printf("[KuCoin WS %d] bootstrap finished %d/%d in %v\n",
		ws.id, completed, total, time.Since(start))
	return nil
}

