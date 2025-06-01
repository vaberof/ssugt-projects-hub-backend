package handlers

import (
	"net/http"
	"ssugt-projects-hub/models"
	"ssugt-projects-hub/pkg/auth"
	"ssugt-projects-hub/pkg/logging/logs"
	"ssugt-projects-hub/pkg/xhttp"
	userservice "ssugt-projects-hub/service/user"
	"strconv"
)

func GetUsersHandler(logs *logs.Logs, userService userservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logs.WithName("get-users-handler")

		ctx, err := auth.GetContext(r)
		if err != nil {
			log.Debug(err.Error())
		}

		userIdsStr := r.URL.Query()["ids"]

		userIds := make([]int, 0, len(userIdsStr))
		for _, s := range userIdsStr {
			id, err := strconv.Atoi(s)
			if err != nil {
				log.Debug(err.Error())
				http.Error(w, "invalid ids param", http.StatusBadRequest)
				return
			}
			userIds = append(userIds, id)
		}

		users, err := userService.GetByIds(ctx, userIds)
		if err != nil {
			log.Error(err.Error())
			xhttp.BadRequest(w)
			return
		}

		response := models.GetUsersResponse{
			Users: models.MapToUsersResponse(users),
		}

		if err = xhttp.WriteResponseJson(w, http.StatusOK, response); err != nil {
			log.Error("Не удалось записать ответ:", err)
			xhttp.BadRequest(w)
		}
	}
}
