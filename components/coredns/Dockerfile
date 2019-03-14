#
# Copyright (C) 2019 Nalej Group - All Rights Reserved
#

FROM iron/go

ADD coredns /coredns

EXPOSE 53 53/udp
ENTRYPOINT ["/coredns"]