package usecase

import (
	"log"
	"os"
	"sync"
	"time"
)

type RetryProxyAuthUseCase struct {
	proxyAuthOutput   *ProxyAuthOutput
	wg                *sync.WaitGroup
	mu                *sync.Mutex
	attemptsLimit     int
	attemptDelay      time.Duration
	refreshTokenDelay time.Duration
}

func NewRetryProxyAuthUseCase(proxyAuthOutput *ProxyAuthOutput, wg *sync.WaitGroup, mu *sync.Mutex, attemptsLimit int, attemptDelay time.Duration, refreshTokenDelay time.Duration) *RetryProxyAuthUseCase {
	return &RetryProxyAuthUseCase{
		proxyAuthOutput:   proxyAuthOutput,
		wg:                wg,
		mu:                mu,
		attemptsLimit:     attemptsLimit,
		attemptDelay:      attemptDelay,
		refreshTokenDelay: refreshTokenDelay,
	}
}

func (r *RetryProxyAuthUseCase) Execute() {
	defer r.wg.Done()

	for {
		for attempt := 1; attempt < r.attemptsLimit; attempt++ {
			log.Printf("%d of %d Proxy Authentication with WebSocket Gateway\n", attempt, r.attemptsLimit)
			out, err := AuthProxy(ProxyAuthInput{
				Host:     os.Getenv("WEBSOCKET_GATEWAY_HTTP"),
				Protocol: "http",
				Path:     "/auth/proxy",
			})
			if err == nil && out != nil && out.Token != "" {
				r.mu.Lock()
				*r.proxyAuthOutput = *out
				r.mu.Unlock()
				log.Println("Proxy authenticated successfully.")
				return
			}
			if err != nil {
				log.Printf("Error on proxy auth (attempt %d): %s", attempt, err.Error())
			} else {
				log.Println("Proxy authentication failed. Retrying...")
			}
			time.Sleep(r.attemptDelay)
		}

		log.Println("Reached max authentication attempts. Retrying after delay...")
		time.Sleep(r.refreshTokenDelay)
	}
}
