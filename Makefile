
all: gpd ppd


.PHONY: gpd
gpd: 
	@cd gpd && go build -v .

.PHONY: ppd
ppd:
	@cd ppd && go build -v . 

.PHONY: ppd-docker 
ppd-docker:
	docker build -f Dockerfile.ppd  -t ppd . 

clean:
	@cd gpd && go clean .
	@cd ppd && go clean . 
	
