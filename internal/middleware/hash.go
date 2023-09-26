package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/ilegorro/almetrics/internal/server"
)

func WithHash(app *server.App) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rHash := r.Header.Get("HashSHA256")
			key := app.Options.Key
			if key != "" && rHash != "" {
				buf, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Error reading body", http.StatusInternalServerError)
					return
				}

				h := hmac.New(sha256.New, []byte(key))
				h.Write(buf)
				sHash := hex.EncodeToString(h.Sum(nil))
				if rHash != sHash {
					http.Error(w, "Invalid hash", http.StatusBadRequest)
					return
				}
				r.Body.Close()
				r.Body = io.NopCloser(bytes.NewBuffer(buf))
			}
			hw := hashResponseWriter{
				ResponseWriter: w,
				key:            key,
			}

			h.ServeHTTP(&hw, r)
		})
	}
}

type hashResponseWriter struct {
	http.ResponseWriter
	key string
}

func (r *hashResponseWriter) Write(b []byte) (int, error) {
	h := hmac.New(sha256.New, []byte(r.key))
	h.Write(b)
	r.ResponseWriter.Header().Set("HashSHA256", hex.EncodeToString(h.Sum(nil)))
	size, err := r.ResponseWriter.Write(b)

	return size, err
}
