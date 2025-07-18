package controllers

import (
	"appsku-golang/app/models"
	"appsku-golang/app/responses"
	"context"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"appsku-golang/app/global-utils/helper"
	"appsku-golang/app/global-utils/model"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gofiber/fiber/v2"
)

type FiberStoreController struct {
	Controller *StoreController
}

func NewFiberStoreController(controller *StoreController) *FiberStoreController {
	return &FiberStoreController{
		Controller: controller,
	}
}

func (c *FiberStoreController) GetAll(ctx *fiber.Ctx) error {
	var result model.Response
	var errorLog *model.ErrorLog

	page := 1
	limit := 10

	pageStr := ctx.Query("page")
	if pageStr != "" {
		pageInt, err := strconv.Atoi(pageStr)
		if err == nil && pageInt > 0 {
			page = pageInt
		}
	}

	limitStr := ctx.Query("limit")
	if limitStr != "" {
		limitInt, err := strconv.Atoi(limitStr)
		if err == nil && limitInt > 0 {
			limit = limitInt
		}
	}

	filter := make(map[string]interface{})

	name := ctx.Query("name")
	if name != "" {
		filter["name"] = primitive.Regex{Pattern: name, Options: "i"}
	}

	storeType := ctx.Query("type")
	if storeType != "" {
		filter["type"] = storeType
	}

	stores, totalCount, err := c.Controller.StoreUseCase.GetAll(context.Background(), filter, page, limit)
	if err != nil {
		errorLog = &model.ErrorLog{
			Message:       "Internal Server Error",
			SystemMessage: err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusInternalServerError
		return responses.SendFiberResponse(ctx, result)
	}

	result.StatusCode = http.StatusOK
	result.Page = page
	result.PerPage = limit
	result.Total = totalCount
	result.MaxPage = int(math.Ceil(float64(totalCount) / float64(limit)))
	result.Data = stores

	return responses.SendFiberResponse(ctx, result)
}

func (c *FiberStoreController) Insert(ctx *fiber.Ctx) error {
	var result model.Response
	var errorLog *model.ErrorLog
	var store models.Store

	if err := ctx.BodyParser(&store); err != nil {
		errorLog = &model.ErrorLog{
			Message:       "Bad Request",
			SystemMessage: err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusBadRequest
		return responses.SendFiberResponse(ctx, result)
	}

	if err := c.Controller.Validator.Struct(&store); err != nil {
		errorLog = &model.ErrorLog{
			Message:       "Bad Request",
			SystemMessage: err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusBadRequest
		return responses.SendFiberResponse(ctx, result)
	}

	createdStore, err := c.Controller.StoreUseCase.Insert(context.Background(), &store)
	if err != nil {
		errorLog := model.ErrorLog{
			Message:       "Internal Server Error",
			SystemMessage: err.Error(),
		}

		result.Error = &errorLog
		result.StatusCode = http.StatusInternalServerError
		return responses.SendFiberResponse(ctx, result)
	}

	result.StatusCode = http.StatusCreated
	result.Data = createdStore

	return responses.SendFiberResponse(ctx, result)
}

func (c *FiberStoreController) InsertWithSetting(ctx *fiber.Ctx) error {
	var result model.Response
	var errorLog *model.ErrorLog
	var storeWithSetting struct {
		Store   *models.Store        `json:"store" validate:"required"`
		Setting *models.StoreSetting `json:"setting" validate:"required"`
	}

	if err := ctx.BodyParser(&storeWithSetting); err != nil {
		errorLog = &model.ErrorLog{
			Message:       "Bad Request",
			SystemMessage: err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusBadRequest
		return responses.SendFiberResponse(ctx, result)
	}

	if err := c.Controller.Validator.Struct(storeWithSetting); err != nil {
		errorLog = &model.ErrorLog{
			Message:       "Bad Request",
			SystemMessage: err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusBadRequest
		return responses.SendFiberResponse(ctx, result)
	}

	if err := c.Controller.Validator.Struct(storeWithSetting.Store); err != nil {
		errorLog = &model.ErrorLog{
			Message:       "Bad Request",
			SystemMessage: "Invalid store data: " + err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusBadRequest
		return responses.SendFiberResponse(ctx, result)
	}

	if err := c.Controller.Validator.Struct(storeWithSetting.Setting); err != nil {
		errorLog = &model.ErrorLog{
			Message:       "Bad Request",
			SystemMessage: "Invalid setting: " + err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusBadRequest
		return responses.SendFiberResponse(ctx, result)
	}

	createdStore, createdSetting, err := c.Controller.StoreUseCase.InsertWithSetting(context.Background(), storeWithSetting.Store, storeWithSetting.Setting)
	if err != nil {
		errorLog := model.ErrorLog{
			Message:       "Internal Server Error",
			SystemMessage: err.Error(),
		}

		result.Error = &errorLog
		result.StatusCode = http.StatusInternalServerError
		return responses.SendFiberResponse(ctx, result)
	}

	storeWithSetting.Setting.ID = createdSetting.ID

	result.StatusCode = http.StatusCreated
	result.Data = map[string]interface{}{
		"store":   createdStore,
		"setting": storeWithSetting.Setting,
	}

	return responses.SendFiberResponse(ctx, result)
}

func (c *FiberStoreController) InsertWithLocation(ctx *fiber.Ctx) error {
	var result model.Response
	var errorLog *model.ErrorLog
	var storeWithLocation struct {
		Name        string                `json:"name" validate:"required"`
		Description string                `json:"description" validate:"required"`
		Type        string                `json:"type" validate:"required,oneof=retail grosir"`
		Location    *models.StoreLocation `json:"location" validate:"required"`
	}

	if err := ctx.BodyParser(&storeWithLocation); err != nil {
		errorLog = &model.ErrorLog{
			Message:       "Bad Request",
			SystemMessage: err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusBadRequest
		return responses.SendFiberResponse(ctx, result)
	}

	if err := c.Controller.Validator.Struct(storeWithLocation); err != nil {
		errorLog = &model.ErrorLog{
			Message:       "Bad Request",
			SystemMessage: err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusBadRequest
		return responses.SendFiberResponse(ctx, result)
	}

	store := &models.Store{
		Name:        storeWithLocation.Name,
		Description: storeWithLocation.Description,
		Type:        storeWithLocation.Type,
		Location:    storeWithLocation.Location,
	}

	createdStore, err := c.Controller.StoreUseCase.Insert(context.Background(), store)
	if err != nil {
		errorLog := model.ErrorLog{
			Message:       "Internal Server Error",
			SystemMessage: err.Error(),
		}

		result.Error = &errorLog
		result.StatusCode = http.StatusInternalServerError
		return responses.SendFiberResponse(ctx, result)
	}

	result.StatusCode = http.StatusCreated
	result.Data = createdStore

	return responses.SendFiberResponse(ctx, result)
}

func (c *FiberStoreController) GetById(ctx *fiber.Ctx) error {
	var result model.Response
	var errorLog *model.ErrorLog

	id := ctx.Params("id")

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		errorLog := model.ErrorLog{
			Message:       "Bad Request",
			SystemMessage: "Invalid ID format",
		}

		result.Error = &errorLog
		result.StatusCode = http.StatusBadRequest
		return responses.SendFiberResponse(ctx, result)
	}

	store, err := c.Controller.StoreUseCase.GetById(context.Background(), objId)
	if err != nil {
		errorLog = &model.ErrorLog{
			Message:       "Internal Server Error",
			SystemMessage: err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusInternalServerError
		return responses.SendFiberResponse(ctx, result)
	}
	result.StatusCode = http.StatusOK
	result.Data = store

	return responses.SendFiberResponse(ctx, result)
}

func (c *FiberStoreController) Update(ctx *fiber.Ctx) error {
	var result model.Response
	var errorLog *model.ErrorLog
	var storeData map[string]interface{}

	id := ctx.Params("id")

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		errorLog := model.ErrorLog{
			Message:       "Bad Request",
			SystemMessage: "Invalid ID format",
		}

		result.Error = &errorLog
		result.StatusCode = http.StatusBadRequest
		return responses.SendFiberResponse(ctx, result)
	}

	if err := ctx.BodyParser(&storeData); err != nil {
		errorLog = &model.ErrorLog{
			Message:       "Bad Request",
			SystemMessage: err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusBadRequest
		return responses.SendFiberResponse(ctx, result)
	}

	delete(storeData, "id")
	delete(storeData, "_id")
	delete(storeData, "created_at")
	delete(storeData, "updated_at")
	delete(storeData, "deleted_at")

	if len(storeData) == 0 {
		errorLog = &model.ErrorLog{
			Message:       "Bad Request",
			SystemMessage: "No valid fields to update",
		}

		result.Error = errorLog
		result.StatusCode = http.StatusBadRequest
		return responses.SendFiberResponse(ctx, result)
	}

	err = c.Controller.StoreUseCase.Update(context.Background(), objId, storeData)
	if err != nil {
		errorLog = &model.ErrorLog{
			Message:       "Internal Server Error",
			SystemMessage: err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusInternalServerError
		return responses.SendFiberResponse(ctx, result)
	}

	result.StatusCode = http.StatusOK
	result.Data = map[string]interface{}{"message": "Store updated successfully"}

	return responses.SendFiberResponse(ctx, result)
}

func (c *FiberStoreController) Delete(ctx *fiber.Ctx) error {
	var result model.Response
	var errorLog *model.ErrorLog

	id := ctx.Params("id")
	hardDelete := false

	hardDeleteStr := ctx.Query("hard_delete")
	if hardDeleteStr == "true" {
		hardDelete = true
	}

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		errorLog := model.ErrorLog{
			Message:       "Bad Request",
			SystemMessage: "Invalid ID format",
		}

		result.Error = &errorLog
		result.StatusCode = http.StatusBadRequest
		return responses.SendFiberResponse(ctx, result)
	}

	err = c.Controller.StoreUseCase.Delete(context.Background(), objId, hardDelete)
	if err != nil {
		errorLog = &model.ErrorLog{
			Message:       "Internal Server Error",
			SystemMessage: err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusInternalServerError
		return responses.SendFiberResponse(ctx, result)
	}

	result.StatusCode = http.StatusOK
	if hardDelete {
		result.Data = map[string]interface{}{"message": "Store permanently deleted"}
	} else {
		result.Data = map[string]interface{}{"message": "Store deleted successfully"}
	}

	return responses.SendFiberResponse(ctx, result)
}

func (c *FiberStoreController) UploadFile(ctx *fiber.Ctx) error {
	var result model.Response
	var errorLog *model.ErrorLog

	file, err := ctx.FormFile("file")
	if err != nil {
		errorLog = &model.ErrorLog{
			Message:       "Bad Request",
			SystemMessage: "File not found in request: " + err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusBadRequest
		return responses.SendFiberResponse(ctx, result)
	}

	uploadsDir := "./files/uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		errorLog = &model.ErrorLog{
			Message:       "Internal Server Error",
			SystemMessage: "Failed to create uploads directory: " + err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusInternalServerError
		return responses.SendFiberResponse(ctx, result)
	}

	fileExt := filepath.Ext(file.Filename)
	uniqueID := helper.GetUniqueImageName()
	fileName := uniqueID + fileExt
	filePath := filepath.Join(uploadsDir, fileName)

	if err := ctx.SaveFile(file, filePath); err != nil {
		errorLog = &model.ErrorLog{
			Message:       "Internal Server Error",
			SystemMessage: "Failed to save file: " + err.Error(),
		}

		result.Error = errorLog
		result.StatusCode = http.StatusInternalServerError
		return responses.SendFiberResponse(ctx, result)
	}

	now := time.Now()
	fileModel := &models.File{
		FileName:     fileName,
		OriginalName: file.Filename,
		FilePath:     filePath,
		FileSize:     file.Size,
		FileType:     helper.GetMime(fileExt[1:]),
		CreatedAt:    &now,
	}

	result.StatusCode = http.StatusOK

	fileModel.ID = primitive.NewObjectID() // we must save the data file into a database

	result.Data = fileModel

	return responses.SendFiberResponse(ctx, result)
}
