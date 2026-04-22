SHELL := /bin/bash

name ?= new-post

.PHONY: dev build new publish

dev:
	hugo server -D

build:
	hugo --gc --minify

new:
	hugo new posts/$(name).md
	@echo "Created: content/posts/$(name).md"
	@echo "Remember to set draft = false before publish."

publish:
	git add -A
	git commit -m "publish: update blog content"
	git push origin master
