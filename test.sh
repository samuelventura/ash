#!/bin/bash -xe

#trace|debug|info
TEST_PKG="${1:-all}"
TEST_SCOPE="${2:-Test}"
export ASH_LOGLEVEL="${3:-info}"
MOD="github.com/samuelventura/ash"
go clean -testcache 

case $TEST_PKG in
    all)
    go test $MOD/pkg/ash -v -run $TEST_SCOPE
    ;;
    ash)
    go test $MOD/pkg/ash -v -run $TEST_SCOPE
    ;;
esac
