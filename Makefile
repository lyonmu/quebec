# VERSION = $(shell git describe --tags --abbrev=0)
AUTHOR = $(shell git log -1 --pretty=format:'%an')
BRANCH = $(shell git branch  --show-current)
OS = $(shell uname)

.PHONY: auth
auth: 
	@echo $(AUTHOR)
	@echo $(BRANCH)
	@echo $(OS)

.PHONY: clean
clean:
	rm -rf bin