// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package standalone

import (
	log "github.com/sirupsen/logrus"
	"go.amzn.com/lambda/core/directinvoke"
	"go.amzn.com/lambda/rapidcore"
	"net/http"
)

func DirectInvokeHandler(w http.ResponseWriter, r *http.Request, s rapidcore.InteropServer) {
	tok := s.CurrentToken()
	if tok == nil {
		log.Errorf("Attempt to call directInvoke without Reserve")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	invoke, err := directinvoke.ReceiveDirectInvoke(w, r, *tok)
	if err != nil {
		log.Errorf("direct invoke error: %s", err)
		return
	}

	if err := s.FastInvoke(w, invoke, true); err != nil {
		switch err {
		case rapidcore.ErrNotReserved:
		case rapidcore.ErrAlreadyReplied:
		case rapidcore.ErrAlreadyInvocating:
			log.Errorf("Failed to set reply stream: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		case rapidcore.ErrInvokeReservationDone:
			w.WriteHeader(http.StatusBadGateway)
		}
	}
}
