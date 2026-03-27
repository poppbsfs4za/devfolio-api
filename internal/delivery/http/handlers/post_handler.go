package handlers

import (
	"strconv"

	"github.com/example/devfolio-api/internal/delivery/http/response"
	"github.com/example/devfolio-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type PostHandler struct {
	usecase *usecase.PostUsecase
}

type postRequest struct {
	Title         string   `json:"title"`
	Summary       string   `json:"summary"`
	Content       string   `json:"content"`
	CoverImageURL string   `json:"cover_image_url"`
	Status        string   `json:"status"`
	Tags          []string `json:"tags"`
}

type PostResponse struct {
	ID            uint   `json:"id"`
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	Summary       string `json:"summary"`
	Content       string `json:"content"`
	CoverImageURL string `json:"cover_image_url"`
	Status        string `json:"status"`
	CreatedBy     uint   `json:"created_by"`
}

func NewPostHandler(usecase *usecase.PostUsecase) *PostHandler {
	return &PostHandler{usecase: usecase}
}

// ListPosts godoc
// @Summary List published posts
// @Description Get list of published posts
// @Tags posts
// @Produce json
// @Success 200 {array} PostResponse
// @Router /posts [get]
func (h *PostHandler) ListPublished(c *fiber.Ctx) error {
	posts, err := h.usecase.ListPublished()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	var res []PostResponse
	for _, p := range posts {
		res = append(res, PostResponse{
			ID:            p.ID,
			Title:         p.Title,
			Slug:          p.Slug,
			Summary:       p.Summary,
			Content:       p.Content,
			CoverImageURL: p.CoverImageURL,
			Status:        p.Status,
			CreatedBy:     p.CreatedBy,
		})
	}

	return response.JSON(c, fiber.StatusOK, res)
}

// GetPostBySlug godoc
// @Summary Get post by slug
// @Description Get single post detail
// @Tags posts
// @Produce json
// @Param slug path string true "Post slug"
// @Success 200 {object} PostResponse
// @Failure 404 {object} map[string]interface{}
// @Router /posts/{slug} [get]
func (h *PostHandler) GetPublishedBySlug(c *fiber.Ctx) error {
	post, err := h.usecase.GetPublishedBySlug(c.Params("slug"))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	if post == nil {
		return response.Error(c, fiber.StatusNotFound, "post not found")
	}

	var res PostResponse
	res = PostResponse{
		ID:            post.ID,
		Title:         post.Title,
		Slug:          post.Slug,
		Summary:       post.Summary,
		Content:       post.Content,
		CoverImageURL: post.CoverImageURL,
		Status:        post.Status,
		CreatedBy:     post.CreatedBy,
	}

	return response.JSON(c, fiber.StatusOK, res)
}

func (h *PostHandler) Create(c *fiber.Ctx) error {
	var req postRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body")
	}
	userID := c.Locals("user_id").(uint)
	post, err := h.usecase.Create(usecase.CreatePostInput{
		Title:         req.Title,
		Summary:       req.Summary,
		Content:       req.Content,
		CoverImageURL: req.CoverImageURL,
		Status:        req.Status,
		TagNames:      req.Tags,
		CreatedBy:     userID,
	})
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return response.JSON(c, fiber.StatusCreated, post)
}

func (h *PostHandler) Update(c *fiber.Ctx) error {
	var req postRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body")
	}
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid post id")
	}
	userID := c.Locals("user_id").(uint)
	post, err := h.usecase.Update(usecase.UpdatePostInput{
		ID:            uint(id),
		Title:         req.Title,
		Summary:       req.Summary,
		Content:       req.Content,
		CoverImageURL: req.CoverImageURL,
		Status:        req.Status,
		TagNames:      req.Tags,
		UpdatedBy:     userID,
	})
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return response.JSON(c, fiber.StatusOK, post)
}

func (h *PostHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid post id")
	}
	if err := h.usecase.Delete(uint(id)); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *PostHandler) AdminList(c *fiber.Ctx) error {
	posts, err := h.usecase.AdminList()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.JSON(c, fiber.StatusOK, posts)
}

func (h *PostHandler) AdminGetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid post id")
	}

	post, err := h.usecase.AdminGetByID(uint(id))
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "post not found")
	}

	return response.JSON(c, fiber.StatusOK, post)
}
