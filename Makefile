
all: gpd ppd


.PHONY: gpd
gpd: 
	@cd gpd && go build -v .

.PHONY: ppd
ppd:
	@cd ppd && go build -v . 

.PHONY: ppd-docker 
ppd-docker:
	@cd ppd && docker build -t ppd . 

clean:
	@cd gpd && go clean .
	@cd ppd && go clean . 
	
