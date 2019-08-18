package service

import (
	"context"
	"github.com/nodias/golang-oauth2.0-common/models"
	"github.com/nodias/golang-oauth2.0-common/shared/logger"
	"github.com/nodias/golang-oauth2.0-common/shared/repository"

	"go.elastic.co/apm"
)

//GetUserInfo is a function that gets specific user information about id.
func GetUserInfo(ctx context.Context, id string) (*models.User, *models.ResponseError) {
	log := logger.New(ctx)
	span, ctx := apm.StartSpan(ctx, "GetUserInfo", "custom")
	defer span.End()

	db := repository.NewOpenDB()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, models.NewResponseError(err, 500)
	}
	user := models.User{}
	defer db.Close()
	row := tx.QueryRowContext(ctx, "SELECT * FROM schema_user.user WHERE id = $1", id)
	err = row.Scan(&user.Id, &user.Name)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		rerr := models.NewResponseError(err, 500)
		log.WithError(rerr).Error("There is no corresponding user information.")
		return nil, rerr
	}
	log.WithField("user", user).Debug("User information retrieval was successful.")
	return &user, nil
}
