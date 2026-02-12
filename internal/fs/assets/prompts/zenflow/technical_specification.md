# Zenflow Step Guide: Technical Specification

Цель этапа: подготовить или обновить спецификацию изменений.

Обязательно:
1. Прочитать `.spec_agent/prompts/spec_rules.md`.
2. Сформировать `{@artifacts_path}/spec.md`.
3. Убедиться, что spec следует обязательной структуре и правилам ссылок.

Проверить в spec:
- marker block `SPEC:*`;
- секции в нужном порядке;
- блоки `Business Logic`, `Flow`, `Links`, `Dependencies`;
- ссылки только в markdown-формате, с якорем `#...`;
- язык спецификации: русский.

Критерий готовности:
- `spec.md` достаточно подробен для реализации без домысливания.
