package media

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"
)

// CropSquare cắt một image.Image về tỷ lệ 1:1 chính giữa
func CropSquare(src image.Image) image.Image {
	bounds := src.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	// Nếu ảnh đã là hình vuông 1:1 thì giữ nguyên
	if w == h {
		return src
	}

	// 1. Xác định kích thước vuông (lấy theo cạnh ngắn hơn)
	size := w
	if h < w {
		size = h
	}

	// 2. Tính tọa độ góc trên-trái để cắt chính giữa ảnh
	x0 := bounds.Min.X + (w-size)/2
	y0 := bounds.Min.Y + (h-size)/2

	// 3. Tạo bức ảnh mới vuông kích thước (size x size)
	dst := image.NewNRGBA(image.Rect(0, 0, size, size))

	// 4. Vẽ vùng ảnh đã cắt từ src vào dst
	draw.Draw(dst, dst.Bounds(), src, image.Pt(x0, y0), draw.Src)

	return dst
}