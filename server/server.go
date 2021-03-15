package server

import (
	"github.com/KhushrajRathod/repl.deploy/signature"
	"io"
	"log"
	"net/http"
	"os"
)

func Listen(handler func() error) {
	http.HandleFunc(sEndpoint, func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			http.Error(w, sBodyParseError, http.StatusBadRequest)
			return
		}

		signatureHeader := req.Header.Get(sSignatureHeaderName)

		if signatureHeader == "" {
			http.Error(w, sMissingSignatureError, http.StatusUnauthorized)
			return
		}

		validationError := signature.ValidateSignatureAndPayload(signatureHeader, body)

		if validationError != nil {
			http.Error(w, validationError.Err, validationError.Status)
			log.Println(sSignatureValidationFailedWarn)
			return
		}

		err = handler()

		if err != nil {
			
		}
	})

	err := http.ListenAndServe(sPort, nil)

	if err != nil {
		log.Fatalln(sUnexpectedHTTPServerCloseError)
		os.Exit(1)
	}
}
