FROM scratch
COPY /ac /ac
ENTRYPOINT ["/ac"]
CMD ["-h"]
