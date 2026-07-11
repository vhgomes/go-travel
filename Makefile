.PHONY: help infra-init infra-plan infra-apply infra-destroy

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

infra-init:
	@cd infra && terraform init

infra-plan:
	@cd infra && terraform plan

infra-apply:
	@cd infra && terraform apply -auto-approve

infra-destroy:
	@cd infra && terraform destroy
