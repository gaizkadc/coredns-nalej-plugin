#
#  Copyright 2018 Nalej
# 

# Name of the target applications to be built
APPS=

publish:
	@echo "This component doesn't have a Docker image"

image:
	@echo "This component doesn't have a Docker image"

# Use global Makefile for common targets
export
%:
	$(MAKE) -f Makefile.golang $@
