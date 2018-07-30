package api

import (
	"net/http"
	"database/sql"
	"github.com/journeymidnight/nier/src/store"
	"github.com/journeymidnight/nier/src/helper"
)

func ListHistory(w http.ResponseWriter, r *http.Request) {
	records, err := store.Db.GetHistoryRecord()
	if err != nil && err != sql.ErrNoRows {
		helper.Logger.Infoln("failed GetSNMP for query:", r)
		WriteErrorResponse(w, r, err)
		return
	}
	WriteSuccessResponse(w, EncodeResponse(records))
}

