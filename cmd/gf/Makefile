.DEFAULT_GOAL := pack

pack: pack.template-single pack.template-mono pack.template-mono-app

pack.template-single:
	@rm -fr temp
	@mkdir temp || exit 0
	@cd temp && git clone https://github.com/gogf/template-single
	@rm -fr temp/template-single/.git
	@cd temp && gf pack template-single ../internal/packed/template-single.go -n=packed -y
	@rm -fr temp

pack.template-mono:
	@rm -fr temp
	@mkdir temp || exit 0
	@cd temp && git clone https://github.com/gogf/template-mono
	@rm -fr temp/template-mono/.git
	@cd temp && gf pack template-mono ../internal/packed/template-mono.go -n=packed -y
	@rm -fr temp

# Note:
# command `sed` only works on MacOS.
# use `grep -irl 'template-single' temp| xargs sed -i'' -e 's/template-single/template-mono-app/g'` on other platforms.
pack.template-mono-app:
	@rm -fr temp
	@mkdir temp || exit 0
	@cd temp && git clone https://github.com/gogf/template-single
	@cd temp && mv template-single template-mono-app
	@rm -fr temp/template-mono-app/.git
	@rm -fr temp/template-mono-app/.gitattributes
	@rm -fr temp/template-mono-app/.gitignore
	@rm -fr temp/template-mono-app/go.mod
	@rm -fr temp/template-mono-app/go.sum
	@grep -irl 'template-single' temp| xargs sed -i '' -e 's/template-single/template-mono-app/g'
	@cd temp && gf pack template-mono-app ../internal/packed/template-mono-app.go -n=packed -y
	@rm -fr temp