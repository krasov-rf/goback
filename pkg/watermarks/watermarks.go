package watermark

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"golang.org/x/image/draw"
)

type WatermarkAlign int

const (
	WATERMARK_BOTTOM_RIGHT WatermarkAlign = iota
	WATERMARK_BOTTOM_LEFT
	WATERMARK_TOP_RIGHT
	WATERMARK_TOP_LEFT
	WATERMARK_DIAGONAL
	WATERMARK_CENTER
)

type Watermark struct {
	Font     string
	TextSize float64
	Color    color.RGBA
	Padding  float64
}

func (w *Watermark) Generate(img image.Image, text string, align WatermarkAlign) (image.Image, error) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	// Создаю контекст и прорисовываю в нем переданное изображение
	dc := gg.NewContext(width, height)
	dc.DrawImage(img, 0, 0)

	// Подгружаю необходимый шрифт, выставляю его размер и цвет
	err := dc.LoadFontFace(w.Font, 12)
	if err != nil {
		return nil, err
	}
	dc.SetColor(w.Color)

	textWidth, textHeight := dc.MeasureString(text)

	switch align {
	case WATERMARK_CENTER:
		// вычисляю масштаб шрифта относительно полотна
		scaleX := float64(width) / textWidth
		scaleY := float64(height) / textHeight

		// беру наименьший масштаб, для того что бы шрифт не вылазил по ней за пределы картинки
		scale := scaleX
		if scaleY < scaleX {
			scale = scaleY
		}

		// временно изменяю масштаб контекста, что бы правильно отрисовать текст
		dc.Push()
		dc.Scale(scale, scale)
		dc.DrawStringAnchored(text, 0, float64(dc.Height())/2/scale, 0, 0.5)
		dc.Pop()
	case WATERMARK_DIAGONAL:
		// Рисую по диагонали множество водяных знаков
		dc.Push()
		dc.Rotate(gg.Radians(-45))
		dc.MoveTo(float64(width+width/2), float64(height))

		spacing := int(textWidth) + 30
		width := img.Bounds().Dx()
		height := img.Bounds().Dy()
		for x := -width; x < width+150; x += spacing {
			for y := 0; y < height+height; y += spacing / 2 {
				dc.DrawStringAnchored(text, float64(x), float64(y), 0.5, 0.5)
			}
		}
	default:
		// Обработка четырех углов
		maxTextWidth := float64(width) * 0.5

		scale := 1.0
		if textWidth > maxTextWidth {
			scale = maxTextWidth / textWidth
		}

		_, newTextHeight := dc.MeasureString(text)

		dc.Push()
		dc.Scale(scale, scale)

		switch align {
		case WATERMARK_TOP_LEFT:
			dc.DrawStringAnchored(text, 0, 20, 0, 1)
		case WATERMARK_TOP_RIGHT:
			x := float64(img.Bounds().Dx()) / scale
			dc.DrawStringAnchored(text, x, 20, 1, 1)
		case WATERMARK_BOTTOM_LEFT:
			y := float64(img.Bounds().Dy())/scale - newTextHeight - 20
			dc.DrawStringAnchored(text, 10, y, 0, 1)
		case WATERMARK_BOTTOM_RIGHT:
			x := float64(img.Bounds().Dx())/scale - 20
			y := float64(img.Bounds().Dy())/scale - newTextHeight - 20
			dc.DrawStringAnchored(text, x, y, 1, 1)
		}
	}

	dc.Pop()

	return dc.Image(), nil
}

// копирует переданную картинку watermark на изображение
func (w *Watermark) GenerateWithWatermark(img image.Image, watermark image.Image, align WatermarkAlign) (image.Image, error) {
	// кэф масштабирования водянного знака
	scale := 0.3
	newW := int(float64(watermark.Bounds().Max.X) * scale)
	newH := int(float64(watermark.Bounds().Max.Y) * scale)

	// новое изображение для водяного знака с новым размером
	wmImgScaled := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.NearestNeighbor.Scale(wmImgScaled, wmImgScaled.Rect, watermark, watermark.Bounds(), draw.Over, nil)

	// новое изображение что и фоновое, с теми же размерами
	dst := image.NewRGBA(img.Bounds())

	// копируем фоновове изображение в новое
	draw.Draw(dst, img.Bounds(), img, image.Point{0, 0}, draw.Src)

	// положение верхнего левого угла для начала размещения
	var watermarkX, watermarkY int

	// рассчитываем позицию для водяного знака
	switch align {
	case WATERMARK_CENTER, WATERMARK_DIAGONAL:
		watermarkX = (img.Bounds().Max.X - wmImgScaled.Bounds().Max.X) / 2
		watermarkY = (img.Bounds().Max.Y - wmImgScaled.Bounds().Max.Y) / 2
	case WATERMARK_BOTTOM_RIGHT:
		watermarkX = img.Bounds().Max.X - wmImgScaled.Bounds().Max.X - 10
		watermarkY = img.Bounds().Max.Y - wmImgScaled.Bounds().Max.Y - 10
	case WATERMARK_BOTTOM_LEFT:
		watermarkX = 10
		watermarkY = img.Bounds().Max.Y - wmImgScaled.Bounds().Max.Y - 10
	case WATERMARK_TOP_RIGHT:
		watermarkX = img.Bounds().Max.X - wmImgScaled.Bounds().Max.X - 10
		watermarkY = 10
	case WATERMARK_TOP_LEFT:
		watermarkX = 10
		watermarkY = 10
	}

	// накладываем водяной знак
	draw.Draw(
		dst,
		image.Rect(
			watermarkX,
			watermarkY,
			watermarkX+wmImgScaled.Bounds().Dx(),
			watermarkY+wmImgScaled.Bounds().Dy(),
		),
		wmImgScaled,
		image.Point{0, 0},
		draw.Over,
	)

	return dst, nil
}
