help: ## ヘルプを表示する
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# init: ## terraform init を実行
# 	@~/bin/terraform init

validate: ## tf ファイルを検証する
	@~/bin/terraform validate && printf "\e[32m***** validated ******\e[m\n"

plan: validate ## terraform plan を実行する
	@~/bin/terraform plan

apply: validate ## terraform apply を実行する
	@~/bin/terraform apply

show: ## terraform show を実行する
	@~/bin/terraform show

destroy: ## terraform destroy を実行する
	@~/bin/terraform destroy

refresh: ## terraform refresh を実行する
	@~/bin/terraform refresh
