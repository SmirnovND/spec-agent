# UserController

## Responsibility
Обрабатывает HTTP-запросы, связанные с управлением пользователями.

## Inputs
- HTTP-запрос (POST/GET/PUT/DELETE)
- Данные пользователя в теле запроса

## Outputs
- HTTP-ответ с результатом операции
- Код ответа (200, 400, 401, 500)

## Business Rules
1. Контроллер не содержит бизнес-логики
2. Все решения делегируются usecase
3. Входные данные валидируются перед передачей
4. Ошибки преобразуются в HTTP-коды

## Flow
1. Принимает HTTP-запрос
2. Парсит и валидирует входные данные
3. Вызывает соответствующий usecase
   → calls: ../usecases/create_user.md
4. Преобразует результат в HTTP-ответ
5. Возвращает ответ клиенту

## Dependencies
- [CreateUserUseCase](../usecases/create_user.md)
- [GetUserUseCase](../usecases/get_user.md)

## Errors
- ErrInvalidRequest — некорректные данные запроса
- ErrUnauthorized — отсутствует авторизация

## Notes
Контроллер является entry-point для HTTP-запросов.
