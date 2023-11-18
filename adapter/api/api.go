package api

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/waffleboot/ddd2/adapter/api/account"
	"github.com/waffleboot/ddd2/domain"
)

var getRegexp = regexp.MustCompile(`\/accounts\/(\d+)`)
var createRegexp = regexp.MustCompile(`\/accounts`)
var depositRegexp = regexp.MustCompile(`\/accounts\/(\d+)\/deposit\/(\d+)`)
var withdrawRegexp = regexp.MustCompile(`\/accounts\/(\d+)\/withdraw\/(\d+)`)

func Start(handler *account.Handler) error {
	s := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(makeHandler(handler)),
	}

	err := s.ListenAndServe()
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	return nil
}

func makeHandler(handler *account.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			s := depositRegexp.FindAllStringSubmatch(r.URL.Path, -1)
			if s != nil {

				accountId, err := parseAccountId(s[0][1])
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				amount, err := parseAmount(s[0][2])
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				err = handler.DepositAccount(r.Context(), accountId, amount)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusOK)
				return
			}

			s = withdrawRegexp.FindAllStringSubmatch(r.URL.Path, -1)
			if s != nil {

				accountId, err := parseAccountId(s[0][1])
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				amount, err := parseAmount(s[0][2])
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				err = handler.WithdrawAccount(r.Context(), accountId, amount)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusOK)
				return
			}

			s = createRegexp.FindAllStringSubmatch(r.URL.Path, -1)
			if s != nil {

				id, err := handler.CreateAccount(r.Context())
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "%d\n", id)
				return
			}
		}

		if r.Method == http.MethodGet {
			s := getRegexp.FindAllStringSubmatch(r.URL.Path, -1)
			if s != nil {

				accountId, err := parseAccountId(s[0][1])
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				account, err := handler.GetAccount(r.Context(), accountId)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "{amount: %d version: %d}\n", account.Amount(), account.Version())
				return
			}
		}

		http.Error(w, "bad request", http.StatusBadRequest)
	}
}

func parseAccountId(s string) (domain.AccountId, error) {
	accountId, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return domain.AccountId(0), fmt.Errorf("parse account id: %w", err)
	}

	return domain.AccountId(accountId), nil
}

func parseAmount(s string) (domain.Amount, error) {
	amount, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return domain.Amount(0), fmt.Errorf("parse amount: %w", err)
	}

	return domain.Amount(amount), nil
}
