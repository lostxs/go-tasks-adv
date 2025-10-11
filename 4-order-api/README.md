Нужно создать новый проект и подключить:

Конфиг, работы с res, валидация
Создать базу данных и подключиться к ней
Сделать модель продукта и провести миграции:
type Product struct {
    gorm.Model
    Name        string
    Description string
    Images      pq.StringArray
    ...
}