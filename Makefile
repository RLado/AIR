BINARY_NAME:=air
DESTDIR:=build
 
build:
	@cd src && go build -ldflags "-s -w" -o ../${BINARY_NAME} main.go

clean:
	@go clean
	@rm -f ${BINARY_NAME}
	@rm -rf ${DESTDIR}

package: build
	@mkdir -p ${DESTDIR}/bin
	@mkdir -p ${DESTDIR}/templates
	@mkdir -p ${DESTDIR}/static
	@mkdir -p ${DESTDIR}/config

	@mv ${BINARY_NAME} ${DESTDIR}/bin
	@cp -r src/templates ${DESTDIR}/templates
	@cp -r src/static ${DESTDIR}/static
	@cp -r config ${DESTDIR}/config