module github.com/user/go-devstack

go 1.24

require (
	github.com/BurntSushi/toml v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/magefile/mage v1.17.0
)

require (
	github.com/sabhiram/go-gitignore v0.0.0-20210923224102-525f6e181f06
	github.com/tree-sitter/go-tree-sitter v0.25.0
	github.com/tree-sitter/tree-sitter-go v0.23.4
	github.com/tree-sitter/tree-sitter-javascript v0.25.0
	github.com/tree-sitter/tree-sitter-typescript v0.23.2
)

require github.com/mattn/go-pointer v0.0.1 // indirect

replace github.com/tree-sitter/tree-sitter-go/bindings/go => github.com/tree-sitter/tree-sitter-go v0.23.4
