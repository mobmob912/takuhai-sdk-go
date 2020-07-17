.PHONY: build run clean kill commit-push

commit-push:
	git add -A
	git commit -m "update"
	git push origin master

build:
	dep ensure
	docker build -t tockn/takuhai-echo .

docker-push: build
	docker push tockn/takuhai-echo

clean:
	-docker rm takuhai-echo

run: clean
	docker run --name takuhai-echo --env managerAddr=192.168.0.140 --env takuhaiJobPort=3333 --env stepID=hoge --env workflowID=wkID -p 3333:3333 tockn/takuhai-echo

all: commit-push clean build run

kill:
	docker kill takuhai-echo
