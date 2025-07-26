package extends

// DataList - обобщенная структура, содержащая какой-либо
// слайс значений и общее колличество объектов в базе,
// для пагинации
type DataList[T1 any] struct {
	Data  T1  `json:"data"`
	Count int `json:"count"`
}

// NewDataList - функция для создания структуры DataList
func NewDataList[T1 any](data T1, count int) DataList[T1] {
	return DataList[T1]{Data: data, Count: count}
}
