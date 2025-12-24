# CreateUserUseCase

## Responsibility
Реализует бизнес-процесс создания нового пользователя со всеми необходимыми проверками и инициализацией.

## Inputs
- Email пользователя
- Пароль
- Имя пользователя

## Outputs
- Созданный объект User с ID
- Статус успешного создания

## Business Rules
1. Email должен быть уникальным
2. Пароль должен быть минимум 8 символов
3. Новый пользователь создается в неактивном состоянии
4. Отправляется email с ссылкой подтверждения
5. Максимум 10 попыток регистрации за час с одного IP

## Flow
1. Валидирует входные данные
   → calls: ../services/validation_service.md#ValidateEmail
2. Проверяет уникальность email
   → reads: ../repositories/user_repository.md#ExistsByEmail
3. Хеширует пароль
   → calls: ../services/crypto_service.md#HashPassword
4. Создает запись в репозитории
   → writes: ../repositories/user_repository.md#Create
5. Отправляет письмо подтверждения
   → calls: ../services/email_service.md#SendConfirmation
6. Возвращает созданного пользователя

## Dependencies
- [UserRepository](../repositories/user_repository.md)
- [ValidationService](../services/validation_service.md)
- [CryptoService](../services/crypto_service.md)
- [EmailService](../services/email_service.md)

## Errors
- ErrEmailExists — email уже зарегистрирован
- ErrInvalidEmail — некорректный формат email
- ErrWeakPassword — пароль не соответствует требованиям
- ErrRateLimitExceeded — превышен лимит попыток регистрации

## Notes
Usecase координирует работу различных сервисов для реализации бизнес-процесса.
Все проверки выполняются до создания пользователя в БД.
