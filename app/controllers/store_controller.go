package controllers

import (
	"appsku-golang/app/constants"
	"appsku-golang/app/models"
	"appsku-golang/app/responses"
	"appsku-golang/app/usecases"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"net/http"
	"strconv"

	"appsku-golang/app/global-utils/model"

	"github.com/gin-gonic/gin"
)

type StoreController struct {
	StoreUseCase usecases.IStoreUseCase
	Validator    *validator.Validate
}

func NewStoreController(usecase usecases.IStoreUseCase) *StoreController {
	return &StoreController{
		StoreUseCase: usecase,
		Validator:    validator.New(),
	}
}

func (c *StoreController) GetAll(r *gin.Context) {
	var result model.Response
	var errorLog *model.ErrorLog

	page := 1
	limit := 10

	pageStr := r.Query("page")
	if pageStr != "" {
		pageInt, err := strconv.Atoi(pageStr)
		if err == nil && pageInt > 0 {
			page = pageInt
		}
	}

	limitStr := r.Query("limit")
	if limitStr != "" {
		limitInt, err := strconv.Atoi(limitStr)
		if err == nil && limitInt > 0 {
			limit = limitInt
		}
	}

	filter := make(map[string]interface{})

	name := r.Query("name")
	if name != "" {
		filter["name"] = primitive.Regex{Pattern: name, Options: "i"}
	}

	storeType := r.Query("type")
	if storeType != "" {
		filter["type"] = storeType
	}

	stores, totalCount, err := c.StoreUseCase.GetAll(r, filter, page, limit)
	if err != nil {
		errorLog = &model.ErrorLog{
			Message:       constants.InternalServerError,
			SystemMessage: err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusInternalServerError
		responses.SendResponse(r, result)
		return
	}

	result.StatusCode = http.StatusOK

	result.Page = page
	result.PerPage = limit
	result.Total = totalCount
	result.MaxPage = int(math.Ceil(float64(totalCount) / float64(limit)))
	result.Data = stores

	responses.SendResponse(r, result)
}

func (c *StoreController) Insert(r *gin.Context) {
	var result model.Response
	var errorLog *model.ErrorLog
	var store *models.Store

	if err := r.ShouldBindJSON(&store); err != nil {
		errorLog = &model.ErrorLog{
			Message:       constants.BadRequest,
			SystemMessage: err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusBadRequest
		responses.SendResponse(r, result)
		return
	}

	if err := c.Validator.Struct(store); err != nil {
		errorLog = &model.ErrorLog{
			Message:       constants.BadRequest,
			SystemMessage: err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusBadRequest
		responses.SendResponse(r, result)
		return
	}

	createdStore, err := c.StoreUseCase.Insert(r, store)
	if err != nil {
		errorLog := model.ErrorLog{
			Message:       constants.InternalServerError,
			SystemMessage: err.Error(),
		}

		result.Error = &errorLog
		result.StatusCode = http.StatusInternalServerError
		responses.SendResponse(r, result)
		return
	}

	result.StatusCode = http.StatusCreated
	result.Data = createdStore

	responses.SendResponse(r, result)
}

type StoreWithSetting struct {
	Store   *models.Store        `json:"store" validate:"required"`
	Setting *models.StoreSetting `json:"setting" validate:"required"`
}

func (c *StoreController) InsertWithSetting(r *gin.Context) {
	var result model.Response
	var errorLog *model.ErrorLog
	var storeWithSetting StoreWithSetting

	if err := r.ShouldBindJSON(&storeWithSetting); err != nil {
		errorLog = &model.ErrorLog{
			Message:       constants.BadRequest,
			SystemMessage: err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusBadRequest
		responses.SendResponse(r, result)
		return
	}

	if err := c.Validator.Struct(storeWithSetting); err != nil {
		errorLog = &model.ErrorLog{
			Message:       constants.BadRequest,
			SystemMessage: err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusBadRequest
		responses.SendResponse(r, result)
		return
	}

	if err := c.Validator.Struct(storeWithSetting.Store); err != nil {
		errorLog = &model.ErrorLog{
			Message:       constants.BadRequest,
			SystemMessage: "Invalid store data: " + err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusBadRequest
		responses.SendResponse(r, result)
		return
	}

	if err := c.Validator.Struct(storeWithSetting.Setting); err != nil {
		errorLog = &model.ErrorLog{
			Message:       constants.BadRequest,
			SystemMessage: "Invalid setting: " + err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusBadRequest
		responses.SendResponse(r, result)
		return
	}

	createdStore, createdSetting, err := c.StoreUseCase.InsertWithSetting(r, storeWithSetting.Store, storeWithSetting.Setting)
	if err != nil {
		errorLog := model.ErrorLog{
			Message:       constants.InternalServerError,
			SystemMessage: err.Error(),
		}

		result.Error = &errorLog
		result.StatusCode = http.StatusInternalServerError
		responses.SendResponse(r, result)
		return
	}

	storeWithSetting.Setting.ID = createdSetting.ID

	result.StatusCode = http.StatusCreated
	result.Data = gin.H{
		"store":   createdStore,
		"setting": storeWithSetting.Setting,
	}

	responses.SendResponse(r, result)
}

func (c *StoreController) GetById(r *gin.Context) {
	var result model.Response
	var errorLog *model.ErrorLog

	id := r.Param("id")

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		errorLog := model.ErrorLog{
			Message:       constants.InternalServerError,
			SystemMessage: err.Error(),
		}

		result.Error = &errorLog
		result.StatusCode = http.StatusInternalServerError
		responses.SendResponse(r, result)
		return
	}

	store, err := c.StoreUseCase.GetById(r, objId)
	if err != nil {
		errorLog = &model.ErrorLog{
			Message:       constants.InternalServerError,
			SystemMessage: err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusInternalServerError
		responses.SendResponse(r, result)
		return
	}
	result.StatusCode = http.StatusOK
	result.Data = store

	responses.SendResponse(r, result)
}

func (c *StoreController) Update(r *gin.Context) {
	var result model.Response
	var errorLog *model.ErrorLog
	var storeData map[string]interface{}

	id := r.Param("id")

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		errorLog := model.ErrorLog{
			Message:       constants.BadRequest,
			SystemMessage: "Invalid ID format",
		}

		result.Error = &errorLog
		result.StatusCode = http.StatusBadRequest
		responses.SendResponse(r, result)
		return
	}

	if err := r.ShouldBindJSON(&storeData); err != nil {
		errorLog = &model.ErrorLog{
			Message:       constants.BadRequest,
			SystemMessage: err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusBadRequest
		responses.SendResponse(r, result)
		return
	}

	delete(storeData, "id")
	delete(storeData, "_id")
	delete(storeData, "created_at")
	delete(storeData, "updated_at")
	delete(storeData, "deleted_at")

	if len(storeData) == 0 {
		errorLog = &model.ErrorLog{
			Message:       constants.BadRequest,
			SystemMessage: "No valid fields to update",
		}

		result.Error = errorLog
		result.StatusCode = http.StatusBadRequest
		responses.SendResponse(r, result)
		return
	}

	err = c.StoreUseCase.Update(r, objId, storeData)
	if err != nil {
		errorLog = &model.ErrorLog{
			Message:       constants.InternalServerError,
			SystemMessage: err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusInternalServerError
		responses.SendResponse(r, result)
		return
	}

	result.StatusCode = http.StatusOK
	result.Data = gin.H{"message": "Store updated successfully"}

	responses.SendResponse(r, result)
}

func (c *StoreController) Delete(r *gin.Context) {
	var result model.Response
	var errorLog *model.ErrorLog

	id := r.Param("id")
	hardDelete := false

	hardDeleteStr := r.Query("hard_delete")
	if hardDeleteStr == "true" {
		hardDelete = true
	}

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		errorLog := model.ErrorLog{
			Message:       constants.BadRequest,
			SystemMessage: "Invalid ID format",
		}

		result.Error = &errorLog
		result.StatusCode = http.StatusBadRequest
		responses.SendResponse(r, result)
		return
	}

	err = c.StoreUseCase.Delete(r, objId, hardDelete)
	if err != nil {
		errorLog = &model.ErrorLog{
			Message:       constants.InternalServerError,
			SystemMessage: err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusInternalServerError
		responses.SendResponse(r, result)
		return
	}

	result.StatusCode = http.StatusOK
	if hardDelete {
		result.Data = gin.H{"message": "Store permanently deleted"}
	} else {
		result.Data = gin.H{"message": "Store deleted successfully"}
	}

	responses.SendResponse(r, result)
}
