# Сервис для подбора одежды

## Идея

Приложение создано с целью создавать подборки одежды, где пользователь сможет сохранять, и просматривать свои пресеты одежды. Сохраняя свой стиль для разных сезонов, а также будет добавлена функция коммерческой подборки пресетов от стилистов.

## Дополнительная информация

Сервис представлен в каестве API.

## Техническая информация

Приложение было покрыто логами
Реальзован github CI workflow
Реализована авторизация jwt
Открытые пути

1. /register
2. /login
3. /refresh
4. /logout

Пользовательские маршруты защищены авторизацией

### Структурная схема базы данных

![Alt text](./screenshots/img.png)

### Маршруты и методы

| Route | Method | Description |
| ----- | ------ | ----------- |
| +text in green localhost:users/register | POST | Register user |
| +text in green localhost:users/login | POST | Login user |
| -text in red localhost:users/:id/presets | POST | Create preset |
| -text in red localhost:users/:id/presets | GET | Get all user`s presets |
| -text in red localhost:users/:id/presets/:id/clothes | POST | Create clothes |
| -text in red localhost:users/:id/presets/:id/clothes | GET | Get clothes(with season able with filter ) |

## Инструкция запуска

1. Для запуска необходимо сказать docker-compose, go(v22.0)
2. Устанавливаем ```go install github.com/Angstreminus/ClothersSelector```
3. Переходим в директорию установленного приложения и запускаем приложения ```go run cmd/main.go``` (приложение пока не готово)

## Todo

1. Swagger doc
2. Интеграционные тесты + добавить покрытие юнитами
