
pack: pack.template-single pack.template-mono

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