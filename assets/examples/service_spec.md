# UserService

## Responsibility
Реализует бизнес-логику работы с пользователями и управляет их состоянием.

## Inputs
- Email пользователя
- Пароль (при регистрации)
- Параметры для поиска и фильтрации

## Outputs
- Объект User с данными пользователя
- Статус операции (создан, обновлен, удален)
- Список пользователей (при поиске)

## Business Rules
1. Email должен быть уникальным в системе
2. Пароль должен быть минимум 8 символов
3. Пароль хранится в хешированном виде
4. Удаленный пользователь не может быть восстановлен
5. Email активируется через подтверждение

## Flow
1. Валидирует входные данные
2. Проверяет уникальность email
   → calls: ../repositories/user_repository.md#ExistsByEmail
3. Хеширует пароль (если требуется)
   → calls: ../services/crypto_service.md#HashPassword
4. Создает или обновляет пользователя в репозитории
   → writes: ../repositories/user_repository.md
5. Возвращает результат операции

## Dependencies
- [UserRepository](../repositories/user_repository.md)
- [CryptoService](../services/crypto_service.md)

## Errors
- ErrEmailExists — email уже зарегистрирован
- ErrInvalidPassword — пароль не соответствует требованиям
- ErrUserNotFound — пользователь не найден

## Notes
Сервис содержит все бизнес-правила и координирует работу с хранилищем данных.
