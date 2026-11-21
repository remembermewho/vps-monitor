package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"syscall"
	"time"
)

type Result struct {
	URL           string        // URL to check
	CheckDate     time.Time     // Date of the check
	ResultOfCheck string        // Result of the check | Might be SUCCESS, SLOW (Status code between 200 and 400 and TimeTaken > 2 sec), FAILURE
	StatusCode    int           // HTTP status code
	TimeTaken     time.Duration // Time taken for the request
	ErrorMessage  string        // Error message if any
}

func main() {
	result := checkUrl("http://127.0.0.1:1")
	fmt.Println(result)
}

func checkUrl(URL string) string {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(URL)
	if err != nil {

		var netTimeoutError net.Error // Проверка на таймаут
		if errors.As(err, &netTimeoutError) && netTimeoutError.Timeout() {
			fmt.Printf("Операция прервана по таймауту: %v\n", err)
			return "FAILURE"
		}

		var netDNSError *net.DNSError // Проверка на DNS
		if errors.As(err, &netDNSError) {
			fmt.Printf("Операция прервана из за DNS: %v\n", err)
			return "FAILURE"
		}

		var netTLSError *tls.CertificateVerificationError // Проверка на сертификат
		if errors.As(err, &netTLSError) {
			fmt.Printf("Операция прервана из за сертификата сервера, подробнее: %v\n", err)
			return "FAILURE"
		}

		if errors.Is(err, syscall.ECONNREFUSED) { // Проверка на подключение
			fmt.Printf("Операция прервана, хост/порт отказался принять подключение. Подробнее: %v\n", err)
			return "FAILURE"
		}

		var netOpErr *net.OpError // Прочие сетевые ошибки
		if errors.As(err, &netOpErr) {
			fmt.Printf("Сетевая ошибка операции '%s' к '%s': '%v'\n", netOpErr.Op, netOpErr.Addr, netOpErr.Err)
			return "FAILURE"
		}

		fmt.Printf("Неизвестная ошибка: %v\n", err)
		return "FAILURE"
	}
	defer resp.Body.Close()
	result := printStatus(resp.StatusCode)

	return result
}

func printStatus(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 400:
		return "SUCCESS"
	default:
		return "FAILURE"
	}
}
