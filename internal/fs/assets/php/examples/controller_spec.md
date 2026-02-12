<!-- SPEC:START -->
<!-- SPEC:FILE=true -->
<!-- SPEC:ID=app/Http/Controllers/PaymentController -->
<!-- SPEC:KIND=controller -->
<!-- SPEC:MENU=true -->
<!-- SPEC:END -->

# PaymentController

## Responsibility
Обрабатывает входящие HTTP-запросы по платежам и передает управление в бизнес-слой.

## Inputs
- Request данные оплаты

## Outputs
- Результат операции оплаты

## Business Logic
1. Контроллер не содержит бизнес-правил.
2. Валидация входа выполняется до вызова usecase/service.

## Flow
1. Принять HTTP запрос.
2. Провалидировать входные данные.
3. Вызвать usecase/service.
4. Вернуть результат клиенту.

## Links
- calls: [CreatePaymentUseCase](../../UseCases/CreatePaymentUseCase.md#execute)

## Dependencies
- [CreatePaymentUseCase](../../UseCases/CreatePaymentUseCase.md)

## Errors
- ValidationError
- BusinessRuleViolation
