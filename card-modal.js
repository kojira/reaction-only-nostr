(() => {
  const feed = document.getElementById('feed');
  const dialog = document.getElementById('cardDialog');
  const expanded = document.getElementById('expandedCard');
  const close = document.getElementById('closeCardDialog');
  let original = null;
  let marker = null;
  let opener = null;

  function open(card) {
    opener = document.activeElement;
    original = card;
    marker = document.createComment('expanded card');
    card.before(marker);
    expanded.appendChild(card);
    dialog.showModal();
    close.focus();
  }
  feed.addEventListener('click', event => {
    const card = event.target.closest('article');
    if (!card || !feed.contains(card) || event.target.closest('button, a, input, select, textarea, summary, label')) return;
    open(card);
  });
  feed.addEventListener('keydown', event => {
    if ((event.key === 'Enter' || event.key === ' ') && event.target.matches('article')) {
      event.preventDefault();
      open(event.target);
    }
  });
  // Timeline is rebuilt as events arrive. Give newly rendered cards a keyboard entry point.
  new MutationObserver(() => {
    feed.querySelectorAll('article:not([tabindex])').forEach(card => {
      card.tabIndex = 0;
      card.setAttribute('aria-label', 'リアクションの詳細を開く');
    });
  }).observe(feed, {childList: true});
  close.addEventListener('click', () => dialog.close());
  dialog.addEventListener('click', event => {
    if (event.target === dialog) dialog.close();
  });
  dialog.addEventListener('close', () => {
    if (marker && marker.isConnected) marker.replaceWith(original);
    else if (original) original.remove();
    if (opener && opener.isConnected) opener.focus();
    original = marker = opener = null;
  });
})();
