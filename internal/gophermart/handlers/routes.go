package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/dnsoftware/gophermart/internal/constants"
	"github.com/dnsoftware/gophermart/internal/gophermart/domain"
	"github.com/dnsoftware/gophermart/internal/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
	"regexp"
	"strconv"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()
	return r
}

// регистрация пользователя
func (h *Server) userRegister(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), constants.DBContextTimeout)
	defer cancel()

	var buf bytes.Buffer

	var user domain.UserItem

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &user); err != nil {
		code, message := constants.StatusData(constants.RegisterBadFormat)
		http.Error(res, message, code)
		return
	}

	token, status, err := h.userMart.AddUser(ctx, user.Login, user.Password)
	if err != nil {
		logger.Log().Error(err.Error())
		code, message := constants.StatusData(status)
		http.Error(res, message+", "+err.Error(), code)
		return
	}

	if status != constants.RegisterOk {
		code, message := constants.StatusData(status)
		http.Error(res, message, code)
		return
	}

	res.Header().Set("Content-Type", constants.ApplicationJSON)
	bearer := "Bearer " + token
	res.Header().Set(constants.HeaderAuthorization, bearer)
	code, message := constants.StatusData(status)
	res.WriteHeader(code)
	res.Write([]byte(message))

}

// вход пользователя
func (h *Server) userLogin(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), constants.DBContextTimeout)
	defer cancel()

	var buf bytes.Buffer

	var user domain.UserItem

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &user); err != nil {
		code, message := constants.StatusData(constants.LoginBadFormat)
		http.Error(res, message, code)
		return
	}

	token, status, err := h.userMart.LoginUser(ctx, user.Login, user.Password)
	if err != nil {
		logger.Log().Error(err.Error())
		code, message := constants.StatusData(status)
		http.Error(res, message+", "+err.Error(), code)
		return
	}

	if status != constants.LoginOk {
		code, message := constants.StatusData(status)
		http.Error(res, message, code)
		return
	}

	res.Header().Set("Content-Type", constants.ApplicationJSON)
	bearer := "Bearer " + token
	res.Header().Set(constants.HeaderAuthorization, bearer)
	code, message := constants.StatusData(status)
	res.WriteHeader(code)
	res.Write([]byte(message))

}

func (h *Server) userOrderUpload(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), constants.DBContextTimeout)
	defer cancel()

	var buf bytes.Buffer

	uid := ctx.Value(constants.UserIDKey)
	userID := uid.(int64)

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	number := buf.String()
	re := regexp.MustCompile(`^\d+$`)
	if !re.MatchString(number) {
		http.Error(res, "", http.StatusBadRequest)
		return
	}

	orderID, _ := strconv.ParseInt(number, 10, 64)

	status, err := h.orderMart.AddOrder(ctx, userID, orderID)
	if err != nil {
		code, message := constants.StatusData(status)
		http.Error(res, message, code)
		return
	}

	res.Header().Set("Content-Type", constants.ApplicationJSON)
	code, message := constants.StatusData(status)
	res.WriteHeader(code)
	res.Write([]byte(message))
}

func (h *Server) userOrdersList(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), constants.DBContextTimeout)
	defer cancel()

	uid := ctx.Value(constants.UserIDKey)
	userID := uid.(int64)

	list, status, err := h.orderMart.OrdersList(ctx, userID)
	if err != nil {
		code, message := constants.StatusData(status)
		http.Error(res, message, code)
		return
	}

	body, err := json.Marshal(list)

	res.Header().Set("Content-Type", constants.ApplicationJSON)
	code, _ := constants.StatusData(status)
	res.WriteHeader(code)
	res.Write([]byte(body))
}

func (h *Server) noMetricType(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "Metric type required!", http.StatusBadRequest)
}

func (h *Server) unrecognized(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "Not found!", http.StatusNotFound)
}
