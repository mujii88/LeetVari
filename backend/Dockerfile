FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y no-install-recommends \
    verilog \
    yosys \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /sandbox


