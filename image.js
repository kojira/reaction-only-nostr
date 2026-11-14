document.getElementById('generate').addEventListener('click', async () => {
  const status = document.getElementById('status');
  try {
    if (!window.nostr) throw Error('NIP-07拡張機能が必要です');
    const input = document.getElementById('imageText').value.trim();
    if (!input) throw Error('画像にする文字を入力してください');
    const canvas = document.createElement('canvas'); canvas.width = canvas.height = 256;
    const ctx = canvas.getContext('2d'); ctx.clearRect(0, 0, 256, 256);
    ctx.textAlign = 'center'; ctx.textBaseline = 'middle'; ctx.font = 'bold 96px sans-serif'; ctx.fillStyle = '#222'; ctx.fillText(input, 128, 128, 240);
    const blob = await new Promise(resolve => canvas.toBlob(resolve, 'image/png'));
    if (!blob) throw Error('画像の生成に失敗しました');
    const digest = await crypto.subtle.digest('SHA-256', await blob.arrayBuffer());
    const hash = Array.from(new Uint8Array(digest), x => x.toString(16).padStart(2, '0')).join('');
    const pubkey = await window.nostr.getPublicKey();
    const signed = await window.nostr.signEvent({kind:24242,pubkey,created_at:Math.floor(Date.now()/1000),tags:[['t','upload'],['expiration',String(Math.floor(Date.now()/1000)+300)],['x',hash],['server','blossom.primal.net']],content:'Upload generated emoji image'});
    const auth = btoa(unescape(encodeURIComponent(JSON.stringify(signed))));
    const response = await fetch('https://blossom.primal.net/upload',{method:'PUT',headers:{Authorization:'Nostr '+auth,'Content-Type':'image/png','X-SHA-256':hash},body:blob});
    if (!response.ok) throw Error('アップロード失敗: HTTP '+response.status);
    const result = await response.json();
    if (result.sha256 !== hash || !result.url || !result.url.startsWith('https://')) throw Error('サーバーの返答を検証できません');
    document.getElementById('emojiURL').value = result.url; status.textContent = '画像を公開しました。ショートコードを指定して登録できます。';
  } catch (e) { status.textContent = '画像生成・公開エラー: '+e.message; }
});
