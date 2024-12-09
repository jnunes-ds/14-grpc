package service

import (
	"context"
	"github.com/jnunes-ds/14-grpc/internal/database"
	"github.com/jnunes-ds/14-grpc/internal/pb"
	"io"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDb database.Category
}

func NewCategoryService(db database.Category) *CategoryService {
	return &CategoryService{CategoryDb: db}
}

func (c *CategoryService) CreateCategory(_ context.Context, in *pb.CreateCategoryRequest) (*pb.Category, error) {
	category, err := c.CategoryDb.Create(in.Name, in.Description)
	if err != nil {
		return nil, err
	}
	CategoryResponse := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return CategoryResponse, nil
}

func (c *CategoryService) ListCategories(context.Context, *pb.Blank) (*pb.CategoryList, error) {
	categories, err := c.CategoryDb.FindAll()
	if err != nil {
		return nil, err
	}

	var CategoryResponses []*pb.Category

	for _, category := range categories {
		CategoryResponses = append(CategoryResponses, &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}

	return &pb.CategoryList{Categories: CategoryResponses}, nil
}

func (c *CategoryService) GetCategory(_ context.Context, in *pb.CategoryGetRequest) (*pb.Category, error) {
	category, err := c.CategoryDb.Find(in.Id)
	if err != nil {
		return nil, err
	}

	CategoryResponse := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return CategoryResponse, nil
}

func (c *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.CategoryList{}

	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(categories)
		}
		if err != nil {
			return err
		}

		categoryResult, err := c.CategoryDb.Create(category.Name, category.Description)
		if err != nil {
			return err
		}
		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		})
	}
}

func (c *CategoryService) CreateCategoryBidirectionalStream(stream pb.CategoryService_CreateCategoryBidirectionalStreamServer) error {
	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		categoryResult, err := c.CategoryDb.Create(category.Name, category.Description)
		if err != nil {
			return err
		}
		err = stream.Send(&pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		})
		if err != nil {
			return err
		}
	}
}
