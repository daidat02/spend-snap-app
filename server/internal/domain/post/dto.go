package post

import "time"


type CreatePostRequest struct {
    Caption       *string `form:"caption"`
    
    // Đặt mặc định là "friends" nếu Client không truyền lên
    Visibility    string  `form:"visibility,default=friends"` 
    
    LocationName  *string `form:"location_name"`
    Amount        *int64  `form:"amount"`
    CategoryID    *string `form:"category_id"`
    Type          *string `form:"type"`
    Note          *string `form:"note"`
}

type PostResponse struct{
	ID string `json:"id"`
	UserID string `json:"user_id"`
	ImageUrl string `json:"image_url"`
	Caption *string `json:"caption,omitempty"`
	Visibility string `json:"visibility"`
	LocationName *string `json:"location_name,omitempty"`
	TransactionID *string `json:"transaction_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
