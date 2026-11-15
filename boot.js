const go = new Go();
WebAssembly.instantiateStreaming(fetch('app.wasm?v=reaction-row1', {cache: 'no-cache'}), go.importObject).then(result => go.run(result.instance)).catch(error => { document.getElementById('status').textContent = '起動エラー: ' + error; });
