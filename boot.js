const go = new Go();
WebAssembly.instantiateStreaming(fetch('app.wasm'), go.importObject).then(result => go.run(result.instance)).catch(error => { document.getElementById('status').textContent = '起動エラー: ' + error; });
