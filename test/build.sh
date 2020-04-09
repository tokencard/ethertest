#!/bin/bash

set -e -o pipefail

SOLC="docker run --rm -u `id -u` -v $PWD:/solidity --workdir /solidity/contracts ethereum/solc:0.6.5 --optimize"

compile_solidity() {
  echo "compiling ${1}"
  ${SOLC} --overwrite --bin --abi ${1}.sol -o /solidity/build/${1} --combined-json bin-runtime,srcmap-runtime,ast,srcmap,bin
}

contract_sources=(
  'test'
)

for c in "${contract_sources[@]}"
do
    compile_solidity $c
done



# Generate Go bindings from solidity contracts.
ABIGEN="docker run --rm -u `id -u` --workdir /contracts -e GOPATH=/go -v $PWD:/contracts ethereum/client-go:alltools-v1.9.12 abigen"


generate_binding() {
  contract=$(echo $1 | awk '{print $1}')
  go_source=$(echo $1 | awk '{print $2}')
  go_type=$(echo $1 | awk '{print $3}')
  package=$(echo $1 | awk '{print $4}')
  echo "Generating binding for ${go_type} (${contract})"
  ${ABIGEN} --abi ./build/${contract}.abi  --bin ./build/${contract}.bin --pkg ${package} --type=${go_type} --out ./bindings/${go_source}
}

contracts=(
  "test/Test test_bindings.go Test bindings"
)

for c in "${contracts[@]}"
do
    generate_binding "$c"
done

echo "done"
