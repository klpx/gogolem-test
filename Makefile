build: compile
	wasm-tools component embed ./wit gogolem_test.module.wasm --output gogolem_test.embed.wasm
	wasm-tools component new gogolem_test.embed.wasm -o gogolem_test.wasm --adapt adapters/tier2/wasi_snapshot_preview1.wasm

bindings:
	wit-bindgen tiny-go --out-dir gogolem_test ./wit

compile: bindings
	tinygo build -target=wasi -o gogolem_test.module.wasm main.go

clean:
	rm -rf gogolem_test
	rm *.wasm
