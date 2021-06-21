FROM k8deployment-base as build
WORKDIR /app
COPY CMakeLists.txt CMakeLists.txt
COPY go go
COPY libs libs
COPY src src
COPY main.c main.c
RUN cd go && go test && go build -buildmode=c-shared -ldflags="-w -s" -gcflags=all=-l -gcflags=all=-B -o ../libs/kubernetes/libkubernetes.so libkubernetes.go
RUN mkdir build && cd build && cmake .. && cmake --build . && staticx k8deployment k8deployment-static
#RUN cd build && ldd k8deployment | tr -s '[:blank:]' '\n' | grep '^/' | xargs -I % sh -c 'mkdir -p $(dirname deps%); cp % deps%;'

FROM gcr.io/distroless/static
#gcr.io/distroless/base
#gcr.io/distroless/static
#COPY --from=build /app/build/deps /
COPY --from=build /app/build/k8deployment-static  /k8deployment
CMD ["./k8deployment"]