package post

import (
	"context"
	"fmt"
	"log"
	"sync"

	postDomain "spendsnap-backend/internal/domain/post"
	"spendsnap-backend/internal/domain/storage"
	transactionDomain "spendsnap-backend/internal/domain/transaction"
	"spendsnap-backend/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

type CreatePostUsecase struct {
	pool               *pgxpool.Pool
	postRepo           postDomain.PostRepository
	uploadUsecase      storage.UploadUsecase
	transactionUsecase transactionDomain.TransactionUsecase
}

func NewCreatePostUsecase(pool *pgxpool.Pool, postRepo postDomain.PostRepository, uploadUC storage.UploadUsecase, transactionUC transactionDomain.TransactionUsecase) *CreatePostUsecase {
	return &CreatePostUsecase{
		pool:               pool,
		postRepo:           postRepo,
		uploadUsecase:      uploadUC,
		transactionUsecase: transactionUC,
	}
}

func (uc *CreatePostUsecase) CreatePost(ctx context.Context, post *postDomain.Post, fileToupload *storage.File, transaction *transactionDomain.Transaction) (*postDomain.PostResponse, error) {
	if uc == nil || uc.postRepo == nil {
		return nil, fmt.Errorf("post usecase not initialized")
	}
	if uc.uploadUsecase == nil {
		return nil, fmt.Errorf("upload usecase not initialized: missing DI wiring")
	}
	if fileToupload == nil {
		return nil, fmt.Errorf("file upload không được để trống")
	}
	if post == nil {
		return nil, fmt.Errorf("post không được để trống")
	}

	post.ID = utils.NewID()
	log.Println("Đang tạo Bài đăng:", post.ID)

	needTransaction := transaction != nil

	// Trường hợp không kèm transaction: chỉ upload song song (1 task) rồi insert post thường
	if !needTransaction {
		uploadedURL, err := uc.uploadUsecase.ProcessUploadFile(ctx, fileToupload)
		if err != nil {
			return nil, fmt.Errorf("lỗi upload ảnh: %w", err)
		}
		post.ImageUrl = uploadedURL
		post.TransactionID = nil
		if err := post.Validate(); err != nil {
			_ = uc.uploadUsecase.DeleteFile(context.Background(), uploadedURL)
			return nil, err
		}
		if _, err := uc.postRepo.Create(ctx, post); err != nil {
			_ = uc.uploadUsecase.DeleteFile(context.Background(), uploadedURL)
			return nil, fmt.Errorf("lỗi khi lưu bài đăng: %w", err)
		}
		return &postDomain.PostResponse{
			ID:            post.ID,
			UserID:        post.UserID,
			ImageUrl:      post.ImageUrl,
			Caption:       post.Caption,
			Visibility:    post.Visibility,
			LocationName:  post.LocationName,
			TransactionID: nil,
		}, nil
	}

	// Có transaction: BEGIN Tx + song song Upload & CreateTransaction
	tx, err := uc.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var (
		mu            sync.Mutex
		uploadedURL   string
		transactionID string
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		log.Println("Đang upload ảnh lên Cloud...")
		url, err := uc.uploadUsecase.ProcessUploadFile(gCtx, fileToupload)
		if err != nil {
			return err
		}
		mu.Lock()
		uploadedURL = url
		mu.Unlock()
		return nil
	})

	g.Go(func() error {
		log.Println("Đang tạo transaction trong DB Tx...")
		resp, err := uc.transactionUsecase.CreateWithTx(gCtx, tx, transaction)
		if err != nil {
			return err
		}
		mu.Lock()
		transactionID = resp.ID
		mu.Unlock()
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Println("Lỗi tiến trình đồng thời:", err)
		mu.Lock()
		url := uploadedURL
		mu.Unlock()
		if url != "" {
			_ = uc.uploadUsecase.DeleteFile(context.Background(), url)
		}
		return nil, fmt.Errorf("lỗi trong quá trình xử lý song song: %w", err)
	}

	mu.Lock()
	post.ImageUrl = uploadedURL
	tID := transactionID
	post.TransactionID = &tID
	mu.Unlock()

	if err := post.Validate(); err != nil {
		_ = uc.uploadUsecase.DeleteFile(context.Background(), post.ImageUrl)
		return nil, err
	}

	if _, err := uc.postRepo.CreateWithTx(ctx, tx, post); err != nil {
		log.Println("Lỗi khi lưu bài đăng trong Tx:", err)
		_ = uc.uploadUsecase.DeleteFile(context.Background(), post.ImageUrl)
		return nil, fmt.Errorf("lỗi khi lưu bài đăng: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		_ = uc.uploadUsecase.DeleteFile(context.Background(), post.ImageUrl)
		return nil, fmt.Errorf("lỗi khi commit tx: %w", err)
	}

	return &postDomain.PostResponse{
		ID:            post.ID,
		UserID:        post.UserID,
		ImageUrl:      post.ImageUrl,
		Caption:       post.Caption,
		Visibility:    post.Visibility,
		LocationName:  post.LocationName,
		TransactionID: post.TransactionID,
	}, nil
}
