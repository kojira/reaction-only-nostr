(() => {
  const menu = document.getElementById('menuPanel');
  const toggle = document.getElementById('menuToggle');
  const select = document.getElementById('fontSize');
  const sizes = {'extra-small': '62.5%', smaller: '75%', small: '87.5%', medium: '100%', large: '112.5%', 'extra-large': '125%'};
  let saved;
  try { saved = localStorage.getItem('reaction-only-font-size'); } catch (_) {}
  select.value = Object.hasOwn(sizes, saved) ? saved : 'medium';
  const apply = () => { document.documentElement.style.fontSize = sizes[select.value]; };
  apply();
  select.addEventListener('change', () => {
    apply();
    try { localStorage.setItem('reaction-only-font-size', select.value); } catch (_) {}
  });
  toggle.addEventListener('click', () => {
    menu.hidden = !menu.hidden;
    toggle.setAttribute('aria-expanded', String(!menu.hidden));
    toggle.setAttribute('aria-label', menu.hidden ? 'メニューを開く' : 'メニューを閉じる');
  });
  document.addEventListener('keydown', e => {
    if (e.key === 'Escape' && !menu.hidden) { toggle.click(); toggle.focus(); }
  });
})();
