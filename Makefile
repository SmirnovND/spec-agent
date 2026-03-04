.PHONY: help eval-openai-shortener-local eval-openai-shortener-local-down eval-sync-draft

help:
	@echo "Available targets:"
	@echo "  make eval-openai-shortener-local      - run OpenAI Codex shortener eval in isolated docker-compose"
	@echo "  make eval-openai-shortener-local-down - stop/remove local eval containers"
	@echo "  make eval-sync-draft                  - sync pinned draft fixture"

eval-sync-draft:
	bash eval/scripts/fetch_draft_fixture.sh

eval-openai-shortener-local: eval-sync-draft
	@mkdir -p eval/results/openai-shortener
	docker compose -f eval/openai/shortener/docker-compose.local.yml up --build --abort-on-container-exit --exit-code-from runner

eval-openai-shortener-local-down:
	docker compose -f eval/openai/shortener/docker-compose.local.yml down -v
