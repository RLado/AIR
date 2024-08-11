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
	@mkdir -p ${DESTDIR}/config

	@cp LICENSE ${DESTDIR}
	@mv ${BINARY_NAME} ${DESTDIR}/bin
	@cp -r src/templates ${DESTDIR}/bin
	@cp -r src/static ${DESTDIR}/bin
	@cp -r config ${DESTDIR}