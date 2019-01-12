export LD_LIBRARY_PATH=$GOPATH/src/github.com/loitd/vabackend/instantclient_12_1:$LD_LIBRARY_PATH
go get github.com/gorilla/mux
go get github.com/loitd/vabackend
go get gopkg.in/goracle.v2